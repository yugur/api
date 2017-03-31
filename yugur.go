// Copyright 2017 Nicholas Brown. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// yugur is a lightweight self-hosted dictionary platform.
package main

import (
  _ "github.com/lib/pq"
  "github.com/yugur/api/endpoint"
  "net/http"
  "log"
)

func main() {
  http.HandleFunc("/entries", endpoint.GetEntries)
  http.HandleFunc("/entries/show", endpoint.GetEntry)
  http.HandleFunc("/entries/create", endpoint.CreateEntry)
  err := http.ListenAndServe(":3000", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
