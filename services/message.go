package services

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	bot "github.com/MixinNetwork/bot-api-go-client"
	number "github.com/MixinNetwork/go-number"
	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/fox-one/mixin-sdk"
	"github.com/gofrs/uuid"
)

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
	user        *mixin.User
	bc          chan WsBroadcastMessage
	recipientID map[string]time.Time
}

func (mc *MessageContext) OnMessage(ctx context.Context, msg *mixin.MessageView, userID string) error {
	if msg.Category == "SYSTEM_ACCOUNT_SNAPSHOT" && msg.UserID != config.AppConfig.Mixin.ClientId {
		data, err := base64.StdEncoding.DecodeString(msg.Data)
		if err != nil {
			return session.BlazeServerError(ctx, err)
		}
		var transfer TransferView
		err = json.Unmarshal(data, &transfer)
		if err != nil {
			return session.BlazeServerError(ctx, err)
		}
		err = handleTransfer(ctx, mc, transfer, msg.UserID)
		if err != nil {
			return session.BlazeServerError(ctx, err)
		}
	} else if msg.ConversationID == models.UniqueConversationId(config.AppConfig.Mixin.ClientId, msg.UserID) {
		if err := handleMessage(ctx, mc, msg, mc.bc); err != nil {
			return err
		}
	}

	return nil
}

func (mc *MessageContext) OnBlazeMessage(ctx context.Context, message *mixin.BlazeMessage, userID string) error {
	if message.Action != "ACKNOWLEDGE_MESSAGE_RECEIPT" {
		return nil
	}

	var msg mixin.MessageView
	if err := json.Unmarshal(message.Data, &msg); err != nil {
		session.Logger(ctx).Error("ACKNOWLEDGE_MESSAGE_RECEIPT json.Unmarshal", err)
		return nil
	}

	if msg.Status != "READ" {
		return nil
	}

	id, err := models.FindDistributedMessageRecipientId(ctx, msg.MessageID)
	if err != nil {
		session.Logger(ctx).Error("ACKNOWLEDGE_MESSAGE_RECEIPT FindDistributedMessageRecipientId", err)
		return nil
	}

	if id == "" {
		return nil
	}

	if time.Since(mc.recipientID[id]) > models.UserActivePeriod {
		if err := models.PingUserActiveAt(ctx, id); err != nil {
			session.Logger(ctx).Error("ACKNOWLEDGE_MESSAGE_RECEIPT PingUserActiveAt", err)
		}
		mc.recipientID[id] = time.Now()
	}

	return nil
}

type TransferMemoInst struct {
	Action string `json:"a"`
	Param1 string `json:"p1"`
	Param2 string `json:"p2"`
}

func (service *MessageService) Run(ctx context.Context, broadcastChan chan WsBroadcastMessage) error {
	go distribute(ctx)
	go loopPendingMessage(ctx)
	go handlePendingParticipants(ctx)
	go handleExpiredPackets(ctx)
	go schedulePluginCronJob(ctx)

	user, err := mixin.NewUser(
		config.AppConfig.Mixin.ClientId,
		config.AppConfig.Mixin.SessionId,
		config.AppConfig.Mixin.SessionKey,
	)
	if err != nil {
		panic(err)
	}

	mc := &MessageContext{
		user:        user,
		bc:          broadcastChan,
		recipientID: map[string]time.Time{},
	}

	for {
		b := mixin.NewBlazeClient(user)
		if err := b.Loop(ctx, mc); err != nil {
			session.Logger(ctx).Error(err)
		}
		session.Logger(ctx).Info("connection loop end")
		time.Sleep(300 * time.Millisecond)
	}
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
		log.Println("Incorrect inst, fallback: ", transfer.TraceId, transfer.Memo, err)
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
	memo := "Rewards from " + strconv.FormatInt(user.IdentityNumber, 10)
	log.Println("Rewards from " + user.FullName + " to " + targetUser.UserId)
	var traceID string
	traceID, err = generateRewardTraceID(transfer.TraceId)
	if err != nil {
		return errors.New("generate trace id failed")
	}
	in := &bot.TransferInput{
		AssetId:     transfer.AssetId,
		RecipientId: targetUser.UserId,
		Amount:      number.FromString(transfer.Amount),
		TraceId:     traceID,
		Memo:        memo,
	}

	if err := bot.CreateTransfer(ctx, in, config.AppConfig.Mixin.ClientId, config.AppConfig.Mixin.SessionId, config.AppConfig.Mixin.SessionKey, config.AppConfig.Mixin.SessionAssetPIN, config.AppConfig.Mixin.PinToken); err != nil {
		log.Println("can't transfer to recipient", err)
		return err
	}

	if user.UserId != targetUser.UserId {
		if err := models.CreateTip(ctx, user.UserId, targetUser.UserId, transfer.AssetId, transfer.Amount, traceID, transfer.CreatedAt); err != nil {
			log.Println("can't record tip", err)
			return err
		}

		if err := models.CreateRewardsMessage(ctx, user, targetUser, transfer.Amount, inst.Param2); err != nil {
			log.Println("can't create rewards message", err)
			return err
		}
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
				session.Logger(ctx).Infof("REFUND ERROR %v, %v\n", id, err)
				break
			}
			if packet != nil {
				session.Logger(ctx).Infof("REFUND %v\n", id)
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

func handleMessage(ctx context.Context, mc *MessageContext, message *mixin.MessageView, broadcastChan chan WsBroadcastMessage) error {
	user, err := models.FindUser(ctx, message.UserID)
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
		return sendTextMessage(ctx, mc, message.ConversationID, config.AppConfig.MessageTemplate.MessageTipsUnsubscribe)
	}
	dataBytes, err := base64.StdEncoding.DecodeString(message.Data)
	if err != nil {
		return session.BadDataError(ctx)
	} else if len(dataBytes) < 10 {
		if strings.ToUpper(string(dataBytes)) == config.AppConfig.MessageTemplate.MessageCommandsInfo {
			if count, err := models.SubscribersCount(ctx); err != nil {
				return err
			} else {
				return sendTextMessage(ctx, mc, message.ConversationID, fmt.Sprintf(config.AppConfig.MessageTemplate.MessageCommandsInfoResp, count))
			}
		}
	}
	// broadcast
	if isBroadcastOn, err := models.ReadBroadcastProperty(ctx); err == nil && isBroadcastOn == "on" {
		go func() {
			bmsg, err := decodeMessage(ctx, user, message)
			if err == nil {
				broadcastChan <- bmsg
			}
		}()
	}
	if _, err := models.CreateMessage(ctx, user, message.MessageID, message.Category, message.QuoteMessageID, message.Data, message.CreatedAt, message.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func sendHelpMessge(ctx context.Context, user *models.User, mc *MessageContext, message *mixin.MessageView) error {
	if err := sendTextMessage(ctx, mc, message.ConversationID, config.AppConfig.MessageTemplate.MessageTipsHelp); err != nil {
		return err
	}
	if err := sendAppButton(ctx, mc, config.AppConfig.MessageTemplate.MessageTipsHelpBtn, message.ConversationID, config.AppConfig.Service.HTTPResourceHost); err != nil {
		return err
	}
	return nil
}

func decodeMessage(ctx context.Context, user *models.User, message *mixin.MessageView) (WsBroadcastMessage, error) {
	var bmsg WsBroadcastMessage
	bmsg.Category = message.Category
	bmsg.MessageId = message.MessageID
	bmsg.CreatedAt = message.UpdatedAt
	bmsg.Data = message.Data
	bmsg.SpeakerId = user.UserId
	bmsg.SpeakerName = user.FullName
	bmsg.SpeakerAvatar = user.AvatarURL

	if message.Category == "PLAIN_TEXT" {
		bytes, _ := base64.StdEncoding.DecodeString(message.Data)
		bmsg.Text = string(bytes)
		return bmsg, nil
	}

	if message.Category != "PLAIN_IMAGE" && message.Category != "PLAIN_VIDEO" && message.Category != "PLAIN_AUDIO" && message.Category != "PLAIN_DATA" {
		return bmsg, nil
	}

	data, err := base64.StdEncoding.DecodeString(message.Data)
	if err != nil {
		log.Println("message data decode error", err)
		return bmsg, err
	}

	att, err := attachmentFromMixinJSON(string(data))
	if err != nil {
		log.Println("decode attachment error", err)
		return bmsg, err
	}
	attResp, err := bot.AttachemntShow(ctx, config.AppConfig.Mixin.ClientId, config.AppConfig.Mixin.SessionId, config.AppConfig.Mixin.SessionKey, att.ID)
	if err != nil {
		log.Println("get attachment details error", err)
	}
	att.ViewUrl = attResp.ViewURL
	bmsg.Attachment = att
	return bmsg, nil
}

func attachmentFromMixinJSON(jsonString string) (att WsBroadcastMessageAttachment, err error) {
	var data struct {
		ID        string  `json:"attachment_id"`
		Size      int     `json:"size"`
		MimeType  string  `json:"mime_type"`
		Name      *string `json:"name"`
		Duration  *uint   `json:"duration"`
		Waveform  *string `json:"waveform"`
		Width     *uint   `json:"width"`
		Height    *uint   `json:"height"`
		Thumbnail *string `json:"thumbnail"`
	}
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return
	}

	att.ID = data.ID
	att.Size = data.Size
	att.MimeType = data.MimeType
	att.Duration = data.Duration
	if data.Waveform != nil {
		att.Waveform, err = base64.StdEncoding.DecodeString(*data.Waveform)
		if err != nil {
			return
		}
	}
	att.Name = data.Name
	att.Width = data.Width
	att.Height = data.Height
	if data.Thumbnail != nil {
		att.Thumbnail, err = base64.StdEncoding.DecodeString(*data.Thumbnail)
		if err != nil {
			return
		}
	}
	return
}

func generateRewardTraceID(originTraceID string) (string, error) {
	h := md5.New()
	io.WriteString(h, originTraceID)
	io.WriteString(h, "REWARD")
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	id, err := uuid.FromBytes(sum)
	return id.String(), err
}
