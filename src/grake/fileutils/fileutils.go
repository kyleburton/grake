// Copyright 2013 Kyle Burton.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package fileutils

import (
    "os"
)

func Exists ( filename string ) bool {
  if _, err := os.Stat(filename); os.IsNotExist(err) {
    return false
  }
  return true
}

