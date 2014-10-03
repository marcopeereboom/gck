package section

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	"github.com/davecgh/go-xdr/xdr2"
)

// XXX Note that this code is basically a copy/paste from variable.
// XXX This needs to be refactored

// Const is an xdr representation
type Const struct {
	Id    uint64      // symbol id (indexes symbol table)
	Name  string      // name of symbol
	Type  uint64      // value type
	Value string      // string representation of value, exported for xdr
	value interface{} // actual value
}

func NewConst(id uint64, name string, value interface{}) (*Const, error) {
	v := Const{
		Id:   id,
		Name: name,
	}
	switch av := value.(type) {
	case *big.Rat:
		v.Type = SymNumId
		v.value = new(big.Rat).Set(av)
		v.Value = v.value.(*big.Rat).String()
	case uint64:
		v.Type = SymLabelId
		v.value = av
		v.Value = fmt.Sprintf("%v", av)
	default:
		return nil, fmt.Errorf("unsuported type %T", value)
	}

	return &v, nil
}

func (v *Const) GetActualValue() interface{} {
	return v.value
}

func encodeConstElement(v *Const) ([]byte, error) {
	vv := v
	if vv.Id < SymReserved {
		return nil, fmt.Errorf("invalid symbol id %x", vv.Id)
	}

	switch val := v.value.(type) {
	case *big.Rat:
		vv.Type = SymNumId
		vv.Value = val.String()
	case uint64:
		vv.Type = SymLabelId
		v.Value = fmt.Sprintf("%v", val)
	default:
		return nil, fmt.Errorf("unsupported const type %T", val)
	}

	var w bytes.Buffer
	_, err := xdr.Marshal(&w, vv)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func encodeConstElements(v []*Const) ([]byte, error) {
	var b bytes.Buffer

	for _, v := range v {
		ve, err := encodeConstElement(v)
		if err != nil {
			return nil, err
		}
		b.Write(ve)
	}
	return b.Bytes(), nil
}

func decodeConstElement(b []byte, consumed *int) (*Const, error) {
	v := Const{}
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
	case SymLabelId:
		newConst, err := strconv.Atoi(v.Value)
		if err != nil {
			return nil, err
		}
		v.value = uint64(newConst)

	default:
		return nil, fmt.Errorf("unsupported type")
	}

	if consumed != nil {
		*consumed = n
	}

	return &v, nil
}

func decodeConstElements(b []byte) ([]*Const, error) {
	var (
		n, at int
		vars  []*Const
	)

	for todo := len(b); todo > 0; {
		v, err := decodeConstElement(b[at:], &n)
		if err != nil {
			return nil, err
		}
		vars = append(vars, v)
		at += n
		todo -= n
	}

	return vars, nil
}

func NewConstSection(args []*Const) (*Section, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("empty const section not allowed")
	}

	seen := make(map[uint64]bool)
	// make sure these are valid
	for _, v := range args {
		_, err := encodeConstElement(v)
		if err != nil {
			return nil, err
		}
		_, found := seen[v.Id]
		if found {
			return nil, fmt.Errorf("duplicate const id %x", v.Id)
		}
		seen[v.Id] = true
	}

	cs := Section{
		Version: Version,
		Name:    Sections[ConstId],
		Id:      ConstId,
		Read:    true,
		Write:   false,
		Execute: false,
		Payload: args,
	}
	return &cs, nil
}
