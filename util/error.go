// Copyright 2017 The Yugur RESTful API Authors. All rights reserved.
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

//---------------------------------------------------------
//---- API Errors 
//---------------------------------------------------------

func Fatal() (int, string, http.ResponseWriter) {
	return FatalError, "Fatal Error", nil
}

//---------------------------------------------------------
//---- HTTP Status Codes
//---------------------------------------------------------

//----
//---- 2xx
//----

// HTTP 200 OK
func OK(w http.ResponseWriter, r *http.Request) (int, string, http.ResponseWriter) {
	return http.StatusOK, http.StatusText(http.StatusOK) + " (" + getRequestMessage(r) + ")", w
}

//----
//---- 4xx
//----

// HTTP 400 Bad Request
func BadRequest(w http.ResponseWriter, r *http.Request) (int, string, http.ResponseWriter) {
	return http.StatusBadRequest, http.StatusText(http.StatusBadRequest) + " (" + getRequestMessage(r) + ")", w
}

// HTTP 404 Not Found
func NotFound(w http.ResponseWriter, r *http.Request) (int, string, http.ResponseWriter) {
	return http.StatusNotFound, http.StatusText(http.StatusNotFound) + " (" + getRequestMessage(r) + ")", w
}

//----
//---- 5xx
//----

// HTTP 500 Internal Server Error
func Internal(w http.ResponseWriter, r *http.Request) (int, string, http.ResponseWriter) {
	return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError) + " (" + getRequestMessage(r) + ")", w
}

// HTTP 501 Not Implemented
func NotImplemented(w http.ResponseWriter, r *http.Request) (int, string, http.ResponseWriter) {
	return http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented) + " (" + getRequestMessage(r) + ")", w
}