package main

import (
	"io"
	"net/http"
	"net/http/httptest"
)

var cache = map[string]map[string]*httptest.ResponseRecorder{}

func cacheHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, ok := cache[r.Method]
		if !ok {
			m = map[string]*httptest.ResponseRecorder{}
			cache[r.Method] = m
		}
		rr, ok := m[r.URL.String()]
		if !ok {
			rr = httptest.NewRecorder()
			m[r.URL.String()] = rr
			next(rr, r)
		}

		resp := shallowCopy(rr).Result()

		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func shallowCopy(r *httptest.ResponseRecorder) *httptest.ResponseRecorder {
	return &httptest.ResponseRecorder{
		Code:      r.Code,
		Flushed:   r.Flushed,
		HeaderMap: r.HeaderMap,
		Body:      r.Body,
	}
}
