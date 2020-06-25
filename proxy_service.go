package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func proxyLoggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func main() {
	r := mux.NewRouter()
	r.Use(proxyLoggingMiddleware)
	r.Handle("/user/profile", authenticateMiddleware(handleUserProfileRequest)).Methods("GET")
	r.HandleFunc("/microservice/name", handleUserMicroserviceNameRequest).Methods("GET")
	err := http.ListenAndServe(":5080", r)
	if err != nil {
		log.Fatalf("Server Failed to Start: %v", err)
	}
}

var handleUserProfileRequest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	proxyReq, err := createProxyRequest(r, "localhost:5082")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	client := &http.Client{}
	proxyRes,err := client.Do(proxyReq)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", proxyRes.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", proxyRes.Header.Get("Content-Length"))
	io.Copy(w, proxyRes.Body)
	proxyRes.Body.Close()
})

func handleUserMicroserviceNameRequest(w http.ResponseWriter, r *http.Request){
	proxyReq, err := createProxyRequest(r, "localhost:5082")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	client := &http.Client{}
	proxyRes,err := client.Do(proxyReq)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", proxyRes.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", proxyRes.Header.Get("Content-Length"))
	io.Copy(w, proxyRes.Body)
	proxyRes.Body.Close()
}

func authenticateMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authRequest, err := http.NewRequest("GET", "http://localhost:5081/auth", nil)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		authRequest.Header.Set("Username", r.Header.Get("Username"))
		client := &http.Client{}
		resp,err := client.Do(authRequest)
		defer resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println(http.StatusText(http.StatusOK))
		if resp.StatusCode == http.StatusOK {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
	})
}

func createProxyRequest(r *http.Request, forwardURL string) (*http.Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	url := fmt.Sprintf("%s://%s%s", "http", forwardURL, r.RequestURI)

	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))

	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}
	return proxyReq, nil

}