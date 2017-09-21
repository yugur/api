// Copyright 2017 The Yugur RESTful API Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

type Entry struct {
  ID         string `json:"id"`
  Headword   string `json:"headword"`
  Wordtype   string `json:"wordtype"`
  Definition string `json:"definition"`

  Headword_Language   string `json:"hw_lang"`
  Definition_Language string `json:"def_lang"`
}

func entrySet(entries ...*Entry) []*Entry {
  var set []*Entry
  for _, i := range entries {
    unique := true
    for _, j := range set {
      if j.ID == i.ID {
        unique = false
        break
      }
    }
    if unique {
      set = append(set, i)
    }
  }
  return set
}