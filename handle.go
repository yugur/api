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
  "time"
  "strconv"

  "github.com/gorilla/sessions"
  "github.com/yugur/api/crypto"
)

type User struct {
  UID      string `json:"uid"`
  Username string `json:"username"`
  Hash     string `json:"hash"`
}

type Entry struct {
  ID         string `json:"id"`
  Headword   string `json:"headword"`
  Wordtype   string `json:"wordtype"`
  Definition string `json:"definition"`

  Headword_Language   string `json:"hw_lang"`
  Definition_Language string `json:"def_lang"`
}

type Tag struct {
  Name     string `json:"name"`
  Tag_Id   int    `json:"tag_id"`
  Entry_Id int    `json:"entry_id"`
}

var store = sessions.NewCookieStore([]byte(conf.Keystore))

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
    session, err := store.Get(r, "uid")
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
    }
    if val, ok := session.Values["uid"].(string); ok {
      log.Println("User ID cookie: ", val)
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

    // Read required fields
    username := r.PostFormValue("username")
    password := r.PostFormValue("password")
    email := r.PostFormValue("email")
    if username == "" || password == "" || email == "" {
      http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
      return
    }

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

    // Create a timestamp for the user's join date
    joindate := time.Now()

    // Insert new user into database
    result, err := db.Exec("INSERT INTO users(username, hash, email, dob, gender, joindate, language, fluency) VALUES($1, $2, $3, $4, $5, $6, $7, $8)", username, hash, email, nil, nil, joindate, nil, nil)
    if err != nil {
      log.Printf(err.Error())
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
    row := db.QueryRow("SELECT uid, username, hash FROM users WHERE username = $1", username)

    user := new(User)
    err = row.Scan(&user.UID, &user.Username, &user.Hash)
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
    session, err := store.Get(r, "uid")
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    session.Values["uid"] = user.UID
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
    id, err := strconv.ParseInt(r.FormValue("q"), 10, 64)
    if err != nil {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    row := db.QueryRow("SELECT * FROM entries WHERE entry_id = $1", id)

    entry := new(Entry)
    err = row.Scan(&entry.ID, &entry.Headword, &entry.Wordtype, &entry.Definition, &entry.Headword_Language, &entry.Definition_Language)
    if err == sql.ErrNoRows {
      http.NotFound(w, r)
      return
    } else if err != nil {
      log.Printf(err.Error())
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

    result, err := db.Exec("INSERT INTO entries (headword, wordtype, definition, hw_lang, def_lang) VALUES($1, $2, $3, $4, $5)", e.Headword, e.Wordtype, e.Definition, e.Headword_Language, e.Definition_Language)
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

    result, err := db.Exec("UPDATE entries SET (headword, wordtype, definition, hw_lang, def_lang) = ($2, $3, $4, $5, $6) WHERE entry_id = $1", e.ID, e.Headword, e.Wordtype, e.Definition, e.Headword_Language, e.Definition_Language)
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
    id := r.FormValue("q")
    if id == "" {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    result, err := db.Exec("DELETE FROM entries WHERE entry_id = $1", id)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    fmt.Fprintf(w, "Entry %s deleted successfully (%d rows affected)\n", id, rowsAffected)
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
      err := rows.Scan(&entry.ID, &entry.Headword, &entry.Wordtype, &entry.Definition, &entry.Headword_Language, &entry.Definition_Language)
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
  letter := r.FormValue("q")
  if letter == "" {
    http.Error(w, http.StatusText(400), 400)
    return
  }
  lower := strings.ToLower(letter) + "%"
  upper := strings.ToUpper(letter) + "%"
  rows, err := db.Query("SELECT headword, definition FROM entries WHERE headword LIKE $1 OR headword LIKE $2", lower, upper)
  if err == sql.ErrNoRows {
    http.NotFound(w, r)
    return
  }
  
  defer rows.Close()
  entries := make([]*Entry, 0)
  for rows.Next() {
    entry := new(Entry)
    err = rows.Scan(&entry.ID, &entry.Headword, &entry.Wordtype, &entry.Definition, &entry.Headword_Language, &entry.Definition_Language)
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
  switch r.Method {
  case http.MethodGet:
  id := r.FormValue("q")
  if id == "" {
    http.Error(w, http.StatusText(400), 400)
    return
  }
  rows, err := db.Query("SELECT entries.entry_id, entries.headword, entries.wordtype, entries.definition, entries.hw_lang, entries.def_lang FROM (SELECT * FROM entry_tags WHERE tag_id = $1) AS entry_tags JOIN entries ON entry_tags.entry_id = entries.entry_id", id)    
  if err == sql.ErrNoRows {
    http.NotFound(w, r)
    return
  }
      
    defer rows.Close()
    entries := make([]*Entry, 0)
    for rows.Next() {
      entry := new(Entry)
      err = rows.Scan(&entry.ID, &entry.Headword, &entry.Wordtype, &entry.Definition, &entry.Headword_Language, &entry.Definition_Language)
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
  case http.MethodPost:
    // Add a new tag relationship
    var t Tag

    if r.Body == nil {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    err := json.NewDecoder(r.Body).Decode(&t)
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    if t.Name == "" {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    result, err := db.Exec("INSERT INTO entry_tags VALUES($1, $2)", t.Tag_Id, t.Entry_Id)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }
    fmt.Fprintf(w, "Tag Id %d added to entry %d successfully (%d rows affected)\n", t.Tag_Id, t.Entry_Id, rowsAffected) 
  case http.MethodDelete:
    // Remove a tag relationship
    tag_id := r.FormValue("tag_id")
    entry_id := r.FormValue("entry_id")
    if tag_id == "" || entry_id == ""  {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    result, err := db.Exec("DELETE FROM entry_tags WHERE tag_id = $1 AND entry_id = $2", tag_id, entry_id)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    if (rowsAffected == 0) {
      fmt.Fprintf(w, "Tag %s doesn't exist (%d rows affected)\n", tag_id, rowsAffected)
    } else {
      fmt.Fprintf(w, "Tag %s deleted successfully (%d rows affected)\n", tag_id, rowsAffected);
    }
  }
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
