package handler

import (
	"fmt"
	"net/http"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/store"
)

// ContentTypeIndex get list of entries
func ContentTypeIndex(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var c, err = s.ContentType.List()

		if err != nil {
			fmt.Printf("Error while getting entries: %v", err)
			common.SendServerError(w, r)
			return
		}

		common.SendJSON(w, r, c, http.StatusOK)
	})
}
