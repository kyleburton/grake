// Copyright 2013 Kyle Burton.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


// grake package of helper functions for the grake library
package grake

import (
	"fmt"
	f "grake/fileutils"
	g "grake/tasks"
	"os/exec"
)

/* TODO: how to have System that returns stdout, or a version that doesn't capture output?
   TODO: move this to a helper namespace
*/


// System executes a command and returns the standard output as a string. If
// the command execution fails, the function panic's
func System(cmd string, args ...string) (res string) {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}


