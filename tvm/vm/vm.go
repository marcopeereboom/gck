// vm package is the Virtual Machine component of the Go Compiler Kit.
//
// The Vm type contains an instance (or context) of a virtual machine.
// Prior to execution it must be primed with an image that contains at least
// a code section.
// The other sections are optional.
//
// Once the machine is primed it can be told to execute via the Run command.
// The VM deals with uint64 only in the code section.
// A single instruction is expressed by <OPCODE>[PARAMETER].
// Parameter is required by some but not all opcodes.
//
// The VM is mostly a RPN (Reverse Polish Notation) type machine with some
// exceptions.
// This means that for example to execute a = b + c the machine does the
// following:
//	PUSH b
//	PUSH c
//	ADD
//	POP  a
//
// Each cycle the VM fetches an opcode and, if required, a parameter.
// Parameters are almost always stored in the symbol table and are indexed
// by the uint64 parameter.
// For example, JSR 0x1234 looks up symbol with index 0x1234, dereferences it
// and then jumps to the location stored in the symbol.
// The bottom 256 (SymReserved) opcodes are reserved to express simple things
// such as TRUE and FALSE.
package vm

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/marcopeereboom/gck/tvm/section"
	"github.com/marcopeereboom/gck/tvm/stdlib"
)

// Opcodes, OP_ABORT must be 0 and OP_INVALID must always be last.
// Opcodes must be consecutive!
const (
	OP_ABORT   = 0
	OP_EXIT    = 1
	OP_NOP     = 2
	OP_PUSH    = 3
	OP_POP     = 4
	OP_ADD     = 5
	OP_SUB     = 6
	OP_MUL     = 7
	OP_DIV     = 8
	OP_NEG     = 9  // unary minus
	OP_JSR     = 10 // jump to subroutine
	OP_EQ      = 11 // ==
	OP_NEQ     = 12 // !=
	OP_LT      = 13 // <
	OP_GT      = 14 // >
	OP_LE      = 15 // <=
	OP_GE      = 16 // >=
	OP_BRT     = 17 // branch if true
	OP_CALL    = 18 // stdlib call
	OP_JMP     = 19 // jump to location
	OP_RET     = 20 // return from subroutine
	OP_INVALID = 21 // must be last
)

const (
	VmInvalidStack = iota
	VmCmdStack
	VmCallStack
)

const (
	vmInitialStackSize     = 1024
	vmInitialCallStackSize = 1024
)

var (
	vmInstructions = []instruction{
		// special one word opcodes
		{1, 0, VmInvalidStack, "abort"},
		{1, 0, VmInvalidStack, "exit"},
		{1, 0, VmInvalidStack, "nop"},

		// require symbol table
		{2, 0, VmCmdStack, "push"},
		{2, 1, VmCmdStack, "pop"},
		{1, 2, VmCmdStack, "add"},
		{1, 2, VmCmdStack, "sub"},
		{1, 2, VmCmdStack, "mul"},
		{1, 2, VmCmdStack, "div"},
		{1, 1, VmCmdStack, "neg"},
		{2, 0, VmCallStack, "jsr"},
		{1, 2, VmCmdStack, "eq"},
		{1, 2, VmCmdStack, "neq"},
		{1, 2, VmCmdStack, "lt"},
		{1, 2, VmCmdStack, "gt"},
		{1, 2, VmCmdStack, "le"},
		{1, 2, VmCmdStack, "ge"},
		{2, 1, VmCmdStack, "brt"},
		{2, 0, VmCmdStack, "call"},

		// don't require symbol table
		{2, 0, VmInvalidStack, "jmp"},
		{1, 1, VmCallStack, "ret"},

		// marks end of opcode list
		{0, 0, VmInvalidStack, "invalid"},
	}
)

type instruction struct {
	size  uint64 // how many uints total
	stack int    // how much stack needed
	which int    // which stack are we manipulating
	name  string // disassembled name
}

// Vm is the Virtual Machine context
type Vm struct {
	sym map[uint64]*section.Symbol // symbol table

	// stacks
	sp        int      // stack pointer
	stack     []uint64 // stack
	cs        int      // call stack pointer
	callStack []uint64 // call stack, contains return addresses

	// cooked sections images
	prog []uint64

	// debug
	singleStep   bool   // set to true to step through code
	trace        bool   // set to true to keep an execution trace
	traceVerbose bool   // set to create more verbose traces
	runTrace     string // runtime trace
}

// generate random uint64
func randomUint64() (uint64, error) {
	x := make([]byte, 8)
	n, err := rand.Read(x)
	if err != nil || n != 8 {
		return 0, fmt.Errorf("could not create random ID")
	}
	id := binary.LittleEndian.Uint64(x)
	return id, nil
}

func (v *Vm) GetId() (uint64, error) {
	var (
		id  uint64
		err error
	)

	// generate new symbol id
	for {
		id, err = randomUint64()
		if err != nil {
			break
		}
		// let's reserve the bottom SymReserved for things such as bool
		// values and stuff
		if id < section.SymReserved {
			continue
		}
		if v.sym[id] == nil {
			break
		}
	}
	return id, err
}

func New(image []byte) (*Vm, error) {
	v := Vm{
		stack:     make([]uint64, vmInitialStackSize),
		callStack: make([]uint64, vmInitialCallStackSize),
		sym:       make(map[uint64]*section.Symbol),
	}

	sections, err := section.SectionsFromImage(image)
	if err != nil {
		return nil, err
	}
	var ok bool
	for _, s := range sections {
		if s.Payload == nil {
			return nil, fmt.Errorf("section %v is nil",
				section.Sections[section.CodeId])
		}
		switch s.Id {
		case section.CodeId:
			v.prog, ok = s.Payload.([]uint64)
			if !ok {
				return nil, fmt.Errorf("invalid type %T "+
					"for code section", s.Payload)
			}
			if len(v.prog) == 0 {
				return nil, fmt.Errorf("empty code section")
			}

		case section.VariableId:
			vars, ok := s.Payload.([]*section.Variable)
			if !ok {
				return nil, fmt.Errorf("invalid type %T "+
					"for variable section", s.Payload)
			}

			for _, newVar := range vars {
				sym, err := section.New(newVar.Id,
					section.VariableId,
					1, // permanent value
					newVar.Name,
					newVar.GetActualValue())
				if err != nil {
					return nil, err
				}
				v.sym[sym.Id] = sym
			}

		case section.ConstId:
			consts, ok := s.Payload.([]*section.Const)
			if !ok {
				return nil, fmt.Errorf("invalid type %T "+
					"for const section", s.Payload)
			}

			for _, newConst := range consts {
				sym, err := section.New(newConst.Id,
					section.ConstId,
					1, // permanent value
					newConst.Name,
					newConst.GetActualValue())
				if err != nil {
					return nil, err
				}
				v.sym[sym.Id] = sym
			}

		case section.OsId:
			oss, ok := s.Payload.([]*section.Os)
			if !ok {
				return nil, fmt.Errorf("invalid type %T "+
					"for os section", s.Payload)
			}

			for _, newOs := range oss {
				sym, err := section.New(newOs.Id,
					section.OsId,
					1, // permanent value
					newOs.Name,
					newOs.GetActualValue())
				if err != nil {
					return nil, err
				}
				v.sym[sym.Id] = sym
			}

		default:
			return nil, fmt.Errorf("invalid section 0x%0x", s.Id)
		}
	}

	// cross reference stdlib symbols
	// XXX we need some sort of reverse lookup from the images
	//funcs := stdlib.GetFunctionNames()
	//for _, fn := range funcs {
	//	id, err := v.GetId()
	//	if err != nil {
	//		return nil, err
	//	}
	//	sym, err := section.New(id, section.OsId, section.SymLabelId,
	//		"os."+fn, fn)
	//	if err != nil {
	//		return nil, err
	//	}
	//	v.sym[sym.Id] = sym
	//}

	return &v, nil
}

func (v *Vm) GC() {
	for k, val := range v.sym {
		if val.RefC > 0 {
			continue
		}
		val.Value = nil
		delete(v.sym, k)
	}
}

func (v *Vm) GetTrace() string {
	return v.runTrace
}

func (v *Vm) Trace(loud bool) {
	v.trace = true
	v.traceVerbose = loud
}

func (v *Vm) SingleStep() {
	v.singleStep = true
}

func (v *Vm) GetSymbols(loud bool) string {
	var s string
	for k := range v.sym {
		s += v.demangle(loud, k) + "\n"
	}

	return s
}

func (v *Vm) GetStack(loud bool, which int) string {
	var (
		sp    int
		stack []uint64
	)
	switch which {
	case VmCmdStack:
		sp = v.sp
		stack = v.stack
	case VmCallStack:
		sp = v.cs
		stack = v.callStack
	default:
		return "INVALID STACK"
	}

	var s string
	for i := 0; i < sp; i++ {
		s += fmt.Sprintf("%016x: %v\n", i, v.demangle(loud, stack[i]))
	}
	return s
}

func (v *Vm) demangle(loud bool, id uint64) string {
	var (
		sym   *section.Symbol
		found bool
	)

	// handle special ids
	switch id {
	case 0:
		return "FALSE"
	case 1:
		return "TRUE"
	case 2:
		return "DISCARD"
	default:
		sym, found = v.sym[id]
		if !found {
			return fmt.Sprintf("%016x", id)
		}

	}

	var val interface{}
	switch valt := sym.Value.(type) {
	case uint64:
		val = fmt.Sprintf("0x%0x", valt)
	default:
		val = sym.Value
	}

	if loud {
		return fmt.Sprintf("%-8v %-8v %3v   %-16v  %v",
			section.Sections[sym.SectionId],
			section.Symbols[sym.TypeId],
			sym.RefC,
			sym.Name,
			val)
	}
	return fmt.Sprintf("%v (%v)", sym.Name, val)
}

func (v *Vm) disassemble(loud bool, pc uint64, prog []uint64) string {
	var (
		args, h string
		i       uint64
	)
	ins := prog[pc]
	for i = 0; i < vmInstructions[ins].size-1; i++ {
		args += " " + v.demangle(loud, prog[pc+i+1])
	}
	if loud {
		todo := 2
		for i = 0; i < vmInstructions[ins].size; i++ {
			h += fmt.Sprintf("%016x  ", prog[pc+i])
			todo--
		}
		for todo != 0 {
			h += fmt.Sprintf("%16s  ", "")
			todo--
		}
	}
	return fmt.Sprintf("%v%-6v%v", h, vmInstructions[ins].name, args)
}

// missing: pause/unpause, step, breakpoints, load/save snapshot, backtrace

func (v *Vm) Run() error {
	if len(v.prog) == 0 {
		return fmt.Errorf("no code section")
	}

	prog := v.prog
	var pc uint64 = 0
	for pc < uint64(len(prog)) {
		i := prog[pc]

		// we try to validate as much as possible up front to keep
		// opcode functions simple
		if pc+vmInstructions[i].size > uint64(len(prog)) {
			return fmt.Errorf("pc out of bounds 0x%0x", pc)
		}

		// make sure stack doesn't underflow
		switch vmInstructions[i].which {
		case VmCmdStack:
			if v.sp-vmInstructions[i].stack < 0 {
				return fmt.Errorf("command stack underflow")
			}
		case VmCallStack:
			if v.cs-vmInstructions[i].stack < 0 {
				return fmt.Errorf("call stack underflow")
			}
		}

		// keep runtime trace
		if v.trace {
			v.runTrace += fmt.Sprintf("%016x: %v\n",
				pc, v.disassemble(v.traceVerbose, pc, prog))
		}

		// jump to command
		switch i {
		case OP_ABORT:
			return fmt.Errorf("aborted at %016x", pc)
		case OP_EXIT:
			return nil
		case OP_NOP:
		case OP_PUSH:
			v.push(pc, prog)
		case OP_POP:
			if err := v.pop(pc, prog); err != nil {
				return err
			}
		case OP_ADD:
			if err := v.add(pc, prog); err != nil {
				return err
			}
		case OP_SUB:
			if err := v.sub(pc, prog); err != nil {
				return err
			}
		case OP_MUL:
			if err := v.mul(pc, prog); err != nil {
				return err
			}
		case OP_DIV:
			if err := v.div(pc, prog); err != nil {
				return err
			}
		case OP_NEG:
			if err := v.neg(pc, prog); err != nil {
				return err
			}
		case OP_EQ:
			if err := v.eq(pc, prog); err != nil {
				return err
			}
		case OP_NEQ:
			if err := v.neq(pc, prog); err != nil {
				return err
			}
		case OP_LT:
			if err := v.lt(pc, prog); err != nil {
				return err
			}
		case OP_GT:
			if err := v.gt(pc, prog); err != nil {
				return err
			}
		case OP_LE:
			if err := v.le(pc, prog); err != nil {
				return err
			}
		case OP_GE:
			if err := v.ge(pc, prog); err != nil {
				return err
			}
		case OP_CALL:
			if err := v.call(pc, prog); err != nil {
				return err
			}
		case OP_BRT:
			if err := v.brt(&pc, prog); err != nil {
				return err
			}
			// note that OP_BRT sets the pc, so continue
			continue
		case OP_JMP:
			if err := v.jmp(&pc, prog); err != nil {
				return err
			}
			// note that OP_JMP sets the pc, so continue
			continue
		case OP_JSR:
			if err := v.jsr(&pc, prog); err != nil {
				return err
			}
			// note that OP_JSR sets the pc, so continue
			continue
		case OP_RET:
			if err := v.ret(&pc, prog); err != nil {
				return err
			}
			// note that OP_RET sets the pc, so continue
			continue
		default:
			return fmt.Errorf("illegal instruction 0x%0x at 0x%0x",
				i, pc)
		}
		pc += vmInstructions[i].size
	}
	return nil
}

func (v *Vm) stackGrow(sp int, oldStack *[]uint64, s string) {
	if sp >= len(*oldStack) {
		// enlarge stack and copy it
		stack := make([]uint64, len(*oldStack)*2)
		for k, v := range *oldStack {
			stack[k] = v
		}
		*oldStack = stack
		if v.singleStep {
			fmt.Printf("enlarge %v stack to %v\n", s, len(stack))
		}
	}
}

// ref adjust the symbols reference counter
// This is the slow path.
func (v *Vm) ref(sym uint64, c int) (int, error) {
	if sym < section.SymReserved {
		return -1, fmt.Errorf("symbol reserved: %v", v)
	}
	s, found := v.sym[sym]
	if !found {
		return -1, fmt.Errorf("symbol not found: %v", v)
	}

	rc, err := s.Ref(c)
	return rc, err
}

// OP_PUSH
func (v *Vm) push(pc uint64, prog []uint64) {
	v.stackGrow(v.sp, &v.stack, "command")
	v.ref(prog[pc+1], 1)
	v.stack[v.sp] = prog[pc+1]
	v.sp++
}

// OP_POP
func (v *Vm) pop(pc uint64, prog []uint64) error {
	defer func() {
		// toss stack value
		v.sp--
	}()

	// discard value
	if prog[pc+1] == 2 {
		src, ok := v.sym[v.stack[v.sp-1]]
		if !ok {
			return fmt.Errorf("discard symbol src not found %016x",
				v.sym[v.stack[v.sp-1]])
		}
		_, err := src.Ref(-1)
		return err
	}

	// if this is a reserved symbol id just toss the stack value
	if prog[pc+1] < section.SymReserved {
		return nil
	}

	// lookup symbols
	dst, ok := v.sym[prog[pc+1]]
	if !ok {
		return fmt.Errorf("symbol dst not found %016x", prog[pc+1])
	}
	src, ok := v.sym[v.stack[v.sp-1]]
	if !ok {
		return fmt.Errorf("symbol src not found %016x",
			v.sym[v.stack[v.sp-1]])
	}

	// check pop section
	if dst.SectionId != section.VariableId {
		return fmt.Errorf("can't pop to %v %016x",
			section.Sections[dst.SectionId],
			dst.Id)
	}

	// overwrite value with a copy
	switch sv := src.Value.(type) {
	case *big.Rat:
		// make a copy
		dst.Value = new(big.Rat).Set(sv)
	default:
		dst.Value = sv
	}
	dst.TypeId = src.TypeId

	// lower ref counter
	_, err := src.Ref(-1)
	return err

}

// generic math operation
func (v *Vm) mathOp(cb func(*big.Rat, *big.Rat) (*big.Rat, error),
	pc uint64, prog []uint64) error {
	s0, found := v.sym[v.stack[v.sp-2]]
	if !found {
		return fmt.Errorf("symbol not found -2 0x%016x",
			v.stack[v.sp-2])
	}
	s1, found := v.sym[v.stack[v.sp-1]]
	if !found {
		return fmt.Errorf("symbol not found -1 0x%016x",
			v.stack[v.sp-1])
	}

	var sym *section.Symbol

	// assert same types
	switch t := s0.Value.(type) {
	case *big.Rat:
		switch t1 := s1.Value.(type) {
		case *big.Rat:
			val, err := cb(t, t1)
			if err != nil {
				return err
			}

			// create new symbol for stack
			id, err := v.GetId()
			if err != nil {
				return err
			}

			sym, err = section.New(id, section.VariableId, 1, "",
				val)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("can't %v %T to %T",
				vmInstructions[prog[pc]].name, t, t1)
		}
	default:
		return fmt.Errorf("%v does not support type: %T",
			vmInstructions[prog[pc]].name, t)
	}

	// insert new symbol
	v.sym[sym.Id] = sym

	// adjust ref counters
	_, err := s0.Ref(-1)
	if err != nil {
		return err
	}
	_, err = s1.Ref(-1)
	if err != nil {
		return err
	}

	// replace 2 stack values with 1 answer
	v.stack[v.sp-2] = sym.Id
	v.sp--
	return nil
}

// OP_ADD
func (v *Vm) add(pc uint64, prog []uint64) error {
	return v.mathOp(func(t, t1 *big.Rat) (*big.Rat, error) {
		return new(big.Rat).Add(t, t1), nil
	}, pc, prog)
}

// OP_SUB
func (v *Vm) sub(pc uint64, prog []uint64) error {
	return v.mathOp(func(t, t1 *big.Rat) (*big.Rat, error) {
		return new(big.Rat).Sub(t, t1), nil
	}, pc, prog)
}

// OP_MUL
func (v *Vm) mul(pc uint64, prog []uint64) error {
	return v.mathOp(func(t, t1 *big.Rat) (*big.Rat, error) {
		return new(big.Rat).Mul(t, t1), nil
	}, pc, prog)
}

// OP_DIV
func (v *Vm) div(pc uint64, prog []uint64) error {
	return v.mathOp(func(t, t1 *big.Rat) (*big.Rat, error) {
		if t1.Sign() == 0 {
			return nil, fmt.Errorf("divide by 0")
		}
		return new(big.Rat).Quo(t, t1), nil
	}, pc, prog)
}

// OP_NEG
func (v *Vm) neg(pc uint64, prog []uint64) error {
	s, found := v.sym[v.stack[v.sp-1]]
	if !found {
		return fmt.Errorf("symbol not found 0x%016x", v.stack[v.sp-1])
	}

	var sym *section.Symbol

	// assert same types
	switch t := s.Value.(type) {
	case *big.Rat:
		val := new(big.Rat).Neg(t)

		// create new symbol for stack
		id, err := v.GetId()
		if err != nil {
			return err
		}
		sym, err = section.New(id, section.VariableId, 1, "", val)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%v does not support type: %T",
			vmInstructions[prog[pc]].name, t)
	}

	// insert new symbol
	v.sym[sym.Id] = sym

	// adjust ref counter of source
	_, err := s.Ref(-1)
	if err != nil {
		return err
	}

	// replace stack value
	v.stack[v.sp-1] = sym.Id

	return nil
}

// generic comparison operation
func (v *Vm) cmpOp(cb func(*big.Rat, *big.Rat) (bool, error),
	pc uint64, prog []uint64) error {

	var rv bool

	s0, found := v.sym[v.stack[v.sp-2]]
	if !found {
		return fmt.Errorf("symbol not found 0x%016x", v.stack[v.sp-2])
	}
	s1, found := v.sym[v.stack[v.sp-1]]
	if !found {
		return fmt.Errorf("symbol not found 0x%016x", v.stack[v.sp-1])
	}

	// assert same types
	switch t := s0.Value.(type) {
	case *big.Rat:
		switch t1 := s1.Value.(type) {
		case *big.Rat:
			var errOp error
			rv, errOp = cb(t, t1)
			if errOp != nil {
				return errOp
			}
		default:
			return fmt.Errorf("can't %v %T to %T",
				vmInstructions[prog[pc]].name, t, t1)
		}
	default:
		return fmt.Errorf("%v does not support type: %T",
			vmInstructions[prog[pc]].name, t)
	}

	// adjust ref counters
	_, err := s0.Ref(-1)
	if err != nil {
		return err
	}
	_, err = s1.Ref(-1)
	if err != nil {
		return err
	}

	v.sp--
	if rv {
		v.stack[v.sp-1] = 1
	} else {
		v.stack[v.sp-1] = 0
	}

	return nil
}

// OP_EQ
func (v *Vm) eq(pc uint64, prog []uint64) error {
	return v.cmpOp(func(t, t1 *big.Rat) (bool, error) {
		return 0 == t.Cmp(t1), nil
	}, pc, prog)
}

// OP_NEQ
func (v *Vm) neq(pc uint64, prog []uint64) error {
	return v.cmpOp(func(t, t1 *big.Rat) (bool, error) {
		return 0 != t.Cmp(t1), nil
	}, pc, prog)
}

// OP_LT
func (v *Vm) lt(pc uint64, prog []uint64) error {
	return v.cmpOp(func(t, t1 *big.Rat) (bool, error) {
		return -1 == t.Cmp(t1), nil
	}, pc, prog)
}

// OP_GT
func (v *Vm) gt(pc uint64, prog []uint64) error {
	return v.cmpOp(func(t, t1 *big.Rat) (bool, error) {
		return 1 == t.Cmp(t1), nil
	}, pc, prog)
}

// OP_LE
func (v *Vm) le(pc uint64, prog []uint64) error {
	return v.cmpOp(func(t, t1 *big.Rat) (bool, error) {
		return 0 >= t.Cmp(t1), nil
	}, pc, prog)
}

// OP_GE
func (v *Vm) ge(pc uint64, prog []uint64) error {
	return v.cmpOp(func(t, t1 *big.Rat) (bool, error) {
		return 0 <= t.Cmp(t1), nil
	}, pc, prog)
}

// OP_JMP
func (v *Vm) jmp(pc *uint64, prog []uint64) error {
	location := prog[*pc+1]
	if location >= uint64(len(prog)) {
		return fmt.Errorf("jmp out of bounds")
	}
	*pc = location
	return nil
}

// OP_CALL
// ABI vars pushed in reverse order
// always pushes success/failure on return
func (v *Vm) call(pc uint64, prog []uint64) error {
	// lookup label in symbol table
	s, found := v.sym[prog[pc+1]]
	if !found {
		return fmt.Errorf("symbol not found 0x%016x", prog[pc+1])
	}

	// validate symbol
	if s.SectionId != section.OsId {
		return fmt.Errorf("call can not jump to section %v",
			section.Sections[s.SectionId])
	}
	if s.TypeId != section.SymLabelId {
		return fmt.Errorf("call can not jump to type %v",
			section.Symbols[s.TypeId])
	}

	// XXX handle args and return values
	rv, err := stdlib.Dispatch(s.Name)
	if err != nil {
		return err
	}

	// push sucess/failure on the stack
	if rv.Error == nil {
		v.stack[v.sp] = 1
	} else {
		v.stack[v.sp] = 0
	}
	v.sp++

	return nil
}

// OP_JSR
func (v *Vm) jsr(pc *uint64, prog []uint64) error {
	// lookup label in symbol table
	s, found := v.sym[prog[*pc+1]]
	if !found {
		return fmt.Errorf("jsr symbol not found 0x%016x",
			prog[*pc+1])
	}

	location, ok := s.Value.(uint64)
	if !ok {
		return fmt.Errorf("jsr invalid label type %T", s.Value)
	}

	// validate symbol
	if s.SectionId != section.ConstId {
		return fmt.Errorf("jsr can not jump using a symbol from "+
			"section %v", section.Sections[s.SectionId])
	}
	if s.TypeId != section.SymLabelId {
		return fmt.Errorf("jsr can not jump using type %v",
			section.Symbols[s.TypeId])
	}

	// validate location
	if location >= uint64(len(prog)) {
		return fmt.Errorf("jsr out of bounds")
	}

	// set return address
	ret := *pc + 2
	if ret >= uint64(len(prog)) {
		return fmt.Errorf("jsr return value out of bounds")
	}
	v.stackGrow(v.cs, &v.callStack, "call")
	v.callStack[v.cs] = ret
	v.cs++

	*pc = location
	return nil
}

// OP_BRT
func (v *Vm) brt(pc *uint64, prog []uint64) error {
	v.sp--
	rv := v.stack[v.sp]
	if rv != 1 {
		*pc += 2
		return nil
	}
	location := prog[*pc+1]
	if location >= uint64(len(prog)) {
		return fmt.Errorf("brt out of bounds")
	}
	*pc = location
	return nil
}

// OP_RET
func (v *Vm) ret(pc *uint64, prog []uint64) error {
	v.cs--
	ret := v.callStack[v.cs]
	if ret >= uint64(len(prog)) {
		return fmt.Errorf("ret return value out of bounds")
	}
	*pc = ret
	return nil
}
