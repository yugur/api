package crypto

import "testing"

func TestHashPassword(t *testing.T) {
  tables := []struct {
    pwd string
  }{
    {"password"},
    {"1234567890"},
    {"!@#$%%^&*()_+-={}[]\\|/,.<>?~`"},
    {"비밀번호"},
  }

  for _, table := range tables {
    hash, err := HashPassword(table.pwd)
    if err != nil {
      t.Errorf("Error occurred on password hash.")
    }
    b := CompareHash(table.pwd, hash)
    if !b {
      t.Errorf(
        `Hash/password comparison failed.
        Expected: %t, got: %t`,
        true, b)
    }
  }
}