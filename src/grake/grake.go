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


