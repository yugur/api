// Copyright 2017 The Yugur RESTful API Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Provides a dictionary entry type and related functions and methods.
package entry

type Entry struct {
  ID         string `json:"id"`
  Headword   string `json:"headword"`
  Wordtype   string `json:"wordtype"`
  Definition string `json:"definition"`

  Headword_Language   string `json:"hw_lang"`
  Definition_Language string `json:"def_lang"`
}

func Set(entries ...*Entry) []*Entry {
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

func (e1 *Entry) Equals(e2 *Entry) bool {
  return e1.ID                  == e2.ID &&
         e1.Headword            == e2.Headword &&
         e1.Wordtype            == e2.Wordtype &&
         e1.Definition          == e2.Definition &&
         e1.Headword_Language   == e2.Headword_Language &&
         e1.Definition_Language == e2.Definition_Language
}