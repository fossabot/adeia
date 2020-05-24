package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	Index(response, request, nil)

	statusCode := response.Result().StatusCode
	assert.Equal(t, statusCode, http.StatusOK, "should return 200")

	got := response.Body.String()
	want := "Welcome\n"
	assert.Equal(t, got, want, "should return welcome message")
}
