# Overview

Inspried by, and draws heavily from: [makengo](https://github.com/remogatto/makengo). Which (at the time of this writing) does not seem to be comptabile with Go 1.1.

# Grake

A [rake](http://rake.rubyforge.org/)-like build tool implemented in [Go](http://golang.org/) with concurrent execution of task dependencies.


# Usage

You must have [go](http://golang.org/) installed and the `go` command must be available on your $PATH.

    git clone https://github.com/kyleburton/grake.git
    cd grake
    export PATH="$PATH:$HOME/grake/bin"

Create a `Grakefile`:

    // vim: ft=go
    
    import(
      "io/ioutil"
    )
    
    func WriteFile (f, s string) {
      ioutil.WriteFile(f, []byte(s), 0644)
    }
    
    g.Desc("not in a namespace")
    g.Task("hello", func (self *g.TaskInfo) {
      fmt.Printf("in task: %s\n", self.Name)
    })
    
    g.Namespace("main", func () {
      g.Desc("This is a test task")
      g.Task("hello", func (self *g.TaskInfo) {
        fmt.Printf("Hello! from: %s\n", self.Name)
      });
      g.Depends("hello")
    
      g.Namespace("two", func () {
        g.Desc("a task with arguments")
        g.Task("hasargs[a,b,c]", func (self *g.TaskInfo) {
          fmt.Printf("%s, self.Args: %q\n", self.Name, self.Args)
          fmt.Printf("%s, and I have arguments: a=%s, b=%s, c=%s\n", self.Name, self.Args["a"], self.Args["b"], self.Args["c"])
        })
        g.Depends("main:hello")
      })
    })
    
    g.Desc("make a file")
    g.Task("createfile", func (self *g.TaskInfo) {
      txt := 
    `this is 
    the 
    file contents`
      WriteFile("output.txt", txt)
    })
    
    
    g.Default("main:hello")
    


Run `grake`

    $ grake -T
    grake hello                     # not in a namespace
    grake main:hello                # This is a test task
    grake main:two:hasargs[a,b,c]   # a task with arguments
    grake createfile                # make a file
    $


# Limitations

* There are bugs, please help me fix them
* Documentation is incomplete
* I'd like to support 'once only' execution of tasks
* Also, support re-enabling and re-execution of tasks

# LICENSE

Copyright (c) 2013 Kyle Burton

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
