package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	Index(response, request, nil)

	if statusCode := response.Result().StatusCode; statusCode != http.StatusOK {
		t.Errorf("should return 200. got %q, want %q", statusCode, http.StatusOK)
	}

	got := response.Body.String()
	want := "Welcome\n"
	if got != want {
		t.Errorf("should return welcome message. got %q, want %q", got, want)
	}
}
