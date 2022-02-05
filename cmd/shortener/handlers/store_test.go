package handlers

import (
	"github.com/genga911/yandexworkshop/cmd/shortener/storages"
	"net/http"
	"testing"
)

func TestStore(t *testing.T) {
	type args struct {
		storage storages.Repository
		w       http.ResponseWriter
		req     *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Store(tt.args.storage, tt.args.w, tt.args.req)
		})
	}
}
