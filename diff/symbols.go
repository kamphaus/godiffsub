package diff

import (
	"go/token"
	"go/parser"
	"go/ast"
	"sort"
	"fmt"
)

func (a *Arguments) readSymbols() {
	a.symbols = make(map[string]struct{})
	for _, src := range a.Src {
		a.readSymbolsForFile(src)
	}
}

func (a *Arguments) readSymbolsForFile(fileName string) error {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}
	ast.Walk(&visitor{a}, f)
	return nil
}

type visitor struct {
	args *Arguments
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return v
	}
	switch n := node.(type) {
	case *ast.FuncDecl:
		v.visitName(n.Name)
		return nil
	case *ast.GenDecl:
		switch n.Tok {
		case token.CONST, token.TYPE, token.VAR:
			for _, spec := range n.Specs {
				v.visitSpecs(spec)
			}
		}
		return nil
	}
	return v
}

func (v *visitor) visitName(name *ast.Ident) {
	v.args.symbols[name.Name] = struct{}{}
}

func (v *visitor) visitSpecs(spec ast.Spec) {
	switch s := spec.(type) {
	case *ast.ValueSpec:
		for _, n := range s.Names {
			v.visitName(n)
		}
	case *ast.TypeSpec:
		v.visitName(s.Name)
	}
}

func (a Arguments) printSymbols() {
	var symbols = make([]string, 0, len(a.symbols))
	for s := range a.symbols {
		symbols = append(symbols, s)
	}
	sort.Strings(symbols)
	for _, s := range symbols {
		fmt.Println(s)
	}
}
