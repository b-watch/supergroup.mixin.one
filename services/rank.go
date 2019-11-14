package services

import (
	"context"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func StartRank(name string, db *durable.Database, connStr string) {
	dsn, err := pq.ParseURL(connStr)
	if err != nil {
		logrus.Error(err)
		return
	}
	ctx := session.WithDatabase(context.Background(), db)
	models.RankManager.Init(ctx, dsn)
}
