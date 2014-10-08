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

// ToyVirtualMachine is the context for the arch.Backend enterface.
type ToyVirtualMachine struct {
	id      uint64
	consts  []*section.Const
	constsL map[string]*section.Const // lookup by value
	varsA   []*section.Variable
	vars    map[string]*section.Variable // lookup by name
	code    []uint64
}

// endure interface is met
var _ arch.Backend = &ToyVirtualMachine{}

// New creates a new ToyVirtualMachine context.
func New() (*ToyVirtualMachine, error) {
	vm := ToyVirtualMachine{
		vars:    make(map[string]*section.Variable),
		constsL: make(map[string]*section.Const),
		id:      1000,
		code:    make([]uint64, 0, 1000),
	}

	return &vm, nil
}

// newId generates a new variable or contant identifier.
// Identifiers must be unique since they are keys to symbol table.
func (t *ToyVirtualMachine) newId() uint64 {
	defer func() { t.id++ }()
	return t.id
}

// EmitCode implements the arch.Backend interface.
// It converts ast.Node n into code, variables, contants etc.
// The result is a Toy Virtual Machine image that can be executed.
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

// addCode adds opcodes and variables from code into the code section.
func (t *ToyVirtualMachine) addCode(code []uint64) {
	// blow me slow
	for _, v := range code {
		t.code = append(t.code, v)
	}
}

// getVar looks up a variable by name and return a new Variable structure if
// the variable name does not exist.
// If the variable name does exist it returns the existing structure instead.
func (t *ToyVirtualMachine) getVar(name string) (*section.Variable, error) {
	var err error

	va, found := t.vars[name]
	if !found {
		// emit empty variable, NUMBER for now
		id := t.newId()
		va, err = section.NewVariable(id, name, new(big.Rat))
		if err != nil {
			return nil, err
		}

		t.varsA = append(t.varsA, va)
		t.vars[name] = va
	}

	return va, nil
}

// getConst looks up a constant by value and returns a new Const structure if
// the value does not exist.
// If the constant values does exist it returns the existing structure instead.
// This eliminates duplicate constant values in the symbol table.
func (t *ToyVirtualMachine) getConst(value *big.Rat) (*section.Const, error) {
	var err error

	v := value.String()
	c, found := t.constsL[v]
	if !found {
		id := t.newId()
		c, err = section.NewConst(id, "c"+strconv.Itoa(int(id)), value)
		if err != nil {
			return nil, err
		}
		t.consts = append(t.consts, c)
		t.constsL[v] = c
	}

	return c, nil
}

// emitCode convert ast.Node into code, variables and constants.
func (t *ToyVirtualMachine) emitCode(ty int, args ...interface{}) {
	switch ty {
	case ast.IDENTIFIER:
		va, err := t.getVar(args[0].(string))
		if err != nil {
			panic(err)
		}
		t.addCode([]uint64{vm.OP_PUSH, va.Id})

	case ast.NUMBER:
		c, err := t.getConst(args[0].(*big.Rat))
		if err != nil {
			panic(err)
		}
		t.addCode([]uint64{vm.OP_PUSH, c.Id})

	case '=':
		va, err := t.getVar(args[0].(string))
		if err != nil {
			panic(err)
		}
		t.addCode([]uint64{vm.OP_POP, va.Id})

	case ast.UMINUS:
		t.addCode([]uint64{vm.OP_NEG})
	case '+':
		t.addCode([]uint64{vm.OP_ADD})
	case '-':
		t.addCode([]uint64{vm.OP_SUB})
	case '*':
		t.addCode([]uint64{vm.OP_MUL})
	case '/':
		t.addCode([]uint64{vm.OP_DIV})
	case ast.DEBUG:
		// ignore
	default:
		panic(fmt.Sprintf("unsuported pseudo opcode %v", t))
	}
}
