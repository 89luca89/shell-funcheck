// Package types contains all the custom types we will use in this program.
package types

import (
	"fmt"

	"mvdan.cc/sh/v3/syntax"
)

// CLIOptions stores options we're passing to cobra's CLI args.
type CLIOptions struct {
	// up options
	Files           []string
	ExcludeComments bool
	WError          bool
}

// ParsedTokens will keep the results of all the parsing of the
// function, and will be then used to apply our rules.
type ParsedTokens struct {
	Arguments   map[string]*syntax.Assign
	Assignments map[string]*syntax.Assign
	Comments    []string
	Iters       map[string]*syntax.WordIter
	Variables   []*syntax.Lit
	Function    *syntax.FuncDecl
}

const (
	Err  = 0
	Warn = 1
	Note = 2
)

var levels = map[int]string{
	0: "error",
	1: "warn",
	2: "note",
}

// ReportedVar is our message to be displayed.
//
//	Subject is what is erroring
//	Function is the context function where Subject is
//	Reason is a brief description of the problem
//	Level is the severity from levels.
type ReportedVar struct {
	Subject  *syntax.Lit
	Function *syntax.FuncDecl
	Reason   string
	Level    int
}

// Print will pretty print our outputs in a Vim compatible way.
func (v *ReportedVar) Print(filepath string) {
	if v.Subject != nil {
		fmt.Printf(
			"%s:%s: %s: %s - %s: %s\n",
			filepath,
			v.Subject.ValuePos,
			levels[v.Level],
			v.Function.Name.Value,
			v.Subject.Value,
			v.Reason)

		return
	}

	fmt.Printf(
		"%s:%s: %s: %s\n",
		filepath,
		"0:0",
		levels[v.Level],
		v.Reason,
	)
}
