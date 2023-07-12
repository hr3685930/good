package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	cmdHTTP "good/cmd/http"
	r "good/internal/logic/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping",  nil)
	gServer := gin.New()
	cmdHTTP.Routes(r.NewRouter(gServer))
	gServer.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}