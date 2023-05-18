package main

import (
	"net/http"
	"testing"
)

func Testmain(m *testing.T) {
	go main()
}

type myHandler struct {
}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
