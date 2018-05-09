package diff

import (
	"go/token"
	"go/parser"
	"golang.org/x/tools/go/ast/astutil"
	"go/ast"
	"fmt"
	"go/format"
	"bytes"
	"io/ioutil"
)

func (a *Arguments) removeSymbols() (int, error) {
	var err error
	var totalDuplicateSymbols int
	for _, from := range a.From {
		dup, e := a.removeSymbolsFromFile(from)
		totalDuplicateSymbols += dup
		if e != nil {
			err = e
			if a.Verbose {
				fmt.Fprintf(a.Stdout, "Error removing symobls from file \"%s\": %v\n", from, e)
			}
		} else if a.Verbose {
			fmt.Fprintf(a.Stdout, "Removed %v duplicate symbols from %s\n", dup, from)
		}
	}
	return totalDuplicateSymbols, err
}

func (a *Arguments) removeSymbolsFromFile(fileName string) (int, error) {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return 0, err
	}
	var duplicateSymbols int
	hasSymbol := func(symbol string) bool {
		_, ok := a.symbols[symbol]
		return ok
	}
	hasIdent := func(ident *ast.Ident) bool {
		if ident == nil {
			return false
		}
		return hasSymbol(ident.Name)
	}
	removeIfSymbolExists := func(cursor *astutil.Cursor, ident *ast.Ident) {
		if hasIdent(ident) {
			duplicateSymbols++
			cursor.Delete()
		}
	}
	removeIfEmptyNames := func(cursor *astutil.Cursor, vs *ast.ValueSpec) {
		if len(vs.Names) == 0 {
			cursor.Delete()
		}
	}
	removeSymbols := func(cursor *astutil.Cursor) bool {
		if cursor == nil {
			return true
		}
		node := cursor.Node()
		switch n := node.(type) {
		case *ast.FuncDecl:
			removeIfSymbolExists(cursor, n.Name)
			return false
		case *ast.GenDecl:
			if n.Tok == token.IMPORT {
				return false
			}
		case *ast.ValueSpec:
			var newNames []*ast.Ident
			for _, n := range n.Names {
				if hasIdent(n) {
					duplicateSymbols++
				} else {
					newNames = append(newNames, n)
				}
			}
			n.Names = newNames
			removeIfEmptyNames(cursor, n)
			return false
		case *ast.TypeSpec:
			removeIfSymbolExists(cursor, n.Name)
			return false
		}
		return true
	}
	astutil.Apply(f, removeSymbols, nil)
	removeEmptyGenDecls := func(cursor *astutil.Cursor) bool {
		if cursor == nil {
			return true
		}
		node := cursor.Node()
		switch n := node.(type) {
		case *ast.GenDecl:
			if n.Tok == token.IMPORT {
				return false
			}
			if len(n.Specs) == 0 {
				cursor.Delete()
			}
		}
		return true
	}
	astutil.Apply(f, removeEmptyGenDecls, nil)

	// write changes to file
	if duplicateSymbols > 0 {
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, f); err != nil {
			handleAstError(fset, f, err)
		}
		err = ioutil.WriteFile(fileName, []byte(buf.String()), 0644)
		if err != nil {
			return 0, fmt.Errorf("writing changed file \"%s\" failed: %v", fileName, err)
		}
	}
	return duplicateSymbols, nil
}

// copied from c2go
func handleAstError(fset *token.FileSet, f *ast.File, err error) {
	// Printing the entire AST will generate a lot of output. However, it is
	// the only way to debug this type of error. Hopefully the error
	// (printed immediately afterwards) will give a clue.
	//
	// You may see an error like:
	//
	//     panic: format.Node internal error (692:23: expected selector or
	//     type assertion, found '[')
	//
	// This means that when Go was trying to convert the Go AST to source
	// code it has come across a value or attribute that is illegal.
	//
	// The line number it is referring to (in this case, 692) is not helpful
	// as it references the internal line number of the Go code which you
	// will never see.
	//
	// The "[" means that there is a bracket in the wrong place. Almost
	// certainly in an identifer, like:
	//
	//     noarch.IntTo[]byte("foo")
	//
	// The "[]" which is obviously not supposed to be in the function name
	// is causing the syntax error. However, finding the original code that
	// produced this can be tricky.
	//
	// The first step is to filter down the AST output to probably lines.
	// In the error message it said that there was a misplaced "[" so that's
	// what we will search for. Using the original command (that generated
	// thousands of lines) we will add two grep filters:
	//
	//     go test ... | grep "\[" | grep -v '{$'
	//     #                   |     |
	//     #                   |     ^ This excludes lines that end with "{"
	//     #                   |       which almost certainly won't be what
	//     #                   |       we are looking for.
	//     #                   |
	//     #                   ^ This is the character we are looking for.
	//
	// Hopefully in the output you should see some lines, like (some lines
	// removed for brevity):
	//
	//     9083  .  .  .  .  .  .  .  .  .  .  Name: "noarch.[]byteTo[]int"
	//     9190  .  .  .  .  .  .  .  .  .  Name: "noarch.[]intTo[]byte"
	//
	// These two lines are clearly the error because a name should not look
	// like this.
	//
	// Looking at the full output of the AST (thousands of lines) and
	// looking at those line numbers should give you a good idea where the
	// error is coming from; by looking at the parents of the bad lines.
	_ = ast.Print(fset, f)

	panic(err)
}
