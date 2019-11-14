package routes

import (
	"net/http"

	"github.com/MixinNetwork/supergroup.mixin.one/middlewares"
	"github.com/MixinNetwork/supergroup.mixin.one/views"
	"github.com/dimfeld/httptreemux"
)

type tiprankImpl struct{}

func registerTipranks(router *httptreemux.TreeMux) {
	impl := &tiprankImpl{}

	router.GET("/tipranks", impl.index)
}

func (impl *tiprankImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if ranks, err := middlewares.CurrentUser(r).ShowTiprank(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, ranks)
	}
}
