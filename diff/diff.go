package diff

import (
	"errors"
	"fmt"
)

// Arguments to the diff-sub algorithm
type Arguments struct {
	Src  []string // the files whose function, constant and variable declarations should be considered
	From []string // the files from where the considered declarations should be removed
	Verbose bool  // whether to output debug statements
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
		fmt.Println("Parsing src files...")
	}
	a.readSymbols()
	if a.Verbose {
		fmt.Println("Found symbols:")
		a.printSymbols()
	}
	return errors.New("not yet implemented")
}
