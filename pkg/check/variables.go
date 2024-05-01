package check

import (
	"github.com/89luca89/shell-funcheck/pkg/types"
	"github.com/89luca89/shell-funcheck/pkg/util"
)

// listUndeclaredVariables will parse the comments and find the "Expected env variables" and
// the "Expected global variables:" section.
// Within those section, it will search for the variables used in the function, for example
// the following script:
//
//		#!/bin/sh
//	 input=$1
//
//		foo() {
//		  echo $input
//		  echo $USER
//		}
//
// will need to declare the argument "input" in the "Expected global variables:" section of the comment
// and "USER" in the "Expected env variables:" section.
func listUndeclaredVariables(parsedTokens types.ParsedTokens) []types.ReportedVar {
	result := []types.ReportedVar{}

	documentedVars := listDocumentedVariables(parsedTokens.Comments)

	for _, variable := range parsedTokens.Variables {
		if util.IsNumber(variable.Value) {
			continue
		}

		if parsedTokens.Assignments[variable.Value] != nil ||
			parsedTokens.Arguments[variable.Value] != nil ||
			parsedTokens.Iters[variable.Value] != nil ||
			documentedVars[variable.Value] != "" {
			// Variable is documented, skip
			continue
		}

		// We also want to encourage to explain what variables do in the comment
		// section, so if no "foo = explanation" comment is found, we do a warning.
		if value, ok := documentedVars[variable.Value]; ok {
			// Variable is documented, but not explained
			if value == "" {
				// An explanation is encouraged only for regular variables,
				// env ones are usually known.
				if util.StringIsAllUpper(variable.Value) {
					continue
				}

				result = append(result, types.ReportedVar{
					Subject:  variable,
					Reason:   "variable is documented, but not explained",
					Function: parsedTokens.Function,
					Level:    types.Warn,
				})

				continue
			}
		}

		result = append(result, types.ReportedVar{
			Subject:  variable,
			Reason:   "variable is not documented",
			Level:    types.Err,
			Function: parsedTokens.Function,
		})
	}

	return result
}
