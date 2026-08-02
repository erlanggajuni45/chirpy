package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	s := http.Server{
		Handler: mux,
		Addr:    ":8081",
	}

	s.ListenAndServe()
}
