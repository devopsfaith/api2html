package http

import (
	"bytes"
	"net/http"

	"github.com/gregjones/httpcache"
)

var (
	cachedTransport  = httpcache.NewMemoryCacheTransport()
	cachedHTTPClient = http.Client{Transport: cachedTransport}
)

type Backend func(map[string]string, map[string]string) (*http.Response, error)

func DefaultClient(URLPattern string) Backend {
	return NewBackend(http.DefaultClient, URLPattern)
}

func CachedClient(URLPattern string) Backend {
	return NewBackend(&cachedHTTPClient, URLPattern)
}

func NewBackend(client *http.Client, URLPattern string) Backend {
	urlPattern := []byte(URLPattern)
	return func(params map[string]string, headers map[string]string) (*http.Response, error) {
		req, err := http.NewRequest("GET", string(replaceParams(urlPattern, params)), nil)
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		return client.Do(req)
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
