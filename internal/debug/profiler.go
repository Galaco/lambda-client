package debug

import (
	"os"
	"runtime/pprof"
)

func StartProfiling(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return pprof.StartCPUProfile(f)
}

func StopProfiling() {
	pprof.StopCPUProfile()
}