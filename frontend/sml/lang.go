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
	integer    int
	number     *big.Rat
	identifier string
	node       ast.Node
}

const PROGRAM = 57346
const INTEGER = 57347
const IDENTIFIER = 57348
const VAR = 57349
const CONST = 57350
const FUNC = 57351
const NUMBER = 57352
const WHILE = 57353
const IF = 57354
const ELSE = 57355
const EOL = 57356
const ASSIGN = 57357
const LE = 57358
const GE = 57359
const NE = 57360
const EQ = 57361
const LT = 57362
const GT = 57363
const UMINUS = 57364

var yyToknames = []string{
	"PROGRAM",
	"INTEGER",
	"IDENTIFIER",
	"VAR",
	"CONST",
	"FUNC",
	"NUMBER",
	"WHILE",
	"IF",
	"ELSE",
	"EOL",
	"ASSIGN",
	"LE",
	"GE",
	"NE",
	"EQ",
	"LT",
	"GT",
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

//line lang.y:123

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 40
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 137

var yyAct = []int{

	20, 19, 16, 41, 70, 32, 33, 34, 35, 49,
	63, 11, 55, 56, 57, 58, 53, 54, 32, 33,
	34, 35, 10, 8, 37, 9, 12, 51, 38, 40,
	42, 42, 44, 7, 62, 45, 46, 47, 48, 36,
	50, 14, 52, 34, 35, 61, 60, 59, 32, 33,
	34, 35, 72, 31, 4, 29, 64, 65, 66, 67,
	68, 69, 22, 24, 6, 28, 1, 23, 27, 28,
	22, 24, 21, 73, 74, 23, 27, 28, 2, 3,
	25, 12, 5, 17, 15, 12, 30, 26, 25, 22,
	39, 71, 15, 12, 23, 26, 32, 33, 34, 35,
	18, 22, 39, 13, 0, 51, 23, 25, 32, 33,
	34, 35, 0, 0, 26, 0, 0, 0, 0, 25,
	0, 0, 0, 0, 0, 0, 43, 55, 56, 57,
	58, 53, 54, 32, 33, 34, 35,
}
var yyPact = []int{

	45, -1000, 45, -1000, 58, -1000, 3, -8, -5, -9,
	-2, -1000, 65, 57, -1000, -1000, 26, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 9, 84, 84, 96, 96, -1000,
	-1000, -1000, 84, 84, 84, 84, -22, 84, -1000, -1000,
	74, -2, 111, 96, -2, 19, 19, -1000, -1000, 7,
	-17, -1000, -1000, 84, 84, 84, 84, 84, 84, -27,
	-4, 39, -1000, -1000, 86, 86, 86, 86, 86, 86,
	-1000, -1000, 53, -1000, -1000,
}
var yyPgo = []int{

	0, 41, 103, 2, 3, 100, 1, 91, 0, 83,
	79, 78, 72, 66,
}
var yyR1 = []int{

	0, 13, 11, 11, 1, 1, 1, 1, 1, 1,
	1, 2, 2, 2, 8, 12, 10, 9, 9, 5,
	6, 7, 7, 7, 4, 4, 4, 4, 4, 4,
	4, 3, 3, 3, 3, 3, 3, 3, 3, 3,
}
var yyR2 = []int{

	0, 1, 1, 2, 1, 2, 1, 1, 1, 1,
	1, 0, 1, 2, 3, 4, 7, 4, 1, 3,
	4, 0, 2, 2, 3, 3, 3, 3, 3, 3,
	3, 1, 1, 1, 2, 3, 3, 3, 3, 3,
}
var yyChk = []int{

	-1000, -13, -11, -10, 9, -10, 6, 30, 31, 30,
	31, -8, 28, -2, -1, 27, -3, -9, -5, -6,
	-8, -12, 5, 10, 6, 23, 30, 11, 12, -1,
	29, 27, 22, 23, 24, 25, 30, 15, -3, 6,
	-3, -4, -3, 30, -4, -3, -3, -3, -3, 31,
	-3, 31, -8, 20, 21, 16, 17, 18, 19, -4,
	-3, -8, 27, 27, -3, -3, -3, -3, -3, -3,
	31, -7, 13, -8, -6,
}
var yyDef = []int{

	0, -2, 1, 2, 0, 3, 0, 0, 0, 0,
	0, 16, 11, 0, 12, 4, 0, 6, 7, 8,
	9, 10, 31, 32, 33, 0, 0, 0, 0, 13,
	14, 5, 0, 0, 0, 0, 0, 0, 34, 33,
	0, 0, 0, 0, 0, 35, 36, 37, 38, 0,
	0, 39, 19, 0, 0, 0, 0, 0, 0, 0,
	0, 21, 15, 17, 24, 25, 26, 27, 28, 29,
	30, 20, 0, 22, 23,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	30, 31, 24, 22, 3, 23, 3, 25, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 27,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 28, 3, 29,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	26,
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
		//line lang.y:49
		{
			d.tree = yyS[yypt-0].node
		}
	case 2:
		//line lang.y:53
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 3:
		//line lang.y:54
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eos, yyS[yypt-1].node, yyS[yypt-0].node)
		}
	case 4:
		//line lang.y:58
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eos)
		}
	case 5:
		//line lang.y:59
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 6:
		//line lang.y:60
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 7:
		//line lang.y:61
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 8:
		//line lang.y:62
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 9:
		//line lang.y:63
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 10:
		//line lang.y:64
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 11:
		//line lang.y:68
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eos)
		}
	case 12:
		//line lang.y:69
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 13:
		//line lang.y:70
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eos, yyS[yypt-1].node, yyS[yypt-0].node)
		}
	case 14:
		//line lang.y:74
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 15:
		//line lang.y:78
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.FunctionCall, ast.NewIdentifier(nil, yyS[yypt-3].identifier))
		}
	case 16:
		//line lang.y:82
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Function, ast.NewIdentifier(nil, yyS[yypt-5].identifier), yyS[yypt-0].node)
		}
	case 17:
		//line lang.y:86
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Assign, ast.NewIdentifier(nil, yyS[yypt-3].identifier), yyS[yypt-1].node)
		}
	case 18:
		//line lang.y:87
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 19:
		//line lang.y:91
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.While, yyS[yypt-1].node, yyS[yypt-0].node)
		}
	case 20:
		//line lang.y:94
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.If, yyS[yypt-2].node, yyS[yypt-1].node, yyS[yypt-0].node)
		}
	case 21:
		//line lang.y:97
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eos)
		}
	case 22:
		//line lang.y:98
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 23:
		//line lang.y:99
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 24:
		//line lang.y:103
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Lt, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 25:
		//line lang.y:104
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Gt, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 26:
		//line lang.y:105
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Le, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 27:
		//line lang.y:106
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Ge, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 28:
		//line lang.y:107
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Ne, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 29:
		//line lang.y:108
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eq, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 30:
		//line lang.y:109
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 31:
		//line lang.y:113
		{
			yyVAL.node = ast.NewInteger(d.d(), yyS[yypt-0].integer)
		}
	case 32:
		//line lang.y:114
		{
			yyVAL.node = ast.NewNumber(d.d(), yyS[yypt-0].number)
		}
	case 33:
		//line lang.y:115
		{
			yyVAL.node = ast.NewIdentifier(d.d(), yyS[yypt-0].identifier)
		}
	case 34:
		//line lang.y:116
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Uminus, yyS[yypt-0].node)
		}
	case 35:
		//line lang.y:117
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Add, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 36:
		//line lang.y:118
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Sub, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 37:
		//line lang.y:119
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Mul, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 38:
		//line lang.y:120
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Div, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 39:
		//line lang.y:121
		{
			yyVAL.node = yyS[yypt-1].node
		}
	}
	goto yystack /* stack new state and value */
}
