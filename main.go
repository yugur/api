// Copyright 2017 The Yugur.io Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// yugur is a lightweight self-hosted dictionary platform.
package main

import (
  "net/http"
  "log"
  "github.com/gorilla/handlers"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", indexHandler)
  mux.HandleFunc("/register", registerHandler)
  mux.HandleFunc("/login", loginHandler)
  mux.HandleFunc("/status", statusHandler)
  mux.HandleFunc("/entry", entryHandler)
  mux.HandleFunc("/fetch", fetchHandler)
  mux.HandleFunc("/search-letter", letterSearchHandler)
  mux.HandleFunc("/search-tag", tagSearchHandler)
  mux.HandleFunc("/search", notImplemented)

  headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
  originsOk := handlers.AllowedOrigins([]string{"*"})
  methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

  err := http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(mux))
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
