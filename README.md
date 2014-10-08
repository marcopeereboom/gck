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
* Frontend - Language to AST translator, components are lexer and parser
* AST - The Abstract Syntax Tree is in this system nothing but plumbing between subsystems
* Optimizer -  Transforms an AST into a better AST
* Backend - Transforms an AST into bytecode
* Virtual Machine - takes the bytecode and executes it

Currently the front and backend are modular and a new language or a new target
can simply be added by adding a driver that adheres to the correct interface.

The provided front and backend examples are:
* sml - Simple Math Language a simple language that does mostly math type operations
* tvm - Toy Virtual Machine is a simple VM that runs tvm binaries

## Examples
First compile the code, we use e1 from the examples directory.
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
.VAR     NUMBER     1   b                 -1/1
.CONST   NUMBER     1   c1001             13/1
.CONST   NUMBER     1   c1003             15/1
.VAR     NUMBER     1   a                 389/1
.CONST   NUMBER     1   c1000             12/1
.CONST   NUMBER     1   c1002             14/1
```
As you can see that generates a lot of stuff.
But that aside the following 2 lines are what matter:
```
.VAR     NUMBER     1   a                 389/1
.VAR     NUMBER     1   b                 -1/1
```
The astute reader can see that the math actually is correct.

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

**Note: unfortunately go does not support running tasks yet.  So be sure to run the Makefile in frontend/ml/ if you change the grammar or tokenizer.**

## License
This code uses the liberal ISC license.
