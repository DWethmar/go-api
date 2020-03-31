package server

import (
	"net/http"

	"github.com/DWethmar/go-api/pkg/common"
	"github.com/DWethmar/go-api/pkg/contententry"
	"github.com/DWethmar/go-api/pkg/request"

	"github.com/gorilla/mux"
)

func requireEntryId(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := contententry.ParseId(vars["id"])

		if err != nil {
			request.SendBadRequestError(w, r, "invalid ID")
			return
		}

		ctx := common.WithUUID(r.Context(), id)

		next(w, r.WithContext(ctx))
	})
}
