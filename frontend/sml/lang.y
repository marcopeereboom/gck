%{

package sml

import (
	"github.com/marcopeereboom/gck/ast"
	"math/big"
)

var d *yylexer // being set so we don't have to type assert all the time

%}

%union{
	number     *big.Rat
	identifier string
	node       ast.Node
}

%token	PROGRAM
%token	IDENTIFIER
%token	VAR
%token	CONST
%token	FUNC
%token	NUMBER
%token	WHILE
%token	EOL

%type	<identifier>	IDENTIFIER
%type	<number>	NUMBER
%type	<node>		statement statementlist expression

%left		'+' '-'
%left		'*' '/'
%nonassoc	UMINUS

%%

program:
           statementlist		{ d.tree = $1 }
        ;

statement:
	  ';'				{ $$ = ast.NewOperand(d.d(), ';') }
	| expression ';'		{ $$ = $1 }
	| IDENTIFIER '=' expression ';'	{ $$ = ast.NewOperand(d.d(), '=', ast.NewIdentifier(nil, $1), $3) }
	| '{' statementlist '}'		{ $$ = $2 }
	;

statementlist:
	  statement			{ $$ = $1 }
	| statementlist statement	{ $$ = ast.NewOperand(d.d(), ';', $1, $2) }
	;

expression:
	  NUMBER			{ $$ = ast.NewNumber(d.d(), $1) }
	| IDENTIFIER			{ $$ = ast.NewIdentifier(d.d(), $1) }
	| '-' expression %prec UMINUS	{ $$ = ast.NewOperand(d.d(), ast.Uminus, $2) }
	| expression '+' expression	{ $$ = ast.NewOperand(d.d(), '+', $1, $3) }
	| expression '-' expression	{ $$ = ast.NewOperand(d.d(), '-', $1, $3) }
	| expression '*' expression	{ $$ = ast.NewOperand(d.d(), '*', $1, $3) }
	| expression '/' expression	{ $$ = ast.NewOperand(d.d(), '/', $1, $3) }
	| '(' expression ')'		{ $$ = $2 }
	;
%%
