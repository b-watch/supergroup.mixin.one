package models

import (
	"context"
	"fmt"

	"github.com/MixinNetwork/supergroup.mixin.one/session"
)

func ReadStatistic(ctx context.Context, user *User) (map[string]interface{}, error) {
	var err error
	s := make(map[string]interface{}, 0)
	count, err := PaidMemberCount(ctx)
	if err != nil {
		return nil, err
	}
	s["users_count"] = count
	s["mode"], _ = ReadGroupModeProperty(ctx)
	s["announcement"], _ = ReadAnnouncementProperty(ctx)
	s["broadcast"], _ = ReadBroadcastProperty(ctx)
	s["pinned_message"], _ = ReadPinnedMessage(ctx)
	return s, nil
}

func ReadBroadcastById(ctx context.Context) ([]*WsBroadcastMessage, error) {
	query := fmt.Sprintf("SELECT * FROM broadcast_message ORDER BY created_at DESC LIMIT 100")
	rows, err := session.Database(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var bmsgs []*WsBroadcastMessage
	for rows.Next() {
		m, err := broadcastMessageFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		bmsgs = append(bmsgs, m)
	}
	return bmsgs, nil
}
