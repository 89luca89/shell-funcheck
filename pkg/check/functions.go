package check

import (
	"bufio"
	"os"

	"github.com/89luca89/shell-funcheck/pkg/types"
	"mvdan.cc/sh/v3/syntax"
)

// listFileFunctions will return the list of functions in the input file.
func listFileFunctions(filepath string, options *types.CLIOptions) ([]*syntax.Stmt, error) {
	var result []*syntax.Stmt

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileReader := bufio.NewReader(file)

	parsedFile, err := syntax.NewParser(syntax.KeepComments(!options.ExcludeComments)).Parse(fileReader, "")
	if err != nil {
		return result, err
	}

	for _, stmt := range parsedFile.Stmts {
		if _, ok := stmt.Cmd.(*syntax.FuncDecl); ok {
			result = append(result, stmt)
		}
	}

	return result, nil
}
