package views

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/models"
)

type UserView struct {
	Type           string `json:"type"`
	UserId         string `json:"user_id"`
	IdentityNumber string `json:"identity_number"`
	FullName       string `json:"full_name"`
	AvatarURL      string `json:"avatar_url"`
	SubscribedAt   string `json:"subscribed_at"`
	Role           string `json:"role"`
}

type AccountView struct {
	UserView
	AuthenticationToken string `json:"authentication_token"`
	TraceId             string `json:"trace_id"`
	State               string `json:"state"`
}

func buildUserView(r *http.Request, user *models.User) UserView {
	return UserView{
		Type:           "user",
		UserId:         user.UserId,
		IdentityNumber: fmt.Sprint(user.IdentityNumber),
		FullName:       user.GetFullName(),
		AvatarURL:      user.AvatarURL,
		SubscribedAt:   user.SubscribedAt.Format(time.RFC3339Nano),
		Role:           user.GetRole(r.Context()),
	}
}

func RenderUsersView(w http.ResponseWriter, r *http.Request, users []*models.User, admins []*models.User, lecturers []*models.User) {
	var payload struct {
		Admins    []UserView `json:"admins"`
		Lecturers []UserView `json:"lecturers"`
		Users     []UserView `json:"users"`
	}
	userViews := make([]UserView, len(users))
	adminViews := make([]UserView, len(admins))
	lectureViews := make([]UserView, len(lecturers))
	for i, user := range users {
		userViews[i] = buildUserView(r, user)
	}
	for i, admin := range admins {
		adminViews[i] = buildUserView(r, admin)
	}
	for i, lecturer := range lecturers {
		lectureViews[i] = buildUserView(r, lecturer)
	}
	payload.Users = userViews
	payload.Admins = adminViews
	payload.Lecturers = lectureViews
	RenderDataResponse(w, r, payload)
}

func RenderUserView(w http.ResponseWriter, r *http.Request, user *models.User) {
	RenderDataResponse(w, r, buildUserView(r, user))
}

func RenderAccount(w http.ResponseWriter, r *http.Request, user *models.User) {
	userView := AccountView{
		UserView:            buildUserView(r, user),
		AuthenticationToken: user.AuthenticationToken,
		TraceId:             user.TraceId,
		State:               user.State,
	}
	RenderDataResponse(w, r, userView)
}
