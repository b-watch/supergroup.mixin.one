package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/lib/pq"
)

const rewards_recipient_DDL = `
CREATE TABLE IF NOT EXISTS rewards_recipients (
	user_id           VARCHAR(36) PRIMARY KEY CHECK (user_id ~* '^[0-9a-f-]{36,36}$'),
	full_name         VARCHAR(512) NOT NULL DEFAULT '',
	avatar_url        VARCHAR(1024) NOT NULL DEFAULT '',
	status			  VARCHAR(16) NOT NULL DEFAULT '',
	created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS rewards_recipients_userx ON rewards_recipients(user_id);
`

type RewardsRecipient struct {
	UserId    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var rewardsRecipientColums = []string{"user_id", "full_name", "avatar_url", "status", "created_at"}

func (c *RewardsRecipient) values() []interface{} {
	return []interface{}{c.UserId, c.Status, c.CreatedAt}
}

func rewardsRecipientFromRow(row durable.Row) (*RewardsRecipient, error) {
	var c RewardsRecipient
	err := row.Scan(&c.UserId, &c.FullName, &c.AvatarURL, &c.Status, &c.CreatedAt)
	return &c, err
}

func GetRewardsRecipients(ctx context.Context) ([]*RewardsRecipient, error) {
	query := fmt.Sprintf("SELECT %s FROM rewards_recipients ORDER BY created_at desc LIMIT 10", strings.Join(rewardsRecipientColums, ","))
	rows, err := session.Database(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var recipients []*RewardsRecipient
	for rows.Next() {
		recipient, err := rewardsRecipientFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		recipients = append(recipients, recipient)
	}
	return recipients, nil
}

func CreateRewardsRecipient(ctx context.Context, identityNumber int64) (*RewardsRecipient, error) {
	var recipient RewardsRecipient
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var users []*User
		var err error
		users, err = findUsersByIdentityNumber(ctx, identityNumber)
		if err != nil {
			return session.TransactionError(ctx, err)
		}
		if len(users) == 0 {
			return session.NotFoundError(ctx)
		}
		user := users[0]
		recipient.UserId = user.UserId
		recipient.FullName = user.FullName
		recipient.AvatarURL = user.AvatarURL
		recipient.Status = "available"
		recipient.CreatedAt = time.Now()

		values := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')", recipient.UserId, recipient.FullName, recipient.AvatarURL, recipient.Status, string(pq.FormatTimestamp(recipient.CreatedAt)))
		query := fmt.Sprintf("INSERT INTO rewards_recipients (user_id, full_name, avatar_url, status, created_at) VALUES %s ON CONFLICT (user_id) DO NOTHING", values)
		_, err = session.Database(ctx).ExecContext(ctx, query)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &recipient, nil
}

func RemoveRewardsRecipient(ctx context.Context, userId string) error {
	_, err := session.Database(ctx).ExecContext(ctx, fmt.Sprintf("DELETE FROM rewards_recipients WHERE user_id=$1"), userId)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func CreateRewardsMessage(ctx context.Context, fromUser, toUser *User, amount, symbol string) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err := createSystemRewardsMessage(ctx, tx, fromUser, toUser, amount, symbol); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
