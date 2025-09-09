package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

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
	var grMem func(w, h, val int) (string, error)
	grCpu, e = types.Graph("CPU Load", 100)
	fatal(e)
	grMem, e = types.Graph("Memory (%)", 100)
	fatal(e)
	var w int
	prnt := func() {
		fmt.Printf("\033[2J\033[1;1H\r%s\n", s)
	}
	for bl {
		w, _ = types.TermSize()
		s = ""
		hw.Update()
		w, _ = types.TermSize()

		// mem
		ss, e = hw.MEM.Bar()
		if errors.Is(e, types.ErrTooNarrow) {
			s = "Terminal too narrow"
			prnt()
			continue
		}
		fatal(e)
		s += "MEM: " + hw.MEM.String() + "\n" + ss + "\n"
		ss, e = grMem(w, 5, int(hw.MEM.Percent.Used+hw.MEM.Percent.Buff))
		fatal(e)
		s += ss + "\n"

		// cpu

		// print cpu info string
		s += "CPU: " + hw.CPU.String() + "\n"
		// make load graph
		ss, e = grCpu(w, 5, int(hw.CPU.LoadAvg))
		fatal(e)
		// print load graph
		s += ss + "\n"

		// cpu bar graphs
		// make bar graph
		ss, e = hw.CPU.LoadAvg.Bar("CPU Avg", w)
		if errors.Is(e, types.ErrTooNarrow) {
			s = "Terminal too narrow"
			prnt()
			continue
		}
		fatal(e)
		// print bar graph
		s += ss + "\n"
		// make core bar graphs
		ss, e = hw.CPU.LoadBar(w)
		if errors.Is(e, types.ErrTooNarrow) {
			s = "Terminal too narrow"
			prnt()
			continue
		}
		fatal(e)
		// print core bar graphs
		s += ss + "\n"

		prnt()
	}
}
