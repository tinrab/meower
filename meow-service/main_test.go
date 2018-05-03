package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
)

func TestInsertMeow(t *testing.T) {
	router := newRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/meows?body=test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	res, _ := simplejson.NewFromReader(w.Body)
	assert.Regexp(t, "[\\w]{10,50}", res.GetPath("id").MustString())
}

func TestInsertMeowInvalidBody(t *testing.T) {
	router := newRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/meows?body=", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	res, _ := simplejson.NewFromReader(w.Body)
	assert.Equal(t, "Invalid body", res.GetPath("error").MustString())
}
