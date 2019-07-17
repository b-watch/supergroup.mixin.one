package models

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/lib/pq"
	"github.com/rs/xid"
)

const invitation_DDL = `
CREATE TABLE IF NOT EXISTS invitations (
	code         			VARCHAR(36) PRIMARY KEY,
	inviter_id        VARCHAR(36) NOT NULL,
	invitee_id	      VARCHAR(36),
	created_at       	TIMESTAMP WITH TIME ZONE NOT NULL,
	used_at        		TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS invitations_inviterx ON invitations(inviter_id);
`

const invitationGroupSize = 3

type Invitation struct {
	Code      string
	InviterID string
	InviteeID sql.NullString
	CreatedAt time.Time
	UsedAt    pq.NullTime
	Invitee   *User
}

var invitationColumns = []string{"code", "inviter_id", "invitee_id", "created_at", "used_at"}

func (r *Invitation) values() []interface{} {
	return []interface{}{r.Code, r.InviterID, r.InviteeID, r.CreatedAt, r.UsedAt}
}

func invitationFromRow(row durable.Row) (*Invitation, error) {
	var r Invitation
	err := row.Scan(&r.Code, &r.InviterID, &r.InviteeID, &r.CreatedAt, &r.UsedAt)
	return &r, err
}

func (user *User) Invitations(ctx context.Context) ([]*Invitation, error) {
	if user.State != PaymentStatePaid {
		return nil, session.ForbiddenError(ctx)
	}
	var invitations []*Invitation
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM invitations WHERE inviter_id = $1 ORDER BY created_at DESC LIMIT $2", strings.Join(invitationColumns, ","))
		rows, queryErr := tx.QueryContext(ctx, query, user.UserId, invitationGroupSize)
		if queryErr != nil {
			return queryErr
		}
		defer rows.Close()
		for rows.Next() {
			invitation, fetchError := invitationFromRow(rows)
			if fetchError != nil {
				return fetchError
			}
			if inviteeID := invitation.InviteeID; inviteeID.Valid {
				user, fetchUserErr := FindUser(ctx, inviteeID.String)
				if fetchUserErr != nil {
					return fetchUserErr
				}
				invitation.Invitee = user
			}
			invitations = append(invitations, invitation)
		}
		return nil
	})
	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return nil, sessionErr
		}
		return nil, session.TransactionError(ctx, err)
	}
	return invitations, nil
}

func (user *User) CreateInvitations(ctx context.Context) ([]*Invitation, error) {
	if user.State != PaymentStatePaid {
		return nil, session.ForbiddenError(ctx)
	}
	var invitations []*Invitation
	currentInvitations, err := user.Invitations(ctx)
	if err != nil {
		return nil, err
	} else if invitationCount := len(currentInvitations); invitationCount > 0 {
		return nil, session.ForbiddenError(ctx)
	} else {
		var values bytes.Buffer
		for i := 1; i <= invitationGroupSize; i++ {
			invitation := &Invitation{InviterID: user.UserId, Code: uniqueInvitationCode(), CreatedAt: time.Now()}
			if i > 1 {
				values.WriteString(",")
			}
			values.WriteString(fmt.Sprintf("('%s', '%s', '%s')", invitation.Code, invitation.InviterID, string(pq.FormatTimestamp(invitation.CreatedAt))))
			invitations = append(invitations, invitation)
		}
		query := fmt.Sprintf("INSERT INTO invitations (code,inviter_id,created_at) VALUES %s", values.String())
		_, err := session.Database(ctx).ExecContext(ctx, query)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		return invitations, nil
	}
}

func (user *User) ApplyInvitation(ctx context.Context, invitationCode string) (*Invitation, error) {
	if user.State != PaymentStateUnverified {
		return nil, session.ForbiddenError(ctx)
	}
	var invitation *Invitation
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		invitation, err = findInvitationByCode(ctx, tx, invitationCode)
		if err != nil {
			return err
		}
		if invitation.UsedAt.Valid {
			return fmt.Errorf("Invitation Code has already been used")
		}

		invitation.InviteeID = sql.NullString{String: user.UserId, Valid: true}
		invitation.UsedAt = pq.NullTime{Time: time.Now(), Valid: true}
		query := fmt.Sprintf("UPDATE invitations SET (invitee_id,used_at)=($1,$2) WHERE code=$3")
		_, err = tx.ExecContext(ctx, query, invitation.InviteeID, invitation.UsedAt, invitationCode)
		if err != nil {
			return err
		}

		query = fmt.Sprintf("UPDATE users SET state=$1 WHERE user_id=$2")
		_, err = tx.ExecContext(ctx, query, PaymentStatePending, user.UserId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return nil, sessionErr
		}
		return nil, session.TransactionError(ctx, err)
	}

	return invitation, nil
}

func findInvitationByCode(ctx context.Context, tx *sql.Tx, code string) (*Invitation, error) {
	query := fmt.Sprintf("SELECT %s FROM invitations WHERE code = $1 FOR UPDATE", strings.Join(invitationColumns, ","))
	row := tx.QueryRowContext(ctx, query, code)
	invitation, err := invitationFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return invitation, err
}

func uniqueInvitationCode() string {
	guid := xid.New()
	return guid.String()
}
