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

// SimpleMathLanguage contains the lexer and parser context.
type SimpleMathLanguage struct {
	lexer *yylexer   // lexer context
	mtx   sync.Mutex // prevent reentrant calls
	src   string     // original source
	code  []string   // generated code
}

// Ensure we are implementing the driver.Frontend interface.
var _ driver.Frontend = &SimpleMathLanguage{}

// New creates a new SimpleMathLanguage context.
func New() (*SimpleMathLanguage, error) {
	return &SimpleMathLanguage{}, nil
}

// Compile lexes and parses src.
// If src compiles it'll return nil.
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

// AST returns the AST representation of the compiled code.
// This is what is subsequently fed into the other layers.
func (s *SimpleMathLanguage) AST() (ast.Node, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// we really should make a copy of the AST
	return s.lexer.tree, nil
}

// Lines returns the original source as an array of strings.
// This is to simplify debugging and enable other human readability tasks.
func (s *SimpleMathLanguage) Lines() ([]string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return nil, fmt.Errorf("Lines not implemented")
}

// Line returns  line l from the original source.
func (s *SimpleMathLanguage) Line(l int) (string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return "", fmt.Errorf("Line not implemented")
}

// yylexer implements the lexer interface.
type yylexer struct {
	src       *bufio.Reader // reader to the code
	buf       []byte        // contains currently lexed bytes
	empty     bool          // indicate if current is valid
	current   byte          // current byte we are lexing
	lastError error         // last error we saw
	line      int           // line we are parsing
	lines     []string      // lines, used for debug etc
	colStart  int           // column where token starts
	colEnd    int           // column where token ends

	tree ast.Node // AST representation of the provided code
}

// newLexer returns a yylexer context.
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

// d generate debug information, short name to keep yacc code readable.
func (y *yylexer) d() *ast.NodeDebugInformation {
	return &ast.NodeDebugInformation{
		LineNo:   y.line,
		ColStart: y.colStart,
		ColEnd:   y.colEnd,
		Line:     y.lines[y.line],
	}
}

// getc returns the next byte from the reader.
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

// Error creates an error structure from a string.
func (y *yylexer) Error(e string) {
	y.lastError = fmt.Errorf("line %v,%v-%v: %v", y.line, y.colStart, y.colEnd, e)
}

// Error creates an error structure using standard formating rules.
func (y *yylexer) Errorf(format string, args ...interface{}) {
	line := fmt.Sprintf("line %v,%v-%v: ", y.line, y.colStart, y.colEnd)
	y.lastError = fmt.Errorf(line+format, args...)
}

// number returns NUMBER and sets the union of the parser to the value of s.
func (y *yylexer) number(val *yySymType, s string) int {
	var ok bool
	val.number, ok = new(big.Rat).SetString(s)
	if !ok {
		log.Fatal("invalid number")
	}
	return NUMBER
}

// number returns IDENTIFIER and sets the union of the parser to the value of s.
func (y *yylexer) identifier(val *yySymType, s string) int {
	val.identifier = string(y.buf)
	return IDENTIFIER
}
