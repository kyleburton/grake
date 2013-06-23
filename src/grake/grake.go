// Copyright 2013 Kyle Burton.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package grake

import (
	"fmt"
	f "grake/fileutils"
	g "grake/tasks"
	"os/exec"
)

var debug = false

// TODO: how to have System that returns stdout, or a version that doesn't capture output?
// TODO: move this to a helper namespace
func System(cmd string, args ...string) (res string) {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}


