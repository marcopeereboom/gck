package vm

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/marcopeereboom/gck/tvm/section"
)

var (
	loud  bool = true // unset for less dumping
	trace bool = false
)

func newImage(prog []uint64) (*section.Image, error) {
	cs := section.NewCodeSection(prog)

	// variables everyone uses
	v1, err := section.NewVariable(1000, "x", new(big.Rat).SetFloat64(2.0))
	if err != nil {
		return nil, err
	}
	v2, err := section.NewVariable(1001, "y", new(big.Rat).SetFloat64(3.0))
	if err != nil {
		return nil, err
	}
	vs, err := section.NewVariableSection([]*section.Variable{v1, v2})
	if err != nil {
		return nil, err
	}

	// constants everyone uses
	c1, err := section.NewConst(1002, "myjsr", uint64(12))
	if err != nil {
		return nil, err
	}
	cos, err := section.NewConstSection([]*section.Const{c1})
	if err != nil {
		return nil, err
	}

	// stdlib cross reference
	o := section.OsCall{
		Id:        1003,
		Name:      "os.true",
		Variables: nil,
		Results:   nil,
	}
	o1, err := section.NewOs(1003, "os.true", o)
	if err != nil {
		return nil, err
	}
	oss, err := section.NewOsSection([]*section.Os{o1})
	if err != nil {
		return nil, err
	}

	// generate image
	i := section.NewImage()
	err = i.AddSection(cs, true)
	if err != nil {
		return nil, err
	}
	err = i.AddSection(vs, true)
	if err != nil {
		return nil, err
	}
	err = i.AddSection(cos, true)
	if err != nil {
		return nil, err
	}
	err = i.AddSection(oss, true)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func execute(prog []uint64, t *testing.T) error {
	i, err := newImage(prog)
	if err != nil {
		return err
	}

	// store a copy of the image
	f, err := ioutil.TempFile(os.TempDir(), "tvm")
	if err != nil {
		return err
	}
	f.Close()
	err = ioutil.WriteFile(f.Name(), i.GetImage(), 0660)
	if err != nil {
		return err
	}
	t.Logf("wrote image to file: %v", f.Name())

	// run
	vm, err := New(i.GetImage())
	if err != nil {
		return err
	}
	vm.Trace(trace)
	err = vm.Run()
	if err != nil {
		return err
	}

	if loud {
		t.Logf("=== run trace  ===\n%v", vm.GetTrace())
		t.Logf("=== cmd stack  ===\n%v", vm.GetStack(false, VmCmdStack))
		t.Logf("=== call stack ===\n%v", vm.GetStack(false, VmCallStack))
		t.Logf("=== symbols    ===\n%v", vm.GetSymbols(true))
	}

	return nil
}

func TestSubr(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_JMP,  // 0
		4,       // 1 JMP to OP_JSR
		OP_NOP,  // 2
		OP_NOP,  // 3
		OP_JSR,  // 4
		1002,    // 5 lookup label in symbol table
		OP_PUSH, // 6
		1000,    // 7
		OP_PUSH, // 8
		1001,    // 9
		OP_MUL,  // 10
		OP_EXIT, // 11
		OP_NOP,  // 12
		OP_RET,  // 13
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestNop(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_NOP,
		OP_NOP,
		OP_NOP,
		OP_NOP,
		OP_NOP,
		OP_NOP,
		OP_NOP,
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPush0(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,
		0,
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPush1(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,
		0,
		OP_PUSH,
		1,
		OP_PUSH,
		2,
		OP_PUSH,
		3,
		OP_PUSH,
		4,
		OP_PUSH,
		5,
		OP_PUSH,
		6,
		OP_PUSH,
		7,
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPushIllegal(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,
	}

	err := execute(prog, t)
	if err == nil {
		t.Errorf("expected out of bounds")
		return
	}
}

func TestAdd(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,
		1000,
		OP_PUSH,
		1001,
		OP_ADD,
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAddIllegal(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,
		1,
		OP_ADD,
	}

	err := execute(prog, t)
	if err == nil {
		t.Error("expected underflow")
		return
	}
}

func TestSub(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,
		1001,
		OP_PUSH,
		1000,
		OP_SUB,
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestNeg(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,
		1001,
		OP_NEG,
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestIf(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,  // 0
		1001,     // 1
		OP_PUSH,  // 2
		1000,     // 3
		OP_GT,    // 4
		OP_BRT,   // 5
		8,        // 6 branch over abort
		OP_ABORT, // 7
		OP_NOP,   // 8
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPop(t *testing.T) {
	// code
	var prog []uint64 = []uint64{
		OP_PUSH, // 0
		1000,    // 1
		OP_POP,  // 2
		1001,    // 3 overwrite symbol 1001 with 1000
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestOs(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_CALL,
		1003,
	}

	err := execute(prog, t)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPopFail(t *testing.T) {
	var prog []uint64 = []uint64{
		OP_PUSH,
		1000,
		OP_POP,
		1001,
		OP_POP,
		1001,
	}

	err := execute(prog, t)
	if err == nil {
		t.Error("expected command stack underflow")
		return
	}
}
