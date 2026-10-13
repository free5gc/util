package httpwrapper

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHttp2Server(t *testing.T) {
	for _, timeout := range []time.Duration{0, time.Millisecond, 30 * time.Second} {
		t.Run(timeout.String(), func(t *testing.T) {
			server, err := NewHttp2Server("127.0.0.1:8000", "", http.NewServeMux(), timeout)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, timeout, server.IdleTimeout)
			assert.True(t, server.Protocols.HTTP1())
			assert.True(t, server.Protocols.HTTP2())
			assert.True(t, server.Protocols.UnencryptedHTTP2())
		})
	}
	server, err := NewHttp2Server("127.0.0.1:8000", "", nil, time.Second)
	assert.Error(t, err)
	assert.Nil(t, server)
}

func TestNewRequest(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(),
		"GET", "http://localhost:8080?name=NCTU&location=Hsinchu", nil)
	if err != nil {
		t.Errorf("TestNewRequest error: %+v", err)
	}
	req.Header.Set("Location", "https://www.nctu.edu.tw/")
	request := NewRequest(req, 1000)
	assert.Equal(t, "https://www.nctu.edu.tw/", request.Header.Get("Location"))
	assert.Equal(t, "NCTU", request.Query.Get("name"))
	assert.Equal(t, "Hsinchu", request.Query.Get("location"))
	assert.Equal(t, 1000, request.Body)
}

func TestNewResponse(t *testing.T) {
	response := NewResponse(http.StatusCreated, map[string][]string{
		"Location": {"https://www.nctu.edu.tw/"},
		"Refresh":  {"url=https://free5gc.org"},
	}, 1000)
	assert.Equal(t, "https://www.nctu.edu.tw/", response.Header.Get("Location"))
	assert.Equal(t, "url=https://free5gc.org", response.Header.Get("Refresh"))
	assert.Equal(t, 1000, response.Body)
}
