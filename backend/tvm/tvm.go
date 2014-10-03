// tvm implements the backend interface to target ToyVM architecture.
package tvm

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	"github.com/marcopeereboom/gck/ast"
	"github.com/marcopeereboom/gck/backend/arch"
	"github.com/marcopeereboom/gck/tvm/section"
	"github.com/marcopeereboom/gck/tvm/vm"
)

type ToyVirtualMachine struct {
	id     uint64
	consts []*section.Const
	varsA  []*section.Variable
	vars   map[string]*section.Variable
	code   []uint64
}

var _ arch.Backend = &ToyVirtualMachine{}

// utility
func New() (*ToyVirtualMachine, error) {
	vm := ToyVirtualMachine{
		vars: make(map[string]*section.Variable),
		id:   1000,
		code: make([]uint64, 0, 1000),
	}

	return &vm, nil
}

func (t *ToyVirtualMachine) NewId() uint64 {
	defer func() { t.id++ }()
	return t.id
}

// interface
func (t *ToyVirtualMachine) EmitCode(n ast.Node) ([]byte, error) {
	var image []byte

	w := bytes.NewBuffer(image)
	err := ast.EmitCode(n, w, t.emitCode)
	if err != nil {
		return nil, err
	}

	// generate sections
	cs := section.NewCodeSection(t.code)

	cos, err := section.NewConstSection(t.consts)
	if err != nil {
		return nil, err
	}

	vs, err := section.NewVariableSection(t.varsA)
	if err != nil {
		return nil, err
	}

	// generate image
	i := section.NewImage()
	err = i.AddSection(cs, true)
	if err != nil {
		return nil, err
	}
	err = i.AddSection(vs, true)
	if err != nil {
		return nil, err
	}
	err = i.AddSection(cos, true)
	if err != nil {
		return nil, err
	}

	return i.GetImage(), nil
}

func (t *ToyVirtualMachine) addCode(code []uint64) {
	// blow me slow
	for _, v := range code {
		t.code = append(t.code, v)
	}
}

func (t *ToyVirtualMachine) emitCode(ty int, args ...interface{}) {
	switch ty {
	case ast.IDENTIFIER:
		fmt.Printf("push id\n")
		//s.addCode("\tpush\t%v\n", args[0].(Node).Value)
	case ast.NUMBER:
		//s.addCode("\tpush\t%v\n", args[0].(*big.Rat))
		id := t.NewId()
		c, err := section.NewConst(id, "c"+strconv.Itoa(int(id)),
			args[0].(*big.Rat))
		if err != nil {
			panic(err)
		}
		t.consts = append(t.consts, c)
		t.addCode([]uint64{vm.OP_PUSH, id})
		//fmt.Printf("push const %v\n", code)
	case '=':
		//fmt.Printf("pop %v\n", args[0].(string))
		//s.addCode("\tpop\t%v\n", args[0].(string))
		var err error
		code := []uint64{vm.OP_POP, 0}
		va, found := t.vars[args[0].(string)]
		if !found {
			id := t.NewId()
			va, err = section.NewVariable(id, args[0].(string),
				new(big.Rat))
			if err != nil {
				panic(err)
			}
			t.varsA = append(t.varsA, va)
		}
		code[1] = va.Id
		t.addCode(code)
		//fmt.Printf("pop %v\n", code)

	case ast.UMINUS:
		//fmt.Printf("umin\n")
		//s.addCode("\tneg\n")
		t.addCode([]uint64{vm.OP_NEG})
		//fmt.Printf("neg %v\n", code)
	case '+':
		//fmt.Printf("plus\n")
		t.addCode([]uint64{vm.OP_ADD})
		//fmt.Printf("add %v\n", code)
		//s.addCode("\tadd\n")
	case '-':
		//fmt.Printf("min\n")
		//s.addCode("\tsub\n")
		t.addCode([]uint64{vm.OP_SUB})
		//fmt.Printf("sub %v\n", code)
	case '*':
		//fmt.Printf("mul\n")
		//s.addCode("\tmul\n")
		t.addCode([]uint64{vm.OP_MUL})
		//fmt.Printf("mul %v\n", code)
	case '/':
		//fmt.Printf("div\n")
		//s.addCode("\tdiv\n")
		t.addCode([]uint64{vm.OP_DIV})
		//fmt.Printf("div %v\n", code)
	case ast.DEBUG:
		// ignore
	default:
		panic(fmt.Sprintf("unsuported pseudo opcode %v", t))
	}
}
