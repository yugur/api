// Copyright 2017 The Yugur.io Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// yugur is a lightweight self-hosted dictionary platform.
package main

import (
  "net/http"
  "log"
  "os"
  "fmt"
  "database/sql"

  "github.com/gorilla/handlers"
  "github.com/yugur/api/config"
)

// Global config values. This should only be changed via a call to config.Load(string)
var conf config.Values

// The primary database instance
var db *sql.DB

func init() {
  var err error
  fmt.Print("Loading configuration...")
  conf, err = config.Load("config/config.json")
  if err != nil {
    log.Fatal(err.Error())
  }
  fmt.Println("done!")

  fmt.Print("Preparing database...")
  db, err = sql.Open("postgres", "postgres://" +
    conf.Database.User     + ":" + 
    conf.Database.Password + "@" +
    conf.Database.Host     + "/" +
    conf.Database.Database)
  if err != nil {
    log.Fatal(err)
  }

  if err = db.Ping(); err != nil {
    log.Fatal(err)
  }
  fmt.Println("done!")
}

func main() {
  fmt.Print("Initialising mux...")
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
  fmt.Println("done!")

  fmt.Println("Listening on port " + conf.Port + "...")
  if conf.CORS {
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

    err := http.ListenAndServe(":" + conf.Port, handlers.CORS(originsOk, headersOk, methodsOk)(mux))
    if err != nil {
      log.Fatal("ListenAndServe: ", err)
    }
  } else {
    err := http.ListenAndServe(":" + conf.Port, handlers.LoggingHandler(os.Stdout, mux))
    if err != nil {
      log.Fatal("ListenAndServe: ", err)
    }
  }
}
