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

{identifier}	return y.identifier(val, string(y.buf))
{number}	return y.number(val, string(y.buf))

%%
		y.empty = true
		return int(c)
}
