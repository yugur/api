// Copyright 2017 The Yugur.io Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// util provides additional logging/benchmarking tools.
package util

import (
	"net/http"
	"log"
)

const (
	// Internal API error codes
	FatalError = 1
)

func getRequestMessage(r *http.Request) string {
	return r.Method + " " + r.URL.String()
}

func Error(e int, msg string, w http.ResponseWriter) {
	if w != nil {
		http.Error(w, http.StatusText(e), e)
		log.Println(msg)
	} else {
		switch (e) {
		case FatalError:
			log.Fatal(msg)
		default:
			log.Println(msg)
		}
	}
}

func Fatal() (int, string, http.ResponseWriter) {
	return FatalError, "Fatal Error", nil
}

func Internal(w http.ResponseWriter, r *http.Request) (int, string, http.ResponseWriter) {
	return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError) + " (" + getRequestMessage(r) + ")", w
}

func BadRequest(w http.ResponseWriter, r *http.Request) (int, string, http.ResponseWriter) {
	return http.StatusBadRequest, http.StatusText(http.StatusBadRequest) + " (" + getRequestMessage(r) + ")", w
}