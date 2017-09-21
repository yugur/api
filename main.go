// Copyright 2017 The Yugur RESTful API Authors. All rights reserved.
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
  psqlInfo := fmt.Sprintf(
    "host=%s port=%d user=%s " + 
    "password=%s dbname=%s sslmode=disable",
    conf.Database.Host,
    conf.Database.Port,
    conf.Database.User,
    conf.Database.Password,
    conf.Database.Name)

  db, err = sql.Open("postgres", psqlInfo)
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

  if conf.Endpoints.Index.Enable {
    mux.HandleFunc(conf.Endpoints.Index.Path, indexHandler)
  }
  if conf.Endpoints.Status.Enable {
    mux.HandleFunc(conf.Endpoints.Status.Path, statusHandler)
  }
  if conf.Endpoints.Search.Enable {
    mux.HandleFunc(conf.Endpoints.Search.Path, searchHandler)
  }
  if conf.Endpoints.Entry.Enable {
    mux.HandleFunc(conf.Endpoints.Entry.Path, entryHandler)
  }
  if conf.Endpoints.Register.Enable {
    mux.HandleFunc(conf.Endpoints.Register.Path, registerHandler)
  }
  if conf.Endpoints.Login.Enable {
    mux.HandleFunc(conf.Endpoints.Login.Path, loginHandler)
  }
  if conf.Endpoints.Tag.Enable {
    mux.HandleFunc(conf.Endpoints.Tag.Path, tagSearchHandler)
  }
  if conf.Endpoints.Fetch.Enable {
    mux.HandleFunc(conf.Endpoints.Fetch.Path, fetchHandler)
  }
  if conf.Endpoints.Random.Enable {
    mux.HandleFunc(conf.Endpoints.Random.Path, notImplemented)
  }
  mux.HandleFunc("/letter", letterSearchHandler)
  fmt.Println("done!")

  fmt.Printf("The API is running at http://%s:%d/\n", conf.Host, conf.Port)
  if conf.CORS {
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

    err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), handlers.CORS(originsOk, headersOk, methodsOk)(mux))
    if err != nil {
      log.Fatal("ListenAndServe: ", err)
    }
  } else {
    err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), handlers.LoggingHandler(os.Stdout, mux))
    if err != nil {
      log.Fatal("ListenAndServe: ", err)
    }
  }
}
