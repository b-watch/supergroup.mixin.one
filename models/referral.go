package models

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/jinzhu/gorm"
)

var gormDB *gorm.DB

type Referral struct {
	InviterID string `gorm:"type:varchar(36);not null"`
	InviteeID string `gorm:"type:varchar(36);not null"`
	Code      string `gorm:"primary_key;type:varchar(36)"`
	IsUsed    bool 	 `gorm:"default:false"`
	CreatedAt time.Time `gorm:"type:time;not null"`
	UsedAt    time.Time `gorm:"type:time;not null"`
}

func (user *User) Referrals() (referrals []Referral) {
	gormDB.Where("inviter_id = ? AND is_used = ?", user.UserId, false).Find(&referrals)
	return
}

func (user *User) CreateReferrals() (referrals []Referral, err error) {
	if referralCount := len(user.Referrals()); referralCount > 0 {
		err = fmt.Errorf("There are %d unused codes, can't create new one", referralCount)
		return
	}

	for i := 1;  i<=3; i++ {
		referral := Referral{InviterID: user.UserId, Code: newReferralCode()}
		gormDB.Create(&referrals)
		referrals = append(referrals, referral)
	}
	return
}

func (user *User) ApplyReferral(referralCode string) (referral Referral, err error) {
	if findErr := gormDB.First(&referral, "code = ?", referralCode).Error; gorm.IsRecordNotFoundError(findErr) {
		err = findErr
		return
	}
	gormDB.Model(&referral).Updates(Referral{InviteeID: user.UserId, IsUsed: true})
	return
}

func newReferralCode() string {
	guid := xid.New()
	return guid.String()
}
