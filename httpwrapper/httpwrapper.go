package httpwrapper

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

type Request struct {
	Params map[string]string
	Header http.Header
	Query  url.Values
	Body   interface{}
	URL    *url.URL
}

func NewRequest(req *http.Request, body interface{}) *Request {
	ret := &Request{}
	ret.Query = req.URL.Query()
	ret.Header = req.Header
	ret.Body = body
	ret.Params = make(map[string]string)
	ret.URL = req.URL
	return ret
}

type Response struct {
	Header http.Header
	Status int
	Body   interface{}
}

func NewResponse(code int, h http.Header, body interface{}) *Response {
	ret := &Response{}
	ret.Status = code
	ret.Header = h
	ret.Body = body
	return ret
}

// NewHttp2Server returns a server instance with HTTP/2.0 and HTTP/2.0 cleartext support.
// Cleartext clients must send the HTTP/2 preface directly: the "h2c" Upgrade handshake is
// deprecated by RFC 9113 section 3.1 and is not accepted.
// If this function cannot open or create the secret log file,
// **it still returns server instance** but without the secret log and error indication
func NewHttp2Server(
	bindAddr string, preMasterSecretLogPath string, handler http.Handler, idleTimeout time.Duration,
) (*http.Server, error) {
	if handler == nil {
		return nil, errors.New("server needs handler to handle request")
	}

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	protocols.SetUnencryptedHTTP2(true)

	server := &http.Server{
		Addr:      bindAddr,
		Handler:   handler,
		Protocols: protocols,
		// TODO: extends the idle time after re-use openapi client
		IdleTimeout: idleTimeout,
	}

	if preMasterSecretLogPath != "" {
		preMasterSecretFile, err := os.OpenFile(
			preMasterSecretLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"create pre-master-secret log [%s] fail: %s", preMasterSecretLogPath, err,
			)
		}
		server.TLSConfig = &tls.Config{
			KeyLogWriter: preMasterSecretFile,
		}
	}

	return server, nil
}
