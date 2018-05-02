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
	req, _ := http.NewRequest("POST", "/v1/meows?body=test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	res, err := simplejson.NewFromReader(w.Body)
	assert.NoError(t, err)

	assert.Regexp(t, "[\\w]{10,50}", res.GetPath("id").MustString())
}
