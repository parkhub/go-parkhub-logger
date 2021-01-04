# go-parkhub-logger

This package provides a singular interface to create logs as well as filtering them out based on level.  It also provides two types of formatting json, pretty.  This logger doesn't ship any logs.

## Installing

### ([godoc](//godoc.org/github.com/parkhub/go-parkhub-logger))

    go get github.com/parkhub/go-parkhub-logger

## Setting up your logger
```
package main

import (
	log "github.com/parkhub/go-parkhub-logger"
)

func main() {
	log.SetupLogger(LogLevelInfo, []string{"test"}, false, Pretty)

	// Print info statement
	log.Infoln("This is an info statement.")

	// Print info statement with data
	type testStruct struct {
		Name string
		Kind string
	}

	test := &testStruct{
		Name: "Logan",
		Kind: "Log",
	}

	log.Infod("This is some text", test)

	// Print debug text
	log.Debugln("This is a debug statement.")

	// Print debug text with additional data.
	log.Debugf("This is a debug statement %d.", 10000)

	// Print info text
	log.Infoln("This is an info statement.")

	// Print info text with additional data.
	log.Infof("This is an info statement %d.", 10000)

	// Print warn text
	log.Warnln("This is a warning.")

	// Print warn text with additional data.
	log.Warnf("This is a warning %d.", 10000)

	// Print error text
	log.Errorln("This is an error.")

	// Print error text with additional data.
	log.Errorf("This is an error %d.", 10000)

	// Print fatal text
	log.Fatalln("This is an error.")

	// Print fatal text with additional data.
	log.Fatalf("This is an error %d.", 10000)
}
```