package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"git.sophuwu.com/statlog/types"
)

func fatal(e error) {
	if e != nil {
		fmt.Println("\033[?1049l\033[?25h")
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

const gpuTempFile = "/sys/class/hwmon/hwmon1/temp1_input"

func getGPU() (int, error) {
	b, err := os.ReadFile(gpuTempFile)
	if err != nil { return 0, err }
	i := 0
	_, err = fmt.Sscanf(string(b), "%d", &i)
	if err != nil { return 0, err }
	return i/1000, err
}

func main() {
	fmt.Println("\033[?25l\033[?1049h\033[2J")
	defer fmt.Println("\033[?1049l\033[?25h")
	ch := make(chan os.Signal, 1)
	bl := true
	go func() {
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		bl = false
	}()
	var s string
	var ss string
	var e error
	var tempGPU func(w, h, val int) (string, error)
	tempGPU, e = types.Graph("GPU Temp", 100)
	fatal(e)
	var w, h, t int
	prnt := func() {
		fmt.Printf("\033[2J\033[1;1H\r%s\n", s)
		time.Sleep(250 * time.Millisecond)
	}
	ERR := func() bool {
		if errors.Is(e, types.ErrTooNarrow) {
			s = "terminal too narrow"
			prnt()
			time.Sleep(100 * time.Millisecond)
			return true
		}
		fatal(e)
		return false
	}
	for bl {
		s = ""
		w, h = types.TermSize()
		t, e = getGPU()
		if ERR() {
			continue
		}
		ss, e = tempGPU(w, h-10, t)
		if ERR() {
			continue
		}
		// print gpu temp graph
		s += ss + "\n"
		prnt()
	}
}
