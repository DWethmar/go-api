package api

import (
	"net/http"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/request"

	"github.com/go-chi/chi"
)

// RequireEntryID requires that a id is provided in the url.
func RequireEntryID(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := common.ParseUUID(chi.URLParam(r, "id"))

		if err != nil {
			request.SendBadRequestError(w, r, "invalid ID")
			return
		}

		ctx := common.WithUUID(r.Context(), id)

		next(w, r.WithContext(ctx))
	})
}
