// Copyright 2017 Nicholas Brown. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// yugur is a lightweight self-hosted dictionary platform.
package main

import (
  "net/http"
  "log"
  "os"
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
  mux.HandleFunc("/search", notImplemented)
  err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, mux))
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
