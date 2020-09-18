package middleware

import (
	"net/http"

	"github.com/dwethmar/go-api/pkg/common"

	"github.com/go-chi/chi"
)

// RequireID requires that a id is provided in the url.
func RequireID(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.StringToID(chi.URLParam(r, "id"))
		if err != nil {
			common.SendBadRequestError(w, r, "invalid ID")
			return
		}

		ctx := common.WithID(r.Context(), id)

		next(w, r.WithContext(ctx))
	})
}
