package views

import (
	"net/http"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/lib/pq"
)

type InvitationView struct {
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

func buildInvitation(invitation *models.Invitation) InvitationView {
	invitee := invitation.Invitee
	inviteeView := InviteeView{
		UserId: invitee.UserId,
		FullName: invitee.FullName,
		AvatarURL: invitee.AvatarURL,
		State: invitee.State,
	}
	
	return InvitationView{
		Type: 	"Invitation",
		Code: invitation.Code,
		Invitee: inviteeView,
		IsUsed: invitation.UsedAt.Valid,
		CreatedAt: invitation.CreatedAt,
		UsedAt:  invitation.UsedAt,
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