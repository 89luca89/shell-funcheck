// Package check package contains the parsing and ruling utilities to report problems.
package check

import (
	"github.com/89luca89/shell-funcheck/pkg/types"
	"github.com/89luca89/shell-funcheck/pkg/util"
	"mvdan.cc/sh/v3/syntax"
)

// listParsedTokens will walk trough the input function and store
//
//	comments
//	arguments
//	assignments
//	iter variables (for i in ..., for x in ...)
//	variables
//
// these will then be stored in a ParsedTokens struct that will later
// be used to apply the rules.
func listParsedTokens(funct *syntax.Stmt) types.ParsedTokens {
	comments := []string{}
	arguments := map[string]*syntax.Assign{}
	assignments := map[string]*syntax.Assign{}
	iters := map[string]*syntax.WordIter{}
	variables := []*syntax.Lit{}

	for _, comment := range funct.Comments {
		comments = append(comments, comment.Text)
	}

	fun, ok := funct.Cmd.(*syntax.FuncDecl)
	if !ok {
		return types.ParsedTokens{}
	}

	syntax.Walk(fun, func(node syntax.Node) bool {
		switch token := node.(type) {
		case *syntax.ParamExp:
			if util.StringIsArgument(token.Param.Value) {
				variables = append(variables, token.Param)
			}
		case *syntax.WordIter:
			if util.StringIsArgument(token.Name.Value) {
				iters[token.Name.Value] = token
			}
		case *syntax.CallExpr:
			for _, assign := range token.Assigns {
				assignments[assign.Name.Value] = assign

				if arg, ok := assign.Value.Parts[0].(*syntax.ParamExp); ok {
					if util.StringIsArgument(arg.Param.Value) {
						arguments[assign.Name.Value] = assign
						delete(assignments, assign.Name.Value)
					}
				} else if arg, ok := assign.Value.Parts[0].(*syntax.DblQuoted); ok {
					if len(arg.Parts) > 0 {
						if nestArg, ok := arg.Parts[0].(*syntax.ParamExp); ok {
							if util.StringIsArgument(nestArg.Param.Value) {
								arguments[assign.Name.Value] = assign
								delete(assignments, assign.Name.Value)
							}
						}
					}
				}
			}
		}

		return true
	})

	return types.ParsedTokens{
		Arguments:   arguments,
		Assignments: assignments,
		Comments:    comments,
		Iters:       iters,
		Variables:   variables,
		Function:    fun,
	}
}

// RunCheck will run checks on an input file.
// this will report:
//
//	uncommented and undocumented functions
//	missing document sections in the comment
//	missing arguments declaration in the comment
//	missing global variables declaration in the comment
//	missing env variables declaration in the comment
//	missing output explanation in the comment
//	warn about incomplete documentation of the previous sections.
func RunCheck(filepath string, options *types.CLIOptions) ([]types.ReportedVar, error) {
	var result []types.ReportedVar

	functions, err := listFileFunctions(filepath, options)
	if err != nil {
		return nil, err
	}

	if len(functions) == 0 {
		return append(result, types.ReportedVar{
			Reason: "no functions found in file",
			Level:  types.Note,
		}), nil
	}

	parsedTokens := []types.ParsedTokens{}

	for _, function := range functions {
		result = append(result, validateFunctionComments(function)...)

		parsedTokens = append(parsedTokens, listParsedTokens(function))
	}

	for _, parsedTokens := range parsedTokens {
		result = append(result, listUndeclaredArguments(parsedTokens)...)
		result = append(result, listUndeclaredVariables(parsedTokens)...)
	}

	return result, nil
}
