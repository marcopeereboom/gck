// c - compiler combines all stages
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/marcopeereboom/gck/ast"
	"github.com/marcopeereboom/gck/backend"
	"github.com/marcopeereboom/gck/frontend"
	"github.com/marcopeereboom/gck/optimizer"
)

var (
	lang     string
	target   string
	in       string
	out      string
	pAST     bool
	optimize bool
)

func langUsage() string {
	return fmt.Sprintf("currently supported languages: %v; default %v",
		frontend.SML,
		frontend.SML)
}

func targetUsage() string {
	return fmt.Sprintf("currently supported architectures: %v; default %v",
		backend.TVM,
		backend.TVM)
}

func init() {
	flag.BoolVar(&pAST, "ast", false, "dump pseudo assembly AST")
	flag.BoolVar(&optimize, "O", false, "enable optimizer")
	flag.StringVar(&lang, "lang", frontend.SML, langUsage())
	flag.StringVar(&target, "target", backend.TVM, targetUsage())
	flag.StringVar(&in, "i", "", "source file")
	flag.StringVar(&out, "o", "-", "output file; default stdout")
}

func _main() error {
	fe, err := frontend.New(lang)
	if err != nil {
		return err
	}

	// read source file
	src, err := ioutil.ReadFile(in)
	if err != nil {
		return err
	}

	// Compile source
	err = fe.Compile(string(src))
	if err != nil {
		return err
	}

	// obtains AST
	a, err := fe.AST()
	if err != nil {
		return err
	}

	// optimize AST
	ao, err := optimizer.Optimize(a)
	if err != nil {
		return err
	}

	// dump AST pseudo asm
	if pAST {
		var w io.Writer
		if out == "-" {
			w = os.Stdout
		} else {
			w, err = os.Create(out)
			if err != nil {
				return err
			}
		}
		ast.DumpPseudoAsm(ao, w)

		// abort compilation
		return nil
	}

	// obtain binary image
	// XXX this New must be done much earlier
	t, err := backend.New(target)
	if err != nil {
		return err
	}
	bi, err := t.EmitCode(ao)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(out, bi, 0660)
	if err != nil {
		return err
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
