package section

import (
	"fmt"
	"math/big"
)

const (
	SymInvalid  = 0
	SymNumId    = 1
	SymLabelId  = 2
	SymReserved = 256 // minimum symbol id
)

var (
	Symbols = map[uint64]string{
		SymInvalid: "INVALID",
		SymNumId:   "NUMBER",
		SymLabelId: "LABEL",
	}
)

type Symbol struct {
	Id        uint64
	Name      string
	SectionId uint64
	TypeId    uint64
	Value     interface{}
}

func New(id, sectionId uint64, name string, val interface{}) (*Symbol,
	error) {

	if id < SymReserved {
		return nil, fmt.Errorf("invalid symbol id %x", id)
	}

	s := Symbol{
		Id:        id,
		Name:      name,
		SectionId: sectionId,
	}
	switch sectionId {
	case VariableId:
		err := variable(&s, val)
		if err != nil {
			return nil, fmt.Errorf("%v for variable symbol %v",
				err, name)
		}

	case ConstId:
		err := constant(&s, val)
		if err != nil {
			return nil, fmt.Errorf("%v for const symbol %v",
				err, name)
		}

	case OsId:
		err := stdlib(&s, val)
		if err != nil {
			return nil, fmt.Errorf("%v for os symbol %v",
				err, name)
		}

	default:
		return nil, fmt.Errorf("invalid symbol section id 0x%0x",
			sectionId)
	}

	if name == "" {
		s.Name = fmt.Sprintf("%016x", id)
	}

	return &s, nil
}

func constant(s *Symbol, val interface{}) error {
	// allowed types in constant section
	switch v := val.(type) {
	case *big.Rat:
		s.TypeId = SymNumId
		s.Value = new(big.Rat).Set(v)
		return nil

	case uint64:
		// this may not be enough of a discriminator
		s.TypeId = SymLabelId
		s.Value = v
		return nil

	default:
		return fmt.Errorf("invalid type %T", val)
	}

	return fmt.Errorf("const, impossible condition")
}

func variable(s *Symbol, val interface{}) error {
	// allowed types in variable section
	switch v := val.(type) {
	case *big.Rat:
		s.TypeId = SymNumId
		s.Value = new(big.Rat).Set(v)
		return nil

	default:
		return fmt.Errorf("invalid type %T", val)
	}

	return fmt.Errorf("variable, impossible condition")
}

func stdlib(s *Symbol, val interface{}) error {
	// allowed types in variable section
	s.TypeId = SymLabelId
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case OsCall:
		s.Value = v
		return nil
	default:
		return fmt.Errorf("invalid type %T", val)
	}

	return fmt.Errorf("os, impossible condition")
}
