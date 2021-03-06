// Copyright 2013 Kyle Burton.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package grake

import (
	"fmt"
	"strings"
	"time"
)

// Controlls verbose output.
var debug = false // true


// TaskInfo defines a task.
type TaskInfo struct {
	Name        string
	Description string
	TaskArgs    []string
	Args        map[string]string

	Block      func(t *TaskInfo)
	Deps       []string
	Err        error
	Tstart     time.Time
	Tdepend    time.Time
	Tdepstart  time.Time
	Tend       time.Time
	DepsElapsd time.Duration
	Elapsd     time.Duration
	Completed  bool
	Running    bool
}

var lastDescription = "no-description"
var lastTask *TaskInfo
var namespace string = ""

// The name of the defailt task for your Grakefile
var DefaultTaskName = ""

// Task Registry type, name to TaskInfo.
type Taskmanager map[string]*TaskInfo

// Registry of tasks declared with Task()
var TaskManager Taskmanager

// Formatter for a TaskInfo.
func (self *TaskInfo) String() string {
	if debug {
		fmt.Printf("Converting to String\n")
	}

	deps := self.Deps
	if nil == deps {
		deps = make([]string, 0)
	}

	if debug {
		fmt.Printf("Converting to String, deps=%q\n", deps)
	}
	return fmt.Sprintf("TaskInfo{descr=\"%s\";name=\"%s\";block=%s;deps=%s}",
		self.Description,
		self.Name,
		self.Block,
		strings.Join(deps, ","))
}

// Set the description for the next task that will be declared.
func Desc(descr string) {
	lastDescription = descr
}

// Prepends the namespace to the name of the task being declared.
func NamespaceTaskName(name string) string {
	if "" != namespace {
		return namespace + ":" + name
	}

	return name
}

// Prepends the namespace to the name of the task being declared, any argument list is removed.
func NamespaceTaskNameNoArgs(name string) string {
	// strip off any arguments
	pos := strings.Index(name, "[")

	if -1 != pos {
		name = name[:pos]
	}

	if "" != namespace {
		return namespace + ":" + name
	}
	return name
}

// Declare a task.
func Task(name string, block func(*TaskInfo)) (t *TaskInfo) {
	t = &TaskInfo{
		Name:        NamespaceTaskName(name),
		Block:       block,
		Description: lastDescription,
	}
	lastDescription = "No Description Provided"
	lastTask = t
	TaskManager[NamespaceTaskNameNoArgs(name)] = t
	return t
}

// Declare the set of dependencies for the current task, each of these will be
// executed in parallel before the current task is executed.
func Depends(deps ...string) {
	lastTask.Deps = deps
	if debug {
		fmt.Printf("Task[%s] deps: %q\n", lastTask.Name, lastTask.Deps)
	}
}

// Invoke a task given a name and an array of arguments.
func InvokeTask(name string, args []string) {
	if t, ok := TaskManager[name]; ok {
		t.TaskArgs = args
		t.ArgsToMap()
		if t.Running {
			panic(fmt.Sprintf("Error: task depends on itself / circular dependency!: %s", name))
		}

		t.Running = true

		if t.Completed {
			if debug {
				fmt.Printf("InvokeTask[%s] task already completed\n", name)
			}
			return
		}

		ok := make(chan bool)

		t.Tdepstart = time.Now()

		// TODO: make sure a task can't depend on itself
		if nil != t.Deps {
			for idx, dep := range t.Deps {
				go func(name string, idx int, dep string) {
					if debug {
						fmt.Printf("InvokeTask[%s] execing dep: %d/%s\n", name, idx, dep)
					}
					InvokeTask(dep, make([]string, 0))
					ok <- true
				}(name, idx, dep)
			}

			if debug {
				fmt.Printf("InvokeTask[%s] waiting for deps\n", name)
			}
			for i := 0; i < len(t.Deps); i++ {
				<-ok
			}
		}

		t.Tdepend = time.Now()
		t.DepsElapsd = t.Tdepend.Sub(t.Tdepstart)

		if debug {
			fmt.Printf("InvokeTask[%s] executing this task\n", name)
		}

		t.Tstart = time.Now()
		t.Block(t)
		t.Tend = time.Now()
		t.Elapsd = t.Tend.Sub(t.Tstart)
		t.Completed = true

		if debug {
			fmt.Printf("InvokeTask[%s] completed in %q, deps took %q\n",
				name, t.Elapsd, t.DepsElapsd)
		}
	} else {
		panic(fmt.Sprintf("Error: no such task: %q", name))
	}
}

// Create the task manager.
func init() {
	TaskManager = make(Taskmanager)
}

// Typename for the namespace body/block.
type NamespaceFn func()

// Declare a namespace for tasks.
func Namespace(ns string, fn NamespaceFn) {
	var oldNs = namespace
	if "" != namespace {
		namespace = namespace + ":" + ns
	} else {
		namespace = ns
	}
	fn()
	namespace = oldNs
}

// Set the default task.
func Default(name string) {
	DefaultTaskName = name
}


// Print's a list of all the registered tasks.
func ShowTasks () {
	maxLen := 0

	for _, t := range TaskManager {
		l := len(t.Name)
		if l > maxLen {
			maxLen = l
		}
	}

	for _, t := range TaskManager {
		fmt.Printf("grake % *s   # %s\n", -1*maxLen, t.Name, t.Description)
	}
}

// Parse a task string into its name and arguments, eg:
//     foo[1,2,3]
// will become:
//     "foo"  ["1","2","3"]
func ParseTaskString(s string) (taskName string, args []string) {
	spos := strings.Index(s, "[")
	args = make([]string, 0)

	if spos == -1 {
		taskName = s
		return
	}

	epos := strings.Index(s, "]")

	if epos == -1 {
		panic(fmt.Sprintf("Error: invalidly formatted argument: %s", s))
	}

	taskName = s[:spos]
	args = strings.Split(s[spos+1:epos], ",")
	return
}

// Convert the task's arguments vector into a map by name.
func (self *TaskInfo) ArgsToMap() {
	_, args := ParseTaskString(self.Name)
	self.Args = make(map[string]string)
	for idx := 0; idx < len(self.TaskArgs) && idx < len(args); idx++ {
		arg := self.TaskArgs[idx]
		argName := args[idx]
		self.Args[argName] = arg
	}
}
