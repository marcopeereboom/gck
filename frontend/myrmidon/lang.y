%{

package myrmidon

import (
	"github.com/marcopeereboom/gck/ast"
	"math/big"
)

var d *yylexer // being set so we don't have to type assert all the time

%}

%union{
	integer	   int
	number     *big.Rat
	identifier string
	node       ast.Node
}

%token	PROGRAM
%token	INTEGER
%token	IDENTIFIER
%token	VAR
%token	CONST
%token	FUNC
%token	NUMBER
%token	WHILE
%token	IF
%token	ELSE
%token	EOL
%token	ASSIGN

%type	<identifier>	IDENTIFIER
%type	<integer>	INTEGER
%type	<number>	NUMBER
%type	<node>		statement statementlist expression boolexpression
%type	<node>		while if else closedstatements identifier function
%type	<node>		functionlist functioncall

%left		LE GE NE EQ LT GT
%left		'+' '-'
%left		'*' '/'
%nonassoc	UMINUS

%%

program:
           functionlist	{ d.tree = $1 }
        ;

functionlist:
	  function		{ $$ = $1 }
	| functionlist function	{ $$ = ast.NewOperand(d.d(), ast.Eos, $1, $2) }
	;

statement:
	  ';'			{ $$ = ast.NewOperand(d.d(), ast.Eos) }
	| expression ';'	{ $$ = $1 }
	| identifier		{ $$ = $1 }
	| while			{ $$ = $1 }
	| if			{ $$ = $1 }
	| closedstatements	{ $$ = $1 }
	| functioncall		{ $$ = $1 }
	;

statementlist:
					{ $$ = ast.NewOperand(d.d(), ast.Eos) }
	| statement			{ $$ = $1 }
	| statementlist statement	{ $$ = ast.NewOperand(d.d(), ast.Eos, $1, $2) }
	;

closedstatements:
	  '{' statementlist '}'	{ $$ = $2 }
	;

functioncall:
	  IDENTIFIER '(' ')' ';'	{ $$ = ast.NewOperand(d.d(), ast.FunctionCall, ast.NewIdentifier(nil, $1)) }
	;

function:
	  FUNC IDENTIFIER '(' ')' '(' ')' closedstatements	{ $$ = ast.NewOperand(d.d(), ast.Function, ast.NewIdentifier(nil, $2), $7) }
	;

identifier:
	  IDENTIFIER ASSIGN expression ';'	{ $$ = ast.NewOperand(d.d(), ast.Assign, ast.NewIdentifier(nil, $1), $3) }
	| functioncall				{ $$ = $1 }
	;

while:
	  WHILE boolexpression closedstatements { $$ = ast.NewOperand(d.d(), ast.While, $2, $3) }
	;
if:
	  IF boolexpression closedstatements else { $$ = ast.NewOperand(d.d(), ast.If, $2, $3, $4) }
	;

else:					{ $$ = ast.NewOperand(d.d(), ast.Eos) }
	| ELSE closedstatements		{ $$ = $2 }
	| ELSE if			{ $$ = $2 }
	;

boolexpression:
	  expression LT expression	{ $$ = ast.NewOperand(d.d(), ast.Lt, $1, $3) }
	| expression GT expression	{ $$ = ast.NewOperand(d.d(), ast.Gt, $1, $3) }
	| expression LE expression	{ $$ = ast.NewOperand(d.d(), ast.Le, $1, $3) }
	| expression GE expression	{ $$ = ast.NewOperand(d.d(), ast.Ge, $1, $3) }
	| expression NE expression	{ $$ = ast.NewOperand(d.d(), ast.Ne, $1, $3) }
	| expression EQ expression	{ $$ = ast.NewOperand(d.d(), ast.Eq, $1, $3) }
	| '(' boolexpression ')'	{ $$ = $2 }
	;

expression:
	  INTEGER			{ $$ = ast.NewInteger(d.d(), $1) }
	| NUMBER			{ $$ = ast.NewNumber(d.d(), $1) }
	| IDENTIFIER			{ $$ = ast.NewIdentifier(d.d(), $1) }
	| '-' expression %prec UMINUS	{ $$ = ast.NewOperand(d.d(), ast.Uminus, $2) }
	| expression '+' expression	{ $$ = ast.NewOperand(d.d(), ast.Add, $1, $3) }
	| expression '-' expression	{ $$ = ast.NewOperand(d.d(), ast.Sub, $1, $3) }
	| expression '*' expression	{ $$ = ast.NewOperand(d.d(), ast.Mul, $1, $3) }
	| expression '/' expression	{ $$ = ast.NewOperand(d.d(), ast.Div, $1, $3) }
	| '(' expression ')'		{ $$ = $2 }
	;
%%
