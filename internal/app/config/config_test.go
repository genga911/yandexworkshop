package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerAddress(t *testing.T) {

	tests := []struct {
		name    string
		want    string
		envHost string
		envPort string
	}{
		{
			name: "Корректно заданный env",
			want: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := GetServerAddress()
			assert.Equal(t, address, tt.want)
		})
	}
}

func TestGetMainLink(t *testing.T) {
	tests := []struct {
		name        string
		want        string
		envHost     string
		envPort     string
		envProtocol string
	}{
		{
			name: "Корректно заданный env",
			want: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := GetMainLink()
			assert.Equal(t, address, tt.want)
		})
	}
}
