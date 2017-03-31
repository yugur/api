// Copyright 2017 Nicholas Brown. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Endpoint provides HTTP/JSON handlers for database manipulation.
package endpoint

import (
  _ "github.com/lib/pq"
  "encoding/json"
  "database/sql"
  "log"
  "net/http"
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

// GetEntries retrieves all entries as a JSON object.
func GetEntries(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, http.StatusText(405), 405)
    return
  }

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
}

// GetEntry utilizes HTTP form values to reply a single entry JSON object.
func GetEntry(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, http.StatusText(405), 405)
    return
  }

  headword := r.FormValue("headword")
  if headword == "" {
    http.Error(w, http.StatusText(400), 400)
    return
  }

  row := db.QueryRow("SELECT * FROM entries WHERE headword = $1", headword)

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
}

// CreateEntry decodes the JSON body of a POST request to create a new entry.
func CreateEntry(w http.ResponseWriter, r *http.Request) {
  var e Entry

  if r.Method != "POST" {
    http.Error(w, http.StatusText(405), 405)
    return
  }

  if r.Body == nil {
    http.Error(w, http.StatusText(400), 400)
    return
  }

  err := json.NewDecoder(r.Body).Decode(&e)
  if err != nil {
    http.Error(w, err.Error(), 400)
    return
  }

  if e.Headword == "" || e.Definition == "" {
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

  fmt.Fprintf(w, "Entry %s created successfully (%d row affected)\n", e.Headword, rowsAffected)
}
