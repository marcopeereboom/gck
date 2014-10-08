package section

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/davecgh/go-xdr/xdr2"
)

// Variable is an xdr representation
type Variable struct {
	Id    uint64      // symbol id (indexes symbol table)
	Name  string      // name of symbol
	Type  uint64      // value type
	Value string      // string representation of value, exported for xdr
	value interface{} // actual value
}

func NewVariable(id uint64, name string, value interface{}) (*Variable, error) {

	v := Variable{
		Id:   id,
		Name: name,
	}
	switch av := value.(type) {
	case *big.Rat:
		v.Type = SymNumId
		v.value = new(big.Rat).Set(av)
		v.Value = v.value.(*big.Rat).String()
	default:
		return nil, fmt.Errorf("unsuported type %T", value)
	}

	return &v, nil
}

func (v *Variable) GetActualValue() interface{} {
	return v.value
}

func encodeVariableElement(v *Variable) ([]byte, error) {
	vv := v
	if vv.Id < SymReserved {
		return nil, fmt.Errorf("invalid symbol id %x", vv.Id)
	}

	switch val := v.value.(type) {
	case *big.Rat:
		vv.Type = SymNumId
		vv.Value = val.String()
	default:
		return nil, fmt.Errorf("unsupported variable type %T", val)
	}

	var w bytes.Buffer
	_, err := xdr.Marshal(&w, vv)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func encodeVariableElements(v []*Variable) ([]byte, error) {
	var b bytes.Buffer

	for _, v := range v {
		ve, err := encodeVariableElement(v)
		if err != nil {
			return nil, err
		}
		b.Write(ve)
	}
	return b.Bytes(), nil
}

func decodeVariableElement(b []byte, consumed *int) (*Variable, error) {
	v := Variable{}
	n, err := xdr.Unmarshal(bytes.NewReader(b), &v)
	if err != nil {
		return nil, err
	}

	switch v.Type {
	case SymNumId:
		var ok bool
		v.value, ok = new(big.Rat).SetString(v.Value)
		if !ok {
			return nil, fmt.Errorf("Value not a big.Rat")
		}
	default:
		return nil, fmt.Errorf("unsupported type")
	}

	if consumed != nil {
		*consumed = n
	}

	return &v, nil
}

func decodeVariableElements(b []byte) ([]*Variable, error) {
	var (
		n, at int
		vars  []*Variable
	)

	for todo := len(b); todo > 0; {
		v, err := decodeVariableElement(b[at:], &n)
		if err != nil {
			return nil, err
		}
		vars = append(vars, v)
		at += n
		todo -= n
	}

	return vars, nil
}

func NewVariableSection(args []*Variable) (*Section, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("empty variable section not allowed")
	}

	seen := make(map[uint64]bool)
	// make sure these are valid
	for _, v := range args {
		_, err := encodeVariableElement(v)
		if err != nil {
			return nil, err
		}
		_, found := seen[v.Id]
		if found {
			return nil, fmt.Errorf("duplicate variable id %x",
				v.Id)
		}
		seen[v.Id] = true
	}

	vs := Section{
		Version: Version,
		Name:    Sections[VariableId],
		Id:      VariableId,
		Read:    true,
		Write:   true,
		Execute: false,
		Payload: args,
	}
	return &vs, nil
}
