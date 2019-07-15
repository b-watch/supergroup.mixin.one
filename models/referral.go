package models

import (
	"context"
	"time"

	"github.com/lib/pq"
)

type Referral struct {
	InviterId string
	InviteeId string
	Code      string
	IsUsed    bool
	CreatedAt time.Time
	UsedAt    pq.NullTime
}

func (user *User) Referrals(ctx context.Context) ([]*Referral, error) {
	var referrals []*Referral
	//
}

func (user *User) CreateReferrals(ctx context.Context) ([]*Referral, error) {
	var referrals []*Referral
	//
}
