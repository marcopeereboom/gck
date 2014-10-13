// arch provides the interface for architectures
package arch

import "github.com/marcopeereboom/gck/ast"

// Backend is the interface that all architectures must adhere to.
type Backend interface {
	EmitCode(ast.Node) ([]byte, error) // return target architecture binary
	Error() error                      // returns errors that are not detected by EmitCode
}
