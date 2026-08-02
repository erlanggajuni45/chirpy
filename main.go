package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	s := http.Server{
		Handler: mux,
		Addr:    ":8081",
	}

	s.ListenAndServe()
}
