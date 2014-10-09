package ast

import (
	"fmt"
	"io"
	"math/big"
	"strings"
)

type astResult struct {
	code []string
	ec   func(int, ...interface{}) error
	lbl  int
}

func (s *astResult) dumpCode(n Node, w io.Writer) error {
	err := s.dumpCodeR(n)
	if err != nil {
		return err
	}

	// emit exit at the end of the code
	err = s.ec(EXIT)
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

func (s *astResult) emitPseudoAsm(t int, args ...interface{}) error {
	switch t {
	case IDENTIFIER:
		s.addCode("\tpush\t%v\n", args[0].(string))
	case INTEGER:
		s.addCode("\tpush\t%v\n", args[0].(int))
	case NUMBER:
		s.addCode("\tpush\t%v\n", args[0].(*big.Rat))
	case Assign:
		s.addCode("\tpop\t%v\n", args[0].(string))
	case Uminus:
		s.addCode("\tneg\n")
	case Add:
		s.addCode("\tadd\n")
	case Sub:
		s.addCode("\tsub\n")
	case Mul:
		s.addCode("\tmul\n")
	case Div:
		s.addCode("\tdiv\n")
	case Le:
		s.addCode("\tle\n")
	case Ge:
		s.addCode("\tge\n")
	case Lt:
		s.addCode("\tlt\n")
	case Gt:
		s.addCode("\tgt\n")
	case Ne:
		s.addCode("\tne\n")
	case Eq:
		s.addCode("\teq\n")
	case LOCATION:
		s.addCode("l%v:\n", args[0])
	case BRT:
		s.addCode("\tbrt\tl%v\n", args[0])
	case JUMP:
		s.addCode("\tjmp\tl%v\n", args[0])
	case DEBUG:
		s.addCode(args[0].(string), args[1:]...)
	case FIXUP:
		// nothing to do for pseudo asm
	case NOP:
		s.addCode("\tnop\n")
	case EXIT:
		s.addCode("\texit\n")
	default:
		return fmt.Errorf("unsuported pseudo opcode %v", t)
	}

	return nil
}

func (s *astResult) emitDebug(v Node) error {
	line := ""
	lineNo := 0
	if v.Debug != nil {
		lineNo = v.Debug.LineNo
		line = v.Debug.Line
	}
	return s.ec(DEBUG, "\n// line %v: %v\n",
		lineNo,
		strings.Trim(line, " \r\t\n"))
}

func (s *astResult) dumpCodeR(n Node) (err error) {
	switch node := n.Value.(type) {
	case NodeIdentifier:
		err = s.ec(IDENTIFIER, node.Value)
	case NodeInteger:
		err = s.ec(INTEGER, node.Value)
	case NodeNumber:
		err = s.ec(NUMBER, node.Value)
	case NodeOperand:
		switch node.Operand {
		case Assign:
			err = s.dumpCodeR(node.Nodes[1])
			if err != nil {
				return
			}
			err = s.ec(Assign, node.Nodes[0].Value.(NodeIdentifier).Value)

		case Eos:
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
			err = s.ec(Uminus)

		case If:
			if len(node.Nodes) == 2 {
				l0 := s.lbl
				s.lbl++
				l1 := s.lbl
				s.lbl++

				// bool expression
				err = s.dumpCodeR(node.Nodes[0])
				if err != nil {
					return
				}
				err = s.ec(BRT, l0)
				if err != nil {
					return
				}

				// jmp past body
				err = s.ec(JUMP, l1)
				if err != nil {
					return
				}

				// body
				err = s.ec(LOCATION, l0)
				if err != nil {
					return
				}
				err = s.dumpCodeR(node.Nodes[1])
				if err != nil {
					return
				}

				// label past body
				err = s.ec(LOCATION, l1)
				if err != nil {
					return
				}

				// fixup labels that didn't exist
				err = s.ec(FIXUP, l0, l1)
				if err != nil {
					return
				}

			} else {
				l0 := s.lbl
				s.lbl++
				l1 := s.lbl
				s.lbl++
				l2 := s.lbl
				s.lbl++

				// bool expression
				err = s.dumpCodeR(node.Nodes[0])
				if err != nil {
					return
				}
				err = s.ec(BRT, l0)
				if err != nil {
					return
				}

				// jmp past body
				err = s.ec(JUMP, l1)
				if err != nil {
					return
				}

				// body
				err = s.ec(LOCATION, l0)
				if err != nil {
					return
				}
				err = s.dumpCodeR(node.Nodes[1])
				if err != nil {
					return
				}
				// jmp past else body
				err = s.ec(JUMP, l2)
				if err != nil {
					return
				}

				// label past body
				err = s.ec(LOCATION, l1)
				if err != nil {
					return
				}

				err = s.dumpCodeR(node.Nodes[2])
				if err != nil {
					return
				}

				// label past else body
				err = s.ec(LOCATION, l2)
				if err != nil {
					return
				}

				// fixup labels that didn't exist
				err = s.ec(FIXUP, l0, l1, l2)
				if err != nil {
					return
				}

			}
		case While:
			l0 := s.lbl // loop label
			s.lbl++
			l1 := s.lbl // loop body label
			s.lbl++
			l2 := s.lbl // past body label
			s.lbl++

			// boolean check
			err = s.ec(LOCATION, l0)
			if err != nil {
				return
			}
			err = s.dumpCodeR(node.Nodes[0])
			if err != nil {
				return
			}
			err = s.ec(BRT, l1)
			if err != nil {
				return
			}

			// jmp past body
			err = s.ec(JUMP, l2)
			if err != nil {
				return
			}

			// body
			err = s.ec(LOCATION, l1)
			if err != nil {
				return
			}
			err = s.dumpCodeR(node.Nodes[1])
			if err != nil {
				return
			}
			err = s.ec(JUMP, l0)
			if err != nil {
				return
			}
			err = s.ec(LOCATION, l2)
			if err != nil {
				return
			}

			// fixup labels that didn't exist
			err = s.ec(FIXUP, l1, l2)
			if err != nil {
				return
			}

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
			case Add:
				err = s.ec(Add)
			case Sub:
				err = s.ec(Sub)
			case Mul:
				err = s.ec(Mul)
			case Div:
				err = s.ec(Div)
			case Lt:
				err = s.ec(Lt)
			case Gt:
				err = s.ec(Gt)
			case Le:
				err = s.ec(Le)
			case Ge:
				err = s.ec(Ge)
			case Eq:
				err = s.ec(Eq)
			case Ne:
				err = s.ec(Ne)
			default:
				err = fmt.Errorf("unknown operand %v%v",
					node.Operand, ExtraDebug(n))
				return
			}
		}
	default:
		err = fmt.Errorf("unknown node type %T%v",
			node, ExtraDebug(n))
		return
	}

	return
}
