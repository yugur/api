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
	FatalError = "FatalError"
	InternalError = "InternalError"
)

func Error(err string, w http.ResponseWriter) {
	if w != nil {
		switch (err) {
		case InternalError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else {
		log.Fatal(err)
	}
}

func Fatal() (string, http.ResponseWriter) {
	return FatalError, nil
}

func Internal(w http.ResponseWriter) (string, http.ResponseWriter) {
	return InternalError, w

}