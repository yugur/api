// Copyright 2017 Nicholas Brown. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// yugur is a lightweight self-hosted dictionary platform.
package main

import (
  _ "github.com/lib/pq"
  "encoding/json"
  "database/sql"
  "net/http"
  "log"
  "fmt"
)

type Entry struct {
  Headword    string  `json:"headword"`
  Definition  string  `json:"definition"`
}

var db *sql.DB

func init() {
  var err error
  db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/dictionary")
  if err != nil {
    log.Fatal(err)
  }

  if err = db.Ping(); err != nil {
    log.Fatal(err)
  }
}

// StatusHandler may be used to confirm the server's current status.
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    w.Write([]byte("OK"))
  default:
    http.Error(w, http.StatusText(405), 405)
  }
})

// NotImplemented is called when a particular endpoint hasn't been provided
// with a handler function.
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    w.Write([]byte("Not Implemented"))
  default:
    http.Error(w, http.StatusText(405), 405)
  }
})

// EntryHandler is responsible for serving, adding, updating and removing
// entries from the dictionary database.
var EntryHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    // Serve the entry
    word := r.FormValue("q")
    if word == "" {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    row := db.QueryRow("SELECT * FROM entries WHERE headword = $1", word)

    entry := new(Entry)
    err := row.Scan(&entry.Headword, &entry.Definition)
    if err == sql.ErrNoRows {
      http.NotFound(w, r)
      return
    } else if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    json.NewEncoder(w).Encode(entry)
  case http.MethodPost:
    // Create a new entry
    var e Entry

    if r.Body == nil {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    err := json.NewDecoder(r.Body).Decode(&e)
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    if e.Headword == "" {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    result, err := db.Exec("INSERT INTO entries VALUES($1, $2)", e.Headword, e.Definition)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    fmt.Fprintf(w, "Entry %s created successfully (%d rows affected)\n", e.Headword, rowsAffected)
  case http.MethodPut:
    // Update an existing entry
    var e Entry

    if r.Body == nil {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    err := json.NewDecoder(r.Body).Decode(&e)
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    if e.Headword == "" {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    result, err := db.Exec("UPDATE entries SET definition = $1 WHERE headword = $2", e.Definition, e.Headword)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    fmt.Fprintf(w, "Entry %s updated successfully (%d rows affected)\n", e.Headword, rowsAffected)
  case http.MethodDelete:
    // Remove an existing entry
    word := r.FormValue("q")
    if word == "" {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    result, err := db.Exec("DELETE FROM entries WHERE headword = $1", word)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    fmt.Fprintf(w, "Entry %s deleted successfully (%d rows affected)\n", word, rowsAffected)
  default:
    http.Error(w, http.StatusText(405), 405)
  }
})

var IndexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    rows, err := db.Query("SELECT * FROM entries")
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }
    defer rows.Close()

    entries := make([]*Entry, 0)
    for rows.Next() {
      entry := new(Entry)
      err := rows.Scan(&entry.Headword, &entry.Definition)
      if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
      }
      entries = append(entries, entry)
    }
    if err = rows.Err(); err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    json.NewEncoder(w).Encode(entries)
  default:
    http.Error(w, http.StatusText(405), 405)
  }
})

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/status", StatusHandler)
  mux.HandleFunc("/entry", EntryHandler)
  mux.HandleFunc("/index", IndexHandler)
  mux.HandleFunc("/search", NotImplemented)
  err := http.ListenAndServe(":8080", mux)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
