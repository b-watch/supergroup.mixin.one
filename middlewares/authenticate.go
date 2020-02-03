package middlewares

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/MixinNetwork/supergroup.mixin.one/views"
)

var whitelist = [][2]string{
	{"GET", "^/$"},
	{"GET", "^/_hc$"},
	{"GET", "^/users"},
	{"GET", "^/config$"},
	{"GET", "^/amount$"},
	{"GET", "^/wechat"},
	{"POST", "^/wechat"},
	{"POST", "^/auth$"},
	{"GET", "^/shortcuts$"},
	{"PUT", "^/invitations/"},
	{"GET", "^/payment/currency$"},
}

var whitelistMutex sync.Mutex

func WhitelistAppend(method, urlRegex string) {
	whitelistMutex.Lock()
	defer whitelistMutex.Unlock()

	whitelist = append(whitelist, [2]string{method, urlRegex})
}

type contextValueKey struct{ int }

var keyCurrentUser = contextValueKey{1000}

func CurrentUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(keyCurrentUser).(*models.User)
	return user
}

func Authenticate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			handleUnauthorized(handler, w, r)
			return
		}

		user, err := models.AuthenticateUserByToken(r.Context(), header[7:])
		if err != nil {
			views.RenderErrorResponse(w, r, err)
		} else if user == nil {
			handleUnauthorized(handler, w, r)
		} else {
			ctx := context.WithValue(r.Context(), keyCurrentUser, user)
			if config.AppConfig.System.InviteToJoin && user.State == models.PaymentStateUnverified {
				handleUnverified(handler, w, r.WithContext(ctx))
			} else {
				handler.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}

func handleUnauthorized(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	for _, pp := range whitelist {
		if pp[0] != r.Method {
			continue
		}
		if matched, err := regexp.MatchString(pp[1], strings.ToLower(r.URL.Path)); err == nil && matched {
			handler.ServeHTTP(w, r)
			return
		}
	}

	views.RenderErrorResponse(w, r, session.AuthorizationError(r.Context()))
}

func handleUnverified(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	for _, pp := range whitelist {
		if pp[0] != r.Method {
			continue
		}
		if matched, err := regexp.MatchString(pp[1], strings.ToLower(r.URL.Path)); err == nil && matched {
			handler.ServeHTTP(w, r)
			return
		}
	}

	views.RenderErrorResponse(w, r, session.UnverifiedError(r.Context()))
}
