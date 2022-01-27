package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func debugServe() {
	if err := http.ListenAndServe("0.0.0.0:45677", http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}
