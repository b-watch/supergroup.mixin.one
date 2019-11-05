package models

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
)

const (
	ProhibitedMessage   = "prohibited-message-property"
	GroupMode           = "group-mode-property"
	AnnouncementMessage = "announcement-message-property"
)

const properties_DDL = `
CREATE TABLE IF NOT EXISTS properties (
	name               VARCHAR(512) PRIMARY KEY,
	value              VARCHAR(2048) NOT NULL,
	created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`

var propertiesColumns = []string{"name", "value", "created_at"}

func (p *Property) values() []interface{} {
	return []interface{}{p.Name, p.Value, p.CreatedAt}
}

func propertyFromRow(row durable.Row) (*Property, error) {
	var p Property
	err := row.Scan(&p.Name, &p.Value, &p.CreatedAt)
	return &p, err
}

type Property struct {
	Name      string
	Value     string
	CreatedAt time.Time
}

func CreateProperty(ctx context.Context, name string, value string) (*Property, error) {
	if utf8.RuneCountInString(value) > 512 {
		return nil, session.BadDataError(ctx)
	}
	property := &Property{
		Name:      name,
		Value:     fmt.Sprint(value),
		CreatedAt: time.Now(),
	}
	params, positions := compileTableQuery(propertiesColumns)
	query := fmt.Sprintf("INSERT INTO properties (%s) VALUES (%s) ON CONFLICT (name) DO UPDATE SET value=EXCLUDED.value", params, positions)
	session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, query, property.values()...)
		if err != nil {
			return err
		}
		data := config.AppConfig
		if name == AnnouncementMessage {
			text := fmt.Sprintf(data.MessageTemplate.MessageAnnouncement, value)
			return createSystemMessage(ctx, tx, "PLAIN_TEXT", base64.StdEncoding.EncodeToString([]byte(text)))
		} else if name == GroupMode {
			text := data.MessageTemplate.MessageGroupModeFree
			if value == "lecture" {
				text = data.MessageTemplate.MessageGroupModeLecture
			} else if value == "mute" {
				text = data.MessageTemplate.MessageGroupModeMute
			}
			text = fmt.Sprintf(text, value)
			return createSystemMessage(ctx, tx, "PLAIN_TEXT", base64.StdEncoding.EncodeToString([]byte(text)))
		}
		return nil
	})
	_, err := session.Database(ctx).ExecContext(ctx, query, property.values()...)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}

	if name == GroupMode {
		plugin.Trigger(plugin.EventTypeGroupModeChanged, value)
	}

	return property, nil
}

func ReadProperty(ctx context.Context, name string) (*Property, error) {
	query := fmt.Sprintf("SELECT %s FROM properties WHERE name=$1", strings.Join(propertiesColumns, ","))
	row := session.Database(ctx).QueryRowContext(ctx, query, name)
	property, err := propertyFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return property, nil
}

func readPropertyAsBool(ctx context.Context, tx *sql.Tx, name string) (bool, error) {
	query := fmt.Sprintf("SELECT %s FROM properties WHERE name=$1", strings.Join(propertiesColumns, ","))
	row := tx.QueryRowContext(ctx, query, name)
	property, err := propertyFromRow(row)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return property.Value == "true", nil
}

func readPropertyAsString(ctx context.Context, tx *sql.Tx, name string) (string, error) {
	query := fmt.Sprintf("SELECT %s FROM properties WHERE name=$1", strings.Join(propertiesColumns, ","))
	row := tx.QueryRowContext(ctx, query, name)
	property, err := propertyFromRow(row)
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return property.Value, nil
}

func ReadGroupModeProperty(ctx context.Context) (string, error) {
	mode := "free"
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		mode, err = readPropertyAsString(ctx, tx, GroupMode)
		return err
	})
	if err != nil {
		return "free", session.TransactionError(ctx, err)
	}
	return mode, nil
}

func ReadAnnouncementProperty(ctx context.Context) (string, error) {
	var b string
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		b, err = readPropertyAsString(ctx, tx, AnnouncementMessage)
		fmt.Println(b)
		return err
	})
	if err != nil {
		return "", session.TransactionError(ctx, err)
	}
	return b, nil
}

func readProhibitedStatus(ctx context.Context, tx *sql.Tx) (bool, error) {
	return readPropertyAsBool(ctx, tx, ProhibitedMessage)
}
