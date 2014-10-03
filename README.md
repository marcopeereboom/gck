Go Compiler Kit
===============

The point of GCK is to implement a working version of all major components that
are required to create, validate and run a scripting language.
If done correctly one should be able to experiment with a single, or several
custom components.

The secondary goal is have a bytecode vm that can run a custom language inside
go programs.

At this time performance is not considered and end goal.
Everything must work first, be correct, secure and modular.

The major pieces of such a system are:
	Frontend - Language to AST translator, components are lexer and parser
	AST - The Abstract Syntax Tree is in this system nothing but plumbing between subsystems
	Optimizer -  Transforms an AST into a better AST
	Backend - Transforms an AST into bytecode
	Virtual Machine - takes the bytecode and executes it

Currently the front and backend are modular and a new language or a new target
can simply be added by adding a driver that adheres to the correct interface.

The provided front and backend examples are:
	ml - Math Language a simple language that does mostly math type operations
	tvm - Toy Virtual Machine is a simple VM that runs tvm binaries

Example use:
First compile the code
```
c -i examples/e1.sml -o /tmp/image.bin
tvm -i /tmp/image.bin -t
```

should result in something like this:
```
=== run trace  ===
0000000000000000: push   c1000 (12/1)
0000000000000002: push   c1001 (13/1)
0000000000000004: push   c1002 (14/1)
0000000000000006: push   c1003 (15/1)
0000000000000008: add   
0000000000000009: mul   
000000000000000a: add   
000000000000000b: pop    a (0/1)
000000000000000d: push   c1005 (12/1)
000000000000000f: push   c1006 (13/1)
0000000000000011: push   c1007 (14/1)
0000000000000013: push   c1008 (15/1)
0000000000000015: neg   
0000000000000016: add   
0000000000000017: mul   
0000000000000018: add   
0000000000000019: pop    b (0/1)
=== cmd stack  ===
=== call stack ===
=== symbols    ===
.CONST   NUMBER   c1008             15/1
.VAR     NUMBER   2d1fe47e779f5f40  -1/1
.VAR     NUMBER   a                 389/1
.CONST   NUMBER   c1000             12/1
.CONST   NUMBER   c1005             12/1
.VAR     NUMBER   5f7b67ffec2b8427  29/1
.VAR     NUMBER   b5b80166f2eb98fb  -15/1
.VAR     NUMBER   6633c6de7e5dfd69  -1/1
.VAR     NUMBER   b                 -1/1
.CONST   NUMBER   c1001             13/1
.VAR     NUMBER   ea4dfe3c6fc397cd  -13/1
.CONST   NUMBER   c1002             14/1
.CONST   NUMBER   c1006             13/1
.VAR     NUMBER   df2b3f1b5df6e944  377/1
.VAR     NUMBER   0bd3ce83c84a1f2a  389/1
.CONST   NUMBER   c1003             15/1
.CONST   NUMBER   c1007             14/1
```

To dump the AST pseudo assembly do this:
```
c -i examples/e1.sml -ast
// intermediary language dump

// line 1: a = 12 + 13 * (14 + 15);
        push    12/1
        push    13/1
        push    14/1
        push    15/1
        add
        mul
        add
        pop     a

// line 2: b = 12 + 13 * (14 + -15);
        push    12/1
        push    13/1
        push    14/1
        push    15/1
        neg
        add
        mul
        add
        pop     b
```

Note: unfortunately go does not support running tasks yet.  So be sure to run the Makefile in frontend/ml/ if you change the grammar or tokenizer.

This code uses the liberal ISC license.
