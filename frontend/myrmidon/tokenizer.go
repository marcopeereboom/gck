package myrmidon

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
	case c == '!':
		goto yystate4
	case c == '.':
		goto yystate6
	case c == '<':
		goto yystate12
	case c == '=':
		goto yystate14
	case c == '>':
		goto yystate16
	case c == '\n' || c == '\r':
		goto yystate3
	case c == '\t' || c == ' ':
		goto yystate2
	case c == 'c':
		goto yystate19
	case c == 'e':
		goto yystate24
	case c == 'f':
		goto yystate28
	case c == 'i':
		goto yystate32
	case c == 'p':
		goto yystate34
	case c == 'v':
		goto yystate41
	case c == 'w':
		goto yystate44
	case c >= '0' && c <= '9':
		goto yystate11
	case c >= 'A' && c <= 'Z' || c == 'a' || c == 'b' || c == 'd' || c == 'g' || c == 'h' || c >= 'j' && c <= 'o' || c >= 'q' && c <= 'u' || c >= 'x' && c <= 'z':
		goto yystate18
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
	case c == '=':
		goto yystate5
	}

yystate5:
	c = y.getc()
	goto yyrule14

yystate6:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate7
	}

yystate7:
	c = y.getc()
	switch {
	default:
		goto yyrule19
	case c == 'E' || c == 'e':
		goto yystate8
	case c >= '0' && c <= '9':
		goto yystate7
	}

yystate8:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '+' || c == '-':
		goto yystate9
	case c >= '0' && c <= '9':
		goto yystate10
	}

yystate9:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate10
	}

yystate10:
	c = y.getc()
	switch {
	default:
		goto yyrule19
	case c >= '0' && c <= '9':
		goto yystate10
	}

yystate11:
	c = y.getc()
	switch {
	default:
		goto yyrule18
	case c == '.':
		goto yystate7
	case c == 'E' || c == 'e':
		goto yystate8
	case c >= '0' && c <= '9':
		goto yystate11
	}

yystate12:
	c = y.getc()
	switch {
	default:
		goto yyrule10
	case c == '=':
		goto yystate13
	}

yystate13:
	c = y.getc()
	goto yyrule12

yystate14:
	c = y.getc()
	switch {
	default:
		goto yyrule16
	case c == '=':
		goto yystate15
	}

yystate15:
	c = y.getc()
	goto yyrule15

yystate16:
	c = y.getc()
	switch {
	default:
		goto yyrule11
	case c == '=':
		goto yystate17
	}

yystate17:
	c = y.getc()
	goto yyrule13

yystate18:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate18
	}

yystate19:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'o':
		goto yystate20
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate18
	}

yystate20:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'n':
		goto yystate21
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate18
	}

yystate21:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 's':
		goto yystate22
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate18
	}

yystate22:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 't':
		goto yystate23
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate18
	}

yystate23:
	c = y.getc()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate18
	}

yystate24:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'l':
		goto yystate25
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate18
	}

yystate25:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 's':
		goto yystate26
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate18
	}

yystate26:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'e':
		goto yystate27
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate18
	}

yystate27:
	c = y.getc()
	switch {
	default:
		goto yyrule9
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate18
	}

yystate28:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'u':
		goto yystate29
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate18
	}

yystate29:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'n':
		goto yystate30
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate18
	}

yystate30:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'c':
		goto yystate31
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate18
	}

yystate31:
	c = y.getc()
	switch {
	default:
		goto yyrule6
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate18
	}

yystate32:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'f':
		goto yystate33
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate18
	}

yystate33:
	c = y.getc()
	switch {
	default:
		goto yyrule8
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate18
	}

yystate34:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'r':
		goto yystate35
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate18
	}

yystate35:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'o':
		goto yystate36
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate18
	}

yystate36:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'g':
		goto yystate37
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate18
	}

yystate37:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'r':
		goto yystate38
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate18
	}

yystate38:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'a':
		goto yystate39
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate18
	}

yystate39:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'm':
		goto yystate40
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'l' || c >= 'n' && c <= 'z':
		goto yystate18
	}

yystate40:
	c = y.getc()
	switch {
	default:
		goto yyrule3
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate18
	}

yystate41:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'a':
		goto yystate42
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate18
	}

yystate42:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'r':
		goto yystate43
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate18
	}

yystate43:
	c = y.getc()
	switch {
	default:
		goto yyrule4
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate18
	}

yystate44:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'h':
		goto yystate45
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate18
	}

yystate45:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'i':
		goto yystate46
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate18
	}

yystate46:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'l':
		goto yystate47
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate18
	}

yystate47:
	c = y.getc()
	switch {
	default:
		goto yyrule17
	case c == 'e':
		goto yystate48
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate18
	}

yystate48:
	c = y.getc()
	switch {
	default:
		goto yyrule7
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate18
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
yyrule7: // "while"
	{
		return WHILE
	}
yyrule8: // "if"
	{
		return IF
	}
yyrule9: // "else"
	{
		return ELSE
	}
yyrule10: // "<"
	{
		return LT
	}
yyrule11: // ">"
	{
		return GT
	}
yyrule12: // "<="
	{
		return LE
	}
yyrule13: // ">="
	{
		return GE
	}
yyrule14: // "!="
	{
		return NE
	}
yyrule15: // "=="
	{
		return EQ
	}
yyrule16: // "="
	{
		return ASSIGN
	}
yyrule17: // {identifier}
	{
		return y.identifier(val, string(y.buf))
	}
yyrule18: // {integer}
	{
		return y.integer(val, string(y.buf))
	}
yyrule19: // {number}
	{
		return y.number(val, string(y.buf))
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	y.empty = true
	return int(c)
}
