package middleware

import (
	"fmt"
	"net/http"

	"github.com/dwethmar/go-api/pkg/request"
	"github.com/dwethmar/go-api/pkg/store"
)

// ContentTypeIndex get list of entries
func ContentTypeIndex(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if (*r).Method == "OPTIONS" {
			return
		}

		var entries, err = s.Content.GetAll()

		if err != nil {
			fmt.Printf("Error while getting entries: %v", err)
			request.SendServerError(w, r)
			return
		}

		request.SendJSON(w, r, entries, http.StatusOK)
	})
}
