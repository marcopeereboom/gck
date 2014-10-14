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
//
// Note that a bunch of documentation isn't visible on godoc.org.
// They do not enable ?m=all in URLs (include unexported doco).
// So make sure to reference the source or run godoc locally to see the
// additional documentation.
package vm

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/marcopeereboom/gck/tvm/section"
	"github.com/marcopeereboom/gck/tvm/stdlib"
)

const (
	// Opcodes, OP_ABORT must be 0 and OP_INVALID must always be last.
	// Opcodes must be consecutive!
	OP_ABORT   = 0  // abort execution, exception
	OP_EXIT    = 1  // exit program, not an exception
	OP_NOP     = 2  // no-op
	OP_PUSH    = 3  // push symbol id or something else onto command stack
	OP_POP     = 4  // pop symbol id or something else from command stack
	OP_ADD     = 5  // add top 2 values on the command stack
	OP_SUB     = 6  // subtract top 2 values on the command stack
	OP_MUL     = 7  // multiply top 2 values on the command stack
	OP_DIV     = 8  // divide top 2 values on the command stack
	OP_NEG     = 9  // unary minus
	OP_JSR     = 10 // jump to subroutine
	OP_EQ      = 11 // ==
	OP_NEQ     = 12 // !=
	OP_LT      = 13 // <
	OP_GT      = 14 // >
	OP_LE      = 15 // <=
	OP_GE      = 16 // >=
	OP_BRT     = 17 // branch if true
	OP_BRF     = 18 // branch if false
	OP_CALL    = 19 // stdlib call
	OP_JMP     = 20 // jump to location
	OP_RET     = 21 // return from subroutine
	OP_INVALID = 22 // must be last
)

const (
	// Constants that define what stack to use.
	VmInvalidStack = iota
	VmCmdStack
	VmCallStack
)

const (
	// Constants that define default stack size.
	// This is denominated in uint64.
	vmInitialStackSize     = 1024
	vmInitialCallStackSize = 1024
)

// instruction describes a VM instruction and its limits.
type instruction struct {
	size  uint64 // how many uint64s total
	stack int    // how much stack needed
	which int    // which stack are we manipulating
	name  string // disassembled name
}

var (
	// vmInstructions is an array used for disassembly and it contains
	// limits for the stacks.
	// This array must match the opcodes.
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
		{2, 1, VmCmdStack, "brf"},
		{2, 0, VmCmdStack, "call"},

		// don't require symbol table
		{2, 0, VmInvalidStack, "jmp"},
		{1, 1, VmCallStack, "ret"},

		// marks end of opcode list
		{0, 0, VmInvalidStack, "invalid"},
	}
)

// Vm is the Virtual Machine context
type Vm struct {
	sym map[uint64]*section.Symbol // symbol table

	// stacks
	sp        int      // stack pointer
	stack     []uint64 // stack
	cs        int      // call stack pointer
	callStack []uint64 // call stack, contains return addresses

	// gc
	zero uint64

	// cooked sections images
	prog []uint64

	// debug
	singleStep   bool   // set to true to step through code
	trace        bool   // set to true to keep an execution trace
	traceVerbose bool   // set to create more verbose traces
	runTrace     string // runtime trace

	// stats
	instructions uint64
}

// randomUint64 generates a random uint64 value.
func randomUint64() (uint64, error) {
	x := make([]byte, 8)
	n, err := rand.Read(x)
	if err != nil || n != 8 {
		return 0, fmt.Errorf("could not create random ID")
	}
	id := binary.LittleEndian.Uint64(x)
	return id, nil
}

// GetId returns a valid random uint64 value that can be used to designate
// an entry in the symbol table.
// It excludes values that are in use and some reserved IDs.
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

// New creates a new VM context for image.
// If the image is invalid the function throws an error.
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

// GC garbage collect symbols that have a reference counter that is less than
// 1.
func (v *Vm) GC() {
	for k, val := range v.sym {
		if val.RefC > 0 {
			continue
		}
		val.Value = nil
		delete(v.sym, k)
		v.zero--
	}
}

// GetTrace returns the runtime trace.
// Note that this is not a traditional backtrace.
// This is a trace of all instructions the machine actually ran.
// In order to be able to obtain the runtime trace one must enable it by
// calling Trace before Run.
// This functionality is not enabled by default due to performance reasons.
func (v *Vm) GetTrace() string {
	return v.runTrace
}

// Trace enables runtime tracing.
// Set loud to true for extra verbosity.
// Enabling this impacts performance negatively.
func (v *Vm) Trace(loud bool) {
	v.trace = true
	v.traceVerbose = loud
}

// SingleStep enables single step mode.
// THIS IS CURRENTLY NOT IMPLEMENTED.
func (v *Vm) SingleStep() {
	v.singleStep = true
}

// GetSymbols returns all symbols from the symbol table.
// Set loud to true for extra verbosity.
func (v *Vm) GetSymbols(loud bool) string {
	var s string
	for k := range v.sym {
		s += v.demangle(loud, k) + "\n"
	}

	return s
}

// GetSstack returns the stack as indicated by which.
// Set loud to true for extra verbosity.
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

// demangle returns a human readable symbol.
// Set loud to true for extra verbosity.
func (v *Vm) demangle(loud bool, id uint64) string {
	var (
		sym   *section.Symbol
		found bool
	)

	// handle special ids
	switch id {
	case section.SymReservedFalse, section.SymReservedTrue,
		section.SymReservedDiscard:
		return section.SymbolsReserved[id]
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

// disassemble returns a human readable opcode.
// Set loud to true for extra verbosity.
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

func (v *Vm) Run() error {
	return v.run(make(chan string), false)
}

// missing: pause/unpause, step, breakpoints, load/save snapshot, backtrace

// Run start executing the image that was provided during New.
// If the program violates any rules it will be aborted and Run will return
// an error.
func (v *Vm) run(c chan string, interactive bool) error {
	if len(v.prog) == 0 {
		return fmt.Errorf("no code section")
	}

	v.instructions = 0
	paused := false

	prog := v.prog
	var pc uint64
	for pc < uint64(len(prog)) {
		// when running interactively collect some stats and some more
		// stuff
		if interactive {
			if paused {
				// paused, block
				cmd := <-c
				fmt.Printf("received paused message %v\n", cmd)
				paused = false
			} else {
				// not paused, don't block
				select {
				case cmd := <-c:
					fmt.Printf("received unpaused message %v\n", cmd)
					paused = true
				default:
				}
			}
			v.instructions++
		}

		// see if we should gc, this is pretty arbitrary
		if v.zero > 5000 {
			v.GC()
		}

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
		case OP_BRF:
			if err := v.brf(&pc, prog); err != nil {
				return err
			}
			// note that OP_BRF sets the pc, so continue
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

// stackGrow validates if the current stack is large enough to handle a push.
// If the stack is not big enough it will be doubled in size.
// Note that stacks never shrink.
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
// This is the slow path and should be avoided if possible.
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

// push handles the OP_PUSH opcode.
// It pushes a reserved or symbol ID onto the stack.
// The stack pointer is incremented by exactly one uint64.
// Push automatically grows the stack if needed.
func (v *Vm) push(pc uint64, prog []uint64) {
	v.stackGrow(v.sp, &v.stack, "command")
	v.ref(prog[pc+1], 1)
	v.stack[v.sp] = prog[pc+1]
	v.sp++
}

// push handles the OP_POP opcode.
// It pops a reserved or symbol ID from the stack into a symbol id.
// Popping a reserved value will result in that value being discarded.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) pop(pc uint64, prog []uint64) error {
	defer func() {
		// toss stack value
		v.sp--
	}()

	// discard value
	if prog[pc+1] == section.SymReservedDiscard {
		src, ok := v.sym[v.stack[v.sp-1]]
		if !ok {
			return fmt.Errorf("discard symbol src not found %016x",
				v.sym[v.stack[v.sp-1]])
		}
		rc, err := src.Ref(-1)
		if rc == 0 {
			v.zero++
		}
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
	rc, err := src.Ref(-1)
	if rc == 0 {
		v.zero++
	}
	return err

}

// mathOp handles generic math operations.
// See individual opcodes for descriptions.
func (v *Vm) mathOp(cb func(int, interface{}, interface{}) (interface{}, error),
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
	// NUMBERS
	case *big.Rat:
		switch t1 := s1.Value.(type) {
		case *big.Rat:
			val, err := cb(section.SymNumId, t, t1)
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

	// INTEGERS
	case int:
		switch t1 := s1.Value.(type) {
		case int:
			val, err := cb(section.SymIntId, t, t1)
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
	rc, err := s0.Ref(-1)
	if err != nil {
		return err
	}
	if rc == 0 {
		v.zero++
	}
	rc, err = s1.Ref(-1)
	if err != nil {
		return err
	}
	if rc == 0 {
		v.zero++
	}

	// replace 2 stack values with 1 answer
	v.stack[v.sp-2] = sym.Id
	v.sp--
	return nil
}

// add handles the OP_ADD opcode.
// It adds the top two values on the stack and replaces them with a single
// result value.
// For example:
//	push x (11)
//	push y (22)
//	add
// Results in 33 which is stored in a symbol.
// The symbol ID resides on top of the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) add(pc uint64, prog []uint64) error {
	return v.mathOp(func(mode int, t, t1 interface{}) (interface{}, error) {
		switch mode {
		case section.SymIntId:
			return t.(int) + t1.(int), nil
		case section.SymNumId:
			return new(big.Rat).Add(t.(*big.Rat), t1.(*big.Rat)), nil
		}
		return nil, fmt.Errorf("invalid add mode %v", mode)
	}, pc, prog)
}

// add handles the OP_SUB opcode.
// It subtracts the top two values on the stack and replaces them with a single
// result value.
// For example:
//	push x (3)
//	push y (2)
//	sub
// Results in 1 which is stored in a symbol.
// The symbol ID resides on top of the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) sub(pc uint64, prog []uint64) error {
	return v.mathOp(func(mode int, t, t1 interface{}) (interface{}, error) {
		switch mode {
		case section.SymIntId:
			return t.(int) - t1.(int), nil
		case section.SymNumId:
			return new(big.Rat).Sub(t.(*big.Rat), t1.(*big.Rat)), nil
		}
		return nil, fmt.Errorf("invalid sub mode %v", mode)
	}, pc, prog)
}

// add handles the OP_MUL opcode.
// It multiplies the top two values on the stack and replaces them with a single
// result value.
// For example:
//	push x (2)
//	push y (5)
//	mul
// Results in 10 which is stored in a symbol.
// The symbol ID resides on top of the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) mul(pc uint64, prog []uint64) error {
	return v.mathOp(func(mode int, t, t1 interface{}) (interface{}, error) {
		switch mode {
		case section.SymIntId:
			return t.(int) * t1.(int), nil
		case section.SymNumId:
			return new(big.Rat).Mul(t.(*big.Rat), t1.(*big.Rat)), nil
		}
		return nil, fmt.Errorf("invalid mul mode %v", mode)
	}, pc, prog)
}

// add handles the OP_DIV opcode.
// It divides the top two values on the stack and replaces them with a single
// result value.
// For example:
//	push x (10)
//	push y (5)
//	div
// Results in 2 which is stored in a symbol.
// The symbol ID resides on top of the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) div(pc uint64, prog []uint64) error {
	return v.mathOp(func(mode int, t, t1 interface{}) (interface{}, error) {
		switch mode {
		case section.SymIntId:
			if t1.(int) == 0 {
				return nil, fmt.Errorf("divide by 0")
			}
			return t.(int) * t1.(int), nil
		case section.SymNumId:
			if t1.(*big.Rat).Sign() == 0 {
				return nil, fmt.Errorf("divide by 0")
			}
			return new(big.Rat).Quo(t.(*big.Rat), t1.(*big.Rat)), nil
		}
		return nil, fmt.Errorf("invalid div mode %v", mode)
	}, pc, prog)
}

// neg handles the OP_NEG opcode.
// It negates the top of the stack value.
// result value.
// For example:
//	push x (10)
//	neg
// Results in -10 which is stored in a symbol.
// The symbol ID resides on top of the stack.
// The symbol ID for x is no longer on the stack.
// The stack pointer is unaltered.
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

	case int:
		val := -t

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
	rc, err := s.Ref(-1)
	if err != nil {
		return err
	}
	if rc == 0 {
		v.zero++
	}

	// replace stack value
	v.stack[v.sp-1] = sym.Id

	return nil
}

// cmpOp is the generic comparison operation.
// See individual opcodes for more information.
func (v *Vm) cmpOp(cb func(int, interface{}, interface{}) (bool, error),
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
			rv, errOp = cb(section.SymNumId, t, t1)
			if errOp != nil {
				return errOp
			}
		default:
			return fmt.Errorf("can't %v %T to %T",
				vmInstructions[prog[pc]].name, t, t1)
		}

	case int:
		switch t1 := s1.Value.(type) {
		case int:
			var errOp error
			rv, errOp = cb(section.SymIntId, t, t1)
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
	rc, err := s0.Ref(-1)
	if err != nil {
		return err
	}
	if rc == 0 {
		v.zero++
	}
	rc, err = s1.Ref(-1)
	if err != nil {
		return err
	}
	if rc == 0 {
		v.zero++
	}

	v.sp--
	if rv {
		v.stack[v.sp-1] = section.SymReservedTrue
	} else {
		v.stack[v.sp-1] = section.SymReservedFalse
	}

	return nil
}

// eq handles the OP_EQ opcode.
// It compares for equality the top two values on the stack and replaces them
// with a single boolean result value.
// For example:
//	push x (10)
//	push y (10)
//	eq
// Results in TRUE which is stored as 0x1 on the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) eq(pc uint64, prog []uint64) error {
	return v.cmpOp(func(mode int, t, t1 interface{}) (bool, error) {
		switch mode {
		case section.SymNumId:
			return 0 == t.(*big.Rat).Cmp(t1.(*big.Rat)), nil
		case section.SymIntId:
			return t.(int) == t1.(int), nil
		}
		return false, fmt.Errorf("invalid == mode %v", mode)
	}, pc, prog)
}

// neq handles the OP_NEQ opcode.
// It compares for inequality the top two values on the stack and replaces them
// with a single boolean result value.
// For example:
//	push x (10)
//	push y (22)
//	neq
// Results in TRUE which is stored as 0x1 on the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) neq(pc uint64, prog []uint64) error {
	return v.cmpOp(func(mode int, t, t1 interface{}) (bool, error) {
		switch mode {
		case section.SymNumId:
			return 0 != t.(*big.Rat).Cmp(t1.(*big.Rat)), nil
		case section.SymIntId:
			return t.(int) != t1.(int), nil
		}
		return false, fmt.Errorf("invalid != mode %v", mode)
	}, pc, prog)
}

// lt handles the OP_LT opcode.
// It compares for less than of the top two values on the stack and replaces
// them with a single boolean result value.
// For example:
//	push x (10)
//	push y (22)
//	lt
// Results in TRUE which is stored as 0x1 on the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) lt(pc uint64, prog []uint64) error {
	return v.cmpOp(func(mode int, t, t1 interface{}) (bool, error) {
		switch mode {
		case section.SymNumId:
			return -1 == t.(*big.Rat).Cmp(t1.(*big.Rat)), nil
		case section.SymIntId:
			return t.(int) < t1.(int), nil
		}
		return false, fmt.Errorf("invalid < mode %v", mode)
	}, pc, prog)
}

// gt handles the OP_GT opcode.
// It compares for greater than of the top two values on the stack and replaces
// them with a single boolean result value.
// For example:
//	push x (10)
//	push y (22)
//	lt
// Results in FALSE which is stored as 0x0 on the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) gt(pc uint64, prog []uint64) error {
	return v.cmpOp(func(mode int, t, t1 interface{}) (bool, error) {
		switch mode {
		case section.SymNumId:
			return 1 == t.(*big.Rat).Cmp(t1.(*big.Rat)), nil
		case section.SymIntId:
			return t.(int) > t1.(int), nil
		}
		return false, fmt.Errorf("invalid > mode %v", mode)
	}, pc, prog)
}

// le handles the OP_LE opcode.
// It compares for greater than of the top two values on the stack and replaces
// them with a single boolean result value.
// For example:
//	push x (10)
//	push y (22)
//	le
// Results in TRUE which is stored as 0x1 on the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) le(pc uint64, prog []uint64) error {
	return v.cmpOp(func(mode int, t, t1 interface{}) (bool, error) {
		switch mode {
		case section.SymNumId:
			return 0 >= t.(*big.Rat).Cmp(t1.(*big.Rat)), nil
		case section.SymIntId:
			return t.(int) <= t1.(int), nil
		}
		return false, fmt.Errorf("invalid <= mode %v", mode)
	}, pc, prog)
}

// ge handles the OP_GE opcode.
// It compares for greater than of the top two values on the stack and replaces
// them with a single boolean result value.
// For example:
//	push x (10)
//	push y (22)
//	le
// Results in FALSE which is stored as 0x0 on the stack.
// The symbol IDs for x and y are no longer on the stack.
// The stack pointer is decremented by exactly one uint64.
func (v *Vm) ge(pc uint64, prog []uint64) error {
	return v.cmpOp(func(mode int, t, t1 interface{}) (bool, error) {
		switch mode {
		case section.SymNumId:
			return 0 <= t.(*big.Rat).Cmp(t1.(*big.Rat)), nil
		case section.SymIntId:
			return t.(int) >= t1.(int), nil
		}
		return false, fmt.Errorf("invalid >= mode %v", mode)
	}, pc, prog)
}

// jmp handles the OP_JMP opcode.
// It jumps to the location that is the opcode argument.
// Jump can only do direct jumps within the code segment.
// It validates jump boundaries before jumping and aborts if it can not
// safely perform the jump.
// For example (pc at 0x00 and jumps to 0x02):
//	0x00	jmp	0x02
//	0x01	0x02	this is the jump location
//	0x02	nop
// The stack pointer is unchanged.
func (v *Vm) jmp(pc *uint64, prog []uint64) error {
	location := prog[*pc+1]
	if location >= uint64(len(prog)) {
		return fmt.Errorf("jmp out of bounds")
	}
	*pc = location
	return nil
}

// call handles the OP_CALL opcode.
// All calls are essentially equivalent to OS standard library calls.
// TODO this is currently not fully implemented.
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
		v.stack[v.sp] = section.SymReservedTrue
	} else {
		v.stack[v.sp] = section.SymReservedFalse
	}
	v.sp++

	return nil
}

// jsr handles the OP_JSR opcode.
// It jumps to the location that is dereferenced from the symbol ID that is its
// argument.
// Jump to subroutine can only do indirect jumps within the code segment.
// It validates jump boundaries before jumping and aborts if it can not
// safely perform the jump.
// For example (pc at 0x00 and jumps to *0x1234, assume it contains 0x80):
//	0x00	jmp	0x1234
//	0x01	0x1234	this is the symbol ID that contains the jump address
//	0x02	nop
//	..
//	0x80	nop
//	0x81	ret
// The command stack pointer is unchanged.
// The call stack pointer is incremented by one and contains the return address
// that ret will jump to.
// In this example it would return to 0x02.
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

// brt handles the OP_BRT opcode.
// Branch true evaluates the top value on the stack and jumps if it is set to
// true (0x01).
// brt can only do direct jumps within the code segment.
// It validates jump boundaries before jumping and aborts if it can not
// safely perform the jump.
// For example:
//	0x00	push	TRUE
//	0x01	0x01	this is the true value
//	0x02	brt	0x05
//	0x03	0x05	this is the branch location
//	0x04	nop
//	0x05	nop
// The stack pointer is decremented by exactly one uint64.
// In this example the brt call would jump over the nop at 0x04.
func (v *Vm) brt(pc *uint64, prog []uint64) error {
	v.sp--
	rv := v.stack[v.sp]
	switch rv {
	case section.SymReservedFalse:
		*pc += 2
		return nil
	case section.SymReservedTrue:
		location := prog[*pc+1]
		if location >= uint64(len(prog)) {
			return fmt.Errorf("brt out of bounds")
		}
		*pc = location
		return nil
	}

	return fmt.Errorf("brt not testing true/false")
}

// brf handles the OP_BRF opcode.
// Branch false evaluates the top value on the stack and jumps if it is set to
// false (0x00).
// brf can only do direct jumps within the code segment.
// It validates jump boundaries before jumping and aborts if it can not
// safely perform the jump.
// For example:
//	0x00	push	FALSE
//	0x01	0x00	this is the false value
//	0x02	brf	0x05
//	0x03	0x05	this is the branch location
//	0x04	nop
//	0x05	nop
// The stack pointer is decremented by exactly one uint64.
// In this example the brf call would jump over the nop at 0x04.
func (v *Vm) brf(pc *uint64, prog []uint64) error {
	v.sp--
	rv := v.stack[v.sp]
	switch rv {
	case section.SymReservedFalse:
		location := prog[*pc+1]
		if location >= uint64(len(prog)) {
			return fmt.Errorf("brf out of bounds")
		}
		*pc = location
		return nil
	case section.SymReservedTrue:
		*pc += 2
		return nil
	}

	return fmt.Errorf("brf not testing true/false")
}

// ret handles the OP_RET opcode.
// It jumps to the location that is on top of the call stack.
// Returns can only do direct jumps within the code segment.
// It validates jump boundaries before jumping and aborts if it can not
// safely perform the jump.
// For example (pc at 0x00 and jumps to *0x1234, assume it contains 0x80):
//	0x00	jmp	0x1234
//	0x01	0x1234	this is the symbol ID that contains the jump address
//	0x02	nop
//	..
//	0x80	nop
//	0x81	ret
// The command stack pointer is unchanged.
// The call stack pointer is incremented by one and contains the return address
// that ret will jump to.
// In this example it would return to 0x02.
func (v *Vm) ret(pc *uint64, prog []uint64) error {
	v.cs--
	ret := v.callStack[v.cs]
	if ret >= uint64(len(prog)) {
		return fmt.Errorf("ret return value out of bounds")
	}
	*pc = ret
	return nil
}
