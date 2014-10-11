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
	lbls    map[int]uint64               // labels by id
	fixup   map[int]uint64               // labels that need fixing up
	code    []uint64
}

// endure interface is met
var _ arch.Backend = &ToyVirtualMachine{}

// New creates a new ToyVirtualMachine context.
func New() (*ToyVirtualMachine, error) {
	vm := ToyVirtualMachine{
		vars:    make(map[string]*section.Variable),
		constsL: make(map[string]*section.Const),
		lbls:    make(map[int]uint64),
		fixup:   make(map[int]uint64),
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
func (t *ToyVirtualMachine) getConst(value interface{}) (*section.Const, error) {
	var (
		err error
		v   string
	)

	switch val := value.(type) {
	case *big.Rat:
		v = val.String()
	case int:
		v = strconv.Itoa(val)
	default:
		return nil, fmt.Errorf("invalid type for .CONST %T", value)
	}

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
func (t *ToyVirtualMachine) emitCode(ty int, args ...interface{}) error {
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

	case ast.INTEGER:
		c, err := t.getConst(args[0].(int))
		if err != nil {
			panic(err)
		}
		t.addCode([]uint64{vm.OP_PUSH, c.Id})

	case ast.Assign:
		va, err := t.getVar(args[0].(string))
		if err != nil {
			panic(err)
		}
		t.addCode([]uint64{vm.OP_POP, va.Id})

	case ast.Uminus:
		t.addCode([]uint64{vm.OP_NEG})

	case ast.Add:
		t.addCode([]uint64{vm.OP_ADD})

	case ast.Sub:
		t.addCode([]uint64{vm.OP_SUB})

	case ast.Mul:
		t.addCode([]uint64{vm.OP_MUL})

	case ast.Div:
		t.addCode([]uint64{vm.OP_DIV})

	case ast.Lt:
		t.addCode([]uint64{vm.OP_LT})

	case ast.Gt:
		t.addCode([]uint64{vm.OP_GT})

	case ast.Le:
		t.addCode([]uint64{vm.OP_LE})

	case ast.Ge:
		t.addCode([]uint64{vm.OP_GE})

	case ast.Ne:
		t.addCode([]uint64{vm.OP_NEQ})

	case ast.Eq:
		t.addCode([]uint64{vm.OP_EQ})

	case ast.BRT:
		i := args[0].(int)
		jl, ok := t.lbls[i]
		if !ok {
			// store fixup location as [label index] memory location
			jl = 0xffffffffffffffff
			t.fixup[i] = uint64(len(t.code)) + 1
		}
		t.addCode([]uint64{vm.OP_BRT, jl})

	case ast.BRF:
		i := args[0].(int)
		jl, ok := t.lbls[i]
		if !ok {
			// store fixup location as [label index] memory location
			jl = 0xffffffffffffffff
			t.fixup[i] = uint64(len(t.code)) + 1
		}
		t.addCode([]uint64{vm.OP_BRF, jl})

	case ast.JUMP:
		i := args[0].(int)
		jl, ok := t.lbls[i]
		if !ok {
			// store fixup location as [label index] memory location
			jl = 0xffffffffffffffff
			t.fixup[i] = uint64(len(t.code)) + 1
		}
		t.addCode([]uint64{vm.OP_JMP, jl})

	case ast.LOCATION:
		_, found := t.lbls[args[0].(int)]
		if found {
			return fmt.Errorf("label already exists %v", args[0])
		}
		t.lbls[args[0].(int)] = uint64(len(t.code))

	case ast.FIXUP:
		for _, v := range args {
			// v = label index
			// t.fixup[v] = memory location that needs to be fixed
			// t.lbls[v] = value for fixup
			t.code[t.fixup[v.(int)]] = t.lbls[v.(int)]
		}

	case ast.NOP:
		t.addCode([]uint64{vm.OP_NOP})

	case ast.EXIT:
		t.addCode([]uint64{vm.OP_EXIT})

	case ast.DEBUG:
		// ignore

	default:
		return fmt.Errorf("unsuported pseudo opcode %v", ty)
	}

	return nil
}
