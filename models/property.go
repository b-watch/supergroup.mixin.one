package models

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/mitchellh/mapstructure"
)

const (
	PropAnnouncementMessage = "announcement-message-property"
	PropBroadcast           = "broadcast-property"
	PropBroadcastOn         = "on"
	PropBroadcastOff        = "off"

	PropGroupRoles         = "roles-property"
	PropGroupRolesAdmin    = "admin"
	PropGroupRolesLecturer = "lecturer"
	PropGroupRolesDefault  = "user"

	PropGroupMode        = "group-mode-property"
	PropGroupModeFree    = "free"
	PropGroupModeLecture = "lecture"
	PropGroupModeMute    = "mute"

	PropPinnedMessage = "pinned-message-property"
)

const properties_DDL = `
CREATE TABLE IF NOT EXISTS properties (
	name               VARCHAR(512) PRIMARY KEY,
	value              VARCHAR(2048) NOT NULL,
	complex_value      JSONB,
	created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`

type RoleSet struct {
	Admins    []string `mapstructure:"admins" json:"admins"`
	Lecturers []string `mapstructure:"lecturers" json:"lecturers"`
}

var propertiesColumns = []string{"name", "value", "complex_value", "created_at"}

func (p *Property) values() []interface{} {
	complexValue, _ := json.Marshal(p.ComplexValue)
	return []interface{}{p.Name, p.Value, string(complexValue), p.CreatedAt}
}

func propertyFromRow(row durable.Row) (*Property, error) {
	var p Property
	var complexValue []byte
	err := row.Scan(&p.Name, &p.Value, &complexValue, &p.CreatedAt)
	json.Unmarshal(complexValue, &p.ComplexValue)
	return &p, err
}

type Property struct {
	Name         string      `json:"name"`
	Value        string      `json:"value,omitempty"`
	ComplexValue interface{} `json:"complex_value,omitempty"`
	CreatedAt    time.Time   `json:"time"`
}

func CreateProperty(ctx context.Context, name string, value string, complexValue interface{}) (*Property, error) {
	property := &Property{
		Name:         name,
		Value:        value,
		ComplexValue: complexValue,
		CreatedAt:    time.Now(),
	}
	if err := property.Validate(); err != nil {
		return nil, session.BadDataError(ctx)
	}

	return overrideProperty(ctx, property)
}

func overrideProperty(ctx context.Context, property *Property) (*Property, error) {
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		params, positions := compileTableQuery(propertiesColumns)
		query := fmt.Sprintf("INSERT INTO properties (%s) VALUES (%s) ON CONFLICT (name) DO UPDATE SET value=EXCLUDED.value, complex_value=EXCLUDED.complex_value", params, positions)
		_, err := tx.ExecContext(ctx, query, property.values()...)
		if err != nil {
			return err
		}

		return property.aroundOverride(ctx, tx)
	})

	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}

	property.afterOverride()
	return property, nil
}

func (p Property) aroundOverride(ctx context.Context, tx *sql.Tx) error {
	switch p.Name {
	case PropAnnouncementMessage:
		msg := fmt.Sprintf(config.AppConfig.MessageTemplate.MessageAnnouncement, p.Value)
		return createSystemMessage(ctx, tx, "PLAIN_TEXT", base64.StdEncoding.EncodeToString([]byte(msg)))
	case PropGroupMode:
		msg := config.AppConfig.MessageTemplate.MessageGroupModeFree
		if p.Value == PropGroupModeLecture {
			msg = config.AppConfig.MessageTemplate.MessageGroupModeLecture
		} else if p.Value == PropGroupModeMute {
			msg = config.AppConfig.MessageTemplate.MessageGroupModeMute
		}
		return createSystemMessage(ctx, tx, "PLAIN_TEXT", base64.StdEncoding.EncodeToString([]byte(msg)))
	case PropBroadcast:
		msg := config.AppConfig.MessageTemplate.MessageBroadcastOff
		if p.Value == PropBroadcastOn {
			msg = fmt.Sprintf(config.AppConfig.MessageTemplate.MessageBroadcastOn, config.AppConfig.Service.HTTPBroadcastHost)
		}
		return createSystemMessage(ctx, tx, "PLAIN_TEXT", base64.StdEncoding.EncodeToString([]byte(msg)))
	}
	return nil
}

func (p Property) afterOverride() {
	if p.Name == PropGroupMode {
		plugin.Trigger(plugin.EventTypeGroupModeChanged, p.Value)
	}
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

func readProperty(ctx context.Context, tx *sql.Tx, name string) (*Property, error) {
	query := fmt.Sprintf("SELECT %s FROM properties WHERE name=$1", strings.Join(propertiesColumns, ","))
	row := tx.QueryRowContext(ctx, query, name)
	property, err := propertyFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return property, nil
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

func ReadPropertyAsString(ctx context.Context, name string) (string, error) {
	var property *Property
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		query := fmt.Sprintf("SELECT %s FROM properties WHERE name=$1", strings.Join(propertiesColumns, ","))
		row := tx.QueryRowContext(ctx, query, name)
		property, err = propertyFromRow(row)
		if err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return PropGroupModeFree, session.TransactionError(ctx, err)
	}
	return property.Value, nil
}

func ReadGroupModeProperty(ctx context.Context) (string, error) {
	mode := PropGroupModeFree
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		mode, err = readPropertyAsString(ctx, tx, PropGroupMode)
		return err
	})
	if err != nil {
		return PropGroupModeFree, session.TransactionError(ctx, err)
	}
	return mode, nil
}

func ReadAnnouncementProperty(ctx context.Context) (string, error) {
	var b string
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		b, err = readPropertyAsString(ctx, tx, PropAnnouncementMessage)
		return err
	})
	if err != nil {
		return "", session.TransactionError(ctx, err)
	}
	return b, nil
}

func ReadBroadcastProperty(ctx context.Context) (string, error) {
	broadcast := "on"
	var err error
	broadcast, err = ReadPropertyAsString(ctx, PropBroadcast)
	if err != nil {
		return broadcast, session.TransactionError(ctx, err)
	}
	return broadcast, nil
}

func readGroupModeProperty(ctx context.Context, tx *sql.Tx) (string, error) {
	return readPropertyAsString(ctx, tx, PropGroupMode)
}

func readBroadcastProperty(ctx context.Context, tx *sql.Tx) (string, error) {
	return readPropertyAsString(ctx, tx, PropBroadcast)
}

func ReadPinnedMessage(ctx context.Context) (string, error) {
	var p *Property
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		p, err = readProperty(ctx, tx, PropPinnedMessage)
		return err
	})

	if err != nil {
		return "", err
	}
	return p.Value, nil
}

func PinMessageProperty(ctx context.Context, msg *Message) error {
	msgStr, _ := json.Marshal(msg)
	_, err := CreateProperty(ctx, PropPinnedMessage, string(msgStr), nil)
	return err
}

func UnpinMessageProperty(ctx context.Context) error {
	_, err := CreateProperty(ctx, PropPinnedMessage, "", nil)
	return err
}

func ReadRolesProperty(ctx context.Context) (RoleSet, error) {
	var r RoleSet
	var p *Property
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		p, err = readProperty(ctx, tx, PropGroupRoles)
		return err
	})

	if err != nil {
		return r, err
	}

	if err := mapstructure.Decode(p.ComplexValue, &r); err != nil {
		return r, errors.New("roleset not in correct form")
	}

	return r, nil
}

func AddRolesProperty(ctx context.Context, userId string, role string) (RoleSet, error) {
	rs, err := ReadRolesProperty(ctx)
	if err != nil {
		return rs, err
	}
	if role == PropGroupRolesAdmin && !rs.HasAdmin(userId) {
		rs.Admins = append(rs.Admins, userId)
		CreateProperty(ctx, PropGroupRoles, "", rs)
		return rs, nil
	} else if role == PropGroupRolesLecturer && !rs.HasLecturer(userId) {
		rs.Lecturers = append(rs.Lecturers, userId)
		CreateProperty(ctx, PropGroupRoles, "", rs)
		return rs, nil
	}
	return rs, err
}

func RemoveRolesProperty(ctx context.Context, userId string, role string) (RoleSet, error) {
	var pos int
	rs, err := ReadRolesProperty(ctx)
	if err != nil {
		return rs, err
	}
	if role == PropGroupRolesAdmin && rs.HasAdmin(userId) {
		pos = sort.Search(len(rs.Admins), func(i int) bool {
			return string(rs.Admins[i]) >= userId
		})
		rs.Admins = append(rs.Admins[:pos], rs.Admins[pos+1:]...)
		CreateProperty(ctx, PropGroupRoles, "", rs)
		return rs, nil
	} else if role == PropGroupRolesLecturer && rs.HasLecturer(userId) {
		pos = sort.Search(len(rs.Lecturers), func(i int) bool {
			return string(rs.Lecturers[i]) >= userId
		})
		rs.Lecturers = append(rs.Lecturers, userId)
		CreateProperty(ctx, PropGroupRoles, "", rs)
		return rs, nil
	}
	return rs, err
}

func (p *Property) Validate() error {
	if utf8.RuneCountInString(p.Value) > 2048 {
		return errors.New("value is too long")
	}

	switch p.Name {
	case PropGroupRoles:
		var roleSet RoleSet
		if err := mapstructure.Decode(p.ComplexValue, &roleSet); err != nil {
			return errors.New("roleset not in correct form")
		}
		p.ComplexValue = roleSet
	}
	return nil
}

func (rs RoleSet) GetRole(user *User) string {
	if user != nil {
		if rs.HasAdmin(user.UserId) {
			return PropGroupRolesAdmin
		} else if rs.HasLecturer(user.UserId) {
			return PropGroupRolesLecturer
		}
	}
	return PropGroupRolesDefault
}

func (rs RoleSet) HasAdmin(userID string) bool {
	if userID == config.AppConfig.Mixin.ClientId {
		return true
	}

	for _, id := range rs.Admins {
		if id == userID {
			return true
		}
	}
	return false
}

func (rs RoleSet) HasLecturer(userID string) bool {
	for _, id := range rs.Lecturers {
		if id == userID {
			return true
		}
	}
	return false
}

func (rs RoleSet) AdminIDs() []string {
	return rs.Admins
}

func (rs RoleSet) LecturerIDs() []string {
	return rs.Lecturers
}

func IsAdmin(ctx context.Context, id string) bool {
	roleSet, _ := ReadRolesProperty(ctx)
	if roleSet.HasAdmin(id) {
		return true
	}
	return false
}

func IsLecturer(ctx context.Context, id string) bool {
	roleSet, _ := ReadRolesProperty(ctx)
	if roleSet.HasLecturer(id) {
		return true
	}
	return false
}
