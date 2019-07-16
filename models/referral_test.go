package models

import (
	"testing"
	"fmt"
	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/stretchr/testify/assert"
)

func TestReferral(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)
	var err error

	user, err := createUser(ctx, "accessToken", bot.UuidNewV4().String(), "1", "Inviter1", "http://localhost")
	err = user.Payment(ctx)
	assert.Nil(err)
	inviter, err := FindUser(ctx, user.UserId)
	assert.Equal(PaymentStatePaid, inviter.State)
	
	invitee, err := createUser(ctx, "accessToken", bot.UuidNewV4().String(), "2", "Invitee1", "http://localhost")
	assert.Equal(PaymentStatePending, invitee.State)

	referrals, _ := inviter.Referrals(ctx)
	assert.Len(referrals, 0)
	referrals, _ = inviter.CreateReferrals(ctx)
	assert.Len(referrals, 3)
	
	referrals, _ = inviter.Referrals(ctx)
	assert.False(referrals[0].UsedAt.Valid)

	referral, err := invitee.ApplyReferral(ctx, referrals[0].Code)
	assert.Nil(err)
	invitee, err = FindUser(ctx, invitee.UserId)
	assert.Equal(PaymentStateInvited, invitee.State)
	assert.True(referral.UsedAt.Valid)
	referrals, _ = inviter.Referrals(ctx)
	assert.Len(referrals, 2)

	// user is able to create referral only when all referral codes are used
	referrals, err = inviter.CreateReferrals(ctx)
	if assert.Error(err) {
		assert.Equal(fmt.Errorf("There are %d unused codes, can't create new one", 2), err)
	}
}