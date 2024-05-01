package check

import (
	"slices"
	"strings"

	"github.com/89luca89/shell-funcheck/pkg/types"
	"mvdan.cc/sh/v3/syntax"
)

// those are the markers for our comment explanation, usually we expect something like:
//
//	# foo_nction does stuff (explanation of the functionality)
//	# Arguments:
//	#	arg1 = string/bool/int and what it does
//	#	arg2 = string/bool/int and what it does
//	# Expected env variables:
//	#	ENV1
//	#	ENV2
//	# Expected global variables:
//	#	globalVar1 = string/bool/int and what it does
//	#	globalVar2 = string/bool/int and what it does
//	# Outputs:
//	#	needed icons in /run/host/$host_home/.local/share/icons
//	#	needed desktop files in /run/host/$host_home/.local/share/applications
//	#	or error code.
//	 foo_nction() {
//	 ...
//	 }
const (
	argumentMarker    = " Arguments:"
	expectedVarMarker = " Expected global variables:"
	expectedEnvMarker = " Expected env variables:"
	outputsMarker     = " Outputs:"
)

// validateFunctionComments will ensure that the comments of the input function
// is conformant to our expected input.
func validateFunctionComments(fn *syntax.Stmt) []types.ReportedVar {
	function, ok := fn.Cmd.(*syntax.FuncDecl)
	if !ok {
		return nil
	}

	// missing at all! this should be an error, at leas an explanation should
	// be expected
	if len(fn.Comments) == 0 {
		return []types.ReportedVar{
			{
				Subject: function.Name,
				Reason: "function " +
					function.Name.Value +
					" should be documented",
				Level:    types.Err,
				Function: function,
			},
		}
	}

	comments := []string{}
	result := []types.ReportedVar{}

	// Explanation should be in the form
	//   foo_nction does the following thing...
	//
	// This idea is transplanted from https://go.dev/wiki/CodeReviewComments
	if !strings.HasPrefix(strings.TrimSpace(fn.Comments[0].Text), function.Name.Value) {
		result = append(result, types.ReportedVar{
			Subject: function.Name,
			Reason: "comments on function " +
				function.Name.Value +
				" should be of the form \"" + function.Name.Value + " ...\"",
			Level:    types.Warn,
			Function: function,
		})
	}

	for _, comment := range fn.Comments {
		comments = append(comments, comment.Text)
	}

	//
	// Check for our markers!
	//
	if !slices.Contains(comments, argumentMarker) {
		result = append(result, types.ReportedVar{
			Subject:  function.Name,
			Reason:   "missing '" + argumentMarker + "' section",
			Level:    types.Err,
			Function: function,
		})
	}

	if !slices.Contains(comments, expectedVarMarker) {
		result = append(result, types.ReportedVar{
			Subject:  function.Name,
			Reason:   "missing '" + expectedVarMarker + "' section",
			Level:    types.Err,
			Function: function,
		})
	}

	if !slices.Contains(comments, expectedEnvMarker) {
		result = append(result, types.ReportedVar{
			Subject:  function.Name,
			Reason:   "missing '" + expectedEnvMarker + "' section",
			Level:    types.Err,
			Function: function,
		})
	}

	if !slices.Contains(comments, outputsMarker) {
		result = append(result, types.ReportedVar{
			Subject:  function.Name,
			Reason:   "missing '" + outputsMarker + "' section",
			Level:    types.Err,
			Function: function,
		})
	}

	return result
}

// listDocumentedVariables will report the variables declared in our comments.
// values from the "Expected env variables:" and "Expected global variables:"
// sections are reported.
func listDocumentedVariables(comments []string) map[string]string {
	documentedVars := map[string]string{}

	skip := false

	for _, comment := range comments {
		if comment == argumentMarker {
			skip = true
		}

		if comment == expectedVarMarker {
			skip = false
		}

		if comment == expectedEnvMarker {
			skip = false
		}

		if comment == outputsMarker {
			skip = true
		}

		if skip {
			continue
		}

		variable := strings.TrimSpace(strings.Split(comment, ":")[0])
		documentedVars[variable] = strings.TrimSpace(
			strings.Join(
				strings.Split(comment, ":")[1:], " "),
		)
	}

	return documentedVars
}

// listDocumentedArguments will report the function arguments declared in our comments.
// values from the "Arguments: section are reported.
func listDocumentedArguments(comments []string) map[string]string {
	documentedArgs := map[string]string{}

	skip := false

	for _, comment := range comments {
		if comment == argumentMarker {
			skip = false
		}

		if comment == expectedVarMarker {
			skip = true
		}

		if comment == outputsMarker {
			skip = true
		}

		if skip {
			continue
		}

		variable := strings.TrimSpace(strings.Split(comment, ":")[0])
		documentedArgs[variable] = strings.TrimSpace(
			strings.Join(
				strings.Split(comment, ":")[1:], " "),
		)
	}

	return documentedArgs
}
