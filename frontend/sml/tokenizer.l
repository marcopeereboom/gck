%{
package sml

func (y *yylexer) Lex(val *yySymType) int {
	c := y.current
	if y.empty {
		c, y.empty = y.getc(), false
	}
%}

%yyc		c
%yyn		c = y.getc()

letter		[A-Za-z]
digit		[0-9]+
integer		{digit}
identifier	{letter}({letter}|{digit}|_)*
exponent	[Ee][-+]?{digit}
number		{digit}"."{digit}?{exponent}?|{digit}{exponent}?|"."{digit}{exponent}?

%%
		y.buf = y.buf[:0]
		y.colStart = y.colEnd

[ \t]+

[\n\r]		y.colStart = 1; y.colEnd = 1; y.line++;

"program"	return PROGRAM
"var"		return VAR
"const"		return CONST
"func"		return FUNC
"while"		return WHILE
"<"		return LT
">"		return GT
"<="		return LE
">="		return GE
"!="		return NE
"=="		return EQ
"="		return ASSIGN

{identifier}	return y.identifier(val, string(y.buf))
{integer}	return y.integer(val, string(y.buf))
{number}	return y.number(val, string(y.buf))

%%
		y.empty = true
		return int(c)
}
