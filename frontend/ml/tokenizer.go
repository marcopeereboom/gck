package sml

func (y *yylexer) Lex(val *yySymType) int {
	c := y.current
	if y.empty {
		c, y.empty = y.getc(), false
	}

yystate0:

	y.buf = y.buf[:0]
	y.colStart = y.colEnd

	goto yystart1

	goto yystate1 // silence unused label error
yystate1:
	c = y.getc()
yystart1:
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate4
	case c == '\n' || c == '\r':
		goto yystate3
	case c == '\t' || c == ' ':
		goto yystate2
	case c == 'c':
		goto yystate11
	case c == 'f':
		goto yystate16
	case c == 'p':
		goto yystate20
	case c == 'v':
		goto yystate27
	case c >= '0' && c <= '9':
		goto yystate9
	case c >= 'A' && c <= 'Z' || c == 'a' || c == 'b' || c == 'd' || c == 'e' || c >= 'g' && c <= 'o' || c >= 'q' && c <= 'u' || c >= 'w' && c <= 'z':
		goto yystate10
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule1
	case c == '\t' || c == ' ':
		goto yystate2
	}

yystate3:
	c = y.getc()
	goto yyrule2

yystate4:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate5
	}

yystate5:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c == 'E' || c == 'e':
		goto yystate6
	case c >= '0' && c <= '9':
		goto yystate5
	}

yystate6:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '+' || c == '-':
		goto yystate7
	case c >= '0' && c <= '9':
		goto yystate8
	}

yystate7:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate8
	}

yystate8:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= '9':
		goto yystate8
	}

yystate9:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c == '.':
		goto yystate5
	case c == 'E' || c == 'e':
		goto yystate6
	case c >= '0' && c <= '9':
		goto yystate9
	}

yystate10:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate11:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'o':
		goto yystate12
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate10
	}

yystate12:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'n':
		goto yystate13
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate10
	}

yystate13:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 's':
		goto yystate14
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate10
	}

yystate14:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 't':
		goto yystate15
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate10
	}

yystate15:
	c = y.getc()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate16:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'u':
		goto yystate17
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate10
	}

yystate17:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'n':
		goto yystate18
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate10
	}

yystate18:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'c':
		goto yystate19
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate10
	}

yystate19:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate20:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'r':
		goto yystate21
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate10
	}

yystate21:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'o':
		goto yystate22
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate10
	}

yystate22:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'g':
		goto yystate23
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate10
	}

yystate23:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'r':
		goto yystate24
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate10
	}

yystate24:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'a':
		goto yystate25
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate10
	}

yystate25:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'm':
		goto yystate26
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate10
	}

yystate26:
	c = y.getc()
	switch {
	default:
		goto yyrule3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yystate27:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'a':
		goto yystate28
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate10
	}

yystate28:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c == 'r':
		goto yystate29
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate10
	}

yystate29:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate10
	}

yyrule1: // [ \t]+

	goto yystate0
yyrule2: // [\n\r]
	{
		y.colStart = 1
		y.colEnd = 1
		y.line++
		goto yystate0
	}
yyrule3: // "program"
	{
		return PROGRAM
	}
yyrule4: // "var"
	{
		return VAR
	}
yyrule5: // "const"
	{
		return CONST
	}
yyrule6: // "func"
	{
		return FUNC
	}
yyrule7: // {identifier}
	{
		return y.identifier(val, string(y.buf))
	}
yyrule8: // {number}
	{
		return y.number(val, string(y.buf))
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	y.empty = true
	return int(c)
}
