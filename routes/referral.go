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
	if referrals, err := user.Referrals(); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderReferralCodes(w, r, referrals)
	}
}

func (impl *referralsImpl) create(w http.ResponseWriter, r *http.Request, params map[string]string) {
	user := middlewares.CurrentUser(r)
	referrals, err := user.CreateReferrals(); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderReferralCodes(w, r, referrals)
	}
}

func (impl *referralsImpl) apply(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// update referral code status if code exists and unused
	user := middlewares.CurrentUser(r)
	referral, err := user.ApplyReferral(params["code"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderReferralCodes(w, r, referral)
	}
}
