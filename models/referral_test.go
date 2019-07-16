package models

import (
	"testing"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/stretchr/testify/assert"
)

func TestReferral(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)
	var err error

	// set up users
	user, err := createUser(ctx, "accessToken", bot.UuidNewV4().String(), "1", "Inviter1", "http://localhost")
	err = user.Payment(ctx)
	assert.Nil(err)
	inviter, err := FindUser(ctx, user.UserId)
	assert.Equal(PaymentStatePaid, inviter.State)
	invitee, err := createUser(ctx, "accessToken", bot.UuidNewV4().String(), "2", "Invitee1", "http://localhost")
	assert.Equal(PaymentStateUnverified, invitee.State)

	// inviter should be able to create valid referral codes
	referrals, _ := inviter.Referrals(ctx)
	assert.Len(referrals, 0)
	referrals, _ = inviter.CreateReferrals(ctx)
	assert.Len(referrals, 3)
	referrals, _ = inviter.Referrals(ctx)
	assert.False(referrals[0].UsedAt.Valid)

	// invitee should be able to join group with valid referral code
	firstReferralCode := referrals[0].Code
	referral, err := invitee.ApplyReferral(ctx, firstReferralCode)
	assert.Nil(err)
	invitee, err = FindUser(ctx, invitee.UserId)
	assert.Equal(PaymentStatePending, invitee.State)
	assert.True(referral.UsedAt.Valid)
	referrals, _ = inviter.Referrals(ctx)
	assert.Len(referrals, 2)
	
	// already joined member can't use referral
	referral, err = invitee.ApplyReferral(ctx, referrals[0].Code)
	assert.NotNil(err)
	assert.Nil(referral)

	// invitee can't join with a used referral code
	invitee2, err := createUser(ctx, "accessToken", bot.UuidNewV4().String(), "3", "Invitee2", "http://localhost")
	referral, err = invitee2.ApplyReferral(ctx, firstReferralCode)
	assert.NotNil(err)
	assert.Nil(referral)

	// inviter can't create valid referral codes when there are unused referral codes presents
	referrals, err = inviter.CreateReferrals(ctx)
	assert.NotNil(err)
	assert.Nil(referrals)
}