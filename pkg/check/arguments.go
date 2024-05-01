package check

import (
	"github.com/89luca89/shell-funcheck/pkg/types"
)

// listUndeclaredArguments will parse the comments and find the "Arguments:" section
// within that section, it will search for the arguments of the function, for example
// the following script:
//
//	#!/bin/sh
//	foo() {
//	  input=$1
//	  echo $input
//	}
//
// will need to declare the argument "input" in the "Arguments:" section of the comment.
func listUndeclaredArguments(parsedTokens types.ParsedTokens) []types.ReportedVar {
	result := []types.ReportedVar{}

	documentedVars := listDocumentedArguments(parsedTokens.Comments)

	for _, argument := range parsedTokens.Arguments {
		if documentedVars[argument.Name.Value] != "" {
			// Argument is documented, skip
			continue
		}

		// We also want to encourage to explain what arguments do in the comment
		// section, so if no "foo = explanation" comment is found, we do a warning.
		if value, ok := documentedVars[argument.Name.Value]; ok {
			// Argument is documented, but not explained
			if value == "" {
				result = append(result, types.ReportedVar{
					Subject:  argument.Name,
					Reason:   "argument is documented, but not explained",
					Level:    types.Warn,
					Function: parsedTokens.Function,
				})

				continue
			}
		}

		result = append(result, types.ReportedVar{
			Subject:  argument.Name,
			Reason:   "argument is not documented",
			Level:    types.Err,
			Function: parsedTokens.Function,
		})
	}

	return result
}
