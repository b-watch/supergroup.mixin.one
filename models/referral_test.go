package models

import (
	"testing"
	"log"
	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/stretchr/testify/assert"
)

func TestReferral(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)
	var err error

	inviter := User{UserId: bot.UuidNewV4().String(), IdentityNumber: 1, FullName: "Inviter1", AccessToken: "accessToken", State: "paid"}
	invitee := User{UserId: bot.UuidNewV4().String(), IdentityNumber: 2, FullName: "Invitee1", AccessToken: "accessToken", State: "pending"}

	referrals, err := inviter.Referrals(ctx)
	if err != nil {
		log.Panicln(err)
	}
	assert.Len(referrals, 0)
	referrals, err = inviter.CreateReferrals(ctx)
	if err != nil {
		log.Panicln(err)
	}
	assert.Len(referrals, 3)
	
	referrals, _ = inviter.Referrals(ctx)

	referral, err := invitee.ApplyReferral(ctx, referrals[0].Code)
	if err == nil {
		assert.NotNil(referral.UsedAt)
		referrals, _ = inviter.Referrals(ctx)
		assert.Len(referrals, 2)
	}
}