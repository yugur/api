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
  mux.HandleFunc("/status", StatusHandler)
  mux.HandleFunc("/entry", EntryHandler)
  mux.HandleFunc("/index", IndexHandler)
  mux.HandleFunc("/search", NotImplemented)
  err := http.ListenAndServe(":3000", mux)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
