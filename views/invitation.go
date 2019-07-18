package views

import (
	"net/http"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/models"
)

type InvitationView struct {
	Type      string 		`json:"type"`
	Code  		string 		`json:"code"`
	Invitee  	*InviteeView 	`json:"invitee"`
	IsUsed   	bool 			`json:"is_used"`
	CreatedAt *time.Time `json:"created_at"`
	UsedAt    *time.Time `json:"used_at"`
}

type InviteeView struct {
	IdentityNumber int64  `json:"identity_number"`
	FullName       string `json:"full_name"`
	AvatarURL      string `json:"avatar_url"`
	State					 string `json:"state"`
}

func buildInvitation(invitation *models.Invitation) InvitationView {
	var inviteeView *InviteeView
	if invitee := invitation.Invitee; invitee != nil {
		inviteeView = &InviteeView{
			IdentityNumber: invitee.IdentityNumber,
			FullName: invitee.FullName,
			AvatarURL: invitee.AvatarURL,
			State: invitee.State,
		}
	}

	var usedAt *time.Time
	if invitation.UsedAt.Valid {
		usedAt = &invitation.UsedAt.Time
	}
	return InvitationView{
		Type: 	"Invitation",
		Code: invitation.Code,
		Invitee: inviteeView,
		IsUsed: invitation.UsedAt.Valid,
		CreatedAt: &invitation.CreatedAt,
		UsedAt: usedAt,
	}
}

func RenderInvitation(w http.ResponseWriter, r *http.Request, invitation *models.Invitation) {
	if invitation != nil {
		RenderDataResponse(w, r, buildInvitation(invitation))
	} else {
		RenderBlankResponse(w, r)
	}
}

func RenderInvitations(w http.ResponseWriter, r *http.Request, invitations []*models.Invitation) {
	views := make([]InvitationView, len(invitations))
	for i, invitation := range invitations {
		views[i] = buildInvitation(invitation)
	}
	RenderDataResponse(w, r, views)
}