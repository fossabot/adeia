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
	Index(response, request)

	statusCode := response.Result().StatusCode
	assert.Equal(t, http.StatusOK, statusCode, "should return 200")

	got := response.Body.String()
	want := "Welcome\n"
	assert.Equal(t, want, got, "should return welcome message")
}

func TestIndex2(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	Index2(response, request)

	statusCode := response.Result().StatusCode
	assert.Equal(t, http.StatusOK, statusCode, "should return 200")

	got := response.Body.String()
	want := "Welcome 2\n"
	assert.Equal(t, want, got, "should return welcome message")
}
