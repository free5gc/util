package httpclient

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

const sbiHttpClientTimeout = 60 * time.Second

var (
	httpClient  http.Client
	httpsClient http.Client
)

func init() {
	httpClient = *http.DefaultClient
	httpClient.Timeout = sbiHttpClientTimeout
	sbiHttpTransport := http2.Transport{
		AllowHTTP: true,
		// ReadIdleTimeout:  0,
		// PingTimeout:      0,
		// WriteByteTimeout: 0,
	}
	httpClient.Transport = &sbiHttpTransport
	if ht, _ := http.DefaultTransport.(*http.Transport); ht != nil && ht.DialContext != nil {
		sbiHttpTransport.DialTLSContext = func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			return ht.DialContext(ctx, network, addr)
		}
	} else {
		sbiHttpTransport.DialTLSContext = func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			d := net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			return d.DialContext(ctx, network, addr)
		}
	}

	httpsClient = *http.DefaultClient
	httpsClient.Timeout = sbiHttpClientTimeout
	sbiHttpsTransport := http2.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		// ReadIdleTimeout:  0,
		// PingTimeout:      0,
		// WriteByteTimeout: 0,
	}
	httpsClient.Transport = &sbiHttpsTransport
}

func GetHttpClient(uri string) *http.Client {
	if strings.HasPrefix(strings.ToLower(uri), "http:") {
		return &httpClient
	} else {
		return &httpsClient
	}
}
