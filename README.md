liner
========

golang library which contains a `io.Writer` for line processing.

Example
---------

Use `NewLiningWriter(lp LineProcessor, eh ErrorHandler) *LiningWriter` to get the `io.Writer`:

~~~~~~~~~~~~~~~~~~~~~~~~~~_example/cmdlogging.go
package main

import (
	"github.com/kui/liner"
	"log"
	"os/exec"
)

func main() {
	// use the default line processor which just print lines with log.Printf
	outw := liner.NewLiningWriter(nil, nil)
	defer outw.Close()

	// use a custom line processor which just print lines with log.Printf
	errw := liner.NewLiningWriter(func(line string) error {
		log.Printf("stderr: %s\n", line)
		return nil
	}, nil)
	defer errw.Close()

	c := exec.Command("bash", "-c", "seq 5; seq 6 10 >&2")
	c.Stdout = outw
	c.Stderr = errw

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
~~~~~~~~~~~~~~~~~~~~~~~~~~

In your terminal:

~~~~~~~~~~~~~~~~~~~~~~~~~bash
$ go run _example/cmdlogging.go
2016/01/31 16:23:23 1
2016/01/31 16:23:23 2
2016/01/31 16:23:23 3
2016/01/31 16:23:23 4
2016/01/31 16:23:23 5
2016/01/31 16:23:23 stderr: 6
2016/01/31 16:23:23 stderr: 7
2016/01/31 16:23:23 stderr: 8
2016/01/31 16:23:23 stderr: 9
2016/01/31 16:23:23 stderr: 10
~~~~~~~~~~~~~~~~~~~~~~~~~
