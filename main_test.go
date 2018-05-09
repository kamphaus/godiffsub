package main

import (
	"testing"
	"io/ioutil"
	"path/filepath"
	"os"
	"io"
	"github.com/kamphaus/godiffsub/diff"
	"path"
	"bytes"
	"github.com/kamphaus/godiffsub/util"
	"strings"
)

const testDir = "./tests"

// TestDiffSub tests the godiffsub algorithm for each set in the tests directory.
// The files are copied into a temp directory.
// Files ending in .src are renamed .go and taken as src argument.
// Files ending in .from are renamed .go and taken as from argument, after execution
// of the algorithm they are compared to the file ending in .dst for equality.
// The stdout output of the algorithm is compared to the content of the out.txt file.
func TestDiffSub(t *testing.T) {
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Error(err)
	}
	var tests []string
	for _, f := range files {
		if f.IsDir() {
			tests = append(tests, f.Name())
		}
	}
	for _, testName := range tests {
		test := testName
		t.Run(testName, func(t *testing.T) {
			testDir := filepath.Join(testDir, test)
			args := prepareTest(t, testDir)
			args.testName = test
			defer os.RemoveAll(args.tempDir) // clean up
			runTest(t, args)
		})
	}
}

type diffTest struct {
	*diff.Arguments
	tempDir      string
	testName     string
	mapFrom2Dest map[string]string
	output       string
}

func runTest(t *testing.T, test *diffTest) {
	err := test.DiffSub()
	if err != nil {
		t.Error(err)
	}
	for _, from := range test.From {
		dest := test.mapFrom2Dest[from]
		compareFiles(t, from, dest)
	}
	if outBuf, ok := test.Stdout.(*bytes.Buffer); ok && test.output != "" {
		out, err := ioutil.ReadFile(test.output)
		if err != nil {
			t.Error(err)
			return
		}
		outStr := outBuf.String()
		outStr = strings.Replace(outStr, test.tempDir, "tests/"+test.testName, -1)
		if outStr != string(out) {
			t.Errorf(util.ShowDiff(outStr, string(out)))
		}
	}
}

func compareFiles(t *testing.T, a string, b string) {
	aStr, err := ioutil.ReadFile(a)
	if err != nil {
		t.Error(err)
		return
	}
	bStr, err := ioutil.ReadFile(b)
	if err != nil {
		t.Error(err)
		return
	}
	if string(aStr) != string(bStr) {
		t.Errorf(util.ShowDiff(string(aStr), string(bStr)))
	}
}

func prepareTest(t *testing.T, testDir string) (a *diffTest) {
	dir, err := ioutil.TempDir("", "godiffsub-test")
	if err != nil {
		t.Fatalf("Cannot create temp folder: %v", err)
	}
	a = &diffTest{
		Arguments: &diff.Arguments{
			Verbose: true,
			Stdout: &bytes.Buffer{},
		},
		mapFrom2Dest: make(map[string]string),
		tempDir: dir,
	}
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Error(err)
	}
	for _, f := range files {
		dstFile := filepath.Join(dir, f.Name())
		ext := path.Ext(dstFile)
		if ext == ".src" {
			dstFile = dstFile[0:len(dstFile)-len(ext)] + ".go"
			a.Src = append(a.Src, dstFile)
		} else if ext == ".from" {
			resultFile := dstFile[0:len(dstFile)-len(ext)] + ".dst"
			dstFile = dstFile[0:len(dstFile)-len(ext)] + ".go"
			a.From = append(a.From, dstFile)
			a.mapFrom2Dest[dstFile] = resultFile
		}
		if f.Name() == "out.txt" {
			a.output = dstFile
		}
		err := copyFileContents(filepath.Join(testDir, f.Name()), dstFile)
		if err != nil {
			t.Error(err)
		}
	}
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
