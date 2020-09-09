package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/go-chi/chi"
	"gotest.tools/v3/assert"
)

func TestRequireID(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := common.UUIDFromContext(r.Context())

		if err != nil {
			t.Error(err)
		}

		if ID.String() != "6f2128ce-e830-4b28-8402-c0795a7b18a2" {
			t.Error("wrong id")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	router := chi.NewRouter()
	router.Get("/{id}", RequireID(nextHandler))

	req := httptest.NewRequest("GET", "/6f2128ce-e830-4b28-8402-c0795a7b18a2", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK, "Status code should be equal")
}
