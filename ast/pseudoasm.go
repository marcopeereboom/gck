package ast

import (
	"fmt"
	"io"
	"math/big"
	"strings"
)

type astResult struct {
	code []string
	ec   func(int, ...interface{})
}

func (s *astResult) dumpCode(n Node, w io.Writer) error {
	err := s.dumpCodeR(n)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "// intermediary language dump\n")
	for k, v := range s.code {
		// this works around printing stray // lines
		// ugly as hell but simpler than screwing with recursive
		// printing for the moment
		if k < len(s.code)-2 &&
			strings.HasPrefix(v, "\n//") &&
			strings.HasPrefix(s.code[k+1], "\n//") {
			continue
		}
		fmt.Fprintf(w, "%v", v)
	}

	return nil
}

func (s *astResult) addCode(f string, args ...interface{}) {
	s.code = append(s.code, fmt.Sprintf(f, args...))
}

const (
	IDENTIFIER = 1
	NUMBER     = 2
	UMINUS     = 3
	DEBUG      = 4
)

func (s *astResult) emitPseudoAsm(t int, args ...interface{}) {
	switch t {
	case IDENTIFIER:
		s.addCode("\tpush\t%v\n", args[0].(Node).Value)
	case NUMBER:
		s.addCode("\tpush\t%v\n", args[0].(*big.Rat))
	case '=':
		s.addCode("\tpop\t%v\n", args[0].(string))
	case UMINUS:
		s.addCode("\tneg\n")
	case '+':
		s.addCode("\tadd\n")
	case '-':
		s.addCode("\tsub\n")
	case '*':
		s.addCode("\tmul\n")
	case '/':
		s.addCode("\tdiv\n")
	case DEBUG:
		s.addCode(args[0].(string), args[1:]...)
	default:
		panic(fmt.Sprintf("unsuported pseudo opcode %v", t))
	}
}

func (s *astResult) emitDebug(v Node) {
	line := ""
	lineNo := 0
	if v.Debug != nil {
		lineNo = v.Debug.LineNo
		line = v.Debug.Line
	}
	s.ec(DEBUG, "\n// line %v: %v\n",
		lineNo,
		strings.Trim(line, " \r\t\n"))
}

func (s *astResult) dumpCodeR(n Node) (err error) {
	switch node := n.Value.(type) {
	case NodeIdentifier:
		s.ec(IDENTIFIER, node.Value)
	case NodeNumber:
		s.ec(NUMBER, node.Value)
	case NodeOperand:
		switch node.Operand {
		case '=':
			err = s.dumpCodeR(node.Nodes[1])
			if err != nil {
				return
			}
			s.ec('=', node.Nodes[0].Value.(NodeIdentifier).Value)
		case ';':
			for _, v := range node.Nodes {
				// the first few of those that are emitted
				// should be ignored; would be nice to fix
				s.emitDebug(v)

				// walk all nodes
				err = s.dumpCodeR(v)
				if err != nil {
					return
				}
			}
		case Uminus:
			err = s.dumpCodeR(node.Nodes[0])
			if err != nil {
				return
			}
			s.ec(UMINUS)
		default:
			err = s.dumpCodeR(node.Nodes[0])
			if err != nil {
				return
			}
			err = s.dumpCodeR(node.Nodes[1])
			if err != nil {
				return
			}
			switch node.Operand {
			case '+':
				s.ec('+')
			case '-':
				s.ec('-')
			case '*':
				s.ec('*')
			case '/':
				s.ec('/')
			default:
				err = fmt.Errorf("unknown operand %v\n",
					node.Operand)
				return
			}
		}
	default:
		err = fmt.Errorf("unknown node type %T\n", n.Value)
		return
	}

	return
}
