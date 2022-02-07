package heplers

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortCode(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   *regexp.Regexp
	}{
		{
			name:   "Генерация короткого кода",
			length: 8,
			want:   regexp.MustCompile(`^[a-zA-Z]{8}$`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortString := ShortCode(tt.length)
			assert.Regexp(t, tt.want, shortString)
		})
	}
}
