package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var cache *Cache

var client *http.Client

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	prepare()

	Info.Println("Ready to serve")

	server := &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		Handler:      http.HandlerFunc(handleGet),
	}

	err := server.ListenAndServe()
	if err != nil {
		Error.Fatal(err.Error())
	}
}

func prepare() {
	Info.Println("Init cache")
	CreateCache()

	client = &http.Client{
		Timeout: time.Second * 30,
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	fullURL := r.URL.Path + "?" + r.URL.RawQuery

	Info.Printf("Requested '%s'\n", fullURL)
	if r.Method != "GET" {
		response, err := gatewayRequest(fullURL, r)
		if err != nil {
			handleError(err, w)
			return
		}

		body, err := ioutil.ReadAll(response.Body)
		_ = response.Body.Close()
		if err != nil {
			handleError(err, w)
			return
		}

		if err != nil {
			Error.Printf("Could not write into cache: %s\n", err)
		}

		copyHeaders(w.Header(), response.Header)
		_, _ = w.Write(body)
	} else if cache.has(fullURL) {
		content, headers, err := cache.get(fullURL)

		var status string
		update := cache.checkUpdate(fullURL)
		if update {
			status = "UPDATING"
		} else {
			status = "HIT"
		}

		if err != nil {
			handleError(err, w)
		} else {
			for key, value := range headers {
				w.Header().Set(key, value)
			}
			w.Header().Set("Cache-Status", status)
			_, _ = w.Write(content)
		}
	} else {
		response, err := gatewayRequest(fullURL, r)
		if err != nil {
			handleError(err, w)
			return
		}

		body, err := ioutil.ReadAll(response.Body)
		_ = response.Body.Close()
		if err != nil {
			handleError(err, w)
			return
		}

		Info.Printf("Creating cache for %s", fullURL)
		err = cache.put(fullURL, response.Header, body)

		if err != nil {
			Error.Printf("Could not write into cache: %s\n", err)
		}

		copyHeaders(w.Header(), response.Header)
		w.Header().Set("Cache-Status", "MISS")
		_, _ = w.Write(body)
	}
}

func gatewayRequest(fullURL string, r *http.Request) (*http.Response, error) {
	URL := getEnv("GATEWAY_HOST",getEnv("GATEWAY_REQUEST", "http://localhost:1333")) + fullURL
	req, _ := http.NewRequest(r.Method, URL, r.Body)
	copyHeaders(req.Header, r.Header)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func gatewayRequestUpdate(fullURL string, headers map[string] string) (*http.Response, error) {
	URL := getEnv("GATEWAY_HOST", "http://localhost:1333") + fullURL
	req, _ := http.NewRequest("GET", URL, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func copyHeaders(dst, src http.Header) {
	for k := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func handleError(err error, w http.ResponseWriter) {
	Error.Println(err.Error())
	w.WriteHeader(500)
	_, _ = fmt.Fprintf(w, err.Error())
}
