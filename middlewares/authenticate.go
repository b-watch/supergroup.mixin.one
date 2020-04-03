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
	"github.com/dgrijalva/jwt-go"
	"github.com/fox-one/pkg/encrypt"
	"github.com/fox-one/pkg/uuid"
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
	{"GET", "^/v1/users"},
	{"GET", "^/shortcuts$"},
	{"PUT", "^/invitations/"},
	// {"GET", "^/me$"},
	{"GET", "^/payment/currency$"},
	// {"POST", "^/payment/create$"},
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

func decodeXuexiAuthToken(publicKey interface{}, tokenString string) (string, bool) {
	var claim jwt.StandardClaims
	if _, err := jwt.ParseWithClaims(tokenString, &claim, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}); err != nil {
		return "", false
	}

	return claim.Id, uuid.IsUUID(claim.Id)
}

func Authenticate(handler http.Handler) http.Handler {
	xuexiPublicKey, err := encrypt.ParsePublicPem(config.AppConfig.Xuexi.PublicKey)
	if err != nil {
		panic("parse xuexi auth key failed")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			handleUnauthorized(handler, w, r)
			return
		}

		tokenString := header[7:]

		var (
			user *models.User
			err  error
		)

		if uid, ok := decodeXuexiAuthToken(xuexiPublicKey, tokenString); ok {
			user, err = models.FindUser(r.Context(), uid)
		} else {
			user, err = models.AuthenticateUserByToken(r.Context(), tokenString)
		}

		if err != nil {
			views.RenderErrorResponse(w, r, err)
		} else if user == nil {
			handleUnauthorized(handler, w, r)
		} else {
			ctx := context.WithValue(r.Context(), keyCurrentUser, user)
			if config.AppConfig.System.InviteToJoin && user.State == models.PaymentStateUnverified {
				handleUnverified(handler, w, r.WithContext(ctx))
				// } else if user.State != models.PaymentStatePaid {
				// 	handleUnpaid(handler, w, r.WithContext(ctx))
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

func handleUnpaid(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	for _, pp := range whitelist {
		if pp[0] != r.Method {
			continue
		}
		if matched, err := regexp.MatchString(pp[1], strings.ToLower(r.URL.Path)); err == nil && matched {
			handler.ServeHTTP(w, r)
			return
		}
	}

	views.RenderErrorResponse(w, r, session.UnpaidError(r.Context()))
}
