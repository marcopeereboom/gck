// frontend generalizes languages.
package frontend

import (
	"fmt"

	"github.com/marcopeereboom/gck/frontend/driver"
	"github.com/marcopeereboom/gck/frontend/myrmidon"
	"github.com/marcopeereboom/gck/frontend/sml"
)

const (
	SML      = "sml"
	MYRMIDON = "myrmidon"
)

// New instantiates a Frontend.
// A supported language name must be provided.
// If the provided language name is not supported the function will throw
// an error.
func New(name string) (driver.Frontend, error) {
	switch name {
	case SML:
		return sml.New()
	case MYRMIDON:
		return myrmidon.New()
	}
	return nil, fmt.Errorf("unsuported language: %v", name)
}
