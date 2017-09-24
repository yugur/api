package entry

import "testing"

func TestEquals(t *testing.T) {
  tables := []struct {
    name string
    e1 Entry
    e2 Entry
    b bool
  }{
    {
      "positive",
      Entry{"1", "dog", "noun", "man's best friend", "en-AU", "en-AU"},
      Entry{"1", "dog", "noun", "man's best friend", "en-AU", "en-AU"},
      true,
    },
    {
      "negative",
      Entry{"1", "dog", "noun", "man's best friend", "en-AU", "en-AU"},
      Entry{"2", "cat", "noun", "not man's best friend", "en-AU", "en-AU"},
      false,
    },
    {
      "nil positive",
      Entry{},
      Entry{},
      true,
    },
    {
      "nil negative",
      Entry{},
      Entry{"1", "dog", "noun", "man's best friend", "en-AU", "en-AU"},
      false,
    },
    {
      "partial negative",
      Entry{"1", "dog", "noun", "man's best friend", "en-AU", "en-AU"},
      Entry{"1", "dg", "noun", "man's best friend", "en-AU", "en-AU"},
      false,
    },
  }

  for _, table := range tables {
    b1 := table.e1.Equals(&table.e2)
    b2 := table.e2.Equals(&table.e1)
    if b1 != table.b || b2 != table.b {
      t.Errorf(
        `Wrong equality for table %q
        Expected: b1:%t == b2:%t, got: b1:%t == b2:%t.`,
        table.name, table.b, table.b, b1, b2)
    }
  }
}

func TestSet(t *testing.T) {
  entry := new(Entry)
  entry.ID = "1"
  entry.Headword = "dog"
  entry.Wordtype = "noun"
  entry.Definition = "man's best friend"
  entry.Headword_Language = "en-AU"
  entry.Definition_Language = "en-AU"

  secondEntry := new(Entry)
  secondEntry.ID = "1"
  secondEntry.Headword = "dog"
  secondEntry.Wordtype = "noun"
  secondEntry.Definition = "man's best friend"
  secondEntry.Headword_Language = "en-AU"
  secondEntry.Definition_Language = "en-AU"


  slice := []*Entry{entry}
  slice = append(slice, entry)
  slice = append(slice, secondEntry)
  result := Set(slice...)

  for _, a := range slice {
    count := 0
    for _, b := range result {
      if a.Equals(b) {
        count++
      }
    }
    if count != 1 {
      t.Errorf("entrySet failed to create a set. Expected: %d occurrence(s), got %d occurrence(s).", 1, count)
    }
  }
}