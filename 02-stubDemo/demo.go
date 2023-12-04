package stubdemo

type Demo struct {
	msg string
}

// Num global variable.
var Num int

// Info global variable.
var Info string

func NewDemo(msg string) *Demo {
	return &Demo{msg: msg}
}

// GetMsg method.
//
//go:noinline
func (d *Demo) GetMsg() string {
	return d.msg
}

func DoStore(d Demo) string {
	return GetResult(d)
}

// GetResult function.
//
//go:noinline
func GetResult(d Demo) string {
	if d.GetMsg() == "store" {
		return "success"
	}
	return "failure"
}

//go:noinline
func getInfo() string {
	return "There are some information for you."
}
