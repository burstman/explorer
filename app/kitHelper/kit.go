package kithelper

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Kit struct {
	Request  *http.Request
	Response http.ResponseWriter
	// ... other fields
}

func (k *Kit) URLParam(key string) string {
	return chi.URLParam(k.Request, key)
}
