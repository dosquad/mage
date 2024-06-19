package helper_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dosquad/mage/helper"
)

func TestIfErrorf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format      string
		args        []any
		expectError bool
	}{
		{"Hello There Friend", nil, false},
		{"Hello There Friend: %w", []any{errors.New("fail")}, true},
		{"Hello There Friend: %w", []any{nil}, false},
		{"Hello There Friend: %s/%s %w", []any{"foor", "bar", errors.New("fail")}, true},
		{"Hello There Friend: %s/%s %w", []any{"foor", "bar", nil}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.format, tt.args...), func(t *testing.T) {
			t.Parallel()

			err := helper.IfErrorf(tt.format, tt.args...)
			if v := (err != nil); v != tt.expectError {
				t.Errorf("helper.IfErrorf(): error == nil : got '%t', want '%t'", v, tt.expectError)
			}
		})
	}
}
