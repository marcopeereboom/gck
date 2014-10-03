//line lang.y:2
package sml

import __yyfmt__ "fmt"

//line lang.y:3
import (
	"github.com/marcopeereboom/gck/ast"
	"math/big"
)

var d *yylexer // being set so we don't have to type assert all the time

//line lang.y:14
type yySymType struct {
	yys        int
	number     *big.Rat
	identifier string
	node       ast.Node
}

const PROGRAM = 57346
const IDENTIFIER = 57347
const VAR = 57348
const CONST = 57349
const FUNC = 57350
const NUMBER = 57351
const EOL = 57352
const UMINUS = 57353

var yyToknames = []string{
	"PROGRAM",
	"IDENTIFIER",
	"VAR",
	"CONST",
	"FUNC",
	"NUMBER",
	"EOL",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"UMINUS",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line lang.y:64

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 16
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 59

var yyAct = []int{

	6, 17, 3, 1, 8, 11, 0, 9, 0, 15,
	16, 4, 0, 7, 27, 10, 13, 14, 15, 16,
	6, 11, 2, 0, 8, 20, 28, 9, 5, 8,
	18, 4, 9, 7, 0, 10, 0, 0, 19, 21,
	10, 0, 22, 23, 24, 25, 26, 13, 14, 15,
	16, 0, 29, 13, 14, 15, 16, 0, 12,
}
var yyPact = []int{

	15, -1000, 15, -1000, -1000, 42, -16, 15, -1000, 20,
	20, -1000, -1000, 20, 20, 20, 20, 20, -5, -1000,
	-1000, 5, -4, -4, -1000, -1000, 36, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 2, 22, 28, 3,
}
var yyR1 = []int{

	0, 4, 1, 1, 1, 1, 2, 2, 3, 3,
	3, 3, 3, 3, 3, 3,
}
var yyR2 = []int{

	0, 1, 1, 2, 4, 3, 1, 2, 1, 1,
	2, 3, 3, 3, 3, 3,
}
var yyChk = []int{

	-1000, -4, -2, -1, 16, -3, 5, 18, 9, 12,
	20, -1, 16, 11, 12, 13, 14, 17, -2, -3,
	5, -3, -3, -3, -3, -3, -3, 19, 21, 16,
}
var yyDef = []int{

	0, -2, 1, 6, 2, 0, 9, 0, 8, 0,
	0, 7, 3, 0, 0, 0, 0, 0, 0, 10,
	9, 0, 11, 12, 13, 14, 0, 5, 15, 4,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	20, 21, 13, 11, 3, 12, 3, 14, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 16,
	3, 17, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 18, 3, 19,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 15,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line lang.y:39
		{
			d.tree = yyS[yypt-0].node
		}
	case 2:
		//line lang.y:43
		{
			yyVAL.node = ast.NewOperand(d.d(), ';')
		}
	case 3:
		//line lang.y:44
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 4:
		//line lang.y:45
		{
			yyVAL.node = ast.NewOperand(d.d(), '=', ast.NewIdentifier(nil, yyS[yypt-3].identifier), yyS[yypt-1].node)
		}
	case 5:
		//line lang.y:46
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 6:
		//line lang.y:50
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 7:
		//line lang.y:51
		{
			yyVAL.node = ast.NewOperand(d.d(), ';', yyS[yypt-1].node, yyS[yypt-0].node)
		}
	case 8:
		//line lang.y:55
		{
			yyVAL.node = ast.NewNumber(d.d(), yyS[yypt-0].number)
		}
	case 9:
		//line lang.y:56
		{
			yyVAL.node = ast.NewIdentifier(d.d(), yyS[yypt-0].identifier)
		}
	case 10:
		//line lang.y:57
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Uminus, yyS[yypt-0].node)
		}
	case 11:
		//line lang.y:58
		{
			yyVAL.node = ast.NewOperand(d.d(), '+', yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 12:
		//line lang.y:59
		{
			yyVAL.node = ast.NewOperand(d.d(), '-', yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 13:
		//line lang.y:60
		{
			yyVAL.node = ast.NewOperand(d.d(), '*', yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 14:
		//line lang.y:61
		{
			yyVAL.node = ast.NewOperand(d.d(), '/', yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 15:
		//line lang.y:62
		{
			yyVAL.node = yyS[yypt-1].node
		}
	}
	goto yystack /* stack new state and value */
}
