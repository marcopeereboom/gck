// backend generalizes architectures.
// All target architectures must conform to this interface.
// Using EmitCode should generate a binary that can be run on a target machine.
package backend

import (
	"fmt"

	"github.com/marcopeereboom/gck/backend/arch"
	"github.com/marcopeereboom/gck/backend/tvm"
)

const (
	TVM = "toyvm"
)

// New instantiates a Backend.
// A supported target architecture name must be provided.
// If the provided architecture name is not supported the function will throw
// an error.
func New(name string) (arch.Backend, error) {
	switch name {
	case TVM:
		return tvm.New()
	}
	return nil, fmt.Errorf("unsuported architecture: %v", name)
}
