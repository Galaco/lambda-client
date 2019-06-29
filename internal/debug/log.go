package debug

import "os"

type stdOut struct{}

func (log *stdOut) Write(data []byte) (n int, err error) {
	line := string(data) + "\n"
	return os.Stdout.Write([]byte(line))
}

func NewStdOut() *stdOut {
	return &stdOut{}
}
