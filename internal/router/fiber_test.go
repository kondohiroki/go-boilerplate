package router

import (
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestNewFiberRouter(t *testing.T) {
	tests := []struct {
		name string
		want *fiber.App
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFiberRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFiberRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
