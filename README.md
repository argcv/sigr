# Signal Register

Package `argcv/sigr` implements a simple signal proxy, which is used to manage the signal easier for multiple tasks.

You can add one or a few functions after interrupt/quit/.... signal comes and befure really quit.

Currently, the process sequence are not in order.

## Install

```bash
go get -u github.com/argcv/sigr
```

## Example

```go
package main

import (
	"fmt"
	"github.com/argcv/sigr"
	"time"
)

func main() {
	// register a new name
	// this is used to
	name1 := sigr.RegisterOnStopFuncAutoName(func() {
		fmt.Println("Hello, World! #1")
	})
	name2 := sigr.RegisterOnStopFuncAutoName(func() {
		fmt.Println("Hello, World! #2")
	})
	sigr.RegisterOnStopFunc("customized name", func() {
		fmt.Println("Hello, World! #3")
	})
	fmt.Println("name1:", name1, "name2:", name2)
	sigr.RegisterOnStopFunc(name1, func() {
		fmt.Println("Hello, World! #1 #overridden")
	})
	sigr.UnregisterOnStopFunc(name2)
	// default: true
	sigr.SetQuitDirectyl(true)
	// get verbose log
	sigr.VerboseLog()
	// without verbose log, your printing is still work here
	//sigr.NoLog()
	// you may type ctrl+C to excute the functions (name1 and "customized name")
	time.Sleep(time.Duration(20 * time.Second))
}

```

The output seems as follow:

```go
$ go run example.go
name1: __auto_1 name2: __auto_2
^C2017/10/26 19:06:17.804737 sigr.go:74: sig: [interrupt] Processing...
2017/10/26 19:06:17.804935 sigr.go:78: Processing task [__auto_1]
Hello, World! #1 #overridden
2017/10/26 19:06:17.804949 sigr.go:78: Processing task [customized name]
Hello, World! #3
2017/10/26 19:06:17.805069 sigr.go:85: sig: [interrupt] Processed, quitting directly
signal: interrupt
```

