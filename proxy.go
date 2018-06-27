package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	req, err := buildRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func copyHeader(dst, src http.Header) {
	for k, v := range src {
		switch k {
		case
			"Connection",
			"Keep-Alive",
			"Proxy-Authentication",
			"Proxy-Authorization",
			"TE",
			"Trailer",
			"Transfer-Encoding",
			"Upgrade":
			log.Printf("stripped header: %q (%v)", k, v)
		case
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
			"Access-Control-Allow-Headers":
			log.Printf("replaced header: %q (%v)", k, v)
			for _, vi := range v {
				dst.Set(k, vi)
			}
		default:
			log.Printf("copy header: %q (%v)", k, v)
			for _, vi := range v {
				dst.Add(k, vi)
			}
		}
	}
}

func transformURL(r *http.Request) (*url.URL, error) {
	q := r.URL.RawQuery
	if len(q) > 0 {
		q = "?" + q
	}

	f := r.URL.Fragment
	if len(f) > 0 {
		f = "#" + f
	}

	p := fmt.Sprintf("%s%s%s", r.URL.Path, q, f)
	if len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	return url.Parse(p)
}

func buildRequest(r *http.Request) (*http.Request, error) {
	dest, err := transformURL(r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.Method, dest.String(), r.Body)
	if err != nil {
		return nil, err
	}
	copyHeader(req.Header, r.Header)
	return req, nil
}
