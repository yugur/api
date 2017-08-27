// Copyright 2017 The Yugur.io Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
  _ "github.com/lib/pq"
  "encoding/json"
  "database/sql"
  "net/http"
  "html/template"
  "fmt"
  "log"
  "strings"

  "github.com/gorilla/sessions"
  "github.com/yugur/api/crypto"
)

type User struct {
  Username string `json:"username"`
  Hash     string `json:"hash"`
}

type Entry struct {
  Headword   string `json:"headword"`
  Definition string `json:"definition"`
}

var store = sessions.NewCookieStore([]byte(config.Keystore))

// statusHandler may be used to confirm the server's current status.
func statusHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    w.Write([]byte("OK\n"))
  default:
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

// notImplemented is a simple stub for incomplete handlers.
func notImplemented(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    w.Write([]byte("Not Implemented"))
  default:
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

// indexHandler serves the root page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    session, err := store.Get(r, "username")
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
    }
    if val, ok := session.Values["username"].(string); ok {
      log.Println("Username cookie: ", val)
      switch val {
        case "": 
          http.Redirect(w, r, "/login", http.StatusFound)
        default:
          // Serve index page (demo)
          render(w, "templates/index.html", nil)
      }
    } else {
      http.Redirect(w, r, "/login", http.StatusFound)
    }
  default:
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

/* 
  registerHandler is responsible for dealing with user registration.
  On GET requests this will serve a basic registration page.
  On POST requests it will attempt to register a user as per
  the provided form values.
*/
func registerHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    // Serve registration page (demo)
    render(w, "templates/register.html", nil)
  case http.MethodPost:
    // Attempt to register new user
    // Parse form values
    err := r.ParseForm()
    if err != nil {
      http.Error(w, http.StatusText(403), 403)
    }
    username := r.PostFormValue("username")
    password := r.PostFormValue("password")

    // Check whether the user already exists in database
    var exists bool
    err = db.QueryRow("SELECT 1 FROM users WHERE username = $1", username).Scan(&exists)
    if err != nil && err != sql.ErrNoRows {
      fmt.Fprintf(w, "User %s already exists.\n", username)
      //http.Redirect(w, r, "/", http.StatusSeeOther)
      return
    }

    // Generate hash for new user
    hash, err := crypto.HashPassword(password)
    if err != nil {
      log.Println(err)
    }

    // Insert new user into database
    result, err := db.Exec("INSERT INTO users VALUES($1, $2)", username, hash)
    if err != nil {
      log.Printf(err.Error())
      http.Error(w, http.StatusText(500), 500)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    fmt.Fprintf(w, "Entry %s created successfully (%d rows affected)\n", username, rowsAffected)
  default:
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

/*
  loginHandler takes care of user login attempts.
  On GET this will serve a demo login page.
  On POST it will attempt to auth using the provided form values.
  In the event of successful authentication, the handler will respond with
  a valid session token. Future requests from the user should include this
  token until it expires or the user logs out.
*/
func loginHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    // Serve login page (demo)
    render(w, "templates/login.html", nil)
  case http.MethodPost:
    // Attempt to login with given credentials
    // Parse form values
    err := r.ParseForm()
    if err != nil {
      http.Error(w, http.StatusText(403), 403)
    }
    username := r.PostFormValue("username")
    password := r.PostFormValue("password")

    // Retrieve the matching user from database
    row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
    user := new(User)
    err = row.Scan(&user.Username, &user.Hash)
    if err == sql.ErrNoRows {
      http.NotFound(w, r)
      return
    } else if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    // Compare existing hash with given credentials
    valid := crypto.CompareHash(password, user.Hash)
    if !valid {
      log.Printf("Failed login attempt: username=%s, password=%s", username, password)
      return
    }
    log.Printf("Successful login attempt: username=%s, password=%s", username, password)
    // fmt.Fprintf(w, "Successfully logged in as user %s\n", username)
    session, err := store.Get(r, "username")
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    session.Values["username"] = username
    err = session.Save(r, w)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
    }

    http.Redirect(w, r, "/", 302)
  default:
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

// entryHandler is responsible for serving, adding, updating and removing
// entries from the dictionary database.
func entryHandler(w http.ResponseWriter, r *http.Request) {
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
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

// fetchHandler provides an index of the entire dictionary for testing purposes.
func fetchHandler(w http.ResponseWriter, r *http.Request) {
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
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}
// Search by letter, returns all entries starting with the requested letter
func letterSearchHandler(w http.ResponseWriter, r *http.Request) {
  word := r.FormValue("q")
  if word == "" {
    http.Error(w, http.StatusText(400), 400)
    return
  }
 /* query := ("SELECT * FROM entries WHERE headword LIKE '" + strings.ToLower(word) + "%%'" + "OR headword LIKE '" + strings.ToUpper(word) + "%%'")*/
  lower := strings.ToLower(word) + "%"
  upper := strings.ToUpper(word) + "%"
  rows, err := db.Query("SELECT * FROM entries WHERE headword LIKE $1 OR headword LIKE $2", lower, upper)
  if err == sql.ErrNoRows {
    http.NotFound(w, r)
    return
  }
  
  defer rows.Close()
  entries := make([]*Entry, 0)
  for rows.Next() {
    entry := new(Entry)
    err = rows.Scan(&entry.Headword, &entry.Definition)
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

// Search by category, returns all entries associated with the requested tag
func tagSearchHandler(w http.ResponseWriter, r *http.Request) {
  word := r.FormValue("q")
  if word == "" {
    http.Error(w, http.StatusText(400), 400)
    return
  }
  rows, err := db.Query("SELECT tags.headword, entries.definition FROM tags JOIN entries ON tags.headword = entries.headword WHERE tags.tag = $1", word)    
  if err == sql.ErrNoRows {
      http.NotFound(w, r)
      return
    }
    
    defer rows.Close()
    entries := make([]*Entry, 0)
    for rows.Next() {
      entry := new(Entry)
      err = rows.Scan(&entry.Headword, &entry.Definition)
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

// render uses html/template to serve a template page.
func render(w http.ResponseWriter, filename string, data interface{}) {
  tmpl, err := template.ParseFiles(filename)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  if err := tmpl.Execute(w, data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
