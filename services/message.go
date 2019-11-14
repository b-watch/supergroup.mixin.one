package services

import (
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	bot "github.com/MixinNetwork/bot-api-go-client"
	number "github.com/MixinNetwork/go-number"
	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/gorilla/websocket"
)

const (
	keepAlivePeriod = 20 * time.Second
	writeWait       = 15 * time.Second
	pongWait        = 10 * time.Second
	pingPeriod      = (pongWait * 9) / 10
)

type BlazeMessage struct {
	Id     string                 `json:"id"`
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
	Data   interface{}            `json:"data,omitempty"`
	Error  *session.Error         `json:"error,omitempty"`
}

type MessageView struct {
	ConversationId string    `json:"conversation_id"`
	UserId         string    `json:"user_id"`
	MessageId      string    `json:"message_id"`
	QuoteMessageId string    `json:"quote_message_id"`
	Category       string    `json:"category"`
	Data           string    `json:"data"`
	Status         string    `json:"status"`
	Source         string    `json:"source"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TransferView struct {
	Type          string    `json:"type"`
	SnapshotId    string    `json:"snapshot_id"`
	CounterUserId string    `json:"counter_user_id"`
	AssetId       string    `json:"asset_id"`
	Amount        string    `json:"amount"`
	TraceId       string    `json:"trace_id"`
	Memo          string    `json:"memo"`
	CreatedAt     time.Time `json:"created_at"`
}

type MessageService struct{}

type MessageContext struct {
	Transactions   *tmap
	ReadDone       chan bool
	WriteDone      chan bool
	DistributeDone chan bool
	ReadBuffer     chan MessageView
	WriteBuffer    chan []byte
	RecipientId    map[string]time.Time
}
type TransferMemoInst struct {
	Action string `json:"a"`
	Param1 string `json:"p1"`
	Param2 string `json:"p2"`
}

func (service *MessageService) Run(ctx context.Context) error {
	go distribute(ctx)
	go loopPendingMessage(ctx)
	go handlePendingParticipants(ctx)
	go handleExpiredPackets(ctx)
	go schedulePluginCronJob(ctx)

	for {
		err := service.loop(ctx)
		if err != nil {
			session.Logger(ctx).Error(err)
		}
		session.Logger(ctx).Info("connection loop end")
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func (service *MessageService) loop(ctx context.Context) error {
	conn, err := ConnectMixinBlaze(config.AppConfig.Mixin.ClientId, config.AppConfig.Mixin.SessionId, config.AppConfig.Mixin.SessionKey)
	if err != nil {
		return err
	}
	defer conn.Close()

	mc := &MessageContext{
		Transactions:   newTmap(),
		ReadDone:       make(chan bool, 1),
		WriteDone:      make(chan bool, 1),
		DistributeDone: make(chan bool, 1),
		ReadBuffer:     make(chan MessageView, 102400),
		WriteBuffer:    make(chan []byte, 102400),
		RecipientId:    make(map[string]time.Time, 0),
	}

	go writePump(ctx, conn, mc)
	go readPump(ctx, conn, mc)

	err = writeMessageAndWait(ctx, mc, "LIST_PENDING_MESSAGES", nil)
	if err != nil {
		return session.BlazeServerError(ctx, err)
	}

	for {
		select {
		case <-mc.ReadDone:
			return nil
		case msg := <-mc.ReadBuffer:
			if msg.Category == "SYSTEM_ACCOUNT_SNAPSHOT" && msg.UserId != config.AppConfig.Mixin.ClientId {
				data, err := base64.StdEncoding.DecodeString(msg.Data)
				if err != nil {
					return session.BlazeServerError(ctx, err)
				}
				var transfer TransferView
				err = json.Unmarshal(data, &transfer)
				if err != nil {
					return session.BlazeServerError(ctx, err)
				}
				err = handleTransfer(ctx, mc, transfer, msg.UserId)
				if err != nil {
					return session.BlazeServerError(ctx, err)
				}
			} else if msg.ConversationId == models.UniqueConversationId(config.AppConfig.Mixin.ClientId, msg.UserId) {
				if err := handleMessage(ctx, mc, &msg); err != nil {
					return err
				}
			}

			params := map[string]interface{}{"message_id": msg.MessageId, "status": "READ"}
			err = writeMessageAndWait(ctx, mc, "ACKNOWLEDGE_MESSAGE_RECEIPT", params)
			if err != nil {
				return session.BlazeServerError(ctx, err)
			}
		}
	}
}

func readPump(ctx context.Context, conn *websocket.Conn, mc *MessageContext) error {
	defer func() {
		conn.Close()
		mc.WriteDone <- true
		mc.ReadDone <- true
		mc.DistributeDone <- true
	}()
	conn.SetReadLimit(1024000 * 128)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return session.BlazeServerError(ctx, err)
		}
		return nil
	})
	for {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return session.BlazeServerError(ctx, err)
		}
		messageType, wsReader, err := conn.NextReader()
		if err != nil {
			return session.BlazeServerError(ctx, err)
		}
		if messageType != websocket.BinaryMessage {
			return session.BlazeServerError(ctx, fmt.Errorf("invalid message type %d", messageType))
		}
		err = parseMessage(ctx, mc, wsReader)
		if err != nil {
			return session.BlazeServerError(ctx, err)
		}
	}
}

func writePump(ctx context.Context, conn *websocket.Conn, mc *MessageContext) error {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		conn.Close()
	}()
	for {
		select {
		case data := <-mc.WriteBuffer:
			err := writeGzipToConn(ctx, conn, data)
			if err != nil {
				return session.BlazeServerError(ctx, err)
			}
		case <-mc.WriteDone:
			return nil
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return session.BlazeServerError(ctx, err)
			}
		}
	}
}

func writeMessageAndWait(ctx context.Context, mc *MessageContext, action string, params map[string]interface{}) error {
	var resp = make(chan BlazeMessage, 1)
	var id = bot.UuidNewV4().String()
	mc.Transactions.set(id, func(t BlazeMessage) error {
		select {
		case resp <- t:
		case <-time.After(2 * time.Second):
			return fmt.Errorf("timeout to hook %s %s", action, id)
		}
		return nil
	})

	blazeMessage, err := json.Marshal(BlazeMessage{Id: id, Action: action, Params: params})
	if err != nil {
		return err
	}
	select {
	case <-time.After(keepAlivePeriod):
		return fmt.Errorf("timeout to write %s %v", action, params)
	case mc.WriteBuffer <- blazeMessage:
	}

	select {
	case <-time.After(keepAlivePeriod):
		mc.Transactions.retrive(id)
		return fmt.Errorf("timeout to wait %s %v", action, params)
	case t := <-resp:
		if t.Error != nil && t.Error.Code != 403 {
			mc.Transactions.retrive(id)
			return writeMessageAndWait(ctx, mc, action, params)
		}
	}
	return nil
}

func writeGzipToConn(ctx context.Context, conn *websocket.Conn, msg []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	wsWriter, err := conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return session.BlazeServerError(ctx, err)
	}
	gzWriter, err := gzip.NewWriterLevel(wsWriter, 3)
	if err != nil {
		return session.BlazeServerError(ctx, err)
	}
	if _, err := gzWriter.Write(msg); err != nil {
		return session.BlazeServerError(ctx, err)
	}

	if err := gzWriter.Close(); err != nil {
		return session.BlazeServerError(ctx, err)
	}
	if err := wsWriter.Close(); err != nil {
		return session.BlazeServerError(ctx, err)
	}
	return nil
}

func parseMessage(ctx context.Context, mc *MessageContext, wsReader io.Reader) error {
	var message BlazeMessage
	gzReader, err := gzip.NewReader(wsReader)
	if err != nil {
		return session.BlazeServerError(ctx, err)
	}
	defer gzReader.Close()
	if err = json.NewDecoder(gzReader).Decode(&message); err != nil {
		return session.BlazeServerError(ctx, err)
	}

	transaction := mc.Transactions.retrive(message.Id)
	if transaction != nil {
		return transaction(message)
	}

	if message.Action == "ACKNOWLEDGE_MESSAGE_RECEIPT" {
		data, err := json.Marshal(message.Data)
		if err != nil {
			session.Logger(ctx).Error("ACKNOWLEDGE_MESSAGE_RECEIPT", err)
			return nil
		}
		var msg MessageView
		err = json.Unmarshal(data, &msg)
		if err != nil {
			session.Logger(ctx).Error("ACKNOWLEDGE_MESSAGE_RECEIPT json.Unmarshal", err)
			return nil
		}
		if msg.Status != "READ" {
			return nil
		}
		id, err := models.FindDistributedMessageRecipientId(ctx, msg.MessageId)
		if err != nil {
			session.Logger(ctx).Error("ACKNOWLEDGE_MESSAGE_RECEIPT FindDistributedMessageRecipientId", err)
			return nil
		}
		if id == "" {
			return nil
		}
		if mc.RecipientId[id].Before(time.Now().Add(-1 * models.UserActivePeriod)) {
			err = models.PingUserActiveAt(ctx, id)
			if err != nil {
				session.Logger(ctx).Error("ACKNOWLEDGE_MESSAGE_RECEIPT PingUserActiveAt", err)
			}
			mc.RecipientId[id] = time.Now()
		}
		return nil
	}

	if message.Action != "CREATE_MESSAGE" {
		return nil
	}

	data, err := json.Marshal(message.Data)
	if err != nil {
		return session.BlazeServerError(ctx, err)
	}
	var msg MessageView
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return session.BlazeServerError(ctx, err)
	}

	select {
	case <-time.After(keepAlivePeriod):
		return fmt.Errorf("timeout to handle %s %s", msg.Category, msg.MessageId)
	case mc.ReadBuffer <- msg:
	}
	return nil
}

func handleTransfer(ctx context.Context, mc *MessageContext, transfer TransferView, userId string) error {
	id, err := bot.UuidFromString(transfer.TraceId)
	if err != nil {
		return nil
	}
	user, err := models.FindUser(ctx, userId)
	if user == nil || err != nil {
		log.Println("No such a user", userId)
		return err
	}
	if inst, err := crackTransferProtocol(ctx, mc, transfer, user); err == nil && inst.Action != "" {
		if inst.Action == "rewards" {
			return handleRewardsPayment(ctx, mc, transfer, user, inst)
		} else {
			log.Println("Unknown instruction", inst)
		}
	} else {
		log.Println("Incorrect inst or err", transfer.TraceId, transfer.Memo, err)
		if user.TraceId == transfer.TraceId {
			log.Println("New legacy payment", userId, transfer.TraceId)
			if transfer.Amount == config.AppConfig.System.PaymentAmount && transfer.AssetId == config.AppConfig.System.PaymentAssetId {
				return user.Payment(ctx)
			}
			for _, asset := range config.AppConfig.System.AccpetPaymentAssetList {
				if number.FromString(transfer.Amount).Equal(number.FromString(asset.Amount).RoundFloor(8)) && transfer.AssetId == asset.AssetId {
					return user.Payment(ctx)
				}
			}
		} else if order, err := models.GetOrder(ctx, transfer.TraceId); err == nil && order != nil {
			log.Println("New order received", userId, transfer.TraceId)
			return handleOrderPayment(ctx, mc, transfer, order)
		} else if packet, err := models.PayPacket(ctx, id.String(), transfer.AssetId, transfer.Amount); err != nil || packet == nil {
			log.Println("New packet paid", userId, transfer.TraceId, id)
			return err
		} else if packet.State == models.PacketStatePaid {
			log.Println("New packet prepared", userId, transfer.TraceId, packet.PacketId)
			return sendAppCard(ctx, mc, packet)
		}
	}
	return nil
}

func crackTransferProtocol(ctx context.Context, mc *MessageContext, transfer TransferView, user *models.User) (*TransferMemoInst, error) {
	var data *TransferMemoInst
	err := json.Unmarshal([]byte(transfer.Memo), &data)
	return data, err
}

func handleRewardsPayment(ctx context.Context, mc *MessageContext, transfer TransferView, user *models.User, inst *TransferMemoInst) error {
	userId := inst.Param1
	targetUser, err := models.FindUser(ctx, userId)
	if err != nil {
		log.Println("can't find user to reward", userId, err)
		return nil
	}
	log.Println("Rewards to ", userId)
	traceID := bot.UuidNewV4().String()
	in := &bot.TransferInput{
		AssetId:     transfer.AssetId,
		RecipientId: targetUser.UserId,
		Amount:      number.FromString(transfer.Amount),
		TraceId:     traceID,
		Memo:        "Rewards from " + user.FullName,
	}

	if user.UserId != targetUser.UserId {
		if err := models.CreateTip(ctx, user.UserId, targetUser.UserId, transfer.AssetId, transfer.Amount, traceID); err != nil {
			log.Println("can't record tip", err)
			return err
		}

		if err := models.CreateRewardsMessage(ctx, user, targetUser, transfer.Amount, inst.Param2); err != nil {
			log.Println("can't create rewards message", err)
			return err
		}
	}

	if err := bot.CreateTransfer(ctx, in, config.AppConfig.Mixin.ClientId, config.AppConfig.Mixin.SessionId, config.AppConfig.Mixin.SessionKey, config.AppConfig.Mixin.SessionAssetPIN, config.AppConfig.Mixin.PinToken); err != nil {
		log.Println("can't transfer to recipient", err)
		return err
	}
	return nil
}

func handleOrderPayment(ctx context.Context, mc *MessageContext, transfer TransferView, order *models.Order) error {
	if order.PayMethod == models.PayMethodMixin &&
		number.FromString(transfer.Amount).Equal(number.FromString(order.Amount).RoundFloor(8)) &&
		order.AssetId == transfer.AssetId {
		_, err := models.MarkOrderAsPaidByOrderId(ctx, order.OrderId)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func sendAppCard(ctx context.Context, mc *MessageContext, packet *models.Packet) error {
	description := fmt.Sprintf(config.AppConfig.MessageTemplate.GroupRedPacketDesc, packet.User.FullName)
	if strings.TrimSpace(packet.User.FullName) == "" {
		description = config.AppConfig.MessageTemplate.GroupRedPacketShortDesc
	}
	if count := utf8.RuneCountInString(description); count > 100 {
		name := string([]rune(packet.User.FullName)[:16])
		description = fmt.Sprintf(config.AppConfig.MessageTemplate.GroupRedPacketDesc, name)
	}
	host := config.AppConfig.Service.HTTPResourceHost
	if config.AppConfig.System.RouterMode == config.RouterModeHash {
		host = host + config.RouterModeHashSymbol
	}
	card, err := json.Marshal(map[string]string{
		"icon_url":    "https://images.mixin.one/X44V48LK9oEBT3izRGKqdVSPfiH5DtYTzzF0ch5nP-f7tO4v0BTTqVhFEHqd52qUeuVas-BSkLH1ckxEI51-jXmF=s256",
		"title":       config.AppConfig.MessageTemplate.GroupRedPacket,
		"description": description,
		"action":      host + "/packets/" + packet.PacketId,
	})
	if err != nil {
		return session.BlazeServerError(ctx, err)
	}
	t := time.Now()
	u := &models.User{UserId: config.AppConfig.Mixin.ClientId, ActiveAt: time.Now()}
	_, err = models.CreateMessage(ctx, u, packet.PacketId, models.MessageCategoryAppCard, "", base64.StdEncoding.EncodeToString(card), t, t)
	if err != nil {
		return session.BlazeServerError(ctx, err)
	}
	return nil
}

func handleExpiredPackets(ctx context.Context) {
	var limit = 100
	for {
		packetIds, err := models.ListExpiredPackets(ctx, limit)
		if err != nil {
			session.Logger(ctx).Error(err)
			time.Sleep(300 * time.Millisecond)
			continue
		}

		for _, id := range packetIds {
			packet, err := models.SendPacketRefundTransfer(ctx, id)
			if err != nil {
				session.Logger(ctx).Error(id, err)
				break
			}
			if packet != nil {
				session.Logger(ctx).Infof("REFUND %v", packet)
			}
		}

		if len(packetIds) < limit {
			time.Sleep(300 * time.Millisecond)
			continue
		}
	}
}

func schedulePluginCronJob(ctx context.Context) {
	plugin.RunCron()
}

func handlePendingParticipants(ctx context.Context) {
	var limit = 100
	for {
		participants, err := models.ListPendingParticipants(ctx, limit)
		if err != nil {
			session.Logger(ctx).Error(err)
			time.Sleep(300 * time.Millisecond)
			continue
		}

		for _, p := range participants {
			err = models.SendParticipantTransfer(ctx, p.PacketId, p.UserId, p.Amount)
			if err != nil {
				session.Logger(ctx).Error(err)
				break
			}
		}

		if len(participants) < limit {
			time.Sleep(300 * time.Millisecond)
			continue
		}
	}
}

func handleMessage(ctx context.Context, mc *MessageContext, message *MessageView) error {
	user, err := models.FindUser(ctx, message.UserId)
	if err != nil {
		return err
	}
	if user == nil || user.State != models.PaymentStatePaid {
		return sendHelpMessge(ctx, user, mc, message)
	}
	if user.ActiveAt.Before(time.Now().Add(-1 * models.UserActivePeriod)) {
		err = models.PingUserActiveAt(ctx, user.UserId)
		if err != nil {
			session.Logger(ctx).Error("handleMessage PingUserActiveAt", err)
		}
	}
	if user.SubscribedAt.IsZero() {
		return sendTextMessage(ctx, mc, message.ConversationId, config.AppConfig.MessageTemplate.MessageTipsUnsubscribe)
	}
	dataBytes, err := base64.StdEncoding.DecodeString(message.Data)
	if err != nil {
		return session.BadDataError(ctx)
	} else if len(dataBytes) < 10 {
		if strings.ToUpper(string(dataBytes)) == config.AppConfig.MessageTemplate.MessageCommandsInfo {
			if count, err := models.SubscribersCount(ctx); err != nil {
				return err
			} else {
				return sendTextMessage(ctx, mc, message.ConversationId, fmt.Sprintf(config.AppConfig.MessageTemplate.MessageCommandsInfoResp, count))
			}
		}
	}
	if _, err := models.CreateMessage(ctx, user, message.MessageId, message.Category, message.QuoteMessageId, message.Data, message.CreatedAt, message.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func sendHelpMessge(ctx context.Context, user *models.User, mc *MessageContext, message *MessageView) error {
	if err := sendTextMessage(ctx, mc, message.ConversationId, config.AppConfig.MessageTemplate.MessageTipsHelp); err != nil {
		return err
	}
	if err := sendAppButton(ctx, mc, config.AppConfig.MessageTemplate.MessageTipsHelpBtn, message.ConversationId, config.AppConfig.Service.HTTPResourceHost); err != nil {
		return err
	}
	return nil
}

type tmap struct {
	mutex sync.Mutex
	m     map[string]mixinTransaction
}

type mixinTransaction func(BlazeMessage) error

func newTmap() *tmap {
	return &tmap{
		m: make(map[string]mixinTransaction),
	}
}

func (m *tmap) retrive(key string) mixinTransaction {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	defer delete(m.m, key)
	return m.m[key]
}

func (m *tmap) set(key string, t mixinTransaction) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m[key] = t
}
