package routes

import (
	"encoding/json"
	"net/http"

	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/MixinNetwork/supergroup.mixin.one/views"
	"github.com/dimfeld/httptreemux"
)

type paymentImpl struct{}

func registerPayment(router *httptreemux.TreeMux) {
	impl := &paymentImpl{}
	router.POST("/payment/create", impl.createPayment)
	router.GET("/payment/:id", impl.checkPayment)
}

func (impl *paymentImpl) createPayment(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var payload struct {
		Method  string `json:"method"`
		AssetID string `json:"asset_id"`
		UserID  string `json:"user_id"`
	}
	var resp struct {
		Order *models.Order `json:"order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if order, err := models.CreateMixinOrder(r.Context(), payload.UserID, payload.AssetID, "0.01"); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		resp.Order = order
		views.RenderDataResponse(w, r, resp)
	}
}

func (impl *paymentImpl) checkPayment(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := params["id"]
	if s, err := models.GetOrder(r.Context(), id); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, s)
	}
}
