package diff

import (
	"errors"
	"fmt"
	"io"
)

// Arguments to the diff-sub algorithm
type Arguments struct {
	Src  []string // the files whose function, constant and variable declarations should be considered
	From []string // the files from where the considered declarations should be removed
	Verbose bool  // whether to output debug statements
	Stdout  io.Writer // where to write the debug statements
	symbols map[string]struct{} // symbols found in src
}

func (a *Arguments) DiffSub() error {
	if len(a.Src) == 0 {
		return NotEnoughSrcFiles
	}
	if len(a.From) == 0 {
		return NotEnoughFromFiles
	}
	if err := a.checkFiles(); err != nil {
		return errors.New("could not read all files")
	}
	if a.Verbose {
		fmt.Fprintf(a.Stdout, "Parsing src files...\n")
	}
	a.readSymbols()
	if a.Verbose {
		fmt.Fprintf(a.Stdout, "Found symbols:\n")
		a.printSymbols()
		fmt.Fprintf(a.Stdout, "Removing duplicate symbols...\n")
	}
	total, err := a.removeSymbols()
	if a.Verbose && len(a.From) > 1 {
		fmt.Fprintf(a.Stdout, "Removed total number of duplicate symbols: %v\n", total)
	}
	return err
}
