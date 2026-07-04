package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"git.sophuwu.com/statlog"
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
	hw := &statlog.HWInfo{}
	hw.Update()
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

	var grCpu func(w, h, val int) (string, error)
	grCpu, e = types.Graph("CPU Load", 100)
	fatal(e)

	var grMem func(w, h, val int) (string, error)
	grMem, e = types.Graph("Memory (%)", 100)
	fatal(e)

	var tempGPU func(w, h, val int) (string, error)
	tempGPU, e = types.Graph("GPU Temp", 100)
	fatal(e)

	var tempCPU func(w, h, val int) (string, error)
	tempCPU, e = types.Graph("CPU Temp", 100)
	fatal(e)

	var w, h, t int
	prnt := func() {
		fmt.Printf("\033[2J\033[1;1H\r%s\n", s)
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
	szfn := func(){
		w, h = types.TermSize()
		h-=30
	}
	szfn()
	for bl {
		s = ""
		szfn()
		hw.Update()

		// mem
		ss, e = hw.MEM.Bar()

		if ERR() {
			continue
		}
		s += "MEM: " + hw.MEM.String() + "\n" + ss + "\n"
		ss, e = grMem(w, h/6, int(hw.MEM.Percent.Used))
		if ERR() {
			continue
		}
		s += ss + "\n"

		// cpu

		// print cpu info string
		s += "CPU: " + hw.CPU.String() + "\n"
		// make load graph
		ss, e = grCpu(w, h/6, int(hw.CPU.LoadAvg))
		if ERR() {
			continue
		}
		// print load graph
		s += ss + "\n"

		// cpu bar graphs
		// make bar graph
		ss, e = hw.CPU.LoadAvg.Bar("CPU Avg", w)
		if ERR() {
			continue
		}
		// print bar graph
		s += ss + "\n"
		// make core bar graphs
		ss, e = hw.CPU.LoadBar(w)
		if ERR() {
			continue
		}
		// print core bar graphs
		s += ss + "\n"

		ss, e = tempCPU(w, h/4, int(hw.CPU.Temp))
		if ERR() {
			continue
		}
		// print cpu temp graph
		s += ss + "\n"


		t, e = getGPU()
		if ERR() {
			continue
		}
		ss, e = tempGPU(w, h/4, t)
		if ERR() {
			continue
		}
		// print gpu temp graph
		s += ss + "\n"

		prnt()
	}
}
