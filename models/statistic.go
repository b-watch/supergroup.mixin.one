package models

import (
	"context"
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
	return s, nil
}
