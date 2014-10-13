// ast is a poor mans abstract syntax tree imlementation.
package ast

import (
	"fmt"
	"io"
	"math/big"
)

// operations
const (
	Uminus       = 65000
	Lt           = 65001
	Gt           = 65002
	Le           = 65003
	Ge           = 65004
	Ne           = 65005
	Eq           = 65006
	Assign       = 65007
	Add          = 65008
	Sub          = 65009
	Mul          = 65010
	Div          = 65011
	Eos          = 65020
	While        = 65030
	If           = 65031
	Function     = 65032
	FunctionCall = 65033
)

// pseudo opcodes
const (
	IDENTIFIER = 0
	NUMBER     = 1
	DEBUG      = 2
	INTEGER    = 3
	LOCATION   = 4
	BRT        = 5
	JUMP       = 6
	FIXUP      = 7
	NOP        = 8
	EXIT       = 9
	BRF        = 10
	JSR        = 11
	RETURN     = 12
)

// NodeDebugInformation contains debug information that can be extracted by
// the backend etc for examination.
type NodeDebugInformation struct {
	LineNo   int    // Line number
	ColStart int    // Token column start on line
	ColEnd   int    // Token column end on line
	Line     string // Raw line text
}

// NodeIdentifier contains a string identifier.
type NodeIdentifier struct {
	Value string
}

// NewIdentifier returns an initialized NodeIdentifier structure.
func NewIdentifier(d *NodeDebugInformation, id string) Node {
	i := NodeIdentifier{
		Value: id,
	}

	return Node{
		Debug: d,
		Value: i,
	}
}

// NodeInteger contains an integer.
type NodeInteger struct {
	Value int
}

// NewInteger returns an initialized NodeInteger structure.
func NewInteger(d *NodeDebugInformation, num int) Node {
	ni := NodeInteger{
		Value: num,
	}

	return Node{
		Debug: d,
		Value: ni,
	}
}

// NodeIdentifier contains a rational number.
type NodeNumber struct {
	Value *big.Rat
}

// NewNumber returns an initialized NodeNumber structure.
func NewNumber(d *NodeDebugInformation, num *big.Rat) Node {
	nu := NodeNumber{
		Value: num,
	}

	return Node{
		Debug: d,
		Value: nu,
	}
}

// NodeIdentifier contains an operand (such as + - ; etc) and its associated
// leaf nodes.
type NodeOperand struct {
	Operand int
	Nodes   []Node
}

// NewNumber returns an initialized NodeOperand structure.
func NewOperand(d *NodeDebugInformation, operand int, args ...Node) Node {
	o := NodeOperand{
		Operand: operand,
		Nodes:   make([]Node, 0, len(args)),
	}

	for _, v := range args {
		o.Nodes = append(o.Nodes, v)
	}

	n := Node{
		Debug: d,
		Value: o,
	}
	return n
}

// Node is the genric container type for all other nodes and is the "currency"
// that is passed around.
type Node struct {
	Value interface{}
	Debug *NodeDebugInformation
}

// DumpPseudoAsm dumps human readable pseudo assembler to w.
func DumpPseudoAsm(n Node, w io.Writer) error {
	a := astResult{}
	a.ec = a.emitPseudoAsm
	err := a.dumpCode(n, w)
	return err
}

// EmitCode dumps a binary image to w.
// This code should be executable by the target architecture.
func EmitCode(n Node, w io.Writer, f func(int, ...interface{}) error) error {
	a := astResult{}
	a.ec = f
	err := a.dumpCode(n, w)
	return err
}

// Clone AST.
// BUG Pointers on debug and leafs are reused for now; that is bad!
func Clone(n Node) Node {
	r := Node{}

	if n.Debug != nil {
		debug := *n.Debug
		r.Debug = &debug
	}

	switch nn := n.Value.(type) {
	case Node:
		r.Value = Clone(nn)
	default:
		r.Value = n.Value
	}
	return r
}

func ExtraDebug(n Node) string {
	if n.Debug == nil {
		return ""
	}

	return fmt.Sprintf(" %v,%v-%v",
		n.Debug.LineNo,
		n.Debug.ColStart,
		n.Debug.ColEnd)
}
