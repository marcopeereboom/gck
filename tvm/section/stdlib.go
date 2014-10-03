package section

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-xdr/xdr2"
)

// XXX Note that this code is basically a copy/paste from variable.
// XXX This needs to be refactored

type OsInOut struct {
	Name string
	Type reflect.Type
}

type OsCall struct {
	Id        uint64
	Name      string
	Variables []OsInOut
	Results   []OsInOut
}

func (o *OsCall) String() string {
	var vars, results string

	// skip types for now
	for _, v := range o.Variables {
		vars += v.Name + ","
	}
	for _, v := range o.Results {
		vars += v.Name + ","
	}

	return fmt.Sprintf("%v(%v)(%v)", o.Name, vars, results)
}

// Os is an xdr representation
type Os struct {
	Id    uint64      // symbol id (indexes symbol table)
	Name  string      // name of symbol
	Type  uint64      // value type
	Value string      // string representation of value, exported for xdr
	value interface{} // actual value
}

func NewOs(id uint64, name string, value interface{}) (*Os, error) {
	v := Os{
		Id:   id,
		Name: name,
	}
	switch av := value.(type) {
	case OsCall:
		v.Type = SymLabelId
		v.value = av
		v.Value = av.String()
	default:
		return nil, fmt.Errorf("unsuported type %T", value)
	}

	return &v, nil
}

func (v *Os) GetActualValue() interface{} {
	return v.value
}

func encodeOsElement(v *Os) ([]byte, error) {
	vv := v
	if vv.Id < SymReserved {
		return nil, fmt.Errorf("invalid symbol id %x", vv.Id)
	}

	switch val := v.value.(type) {
	case OsCall:
		vv.Type = SymLabelId
		vv.Value = val.String()
	default:
		return nil, fmt.Errorf("unsupported os type %T", val)
	}

	var w bytes.Buffer
	_, err := xdr.Marshal(&w, vv)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func encodeOsElements(v []*Os) ([]byte, error) {
	var b bytes.Buffer

	for _, v := range v {
		ve, err := encodeOsElement(v)
		if err != nil {
			return nil, err
		}
		b.Write(ve)
	}
	return b.Bytes(), nil
}

func decodeOsElement(b []byte, consumed *int) (*Os, error) {
	v := Os{}
	n, err := xdr.Unmarshal(bytes.NewReader(b), &v)
	if err != nil {
		return nil, err
	}

	switch v.Type {
	case SymLabelId:
		if strings.HasPrefix(v.Value, v.Name) {
			// XXX good enough for now
			// XXX enforce args/results count + type
		}

	default:
		return nil, fmt.Errorf("unsupported type")
	}

	if consumed != nil {
		*consumed = n
	}

	return &v, nil
}

func decodeOsElements(b []byte) ([]*Os, error) {
	var (
		n, at int
		vars  []*Os
	)

	for todo := len(b); todo > 0; {
		v, err := decodeOsElement(b[at:], &n)
		if err != nil {
			return nil, err
		}
		vars = append(vars, v)
		at += n
		todo -= n
	}

	return vars, nil
}

func NewOsSection(args []*Os) (*Section, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("empty os section not allowed")
	}

	seen := make(map[uint64]bool)
	// make sure these are valid
	for _, v := range args {
		_, err := encodeOsElement(v)
		if err != nil {
			return nil, err
		}
		_, found := seen[v.Id]
		if found {
			return nil, fmt.Errorf("duplicate os id %x", v.Id)
		}
		seen[v.Id] = true
	}

	s := Section{
		Version: Version,
		Name:    Sections[OsId],
		Id:      OsId,
		Read:    true,
		Write:   false,
		Execute: true,
		Payload: args,
	}
	return &s, nil
}
