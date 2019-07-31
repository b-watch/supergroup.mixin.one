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
	"github.com/lithammer/shortuuid"
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

var invitationColumns = []string{"invitations.code", "invitations.inviter_id", "invitations.invitee_id", "invitations.created_at", "invitations.used_at"}

func (r *Invitation) values() []interface{} {
	return []interface{}{r.Code, r.InviterID, r.InviteeID, r.CreatedAt, r.UsedAt}
}

func invitationFromRow(row durable.Row) (*Invitation, error) {
	var r Invitation
	err := row.Scan(&r.Code, &r.InviterID, &r.InviteeID, &r.CreatedAt, &r.UsedAt)
	return &r, err
}

func (user *User) Invitations(ctx context.Context) ([]*Invitation, error) {
	return user.invitations(ctx, false)
}

func (user *User) InvitationsHistory(ctx context.Context) ([]*Invitation, error) {
	return user.invitations(ctx, true)
}

func (user *User) invitations(ctx context.Context, historyFlag bool) ([]*Invitation, error) {
	if user.State != PaymentStatePaid {
		return nil, session.ForbiddenError(ctx)
	}
	var invitations []*Invitation
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var (
			err   error
			rows  *sql.Rows
			query string
		)
		if historyFlag {
			query = fmt.Sprintf("SELECT %s FROM invitations INNER JOIN users on invitations.invitee_id = users.user_id WHERE invitations.inviter_id = $1 AND users.state = $2 ORDER BY invitations.created_at DESC", strings.Join(invitationColumns, ","))
			rows, err = tx.QueryContext(ctx, query, user.UserId, "paid")
		} else {
			query = fmt.Sprintf("SELECT %s FROM invitations WHERE inviter_id = $1 ORDER BY created_at DESC LIMIT $2", strings.Join(invitationColumns, ","))
			rows, err = tx.QueryContext(ctx, query, user.UserId, invitationGroupSize)
		}
		if err != nil {
			return err
		}

		var userIDs []string
		defer rows.Close()
		for rows.Next() {
			invitation, err := invitationFromRow(rows)
			if err != nil {
				return err
			}
			if inviteeID := invitation.InviteeID; inviteeID.Valid {
				userIDs = append(userIDs, inviteeID.String)
			}
			invitations = append(invitations, invitation)
		}

		if len(userIDs) > 0 {
			users, err := findUsersByIds(ctx, tx, userIDs)
			if err != nil {
				return err
			}
			userMap := make(map[string]*User)
			for _, user := range users {
				userMap[user.UserId] = user
			}
			for _, invitation := range invitations {
				if inviteeID := invitation.InviteeID; inviteeID.Valid {
					invitation.Invitee = userMap[inviteeID.String]
				}
			}
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
	currentInvitations, err := user.Invitations(ctx)
	if err != nil {
		return nil, err
	} else if len(currentInvitations) > 0 {
		for _, invitation := range currentInvitations {
			if invitee := invitation.Invitee; invitee != nil {
				if invitee.State != PaymentStatePaid {
					return nil, session.ForbiddenError(ctx)
				}
			} else {
				return nil, session.ForbiddenError(ctx)
			}
		}
	}

	var invitations []*Invitation
	var values bytes.Buffer
	createTime := time.Now()
	for i := 1; i <= invitationGroupSize; i++ {
		invitation := &Invitation{InviterID: user.UserId, Code: uniqueInvitationCode(), CreatedAt: createTime}
		if i > 1 {
			values.WriteString(",")
		}
		values.WriteString(fmt.Sprintf("('%s', '%s', '%s')", invitation.Code, invitation.InviterID, string(pq.FormatTimestamp(invitation.CreatedAt))))
		invitations = append(invitations, invitation)
	}
	query := fmt.Sprintf("INSERT INTO invitations (code,inviter_id,created_at) VALUES %s", values.String())
	_, err = session.Database(ctx).ExecContext(ctx, query)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return invitations, nil
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

func (user *User) CleanUnpaidUser(ctx context.Context) (int, error) {
	var pendingUserIDs []string
	var invitationCodes []string
	currentInvitations, err := user.Invitations(ctx)
	if err != nil {
		return 0, err
	} else if len(currentInvitations) > 0 {
		for _, invitation := range currentInvitations {
			if invitee := invitation.Invitee; invitee != nil {
				if invitee.State != PaymentStatePaid {
					pendingUserIDs = append(pendingUserIDs, invitee.UserId)
					invitationCodes = append(invitationCodes, invitation.Code)
				}
			}
		}
	}

	err = session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err = deleteUsersByIds(ctx, tx, pendingUserIDs); err != nil {
			return err
		}
		if err = deleteInvitationsByCodes(ctx, tx, invitationCodes); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return 0, sessionErr
		}
		return 0, session.TransactionError(ctx, err)
	}

	return len(pendingUserIDs), nil
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

func deleteInvitationsByCodes(ctx context.Context, tx *sql.Tx, codes []string) error {
	for i, code := range codes {
		codes[i] = fmt.Sprintf("'%s'", code)
	}
	query := fmt.Sprintf("DELETE FROM invitations WHERE code IN (%s)", strings.Join(codes, ","))
	_, err := tx.QueryContext(ctx, query)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func uniqueInvitationCode() string {
	return shortuuid.New()
}
