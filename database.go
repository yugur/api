// Copyright 2017 Nicholas Brown. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
  _ "github.com/lib/pq"
  "database/sql"
  "log"
)

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
