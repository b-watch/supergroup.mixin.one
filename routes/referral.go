package routes

import (
	"net/http"

	"github.com/MixinNetwork/supergroup.mixin.one/middlewares"
	"github.com/MixinNetwork/supergroup.mixin.one/views"
	"github.com/dimfeld/httptreemux"
)

type referralsImpl struct{}

func registerReferrals(router *httptreemux.TreeMux) {
	impl := &referralsImpl{}
	router.GET("/referral_codes", impl.index)
	router.POST("/referral_codes", impl.create)
	router.PUT("/referral_codes/:code", impl.apply)
}

func (impl *referralsImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	user := middlewares.CurrentUser(r)
	if referrals, err := user.Referrals(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderReferrals(w, r, referrals)
	}
}

func (impl *referralsImpl) create(w http.ResponseWriter, r *http.Request, params map[string]string) {
	user := middlewares.CurrentUser(r)
	if referrals, err := user.CreateReferrals(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderReferrals(w, r, referrals)
	}
}

func (impl *referralsImpl) apply(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// update referral code status if code exists and unused
	user := middlewares.CurrentUser(r)
	if referral, err := user.ApplyReferral(r.Context(), params["code"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderReferral(w, r, referral)
	}
}
