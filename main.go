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
  "encoding/json"
  "database/sql"
  "github.com/gorilla/handlers"
)

var config Config

type Config struct {
  Database struct {
    Host     string `json:"host"`
    User     string `json:"user"`
    Password string `json:"password"`
    Database string `json:"database"`
  }
  Host     string `json:"host"`
  Port     string `json:"port"`
  Keystore string `json:"keystore"`
  CORS     bool   `json:"cors"`
  Verbose  bool   `json:"verbose"`
}

func LoadConfiguration(file string) (config Config, err error) {
  var conf Config
  configFile, err := os.Open(file)
  defer configFile.Close()
  if err != nil {
    fmt.Println(err.Error())
  }
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&conf)
  return conf, nil
}

// The primary database instance
var db *sql.DB

func main() {
  var err error
  fmt.Print("Loading configuration...")
  config, err = LoadConfiguration("config/config.json")
  if err != nil {
    log.Fatal(err.Error())
  }
  fmt.Println("done!")

  fmt.Print("Preparing database...")
  db, err = sql.Open("postgres", "postgres://" +
    config.Database.User     + ":" + 
    config.Database.Password + "@" +
    config.Database.Host     + "/" +
    config.Database.Database)
  if err != nil {
    log.Fatal(err)
  }

  if err = db.Ping(); err != nil {
    log.Fatal(err)
  }
  fmt.Println("done!")

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

  fmt.Println("Listening on port " + config.Port + "...")
  if config.CORS {
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
    err = http.ListenAndServe(":" + config.Port, handlers.CORS(originsOk, headersOk, methodsOk)(mux))
  } else {
    err = http.ListenAndServe(":" + config.Port, handlers.LoggingHandler(os.Stdout, mux))
  }
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
