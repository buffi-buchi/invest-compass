package middleware

import (
	"net/http"
	"strings"

	"github.com/buffi-buchi/invest-compass/backend/internal/api"
	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

type JWTProvider interface {
	Validate(token string) (model.AuthClaims, error)
}

func NewAuthMiddleware(jwtProvider JWTProvider) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				api.EncodeErrorf(w, http.StatusUnauthorized, "Authorization header is required")

				return
			}

			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				api.EncodeErrorf(w, http.StatusUnauthorized, "Authorization header format is invalid")

				return
			}

			token := parts[1]

			claims, err := jwtProvider.Validate(token)
			if err != nil {
				api.EncodeErrorf(w, http.StatusUnauthorized, "Invalid token")

				return
			}

			// TODO: Check that user exists.

			ctx := model.WithAuthClaims(r.Context(), claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
