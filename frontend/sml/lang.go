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

//line lang.y:85

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 26
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 148

var yyAct = []int{

	3, 55, 2, 14, 16, 17, 18, 19, 40, 33,
	20, 5, 25, 42, 18, 19, 10, 6, 54, 22,
	22, 11, 7, 8, 26, 27, 14, 1, 28, 29,
	30, 31, 32, 21, 12, 0, 44, 0, 4, 9,
	57, 13, 24, 51, 0, 14, 45, 46, 47, 48,
	49, 50, 14, 0, 10, 6, 0, 14, 56, 11,
	7, 8, 10, 6, 0, 10, 23, 11, 7, 8,
	11, 0, 12, 0, 0, 0, 4, 9, 53, 13,
	12, 0, 0, 12, 4, 9, 52, 13, 10, 6,
	13, 0, 0, 11, 7, 8, 10, 6, 0, 0,
	0, 11, 7, 8, 0, 0, 12, 0, 0, 0,
	4, 9, 41, 13, 12, 0, 0, 0, 4, 9,
	0, 13, 36, 37, 38, 39, 34, 35, 16, 17,
	18, 19, 16, 17, 18, 19, 0, 43, 16, 17,
	18, 19, 0, 15, 16, 17, 18, 19,
}
var yyPact = []int{

	91, -1000, 91, -1000, -1000, 116, -5, 60, 60, 91,
	-1000, -1000, 60, 60, -1000, -1000, 60, 60, 60, 60,
	60, -19, 106, -1000, -20, 83, -1000, -18, -10, -10,
	-1000, -1000, 110, 91, 60, 60, 60, 60, 60, 60,
	91, -1000, -1000, -1000, 57, 122, 122, 122, 122, 122,
	122, 49, -1000, 5, -27, 91, 11, -1000,
}
var yyPgo = []int{

	0, 0, 2, 11, 33, 27,
}
var yyR1 = []int{

	0, 5, 1, 1, 1, 1, 1, 1, 1, 2,
	2, 4, 4, 4, 4, 4, 4, 3, 3, 3,
	3, 3, 3, 3, 3, 3,
}
var yyR2 = []int{

	0, 1, 1, 2, 4, 5, 5, 9, 3, 1,
	2, 3, 3, 3, 3, 3, 3, 1, 1, 1,
	2, 3, 3, 3, 3, 3,
}
var yyChk = []int{

	-1000, -5, -2, -1, 27, -3, 6, 11, 12, 28,
	5, 10, 23, 30, -1, 27, 22, 23, 24, 25,
	15, -4, -3, 6, -4, -2, -3, -3, -3, -3,
	-3, -3, -3, 28, 20, 21, 16, 17, 18, 19,
	28, 29, 31, 27, -2, -3, -3, -3, -3, -3,
	-3, -2, 29, 29, 13, 28, -2, 29,
}
var yyDef = []int{

	0, -2, 1, 9, 2, 0, 19, 0, 0, 0,
	17, 18, 0, 0, 10, 3, 0, 0, 0, 0,
	0, 0, 0, 19, 0, 0, 20, 0, 21, 22,
	23, 24, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 8, 25, 4, 0, 11, 12, 13, 14, 15,
	16, 0, 5, 6, 0, 0, 0, 7,
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
		//line lang.y:47
		{
			d.tree = yyS[yypt-0].node
		}
	case 2:
		//line lang.y:51
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eos)
		}
	case 3:
		//line lang.y:52
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 4:
		//line lang.y:53
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Assign, ast.NewIdentifier(nil, yyS[yypt-3].identifier), yyS[yypt-1].node)
		}
	case 5:
		//line lang.y:54
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.While, yyS[yypt-3].node, yyS[yypt-1].node)
		}
	case 6:
		//line lang.y:55
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.If, yyS[yypt-3].node, yyS[yypt-1].node)
		}
	case 7:
		//line lang.y:56
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.If, yyS[yypt-7].node, yyS[yypt-5].node, yyS[yypt-1].node)
		}
	case 8:
		//line lang.y:57
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 9:
		//line lang.y:61
		{
			yyVAL.node = yyS[yypt-0].node
		}
	case 10:
		//line lang.y:62
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eos, yyS[yypt-1].node, yyS[yypt-0].node)
		}
	case 11:
		//line lang.y:66
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Lt, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 12:
		//line lang.y:67
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Gt, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 13:
		//line lang.y:68
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Le, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 14:
		//line lang.y:69
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Ge, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 15:
		//line lang.y:70
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Ne, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 16:
		//line lang.y:71
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Eq, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 17:
		//line lang.y:75
		{
			yyVAL.node = ast.NewInteger(d.d(), yyS[yypt-0].integer)
		}
	case 18:
		//line lang.y:76
		{
			yyVAL.node = ast.NewNumber(d.d(), yyS[yypt-0].number)
		}
	case 19:
		//line lang.y:77
		{
			yyVAL.node = ast.NewIdentifier(d.d(), yyS[yypt-0].identifier)
		}
	case 20:
		//line lang.y:78
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Uminus, yyS[yypt-0].node)
		}
	case 21:
		//line lang.y:79
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Add, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 22:
		//line lang.y:80
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Sub, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 23:
		//line lang.y:81
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Mul, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 24:
		//line lang.y:82
		{
			yyVAL.node = ast.NewOperand(d.d(), ast.Div, yyS[yypt-2].node, yyS[yypt-0].node)
		}
	case 25:
		//line lang.y:83
		{
			yyVAL.node = yyS[yypt-1].node
		}
	}
	goto yystack /* stack new state and value */
}
