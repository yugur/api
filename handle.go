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

  "github.com/gorilla/sessions"
  "github.com/yugur/api/crypto"
  "github.com/yugur/api/util"
)

//---------------------------------------------------------
//---- Database Structs
//---------------------------------------------------------

type User struct {
  UID      string `json:"uid"`
  Username string `json:"username"`
  Hash     string `json:"hash"`
}

type Tag struct {
  ID       string    `json:"id"`
  Name     string `json:"name"`
}

var store = sessions.NewCookieStore([]byte(conf.Keystore))

//---------------------------------------------------------
//---- Endpoint Handlers
//---------------------------------------------------------

//----
//---- General Handlers
//----

// statusHandler may be used to confirm the server's current status.
func statusHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    return
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

// indexHandler serves the root page. Can be ignored if you bring your own front end.
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

//----
//---- User Authentication
//----

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
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }

    session.Values["uid"] = user.UID
    err = session.Save(r, w)
    if err != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }

    http.Redirect(w, r, "/", 302)
  default:
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

//----
//---- Dictionary Handlers
//----

// searchHandler returns a collection of unique entries given some query 'q'.
func searchHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    var entries []*Entry

    query := r.FormValue("q")

    headwordResults, err := headwordSearch(query)
    if err != nil {
      headwordResults = nil
    }
    entries = append(entries, headwordResults...)

    tagResults, err := tagSearch(query)
    if err != nil {
      tagResults = nil
    }
    entries = append(entries, tagResults...)

    wordtypeResults, err := wordtypeSearch(query)
    if err != nil {
      wordtypeResults = nil
    }
    entries = append(entries, wordtypeResults...)

    definitionResults, err := definitionSearch(query)
    if err != nil {
      log.Println(err)
      wordtypeResults = nil
    }
    entries = append(entries, definitionResults...)

    entries = entrySet(entries...)

    response, err := asOutgoing(entries...)
    if err != nil {
      util.Error(util.Internal(w, r))
      return
    }
    
    json.NewEncoder(w).Encode(response)
  default:
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

// entryHandler provides Create, Read, Update and Delete access to entries.
func entryHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    // Serve the entry
    query := r.FormValue("q")

    entry, err := idSearch(query)
    if err != nil {
      util.Error(util.NotFound(w, r))
      break
    }

    response, err := asOutgoing(entry...)
    if err != nil {
      util.Error(util.Internal(w, r))
      break
    }
    
    json.NewEncoder(w).Encode(response)
  case http.MethodPost:
    // Create a new entry
    e := new(Entry)

    err := json.NewDecoder(r.Body).Decode(&e)
    if err != nil {
      util.Error(util.BadRequest(w, r))
      break
    }

    // Entry headwords cannot be nil
    if e.Headword == "" {
      util.Error(util.BadRequest(w, r))
      break
    }

    request, err := asIncoming(e)
    if err != nil {
      util.Error(util.BadRequest(w, r))
      break
    }

    _, err = insertEntry(request...)
    if err != nil {
      util.Error(util.BadRequest(w, r))
    }
  case http.MethodPut:
    // Update an existing entry
    e := new(Entry)

    err := json.NewDecoder(r.Body).Decode(&e)
    if err != nil {
      util.Error(util.BadRequest(w, r))
      break
    }

    // Entry headwords cannot be nil
    if e.Headword == "" {
      util.Error(util.BadRequest(w, r))
      break
    }

    request, err := asIncoming(e)
    if err != nil {
      util.Error(util.BadRequest(w, r))
      break
    }

    _, err = insertEntry(request...)
    if err != nil {
      util.Error(util.BadRequest(w, r))
    }
  case http.MethodDelete:
    // Remove an existing entry
    query := r.FormValue("q")
    if query == "" {
      util.Error(util.BadRequest(w, r))
      break
    }

    _, err := deleteEntry(query)
    if err != nil {
      util.Error(util.Internal(w, r))
    }
  default:
    // Unsupported method
    http.Error(w, http.StatusText(405), 405)
  }
}

// fetchHandler provides an index of the entire dictionary for testing purposes.
func fetchHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    entries, err := index()
    if err != nil {
      util.Error(util.Internal(w, r))
    }

    response, err := asOutgoing(entries...)
    if err != nil {
      util.Error(util.Internal(w, r))
      return
    }

    json.NewEncoder(w).Encode(response)
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
  rows, err := db.Query("SELECT * FROM entries WHERE headword LIKE $1 OR headword LIKE $2", lower, upper)
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

  response, err := asOutgoing(entries...)
  if err != nil {
    util.Error(util.Internal(w, r))
    return
  }

  json.NewEncoder(w).Encode(response)
}

// Search by category, returns all entries associated with the requested tag
func tagSearchHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    var entries []*Entry

    query := r.FormValue("q")

    entries, err := tagSearch(query)
    if err != nil {
      entries = nil
    }

    response, err := asOutgoing(entries...)
    if err != nil {
      util.Error(util.Internal(w, r))
      return
    }
    
    json.NewEncoder(w).Encode(response)
  case http.MethodPost:
    // Add a new tag relationship
    entryID := r.FormValue("entry")
    tagID, err := getTagID(r.FormValue("tag"))
    if err != nil {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    result, err := db.Exec("INSERT INTO entry_tags VALUES($1, $2)", tagID, entryID)
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }
    fmt.Fprintf(w, "Tag Id %d added to entry %d successfully (%d rows affected)\n", tagID, entryID, rowsAffected) 
  case http.MethodDelete:
    // Remove a tag relationship
    entryID := r.FormValue("entry")
    tagID, err := getTagID(r.FormValue("tag"))
    if err != nil {
      http.Error(w, http.StatusText(400), 400)
      return
    }

    result, err := db.Exec("DELETE FROM entry_tags WHERE tag_id = $1 AND entry_id = $2", tagID, entryID)
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
      fmt.Fprintf(w, "Tag %s doesn't exist in entry %s (%d rows affected)\n", tagID, entryID, rowsAffected)
    } else {
      fmt.Fprintf(w, "Tag %s deleted successfully from entry %s (%d rows affected)\n", tagID, entryID, rowsAffected);
    }
  }
}

//---------------------------------------------------------
//---- HTTP Helper Functions
//---------------------------------------------------------

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
