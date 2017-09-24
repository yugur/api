// Copyright 2017 The Yugur RESTful API Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// util provides additional logging/benchmarking tools.
package util

import "fmt"

// returns the type of v as a string
func Type(v interface{}) string {
    return fmt.Sprintf("%T", v)
}