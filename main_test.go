package main

import (
	"log"
	"net/http"
	"testing"

	_ "net/http/pprof"
)

func TestPprof(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:45678", http.DefaultServeMux))
	}()
	main()
}
