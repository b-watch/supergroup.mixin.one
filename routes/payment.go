package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
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
	router.GET("/payment/currency", impl.getCurrency)
}

func calculateAmount(price, base string) (string, error) {
	var realPrice float64
	var realBase float64
	var err error
	realPrice, err = strconv.ParseFloat(price, 64)
	if err != nil {
		return "", err
	}
	realBase, err = strconv.ParseFloat(base, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.8f", realBase/realPrice), nil
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
	var foundAsset *config.PaymentAsset
	for _, asset := range config.AppConfig.System.AccpetPaymentAssetList {
		if payload.AssetID == asset.AssetId {
			foundAsset = &asset
			break
		}
	}
	if foundAsset != nil {
		amount := foundAsset.Amount
		var rate *models.CurrencyRate
		var err error
		if foundAsset.Amount == "auto" {
			if rate, err = models.GetCurrencyRate(r.Context(), foundAsset.Symbol); err != nil {
				views.RenderErrorResponse(w, r, err)
			} else {
				if config.AppConfig.System.AutoEstimateCurrency == "usd" {
					amount, err = calculateAmount(rate.PriceUsd, config.AppConfig.System.AutoEstimateBase)
				} else {
					amount, err = calculateAmount(rate.PriceCny, config.AppConfig.System.AutoEstimateBase)
				}
			}
		}
		if order, err := models.CreateMixinOrder(r.Context(), payload.UserID, payload.AssetID, amount); err != nil {
			views.RenderErrorResponse(w, r, err)
		} else {
			resp.Order = order
			views.RenderDataResponse(w, r, resp)
		}
	} else {
		views.RenderErrorResponse(w, r, nil)
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

func (impl *paymentImpl) getCurrency(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if s, err := models.GetCurrencyRates(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, s)
	}

}
