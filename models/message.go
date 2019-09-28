package models

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/gofrs/uuid"
)

const (
	MessageStatePending = "pending"
	MessageStateSuccess = "success"

	MessageCategoryMessageRecall = "MESSAGE_RECALL"
	MessageCategoryPlainText     = "PLAIN_TEXT"
	MessageCategoryPlainImage    = "PLAIN_IMAGE"
	MessageCategoryPlainVideo    = "PLAIN_VIDEO"
	MessageCategoryPlainData     = "PLAIN_DATA"
	MessageCategoryPlainSticker  = "PLAIN_STICKER"
	MessageCategoryPlainContact  = "PLAIN_CONTACT"
	MessageCategoryPlainAudio    = "PLAIN_AUDIO"
	MessageCategoryAppCard       = "APP_CARD"
)

const messages_DDL = `
CREATE TABLE IF NOT EXISTS messages (
	message_id            VARCHAR(36) PRIMARY KEY CHECK (message_id ~* '^[0-9a-f-]{36,36}$'),
	user_id	              VARCHAR(36) NOT NULL CHECK (user_id ~* '^[0-9a-f-]{36,36}$'),
	category              VARCHAR(512) NOT NULL,
	quote_message_id      VARCHAR(36) NOT NULL DEFAULT '',
	data                  TEXT NOT NULL,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	state                 VARCHAR(128) NOT NULL,
	last_distribute_at    TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS messages_state_updatedx ON messages(state, updated_at);
`

var messagesCols = []string{"message_id", "user_id", "category", "quote_message_id", "data", "created_at", "updated_at", "state", "last_distribute_at"}

func (m *Message) values() []interface{} {
	return []interface{}{m.MessageId, m.UserId, m.Category, m.QuoteMessageId, m.Data, m.CreatedAt, m.UpdatedAt, m.State, m.LastDistributeAt}
}

func messageFromRow(row durable.Row) (*Message, error) {
	var m Message
	err := row.Scan(&m.MessageId, &m.UserId, &m.Category, &m.QuoteMessageId, &m.Data, &m.CreatedAt, &m.UpdatedAt, &m.State, &m.LastDistributeAt)
	return &m, err
}

type Message struct {
	MessageId        string
	UserId           string
	Category         string
	QuoteMessageId   string
	Data             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	State            string
	LastDistributeAt time.Time

	FullName sql.NullString
}

func CreateMessage(ctx context.Context, user *User, messageId, category, quoteMessageId, data string, createdAt, updatedAt time.Time) (*Message, error) {
	if len(data) > 5*1024 {
		return nil, nil
	}
	if user.UserId != config.AppConfig.Mixin.ClientId && !user.isAdmin() {
		if category != MessageCategoryMessageRecall && !durable.Allow(user.UserId) {
			text := base64.StdEncoding.EncodeToString([]byte(config.AppConfig.MessageTemplate.MessageTipsTooMany))
			if err := CreateSystemDistributedMessage(ctx, user, MessageCategoryPlainText, text); err != nil {
				return nil, err
			}
			return nil, nil
		}
	}
	if !user.isAdmin() {
		b, err := ReadProhibitedProperty(ctx)
		if err != nil {
			return nil, err
		} else if b {
			return nil, nil
		}
	}
	if category == MessageCategoryPlainAudio {
		if !user.isAdmin() {
			return nil, nil
		}
		if !config.AppConfig.System.AudioMessageEnable {
			return nil, nil
		}
	}
	if category == MessageCategoryPlainImage {
		if !user.isAdmin() && !config.AppConfig.System.ImageMessageEnable {
			return nil, nil
		}
	}
	if category == MessageCategoryPlainVideo {
		if !user.isAdmin() && !config.AppConfig.System.VideoMessageEnable {
			return nil, nil
		}
	}
	if category == MessageCategoryPlainContact {
		if !user.isAdmin() && !config.AppConfig.System.ContactMessageEnable {
			return nil, nil
		}
	}
	message := &Message{
		MessageId:        messageId,
		UserId:           user.UserId,
		Category:         category,
		Data:             data,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		State:            MessageStatePending,
		LastDistributeAt: genesisStartedAt(),
	}

	if quoteMessageId != "" {
		if id, _ := uuid.FromString(quoteMessageId); id.String() == quoteMessageId {
			message.QuoteMessageId = quoteMessageId
			dm, err := FindDistributedMessage(ctx, quoteMessageId)
			if err != nil {
				return nil, err
			}
			if dm != nil {
				message.QuoteMessageId = dm.ParentId
			}
		}
	}
	if category == MessageCategoryMessageRecall {
		bytes, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, session.BadDataError(ctx)
		}
		var recallMessage RecallMessage
		err = json.Unmarshal(bytes, &recallMessage)
		if err != nil {
			return nil, session.BadDataError(ctx)
		}
		m, err := FindMessage(ctx, recallMessage.MessageId)
		if err != nil || m == nil {
			return nil, err
		}
		if m.UserId != user.UserId && !user.isAdmin() {
			return nil, session.ForbiddenError(ctx)
		}
		if user.isAdmin() {
			message.UserId = m.UserId
		}
	}
	params, positions := compileTableQuery(messagesCols)
	query := fmt.Sprintf("INSERT INTO messages (%s) VALUES (%s) ON CONFLICT (message_id) DO NOTHING", params, positions)
	_, err := session.Database(ctx).ExecContext(ctx, query, message.values()...)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	plugin.Trigger(plugin.EventTypeMessageCreated, *message)
	return message, nil
}

func createSystemMessage(ctx context.Context, tx *sql.Tx, category, data string) error {
	mixin := config.AppConfig.Mixin
	t := time.Now()
	message := &Message{
		MessageId:        bot.UuidNewV4().String(),
		UserId:           mixin.ClientId,
		Category:         category,
		Data:             data,
		CreatedAt:        t,
		UpdatedAt:        t,
		State:            MessageStatePending,
		LastDistributeAt: genesisStartedAt(),
	}
	params, positions := compileTableQuery(messagesCols)
	query := fmt.Sprintf("INSERT INTO messages (%s) VALUES (%s) ON CONFLICT (message_id) DO NOTHING", params, positions)
	_, err := tx.ExecContext(ctx, query, message.values()...)
	return err
}

func createSystemJoinMessage(ctx context.Context, tx *sql.Tx, user *User) error {
	b, err := readProhibitedStatus(ctx, tx)
	prohibited := err != nil || b
	if prohibited {
		// send MessageTipsJoinUserProhibited to joined user
		CreateSystemDistributedMessage(ctx, user, "PLAIN_TEXT", base64.StdEncoding.EncodeToString([]byte(config.AppConfig.MessageTemplate.MessageTipsJoinUserProhibited)))
	} else {
		// send MessageTipsJoinUser to joined user while send MessageTipsJoin to all users
		CreateSystemDistributedMessage(ctx, user, "PLAIN_TEXT", base64.StdEncoding.EncodeToString([]byte(config.AppConfig.MessageTemplate.MessageTipsJoinUser)))
		err = createSystemMessage(ctx, tx, "PLAIN_TEXT", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(config.AppConfig.MessageTemplate.MessageTipsJoin, user.FullName))))
		if err != nil {
			return err
		}
	}
	return nil
}

func generateRandomColor() string {
	colors := []string{
		"#0C9C9C",
		"#80C748",
		"#E57C00",
		"#E50000",
		"#00D8E5",
		"#4C7EFE",
		"#854CFE",
		"#E54CFE",
		"#FE4C82",
		"#FFA800",
		"#3672CC",
		"#88462A",
	}
	ix := rand.Intn(len(colors))
	return colors[ix]
}

func createSystemRewardsMessage(ctx context.Context, tx *sql.Tx, fromUser *User, toUser *User, amount, symbol string) error {
	b, err := readProhibitedStatus(ctx, tx)
	if err != nil || b {
		return nil
	}

	label := fmt.Sprintf(config.AppConfig.MessageTemplate.MessageTipsRewards, fromUser.FullName, toUser.FullName, amount, symbol)
	t := time.Now()
	host := config.AppConfig.Service.HTTPResourceHost
	actionURL := host
	// @TODO uncomment followed lines when new Messenger iOS fix the transpile bug.
	// if config.AppConfig.System.RouterMode == config.RouterModeHash {
	// 	host = host + config.RouterModeHashSymbol
	// }
	// actionURL := fmt.Sprintf(host + "/rewards")
	if utf8.RuneCountInString(label) > 36 {
		label = string([]rune(label)[:36])
	}
	btns, err := json.Marshal([]interface{}{map[string]string{
		"label":  label,
		"action": actionURL,
		"color":  generateRandomColor(),
	}})
	message := &Message{
		MessageId: bot.UuidNewV4().String(),
		UserId:    config.AppConfig.Mixin.ClientId,
		Category:  "APP_BUTTON_GROUP",
		Data: base64.StdEncoding.EncodeToString(
			btns),
		CreatedAt: t,
		UpdatedAt: t,
		State:     MessageStatePending,
	}

	params, positions := compileTableQuery(messagesCols)
	query := fmt.Sprintf("INSERT INTO messages (%s) VALUES (%s)", params, positions)
	_, err = tx.ExecContext(ctx, query, message.values()...)
	return err
}

func PendingMessages(ctx context.Context, limit int64) ([]*Message, error) {
	var messages []*Message
	query := fmt.Sprintf("SELECT %s FROM messages WHERE state=$1 ORDER BY state,updated_at LIMIT $2", strings.Join(messagesCols, ","))
	rows, err := session.Database(ctx).QueryContext(ctx, query, MessageStatePending, limit)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	for rows.Next() {
		m, err := messageFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func FindMessage(ctx context.Context, id string) (*Message, error) {
	query := fmt.Sprintf("SELECT %s FROM messages WHERE message_id=$1", strings.Join(messagesCols, ","))
	row := session.Database(ctx).QueryRowContext(ctx, query, id)
	message, err := messageFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return message, nil
}

func LastestMessageWithUser(ctx context.Context, limit int64) ([]*Message, error) {
	query := "SELECT messages.message_id,messages.category,messages.data,messages.created_at,users.full_name FROM messages LEFT JOIN users ON messages.user_id=users.user_id ORDER BY updated_at DESC LIMIT $1"
	rows, err := session.Database(ctx).QueryContext(ctx, query, limit)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.MessageId, &m.Category, &m.Data, &m.CreatedAt, &m.FullName)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		if m.Category == "PLAIN_TEXT" {
			data, _ := base64.StdEncoding.DecodeString(m.Data)
			m.Data = string(data)
		} else {
			m.Data = ""
		}
		messages = append(messages, &m)
	}
	return messages, nil
}

func readLastestMessages(ctx context.Context, limit int64) ([]*Message, error) {
	var messages []*Message
	query := fmt.Sprintf("SELECT %s FROM messages WHERE state=$1 ORDER BY updated_at DESC LIMIT $2", strings.Join(messagesCols, ","))
	rows, err := session.Database(ctx).QueryContext(ctx, query, MessageStateSuccess, limit)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	for rows.Next() {
		m, err := messageFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func readLastestMessagesInTx(ctx context.Context, tx *sql.Tx, limit int64) ([]*Message, error) {
	var messages []*Message
	query := fmt.Sprintf("SELECT %s FROM messages WHERE state=$1 ORDER BY updated_at DESC LIMIT $2", strings.Join(messagesCols, ","))
	rows, err := tx.QueryContext(ctx, query, MessageStateSuccess, limit)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	for rows.Next() {
		m, err := messageFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}

type RecallMessage struct {
	MessageId string `json:"message_id"`
}
