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
* Myrmidon - Simple scripting language; unlike sml Myrmidon suports functions.
* tvm - Toy Virtual Machine is a simple VM that runs tvm binaries

## Examples
First compile the code, we use e1 from the examples directory.
```
c -i examples/sml/e1.sml -o /tmp/image.bin
tvm -i /tmp/image.bin -t
```

should result in something like this:
```
=== run trace  ===
0000000000000000: jsr    main (0x3)
0000000000000003: push   c1001 (12)
0000000000000005: push   c1002 (13)
0000000000000007: push   c1003 (14)
0000000000000009: push   c1004 (15)
000000000000000b: add   
000000000000000c: mul   
000000000000000d: add   
000000000000000e: pop    a (0/1)
0000000000000010: push   c1001 (12)
0000000000000012: push   c1002 (13)
0000000000000014: push   c1003 (14)
0000000000000016: push   c1004 (15)
0000000000000018: neg   
0000000000000019: add   
000000000000001a: mul   
000000000000001b: add   
000000000000001c: pop    b (0/1)
000000000000001e: ret   
0000000000000002: exit  
=== cmd stack  ===
=== call stack ===
=== symbols    ===
.VAR     INTEGER    1   a                 389
.CONST   INTEGER    1   c1001             12
.VAR     INTEGER    1   b                 -1
.CONST   INTEGER    1   c1002             13
.CONST   INTEGER    1   c1003             14
.CONST   LABEL      1   main              0x3
.CONST   INTEGER    1   c1004             15
```
As you can see that generates a lot of stuff.
But that aside the following 2 lines are what matter:
```
.VAR     INTEGER    1   a                 389
.VAR     INTEGER    1   b                 -1
```
The astute reader can see that the math actually is correct.

To dump the pseudo assembly do this:
```
c -i examples/sml/e1.sml -asm
// intermediary language dump
        jsr     main
        exit
main:

// line 2: a = 12 + 13 * (14 + 15);
        push    12
        push    13
        push    14
        push    15
        add
        mul
        add
        pop     a

// line 3: b = 12 + 13 * (14 + -15);
        push    12
        push    13
        push    14
        push    15
        neg
        add
        mul
        add
        pop     b
        ret
```
To dump the AST do this:
c -i examples/sml/e1.sml -ast
```
= \
   | a
   | + \
   |    | 12
   |    | * \
   |    |    | 13
   |    |    | + \
   |    |    |    | 14
   |    |    |    | 15
= \
   | b
   | + \
   |    | 12
   |    | * \
   |    |    | 13
   |    |    | + \
   |    |    |    | 14
   |    |    |    | - \
   |    |    |    |    | 15
```
**Note: unfortunately go does not support running tasks yet.  So be sure to run
the Makefile in frontend/sml/ or frontend/myrmidon if you change the grammar or
tokenizer.**

## License
This code uses the liberal ISC license.
