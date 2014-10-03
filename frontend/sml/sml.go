// sml is the "simple math language" interface, lexer and parser.
package sml

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/marcopeereboom/gck/ast"
	"github.com/marcopeereboom/gck/frontend/driver"
)

type SimpleMathLanguage struct {
	lexer *yylexer   // lexer context
	mtx   sync.Mutex // prevent reentrant calls
	src   string     // original source
	code  []string   // generated code
}

var _ driver.Frontend = &SimpleMathLanguage{}

// utility
func New() (*SimpleMathLanguage, error) {
	return &SimpleMathLanguage{}, nil
}

// interface
func (s *SimpleMathLanguage) Compile(src string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// slice up src
	var err error
	s.src = src
	lines, err := driver.LineGenerator(src)
	if err != nil {
		return err
	}

	// compile code
	r := bufio.NewReader(strings.NewReader(s.src))
	s.lexer = newLexer(r)
	s.lexer.lines = lines
	result := yyParse(s.lexer)
	if result == 0 {
		return nil
	}

	return nil
}

func (s *SimpleMathLanguage) AST() (ast.Node, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// we really should make a copy of the AST
	return s.lexer.tree, nil
}

func (s *SimpleMathLanguage) Lines() ([]string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return nil, fmt.Errorf("Lines not implemented")
}

func (s *SimpleMathLanguage) Line(int) (string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return "", fmt.Errorf("Line not implemented")
}

// implement lexer
type yylexer struct {
	src       *bufio.Reader
	buf       []byte
	empty     bool
	current   byte
	lastError error
	line      int      // line we are parsing
	lines     []string // lines, used for debug etc
	colStart  int      // column where token starts
	colEnd    int      // column where token ends

	tree ast.Node // AST
}

func newLexer(src *bufio.Reader) *yylexer {
	y := yylexer{
		line: 1,
		src:  src,
	}

	d = &y // hack around having to type asser yylex.(*yyLexer)

	if b, err := src.ReadByte(); err == nil {
		y.current = b
		y.colEnd++
	}

	return &y
}

// generate debug information, short name to keep yacc code readable.
func (y *yylexer) d() *ast.NodeDebugInformation {
	return &ast.NodeDebugInformation{
		LineNo:   y.line,
		ColStart: y.colStart,
		ColEnd:   y.colEnd,
		Line:     y.lines[y.line],
	}
}

func (y *yylexer) getc() byte {
	if y.current != 0 {
		y.buf = append(y.buf, y.current)
	}
	y.current = 0
	if b, err := y.src.ReadByte(); err == nil {
		y.current = b
		y.colEnd++
	}
	return y.current
}

func (y *yylexer) Error(e string) {
	y.lastError = fmt.Errorf("line %v,%v-%v: %v", y.line, y.colStart, y.colEnd, e)
}

func (y *yylexer) Errorf(format string, args ...interface{}) {
	line := fmt.Sprintf("line %v,%v-%v: ", y.line, y.colStart, y.colEnd)
	y.lastError = fmt.Errorf(line+format, args...)
}

func (y *yylexer) number(val *yySymType, s string) int {
	var ok bool
	val.number, ok = new(big.Rat).SetString(s)
	if !ok {
		log.Fatal("invalid number")
	}
	return NUMBER
}

func (y *yylexer) identifier(val *yySymType, s string) int {
	val.identifier = string(y.buf)
	return IDENTIFIER
}
