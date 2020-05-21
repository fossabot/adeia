package server

import (
	"adeia-api/internal/config"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestNewAPIServer(t *testing.T) {
	want := &APIServer{Srv: httprouter.New(), Config: &config.Config{}}

	got := NewAPIServer(&config.Config{})

	if !reflect.DeepEqual(got, want) {
		t.Error("should return new APIServer")
	}
}
