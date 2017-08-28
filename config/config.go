// Copyright 2017 The Yugur.io Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// config is a simple config management package.
package config

import (
  "os"
  "log"
  "encoding/json"
)

// Configuration values
type Values struct {
  Database struct {
    Host     string `json:"host"`
    User     string `json:"user"`
    Password string `json:"password"`
    Database string `json:"database"`
  }
  Host     string `json:"host"`
  Port     string `json:"port"`
  Keystore string `json:"keystore"`
  CORS     bool   `json:"cors"`
  Verbose  bool   `json:"verbose"`
}

// Demarshals the provided JSON object into a Values struct
func Load(file string) (config Values, err error) {
  var conf Values
  configFile, err := os.Open(file)
  defer configFile.Close()
  if err != nil {
    log.Fatal(err)
  }
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&conf)
  return conf, nil
}