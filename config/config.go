// Copyright 2017 The Yugur RESTful API Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// config is a simple config management package.
package config

import (
  "os"
  "log"
  "encoding/json"
)

type Endpoint struct {
  Path   string `json:"path"`
  Enable bool   `json:"enable"`
}

// Configuration values
type Values struct {
  Database struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
    Name     string `json:"name"`
  }

  Host     string `json:"host"`
  Port     int    `json:"port"`
  Keystore string `json:"keystore"`
  CORS     bool   `json:"cors"`
  Verbose  bool   `json:"verbose"`

  Endpoints struct {
    Index    Endpoint
    Status   Endpoint
    Search   Endpoint
    Entry    Endpoint
    Register Endpoint
    Login    Endpoint
    Tag      Endpoint
    Fetch    Endpoint
    Random   Endpoint
  }
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