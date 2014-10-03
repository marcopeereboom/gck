// driver provides the interface for languages
package driver

import (
	"bufio"
	"io"
	"strings"

	"github.com/marcopeereboom/gck/ast"
)

// Frontend is the interface that all languages must adhere to.
type Frontend interface {
	Compile(string) error     // compile code
	AST() (ast.Node, error)   // return AST of compiled code
	Lines() ([]string, error) // return the script in array of strings
	Line(int) (string, error) // return an individual line of script
}

// LineGenerator slices the source file up in individual lines.
func LineGenerator(src string) ([]string, error) {
	lines := make([]string,
		1,    // skip line 0
		1024, // guess something
	)

	r := bufio.NewReader(strings.NewReader(src))
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		lines = append(lines, string(line))
	}

	return lines, nil
}
