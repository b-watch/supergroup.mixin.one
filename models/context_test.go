package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
)

const (
	testEnvironment = "test"
	testDatabase    = "group_test"
)

const (
	dropCouponsDDL             = `DROP TABLE IF EXISTS coupons;`
	dropPropertiesDDL          = `DROP TABLE IF EXISTS properties;`
	dropParticipantsDDL        = `DROP TABLE IF EXISTS participants;`
	dropPacketsDDL             = `DROP TABLE IF EXISTS packets;`
	dropAssetsDDL              = `DROP TABLE IF EXISTS assets;`
	dropBlacklistsDDL          = `DROP TABLE IF EXISTS blacklists;`
	dropDistributedMessagesDDL = `DROP TABLE IF EXISTS distributed_messages;`
	dropMessagesDDL            = `DROP TABLE IF EXISTS messages;`
	dropUsersDDL               = `DROP TABLE IF EXISTS users;`
	dropInvitationsDDL           = `DROP TABLE IF EXISTS invitations;`
)

func TestClear(t *testing.T) {
	ctx := setupTestContext()
	teardownTestContext(ctx)
}

func teardownTestContext(ctx context.Context) {
	db := session.Database(ctx)
	tables := []string{
		dropUsersDDL,
		dropMessagesDDL,
		dropDistributedMessagesDDL,
		dropBlacklistsDDL,
		dropAssetsDDL,
		dropParticipantsDDL,
		dropPacketsDDL,
		dropPropertiesDDL,
		dropCouponsDDL,
		dropInvitationsDDL,
	}
	for _, q := range tables {
		if _, err := db.Exec(q); err != nil {
			log.Panicln(err)
		}
	}
}

func setupTestContext() context.Context {
	config.LoadConfig("../")
	if config.AppConfig.Service.Environment != testEnvironment || config.AppConfig.Database.DatabaseName != testDatabase {
		log.Panicln(config.AppConfig.Service.Environment, config.AppConfig.Database.DatabaseName)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.AppConfig.Database.DatebaseUser, config.AppConfig.Database.DatabasePassword, config.AppConfig.Database.DatabaseHost, config.AppConfig.Database.DatabasePort, config.AppConfig.Database.DatabaseName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panicln(err)
	}
	tables := []string{
		users_DDL,
		messages_DDL,
		distributed_messages_DDL,
		assets_DDL,
		blacklist_DDL,
		packets_DDL,
		participants_DDL,
		properties_DDL,
		coupons_DDL,
		invitation_DDL,
	}
	for _, q := range tables {
		if _, err := db.Exec(q); err != nil {
			log.Panicln(err)
		}
	}
	database, err := durable.NewDatabase(context.Background(), db)
	if err != nil {
		log.Panicln(err)
	}
	return session.WithDatabase(context.Background(), database)
}
