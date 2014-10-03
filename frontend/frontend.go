package frontend

import (
	"fmt"

	"github.com/marcopeereboom/gck/frontend/driver"
	"github.com/marcopeereboom/gck/frontend/ml"
)

const (
	SML = "sml"
)

// New instantiates a Frontend.
// A supported language name must be provided.
// If the provided language name is not supported the function will throw
// an error.
func New(name string) (driver.Frontend, error) {
	switch name {
	case SML:
		return sml.New()
	}
	return nil, fmt.Errorf("unsuported language: %v", name)
}
