// hide some of the terrible rust that is interactivity in its own file
package vm

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
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
		case "pause":
			r.rv = "paused"
			v.paused = true
			v.tainted = true // mark stats as tainted
		case "unpause":
			r.rv = "unpaused"
			v.paused = false
		case "pc":
			r.rv = fmt.Sprintf("PC: %016x", v.pc)
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
			switch l {
			case "":
				continue
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
					v.GC() // reset old symbols and stats
					t1 := time.Now()
					err := v.run(cmd, response, true)
					t2 := time.Now()
					if err != nil {
						fmt.Printf("run error: %v\n",
							err)
					} else {
						s := ""
						if v.tainted {
							s = "tainted statistics "
						}
						fmt.Printf("program exited "+
							"normally (%vruntime %v %v MIPS)\n",
							s,
							t2.Sub(t1),
							float64(v.instructions)/t2.Sub(t1).Seconds()/1000000)
					}
					running = false
					line <- ""
				}()

			case "sym", "symbols":
				fmt.Printf("%v", v.GetSymbols(true))
			case "s", "stack":
				fmt.Printf("%v", v.GetStack(true, VmCmdStack))
			case "cs", "callstack":
				fmt.Printf("%v", v.GetStack(true, VmCallStack))
			case "gc", "garbagecollect":
				v.GC()
			case "pause":
				fmt.Printf("vm paused\n")
			case "c", "continue":
				cmd <- vmCommand{cmd: "unpause"}
			case "pc":
				if running {
					cmd <- vmCommand{cmd: "pc"}
				} else {
					fmt.Printf("PC: %016x\n", v.pc)
				}
			default:
				fmt.Printf("invalid command %v\n", l)
			}
		case <-interrupt:
			if running == false {
				fmt.Printf("vm not running\n")
				continue
			}
			cmd <- vmCommand{cmd: "pause"}
		case r := <-response:
			if r.err != nil {
				fmt.Printf("%v\n", r.err)
				continue
			}

			// handle response
			switch res := r.rv.(type) {
			case string:
				fmt.Printf("%v\n", res)
			default:
				fmt.Printf("invalid response type %T\n", r.rv)
			}
		}
	}

	return nil
}
