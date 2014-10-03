// stdlib contains callouts from tvm to go.
// This is where you put performance critical things and OS type functionality.
// For example, a print function could go in here.
package stdlib

import "fmt"

type Result struct {
	Error error         // indicate if call failed or succeeded
	Rv    []interface{} // return values
}

const (
	Version = 1

	RetTrue  = "os.true"
	RetFalse = "os.false"
	RetError = "os.error"
)

var (
	jumpTable = map[string]func(...interface{}) (*Result, error){
		// test functions
		RetTrue:  retTrue,
		RetFalse: retFalse,
		RetError: retError,

		// actual functions
	}
)

func GetFunctionNames() []string {
	fn := make([]string, 0, len(jumpTable))
	for k := range jumpTable {
		fn = append(fn, k)
	}
	return fn
}

func Dispatch(name string) (*Result, error) {
	f, found := jumpTable[name]
	if !found {
		return nil, fmt.Errorf("stdlib function not found: %v", name)
	}

	return f()
}

func retTrue(args ...interface{}) (*Result, error) {
	return &Result{Rv: []interface{}{true}}, nil
}

func retFalse(args ...interface{}) (*Result, error) {
	return &Result{Rv: []interface{}{false}}, nil
}

func retError(args ...interface{}) (*Result, error) {
	return nil, fmt.Errorf("returned error")
}
