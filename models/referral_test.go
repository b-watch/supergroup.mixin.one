package models

import (
	"github.com/jinzhu/gorm"
	"testing"

	bot "github.com/MixinNetwork/bot-api-go-client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	var err error
	gormDB, err = gorm.Open("postgres", "postgres://postgres:@localhost/mixin_test?sslmode=disable")
	if err != nil {
		log.Panicln(err)
	}
	gormDB.AutoMigrate(&Referral{})
}

func TestReferral(t *testing.T) {
	gormDB.Unscoped().Delete(Referral{})

	inviter := User{UserId: bot.UuidNewV4().String(), IdentityNumber: 1, FullName: "Inviter1", AccessToken: "accessToken", State: "paid"}
	invitee := User{UserId: bot.UuidNewV4().String(), IdentityNumber: 2, FullName: "Invitee1", AccessToken: "accessToken", State: "pending"}

	assert.Len(t, inviter.Referrals(), 0)
	referrals, _ := inviter.CreateReferrals()
	assert.Len(t, referrals, 3)
	
	referralCodes := inviter.Referrals()
	invitee.ApplyReferral(referralCodes[0].Code)
	assert.Len(t, referralCodes, 2)
}