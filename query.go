// Copyright 2017 The Yugur RESTful API Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
  "database/sql"

  d "github.com/yugur/api/entry"
)
//---------------------------------------------------------
//---- Search Queries
//---------------------------------------------------------

func index() ([]*d.Entry, error) {
  query := `SELECT *
            FROM entries`

  rows, err := db.Query(query)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  entries, err := scanRows(rows)
  if err != nil {
    return entries, err
  }

  return entries, nil
}

// idSearch returns the matching entries for each id.
// Raises sql.ErrNoRows if there is no matching entry for any provided id.
func idSearch(ids ...string) ([]*d.Entry, error) {
  var errNoRows error
  entries := make([]*d.Entry, 0)

  for _, id := range ids {
    e := new(d.Entry)

    row := db.QueryRow("SELECT * FROM entries WHERE entry_id = $1", id)
    err := row.Scan(&e.ID, &e.Headword, &e.Wordtype, &e.Definition, &e.Headword_Language, &e.Definition_Language)
    if err != nil {
      errNoRows = err
      continue
    }
    entries = append(entries, e)
  }
  return entries, errNoRows
}

func headwordSearch(word string) ([]*d.Entry, error) {
  rows, err := db.Query("SELECT * FROM entries WHERE headword = $1", word)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  entries, err := scanRows(rows)
  if err != nil {
    return entries, err
  }

  return entries, nil
}

func tagSearch(tag string) ([]*d.Entry, error) {
  tagID, err := getTagID(tag)
  if err != nil {
    return nil, err
  }

  rows, err := db.Query("SELECT entries.entry_id, entries.headword, entries.wordtype, entries.definition, entries.hw_lang, entries.def_lang FROM (SELECT * FROM entry_tags WHERE tag_id = $1) AS entry_tags JOIN entries ON entry_tags.entry_id = entries.entry_id", tagID)    
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  entries, err := scanRows(rows)
  if err != nil {
    return entries, err
  }

  return entries, nil
}

func wordtypeSearch(wordtype string) ([]*d.Entry, error) {
  id, err := getWordtypeID(wordtype)
  if err != nil {
    return nil, err
  }

  query := `SELECT * FROM entries
            WHERE wordtype = $1`
  rows, err := db.Query(query, id)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  entries, err := scanRows(rows)
  if err != nil {
    return entries, err
  }

  return entries, nil
}

func definitionSearch(token string) ([]*d.Entry, error) {
  if token == "" {
    return nil, nil
  }

  query := `SELECT * FROM entries
            WHERE STRPOS(definition, $1) > 0`
  rows, err := db.Query(query, token)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  entries, err := scanRows(rows)
  if err != nil {
    return entries, err
  }

  return entries, nil
}

//---------------------------------------------------------
//---- Executable Queries
//---------------------------------------------------------

func insertEntry(entries ...*d.Entry) (int64, error) {
  var rowsAffected int64
  for _, entry := range entries {
    var query string
    var result sql.Result
    var err error

    if entry.ID == "" {
      query = `INSERT INTO entries (headword, wordtype, definition, hw_lang, def_lang) 
                VALUES($1, $2, $3, $4, $5)`
      result, err = db.Exec(
        query,
        entry.Headword,
        entry.Wordtype,
        entry.Definition,
        entry.Headword_Language,
        entry.Definition_Language)
    } else {
      query = `UPDATE entries
                SET headword = $1, wordtype = $2, definition = $3, hw_lang = $4, def_lang = $5
                WHERE entry_id = $6`
      result, err = db.Exec(
        query,
        entry.Headword,
        entry.Wordtype,
        entry.Definition,
        entry.Headword_Language,
        entry.Definition_Language,
        entry.ID)
    }

    if err != nil {
      return rowsAffected, err
    }

    r, err := result.RowsAffected()
    if err != nil {
      return rowsAffected, err
    }
    rowsAffected += r
  }
  return rowsAffected, nil
}

func deleteEntry(ids ...string) (int64, error) {
  var rowsAffected int64
  for _, id := range ids {
    query := `DELETE FROM entries
              WHERE entry_id = $1`
    result, err := db.Exec(query, id)
    if err != nil {
      return rowsAffected, err
    }

    r, err := result.RowsAffected()
    if err != nil {
      return rowsAffected, err
    }
    rowsAffected += r
  }
  return rowsAffected, nil
}

//---------------------------------------------------------
//---- Helper Functions
//---------------------------------------------------------

func getTagID(tag string) (string, error) {
  var result string

  row := db.QueryRow("SELECT tag_id FROM tags WHERE name = $1", tag)
  err := row.Scan(&result)

  return result, err
}

func getTagName(id string) (string, error) {
  var result string

  row := db.QueryRow("SELECT name FROM tags WHERE tag_id = $1", id)
  err := row.Scan(&result)

  return result, err
}

func getWordtypeID(name string) (string, error) {
  var result string

  row := db.QueryRow("SELECT wordtype_id FROM wordtypes WHERE name = $1", name)
  err := row.Scan(&result)

  return result, err
}

func getWordtypeName(id string) (string, error) {
  var result string

  row := db.QueryRow("SELECT name FROM wordtypes WHERE wordtype_id = $1", id)
  err := row.Scan(&result)

  return result, err
}

func getLocaleID(code string) (string, error) {
  var result string

  row := db.QueryRow("SELECT lang_id FROM languages WHERE code = $1", code)
  err := row.Scan(&result)

  return result, err
}

func getLocaleCode(id string) (string, error) {
  var result string

  row := db.QueryRow("SELECT code FROM languages WHERE lang_id = $1", id)
  err := row.Scan(&result)

  return result, err
}

// Given a variadic d.Entry(s) with database identifiers,
// returns list of same entries with human names instead
func asOutgoing(entries ...*d.Entry) ([]*d.Entry, error) {
  for _, entry := range entries {
    wordtype, err := getWordtypeName(entry.Wordtype)
    if err != nil {
      return entries, err
    }
    headwordLanguage, err := getLocaleCode(entry.Headword_Language)
    if err != nil {
      return entries, err
    }
    definitionLanguage, err := getLocaleCode(entry.Definition_Language)
    if err != nil {
      return entries, err
    }
    entry.Wordtype = wordtype
    entry.Headword_Language = headwordLanguage
    entry.Definition_Language = definitionLanguage
  }
  return entries, nil
}

// Given a variadic d.Entry(s) with human names,
// returns list of same entries with database identifiers instead
func asIncoming(entries ...*d.Entry) ([]*d.Entry, error) {
  for _, entry := range entries {
    wordtype, err := getWordtypeID(entry.Wordtype)
    if err != nil {
      return entries, err
    }
    headwordLanguage, err := getLocaleID(entry.Headword_Language)
    if err != nil {
      return entries, err
    }
    definitionLanguage, err := getLocaleID(entry.Definition_Language)
    if err != nil {
      return entries, err
    }
    entry.Wordtype = wordtype
    entry.Headword_Language = headwordLanguage
    entry.Definition_Language = definitionLanguage
  }
  return entries, nil
}

func scanRows(rows *sql.Rows) ([]*d.Entry, error) {
  var entries []*d.Entry

  for rows.Next() {
    entry := new(d.Entry)

    err := rows.Scan(
      &entry.ID,
      &entry.Headword,
      &entry.Wordtype,
      &entry.Definition,
      &entry.Headword_Language,
      &entry.Definition_Language)
    if err != nil {
      return entries, err
    }
    entries = append(entries, entry)
  }
  if err := rows.Err(); err != nil {
    return entries, err
  }
  return entries, nil
}