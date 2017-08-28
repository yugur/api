// Copyright 2017 The Yugur.io Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// crypto provides cryptographic support such as hashing functions.
package crypto

import (
  "time"
  "golang.org/x/crypto/bcrypt"
  "github.com/yugur/api/util"
)

func HashPassword(password string) (string, error) {
  defer util.TrackTime(time.Now(), "HashPassword")
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

func CompareHash(password, hash string) bool {
  defer util.TrackTime(time.Now(), "CompareHash")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
