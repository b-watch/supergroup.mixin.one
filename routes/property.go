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

type propertyImpl struct{}

func registerProperties(router *httptreemux.TreeMux) {
	impl := propertyImpl{}

	router.POST("/properties", impl.create)
	router.GET("/properties/:name", impl.show)
}

func (impl *propertyImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if middlewares.CurrentUser(r).GetRole(r.Context()) != "admin" {
		views.RenderErrorResponse(w, r, session.ForbiddenError(r.Context()))
		return
	}

	p, err := models.ReadProperty(r.Context(), params["name"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, p)
	}
}

func (impl *propertyImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if middlewares.CurrentUser(r).GetRole(r.Context()) != models.PropGroupRolesAdmin || middlewares.CurrentUser(r).GetRole(r.Context()) != models.PropGroupRolesLecturer {
		views.RenderErrorResponse(w, r, session.ForbiddenError(r.Context()))
		return
	}

	var body struct {
		Key          string      `json:"key"`
		Value        string      `json:"value"`
		ComplexValue interface{} `json:"complex_value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	// @TODO need to implement access control for lecturers
	// lecturers allow changing following props:
	// - group mode
	// - broadcast
	// - announcement-message-property
	p, err := models.CreateProperty(r.Context(), body.Key, body.Value, body.ComplexValue)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, p)
	}
}
