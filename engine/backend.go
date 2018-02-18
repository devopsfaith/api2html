package engine

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gregjones/httpcache"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

var (
	cachedTransport  = httpcache.NewMemoryCacheTransport()
	cachedHTTPClient = http.Client{Transport: cachedTransport}
)

// DefaultClient returns a Dackend to the received URLPattern with the default http client
// from the stdlib
func DefaultClient(URLPattern string) Backend {
	return NewBackend(http.DefaultClient, URLPattern)
}

// CachedClient returns a Dackend to the received URLPattern with a in-memory cache aware
// http client
func CachedClient(URLPattern string) Backend {
	return NewBackend(&cachedHTTPClient, URLPattern)
}

// NewBackend creates a Backend with the received http client and url pattern
func NewBackend(client *http.Client, URLPattern string) Backend {
	urlPattern := []byte(URLPattern)
	return func(params map[string]string, headers map[string]string, c *gin.Context) (*http.Response, error) {
		req, err := http.NewRequest("GET", string(replaceParams(urlPattern, params)), nil)
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		var s newrelic.ExternalSegment
		if newrelicApp != nil {
			s = newrelic.StartExternalSegment(nrgin.Transaction(c), req)
		}
		resp, err := client.Do(req)
		if newrelicApp != nil {
			s.Response = resp
			s.End()
		}
		return resp, err
	}
}

func replaceParams(URLPattern []byte, params map[string]string) []byte {
	if len(params) == 0 {
		return URLPattern
	}
	buff := URLPattern
	for k, v := range params {
		key := []byte{}
		key = append(key, ":"...)
		key = append(key, k...)
		buff = bytes.Replace(buff, key, []byte(v), -1)
	}
	return buff
}
