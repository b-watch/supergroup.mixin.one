package views

import (
	"net/http"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/lib/pq"
)

type ReferralView struct {
	Type      string 		`json:"type"`
	Code  		string 		`json:"code"`
	Invitee  	InviteeView 	`json:"invitee"`
	IsUsed   	bool 			`json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
	UsedAt    pq.NullTime `json:"used_at"`
}

type InviteeView struct {
	UserId         string `json:"user_id"`
	FullName       string `json:"full_name"`
	AvatarURL      string `json:"avatar_url"`
	State					 string `json:"state"`
}

func buildReferral(referral *models.Referral) ReferralView {
	invitee := referral.Invitee
	inviteeView := InviteeView{
		UserId: invitee.UserId,
		FullName: invitee.FullName,
		AvatarURL: invitee.AvatarURL,
		State: invitee.State,
	}
	
	return ReferralView{
		Type: 	"Referral",
		Code: referral.Code,
		Invitee: inviteeView,
		IsUsed: referral.IsUsed,
		CreatedAt: referral.CreatedAt,
		UsedAt:  referral.UsedAt,
	}
}

func RenderReferral(w http.ResponseWriter, r *http.Request, referral *models.Referral) {
	if referral != nil {
		RenderDataResponse(w, r, buildReferral(referral))
	} else {
		RenderBlankResponse(w, r)
	}
}

func RenderReferrals(w http.ResponseWriter, r *http.Request, referrals []*models.Referral) {
	views := make([]ReferralView, len(referrals))
	for i, referral := range referrals {
		views[i] = buildReferral(referral)
	}
	RenderDataResponse(w, r, views)
}