// Copyright 2017 Nicholas Brown. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// yugur is a lightweight self-hosted dictionary platform.
package main

import (
  "net/http"
  "log"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", index)
  mux.HandleFunc("/register", register)
  mux.HandleFunc("/login", login)
  mux.HandleFunc("/status", statusHandler)
  mux.HandleFunc("/entry", entryHandler)
  mux.HandleFunc("/fetch", fetchHandler)
  mux.HandleFunc("/search", notImplemented)
  err := http.ListenAndServe(":3000", mux)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
