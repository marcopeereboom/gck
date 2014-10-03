// optimizer walks an AST tree and tries to do optimizations on it.
package optimizer

import "github.com/marcopeereboom/gck/ast"

// Optimize transforms n AST and returns an optimized version of it.
func Optimize(n ast.Node) (ast.Node, error) {
	return ast.Clone(n), nil
}
