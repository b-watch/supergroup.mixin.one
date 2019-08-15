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
	user_id           VARCHAR(36) PRIMARY KEY CHECK (coupon_id ~* '^[0-9a-f-]{36,36}$'),
	full_name         VARCHAR(512) NOT NULL DEFAULT '',
	avatar_url        VARCHAR(1024) NOT NULL DEFAULT '',
	status			  VARCHAR(16) NOT NULL DEFAULT '',
	created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS rewards_recipients_userx ON rewards_recipients(user_id);
`

type RewardsRecipient struct {
	UserId    string
	FullName  string
	AvatarURL string
	Status    string
	CreatedAt time.Time
}

var rewardsRecipientColums = []string{"user_id", "full_name", "avatar_url", "status", "created_at"}

func (c *RewardsRecipient) values() []interface{} {
	return []interface{}{c.UserId, c.Status, c.CreatedAt}
}

func rewardsRecipientFromRow(row durable.Row) (*RewardsRecipient, error) {
	var c RewardsRecipient
	err := row.Scan(&c.UserId, &c.Status, &c.CreatedAt)
	return &c, err
}

func GetRewardsRecipients(ctx context.Context) ([]*RewardsRecipient, error) {
	query := fmt.Sprintf("SELECT %s FROM rewards_recipient WHERE created_at IS NULL LIMIT 10", strings.Join(rewardsRecipientColums, ","))
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

func CreateRewardsRecipient(ctx context.Context, userId string) (*RewardsRecipient, error) {
	var recipient RewardsRecipient
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var user *User
		var err error
		user, err = findUserById(ctx, tx, userId)
		if err != nil {
			return session.TransactionError(ctx, err)
		}
		recipient.UserId = user.UserId
		recipient.FullName = user.FullName
		recipient.AvatarURL = user.AvatarURL
		recipient.Status = "available"
		recipient.CreatedAt = time.Now()

		values := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')", recipient.UserId, recipient.FullName, recipient.AvatarURL, recipient.Status, string(pq.FormatTimestamp(recipient.CreatedAt)))
		query := fmt.Sprintf("INSERT INTO rewards_recipients (user_id, full_name, avatar_url, status, created_at) VALUES %s", values)
		_, err = session.Database(ctx).ExecContext(ctx, query)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
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
