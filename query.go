// Copyright 2017 The Yugur.io Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
  "database/sql"
)
//---------------------------------------------------------
//---- Search Queries
//---------------------------------------------------------

func index() ([]*Entry, error) {
  entries := make([]*Entry, 0)

  query := `SELECT *
            FROM entries`

  rows, err := db.Query(query)
  if err != nil {
    return entries, err
  }
  defer rows.Close()

  for rows.Next() {
    entry := new(Entry)

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
  if err = rows.Err(); err != nil {
    return entries, err
  }

  return entries, nil
}

// idSearch returns the matching entries for each id.
// Raises sql.ErrNoRows if there is no matching entry for any provided id.
func idSearch(ids ...string) ([]*Entry, error) {
  var errNoRows error
  entries := make([]*Entry, 0)

  for _, id := range ids {
    e := new(Entry)

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

func headwordSearch(word string) ([]*Entry, error) {
  entries := make([]*Entry, 0)

  rows, err := db.Query("SELECT * FROM entries WHERE headword = $1", word)
  if err != nil {
    return entries, err
  }
  defer rows.Close()

  for rows.Next() {
    entry := new(Entry)

    err = rows.Scan(&entry.ID, &entry.Headword, &entry.Wordtype, &entry.Definition, &entry.Headword_Language, &entry.Definition_Language)
    if err != nil {
      return entries, err
    }
    entries = append(entries, entry)
  }
  if err = rows.Err(); err != nil {
    return entries, err
  }

  return entries, nil
}

func tagSearch(tag string) ([]*Entry, error) {
  entries := make([]*Entry, 0)

  tagID, err := getTagID(tag)
  if err != nil {
    return entries, err
  }

  rows, err := db.Query("SELECT entries.entry_id, entries.headword, entries.wordtype, entries.definition, entries.hw_lang, entries.def_lang FROM (SELECT * FROM entry_tags WHERE tag_id = $1) AS entry_tags JOIN entries ON entry_tags.entry_id = entries.entry_id", tagID)    
  if err != nil {
    return entries, err
  }
  defer rows.Close()

  for rows.Next() {
    entry := new(Entry)

    err = rows.Scan(&entry.ID, &entry.Headword, &entry.Wordtype, &entry.Definition, &entry.Headword_Language, &entry.Definition_Language)
      if err != nil {
        return entries, err
      }
      entries = append(entries, entry)
    }
    if err = rows.Err(); err != nil {
      return entries, err
    }

    return entries, nil
}

//---------------------------------------------------------
//---- Executable Queries
//---------------------------------------------------------

func insertEntry(entries ...*Entry) (int64, error) {
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

// Given a variadic Entry(s) with database identifiers,
// returns list of same entries with human names instead
func asOutgoing(entries ...*Entry) ([]*Entry, error) {
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

// Given a variadic Entry(s) with human names,
// returns list of same entries with database identifiers instead
func asIncoming(entries ...*Entry) ([]*Entry, error) {
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