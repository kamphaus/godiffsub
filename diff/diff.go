package diff

import (
	"fmt"
	"errors"
)

// Arguments to the diff-sub algorithm
type Arguments struct {
	Src  []string // the files whose function, constant and variable declarations should be considered
	From []string // the files from where the considered declarations should be removed
	Verbose bool  // whether to output debug statements
	symbols map[string]struct{} // symbols found in src
}

func (a Arguments) DiffSub() error {
	if a.Verbose {
		for _, src := range a.Src {
			fmt.Println(fmt.Sprintf("Considering %s src file.", src))
		}
		for _, from := range a.From {
			fmt.Println(fmt.Sprintf("Considering %s from file.", from))
		}
	}
	return errors.New("not yet implemented")
}

