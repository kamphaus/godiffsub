// Package godiffsub contains the main function for running the executable.
//
// Installation
//
//     go get -u github.com/kamphaus/godiffsub
//
// Usage
//
//     godiffsub -src filea.go -from fileb.go
//
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kamphaus/godiffsub/diff"
	"github.com/kamphaus/godiffsub/program"
)

var stderr io.Writer = os.Stderr

type inputDataFlags []string

func (i *inputDataFlags) String() (s string) {
	for pos, item := range *i {
		s += fmt.Sprintf("Flag %d. %s\n", pos, item)
	}
	return
}

func (i *inputDataFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	versionFlag = flag.Bool("v", false, "print the version and exit")
	verboseFlag = flag.Bool("V", false, "print progress as comments")
	helpFlag    = flag.Bool("h", false, "print help information")
	srcFlags    inputDataFlags
	fromFlags   inputDataFlags
)

func init() {
	flag.Var(&srcFlags, "src", "Files whose functions, variables and constants should be considered.")
	flag.Var(&fromFlags, "from", "Files from which any functions, variables and constants having the same name as the considered ones should be removed.")
}

func main() {
	code := runCommand()
	if code != 0 {
		os.Exit(code)
	}
}

func runCommand() int {

	flag.Usage = func() {
		usage := "Usage: %s [<flags>]\n\n"
		usage += "Flags:\n"
		fmt.Fprintf(stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println(program.Version)
		return 0
	}

	if *helpFlag {
		flag.Usage()
		return 1
	}

	args := &diff.Arguments{
		Src:     srcFlags,
		From:    fromFlags,
		Verbose: *verboseFlag,
	}
	err := args.DiffSub()
	if err != nil {
		fmt.Fprintf(stderr, "Error performing diff-sub operation: %v\n", err)
		return 1
	}

	return 0
}
