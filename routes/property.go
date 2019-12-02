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
}

func (impl *propertyImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body struct {
		Key          string      `json:"key"`
		Value        string      `json:"value"`
		ComplexValue interface{} `json:"complex_value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if middlewares.CurrentUser(r).GetRole(r.Context()) != "admin" {
		views.RenderErrorResponse(w, r, session.ForbiddenError(r.Context()))
		return
	}
	var p *models.Property
	var err error
	if body.ComplexValue != nil {
		p, err = models.CreateComplexProperty(r.Context(), body.Key, body.ComplexValue)
	} else {
		p, err = models.CreateProperty(r.Context(), body.Key, body.Value)
	}
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, p)
	}
}
