package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/lib/pq"
)

const broadcast_message_DDL = `
CREATE TABLE IF NOT EXISTS broadcast_message (
	id VARCHAR(36) NOT NULL PRIMARY KEY,
  quote_message_id VARCHAR(36) NOT NULL,
	speaker_name VARCHAR(512) NOT NULL DEFAULT '',
	speaker_avatar VARCHAR(1024) NOT NULL DEFAULT '',
	speaker_id VARCHAR(36) NOT NULL,
	category VARCHAR(512) NOT NULL,
	data TEXT NOT NULL,
	text TEXT,
	attachment JSONB,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`

func SaveBroadcastMessage(ctx context.Context, message *WsBroadcastMessage) {
	ma, _ := json.Marshal(message.Attachment)
	query := fmt.Sprintf("INSERT INTO broadcast_message(id,quote_message_id,speaker_name,speaker_avatar,speaker_id,category,data,text,attachment,created_at) VALUES('%s', '%s','%s', '%s','%s', '%s','%s', '%s','%s', '%s')",
		message.MessageId, message.QuoteMessageId, message.SpeakerName, message.SpeakerAvatar, message.SpeakerId, message.Category, message.Data, message.Text, string(ma), string(pq.FormatTimestamp(message.CreatedAt)))
	session.Database(ctx).ExecContext(ctx, query)
	return
}

func broadcastMessageFromRow(row durable.Row) (*WsBroadcastMessage, error) {
	var m WsBroadcastMessage
	var mas string
	var ma WsBroadcastMessageAttachment
	err := row.Scan(&m.MessageId, &m.QuoteMessageId, &m.SpeakerName, &m.SpeakerAvatar, &m.SpeakerId, &m.Category, &m.Data, &m.Text, &mas, &m.CreatedAt)
	json.Unmarshal([]byte(mas), &ma)
	m.Attachment = ma
	return &m, err
}
