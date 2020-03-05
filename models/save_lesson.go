package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/gofrs/uuid"
)

const lesson_messages_DDL = `
CREATE TABLE IF NOT EXISTS lesson_record_message (
  lesson_id VARCHAR(36) NOT NULL,
	id VARCHAR(36) NOT NULL,
  quote_message_id VARCHAR(36) NOT NULL,
	speaker_name VARCHAR(512) NOT NULL DEFAULT '',
	speaker_avatar VARCHAR(1024) NOT NULL DEFAULT '',
	speaker_id VARCHAR(36) NOT NULL,
	category VARCHAR(512) NOT NULL,
	data TEXT NOT NULL,
	text TEXT,
	attachment JSONB,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	PRIMARY KEY(lesson_id, id)
);
`

type LessonMessage struct {
	LessonId       string `json:"lesson_id"`
	MessageId      string `json:"id"`
	QuoteMessageId string `json:"quote_message_id"`
	SpeakerName    string `json:"speaker_name"`
	SpeakerAvatar  string `json:"speaker_avatar"`
	SpeakerId      string `json:"speaker_id"`
	Category       string `json:"category"`
	Data           string `json:"data"`
	Text           string `json:"text"`
	Attachment     string `json:"attachment"`
	CreatedAt      string `json:"created_at"`
}

func SaveLessonMessage(ctx context.Context, message *WsBroadcastMessage) {
	value, err := ReadPropertyAsString(ctx, "lesson-id")
	if err != nil {
		log.Panicln("SaveLessonMessage err", err)
	}
	if value == "" {
		value = uuid.Must(uuid.NewV4()).String()
		CreateProperty(ctx, "lesson-id", value, "")
	}
	tA, _ := json.Marshal(message.Attachment)
	attachment := string(tA)
	var bmsg LessonMessage
	bmsg.LessonId = value
	bmsg.MessageId = message.MessageId
	bmsg.QuoteMessageId = message.QuoteMessageId
	bmsg.SpeakerName = message.SpeakerName
	bmsg.SpeakerAvatar = message.SpeakerAvatar
	bmsg.SpeakerId = message.SpeakerId
	bmsg.Category = message.Category
	bmsg.Data = message.Data
	bmsg.Text = message.Text
	bmsg.Attachment = attachment
	bmsg.CreatedAt = message.CreatedAt.Format("2006-01-02T15:04:05.000000Z")
	query := fmt.Sprintf("INSERT INTO lesson_record_message(lesson_id,id,quote_message_id,speaker_name,speaker_avatar,speaker_id,category,data,text,attachment,created_at) VALUES('%s', '%s','%s', '%s','%s', '%s','%s', '%s','%s', '%s','%s')",
		bmsg.LessonId, bmsg.MessageId, bmsg.QuoteMessageId, bmsg.SpeakerName, bmsg.SpeakerAvatar, bmsg.SpeakerId, bmsg.Category, bmsg.Data, bmsg.Text, bmsg.Attachment, bmsg.CreatedAt)
	session.Database(ctx).ExecContext(ctx, query)
	return
}

func LessonFinished(ctx context.Context) {
	CreateProperty(ctx, "lesson-id", "", "")
}
