package services

import (
	"context"
	"encoding/json"

	"github.com/MixinNetwork/supergroup.mixin.one/broker"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
)

func SyncUsers(db *durable.Database) {
	ctx := session.WithDatabase(context.Background(), db)
	_ = broker.Sub(ctx, broker.Connect(), func(e *broker.Event) error {
		switch plugin.EventType(e.Topic) {
		case plugin.EventTypeUserCreated, plugin.EventTypeUserPaid:
			user := decodeUserFromEvent(e)
			return models.SyncUser(ctx, user)
		default:
			return nil
		}
	})
}

func decodeUserFromEvent(e *broker.Event) *models.User {
	var user models.User
	_ = json.Unmarshal(e.Body, &user)
	return &user
}
