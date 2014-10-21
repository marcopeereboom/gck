// hide some of the terrible rust that is interactivity in its own file
package vm

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type vmCommand struct {
	cmd interface{}
}

type vmResponse struct {
	rv  interface{}
	err error
}

// readline
func readline(line chan string) {
	r := bufio.NewReader(os.Stdin)

	for {
		l, err := r.ReadBytes('\n')
		if err != nil {
			line <- "quit"
			return
		}
		line <- strings.Trim(string(l), "\n\r\t ")
	}
}

// cmd handles incoming interactive commands and produces a reply.
func (v *Vm) cmd(cmd vmCommand) vmResponse {
	r := vmResponse{}

	switch c := cmd.cmd.(type) {
	case string:
		// simple string commands
		switch c {
		case "break":
			v.paused = true
			v.tainted = true // mark stats as tainted
			r.rv = fmt.Sprintf("break point PC: %016x", v.pc)
		case "pause":
			r.rv = "paused"
			v.paused = true
			v.tainted = true // mark stats as tainted
		case "unpause":
			r.rv = "unpaused"
			v.paused = false
		case "pc":
			r.rv = fmt.Sprintf("PC: %016x", v.pc)
		case "gc":
			s := len(v.sym)
			v.GC()
			e := len(v.sym)
			r.rv = fmt.Sprintf("reaped %v symbols", s-e)
		case "sym":
			r.rv = fmt.Sprintf("%v", v.GetSymbols(true))
		case "s":
			r.rv = fmt.Sprintf("%v", v.GetStack(true, VmCmdStack))
		case "cs":
			r.rv = fmt.Sprintf("%v", v.GetStack(true, VmCallStack))
		}

	default:
		r.err = fmt.Errorf("invalid interactive command type %T", cmd)
	}

	return r
}

func (v *Vm) RunInteractive() error {
	cmd := make(chan vmCommand, 1)
	response := make(chan vmResponse, 1)
	line := make(chan string, 1)
	interrupt := make(chan bool, 1)

	fmt.Printf("=== TVM V1.0 ===\n\npress h for help\n\n")

	go readline(line)

	// catch ctrl-c to pause execution
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			interrupt <- true
		}
	}()

	running := false // may race, doesn't matter
	for {
		fmt.Printf("> ")
		select {
		case l := <-line:
			s := strings.Split(l, " ")
			if len(s) == 0 {
				continue
			}
			switch s[0] {
			default:
				fmt.Printf("invalid command %v\n", s[0])

			case "":
				continue

			case "h", "help":
				fmt.Printf("h, help - this help\n")
				fmt.Printf("q, quit - exit tvm\n")
				fmt.Printf("r, run - run image\n")
				fmt.Printf("pc - print program counter\n")
				fmt.Printf("sym, symbols - dump symbol table\n")
				fmt.Printf("s, stack - dump stack\n")
				fmt.Printf("cs, callstack - dump call stack\n")
				fmt.Printf("gc, garbagecollect - run GC\n")
				fmt.Printf("c, continue - resume execution\n")
				fmt.Printf("ctrl-c - pause execution\n")
				fmt.Printf("d, disassemble <start> <count>" +
					" - disassemble code segment; " +
					"use D for extra verbosity\n")
			case "q", "quit":
				return nil
			case "r", "run":
				if running == true {
					fmt.Printf("program already running\n")
					continue
				}

				go func() {
					fmt.Printf("program started\n")
					line <- ""
					running = true
					t1 := time.Now()
					v.run(cmd, response, true)
					t2 := time.Now()
					s := ""
					if v.tainted {
						s = "tainted statistics "
					}
					d := t2.Sub(t1)
					df := d.Seconds()
					insf := float64(v.instructions)
					fmt.Printf("program exited, "+
						"%vruntime %v instructions %v %v MIPS\n",
						s,
						d,
						v.instructions,
						insf/df/1e6)
					running = false
				}()

			case "pc":
				if running {
					cmd <- vmCommand{cmd: "pc"}
				} else {
					fmt.Printf("PC: %016x\n", v.pc)
				}
			case "sym", "symbols":
				if running {
					cmd <- vmCommand{cmd: "sym"}
				} else {
					fmt.Printf("%v", v.GetSymbols(true))
				}
			case "s", "stack":
				if running {
					cmd <- vmCommand{cmd: "s"}
				} else {
					fmt.Printf("%v", v.GetStack(true, VmCmdStack))
				}
			case "cs", "callstack":
				if running {
					cmd <- vmCommand{cmd: "cs"}
				} else {
					fmt.Printf("%v", v.GetStack(true, VmCallStack))
				}
			case "gc", "garbagecollect":
				if running {
					cmd <- vmCommand{cmd: "gc"}
				} else {
					v.GC()
				}
			case "c", "continue":
				if running {
					cmd <- vmCommand{cmd: "unpause"}
				} else {
					fmt.Printf("vm not running\n")
				}
			case "b", "break":
				var (
					brk uint64
					err error
				)

				if len(s) > 1 {
					brk, err = strconv.ParseUint(s[1], 0, 64)
					if err != nil {
						fmt.Printf("break: %v", err)
						continue
					}
					if brk >= uint64(len(v.prog)) {
						fmt.Printf("out of bounds\n")
						continue
					}
				}
				v.SetBreak(brk)

			case "d", "D", "disassemble":
				var (
					start, pc uint64
					count     int = 20
					err       error
				)

				if len(s) > 1 {
					start, err = strconv.ParseUint(s[1], 0, 64)
					if err != nil {
						fmt.Printf("start: %v", err)
						continue
					}
					if start >= uint64(len(v.prog)) {
						fmt.Printf("out of bounds\n")
						continue
					}
				}
				if len(s) > 2 {
					cnt, err := strconv.ParseInt(s[2], 0, 64)
					if err != nil {
						fmt.Printf("count: %v", err)
						continue
					}
					count = int(cnt)
				}
				// code segment is readonly
				pc = start
				for i := 0; i < count; i++ {
					ins := v.disassemble(s[0] == "D",
						pc, v.prog)
					fmt.Printf("%016x: %v\n", pc, ins)
					if v.prog[pc] >= OP_INVALID {
						pc += 1
					} else {
						pc += vmInstructions[v.prog[pc]].size
					}
					if pc >= uint64(len(v.prog)) {
						fmt.Printf("--- end of " +
							"image ---\n")
						break
					}
				}
			}

		case <-interrupt:
			if running == false {
				fmt.Printf("vm not running\n")
				continue
			}
			cmd <- vmCommand{cmd: "pause"}
		case r := <-response:
			if r.err != nil {
				fmt.Printf("return value: %v\n", r.err)
				if r.err != ErrExit {
					continue
				}
			}

			// handle response
			if r.rv == nil {
				continue
			}
			switch res := r.rv.(type) {
			case string:
				fmt.Printf("%v\n", strings.Trim(res, "\r\n"))
			default:
				fmt.Printf("invalid response type %T\n", r.rv)
			}
		}
	}

	return nil
}
