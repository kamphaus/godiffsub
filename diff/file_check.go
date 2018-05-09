package diff

import (
	"os"
	"fmt"
	"errors"
)

func checkFile(file string) error {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("could not find file: %s", file)
		}
		return fmt.Errorf("file error: %v", err)
	}
	if file, err := os.Open(file); err != nil {
		return fmt.Errorf("could not open file: %v", err)
	} else {
		if err = file.Close(); err != nil {
			return fmt.Errorf("could not close file: %v", err)
		}
	}
	return nil
}

func (a Arguments) checkFiles() (err error) {
	for _, src := range a.Src {
		if a.Verbose {
			fmt.Println(fmt.Sprintf("Considering src file: %s", src))
		}
		if e := checkFile(src); e != nil {
			err = e
			if a.Verbose {
				fmt.Println(fmt.Sprintf("%v", e))
			}
		}
	}
	for _, from := range a.From {
		if a.Verbose {
			fmt.Println(fmt.Sprintf("Considering from file: %s", from))
		}
		if e := checkFile(from); e != nil {
			err = e
			if a.Verbose {
				fmt.Println(fmt.Sprintf("%v", e))
			}
		}
	}
	return
}

var (
	NotEnoughSrcFiles error
	NotEnoughFromFiles error
)

func init() {
	NotEnoughSrcFiles = errors.New("not enough src files")
	NotEnoughFromFiles = errors.New("not enough from files")
}
