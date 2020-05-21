package server

import (
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestNewAPIServer(t *testing.T) {
	want := &APIServer{httprouter.New()}

	got := NewAPIServer()

	if !reflect.DeepEqual(got, want) {
		t.Error("should return new APIServer")
	}
}
