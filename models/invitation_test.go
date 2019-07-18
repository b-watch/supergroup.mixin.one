package models

import (
	"testing"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/stretchr/testify/assert"
)

func TestInvitation(t *testing.T) {
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

	// paid user should be able to create valid invitation codes
	invitations, err := inviter.Invitations(ctx)
	assert.Nil(err)
	assert.Len(invitations, 0)
	_, err = inviter.CreateInvitations(ctx)
	assert.Nil(err)
	invitations, err = inviter.Invitations(ctx)
	assert.Nil(err)
	assert.Len(invitations, 3)
	assert.False(invitations[0].UsedAt.Valid)

	firstInvitationCode := invitations[0].Code

	// already joined member can't use invitation
	invitation, err := inviter.ApplyInvitation(ctx, firstInvitationCode)
	assert.NotNil(err)
	assert.Nil(invitation)

	// unverified user should be able to join group with valid invitation code
	invitation, err = invitee.ApplyInvitation(ctx, firstInvitationCode)
	assert.Nil(err)
	invitee, err = FindUser(ctx, invitee.UserId)
	assert.Equal(PaymentStatePending, invitee.State)
	invitations, _ = inviter.Invitations(ctx)
	assert.Len(invitations, 3)
	assert.True(invitations[0].UsedAt.Valid)

	// invitee can't join with a used invitation code
	invitee2, err := createUser(ctx, "accessToken", bot.UuidNewV4().String(), "3", "Invitee2", "http://localhost")
	invitation, err = invitee2.ApplyInvitation(ctx, firstInvitationCode)
	assert.NotNil(err)
	assert.Nil(invitation)

	// inviter can't create valid invitation codes when there are unused invitation codes presents
	invitations, err = inviter.CreateInvitations(ctx)
	assert.NotNil(err)
	assert.Nil(invitations)

	// inviter should be able to list all users which invite by itself
	invitations, _ = inviter.InvitationsHistory(ctx)
	for _, invitation = range invitations {
		assert.NotNil(invitation.Invitee)
	}
}