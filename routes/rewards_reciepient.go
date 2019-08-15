package routes

import (
	"encoding/json"
	"net/http"

	"github.com/MixinNetwork/supergroup.mixin.one/middlewares"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/MixinNetwork/supergroup.mixin.one/views"
	"github.com/dimfeld/httptreemux"
)

type rewardsImpl struct{}

type rewardsRequest struct {
	UserId string `json:"user_id"`
}

func registerRewardsRecipients(router *httptreemux.TreeMux) {
	impl := &rewardsImpl{}

	router.GET("/rewards/recipients", impl.index)
	router.POST("/rewards/recipients", impl.create)
	router.DELETE("/rewards/recipients/:id", impl.delete)
}

func (impl *rewardsImpl) create(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if middlewares.CurrentUser(r).GetRole() != "admin" {
		views.RenderErrorResponse(w, r, session.ForbiddenError(r.Context()))
		return
	}

	var body rewardsRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}

	recipient, err := models.CreateRewardsRecipient(r.Context(), body.UserId)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, recipient)
	}
}

func (impl *rewardsImpl) delete(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if middlewares.CurrentUser(r).GetRole() != "admin" {
		views.RenderErrorResponse(w, r, session.ForbiddenError(r.Context()))
		return
	}
	err := models.RemoveRewardsRecipient(r.Context(), params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, nil)
	}
}

func (impl *rewardsImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	recipients, err := models.GetRewardsRecipients(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, recipients)
	}
}
