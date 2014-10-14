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

type vmResponse struct{}

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

func (v *Vm) RunInteractive() error {
	vmCmd := make(chan string, 1)
	line := make(chan string)
	interrupt := make(chan string, 1)

	fmt.Printf("=== TVM V1.0 ===\n\npress h for help\n\n")

	go readline(line)

	// catch ctrl-c to pause execution
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			interrupt <- "interrupt"
		}
	}()

	running := false
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
				go func() {
					fmt.Printf("program started\n")
					line <- ""
					running = true
					v.GC() // reset old symbols and stats
					t1 := time.Now()
					err := v.run(vmCmd, true)
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
				running = false
				fmt.Printf("vm paused\n")
			default:
				fmt.Printf("invalid command %v\n", l)
			}
		case <-interrupt:
			//fmt.Printf("pc %016x\n", v.pc)
			fmt.Printf("interrupt!\n")
			if running == false {
				fmt.Printf("vm not running\n")
				continue
			}
			vmCmd <- "pause"
		}
	}

	return nil
}
