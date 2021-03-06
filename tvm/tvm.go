// tvm - Toy Virtual Machine
// tvm wraps the tvm/vm package and makes it an executable command.
// For example, to run a compiled image:
//	tvm -t -i image.bin
// This would execute the code contained in image.bin and it keeps a runtime
// trace.
// See tvm/vm package for detailed information of the virtual machine
// implementation.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/marcopeereboom/gck/tvm/vm"
)

var (
	interactive bool
	trace       bool
	in          string
)

func init() {
	flag.BoolVar(&interactive, "I", false, "run tvm in interactive mode")
	flag.BoolVar(&trace, "t", false, "dump runtime trace")
	flag.StringVar(&in, "i", "", "binary image")
}

func _main() error {
	// read image file
	image, err := ioutil.ReadFile(in)
	if err != nil {
		return err
	}

	// run
	v, err := vm.New(image)
	if err != nil {
		return err
	}

	// see if we want a runtime trace
	if trace {
		v.Trace(false)
	}

	// execute
	if interactive {
		err = v.RunInteractive()
	} else {
		err = v.Run()
	}
	if err != nil {
		if err != vm.ErrExit {
			return err
		}
	}

	// garbage collect
	v.GC()

	// dump results
	if trace {
		fmt.Printf("=== run trace  ===\n%v", v.GetTrace())
		fmt.Printf("=== cmd stack  ===\n%v",
			v.GetStack(false, vm.VmCmdStack))
		fmt.Printf("=== call stack ===\n%v",
			v.GetStack(false, vm.VmCallStack))
		fmt.Printf("=== symbols    ===\n%v", v.GetSymbols(true))

	}

	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// command line
	flag.Parse()

	// check required flags
	if in == "" {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n-i must be provided\n")
		os.Exit(1)
	}

	err := _main()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
