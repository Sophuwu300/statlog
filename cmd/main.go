package main

import (
	"fmt"
	"git.sophuwu.com/statlog"
	"git.sophuwu.com/statlog/types"
	"os"
	"os/signal"
)

func fatal(e error) {
	if e != nil {
		fmt.Println("\033[?1049l\033[?25h")
		os.Exit(1)
	}
}

func main() {
	fmt.Println("\033[?25l\033[?1049h\033[2J")
	defer fmt.Println("\033[?1049l\033[?25h")
	hw := &statlog.HWInfo{}
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
	var fn func(w, h, val int) (string, error)
	fn, e = types.Graph("CPU Load", 100)
	w, _ := types.TermSize()
	fatal(e)
	for bl {
		s = ""
		hw.Update()
		w, _ = types.TermSize()
		ss, e = fn(w/2, 4, int(hw.CPU.LoadAvg))
		fatal(e)
		s += ss + "\n\n"

		// s += "MEM: " + hw.MEM.String()
		// s += "\nCPU: " + hw.CPU.MHzAvg.String() + " MHz  " + hw.CPU.Temp.String() + "\n\n"
		// ss, e = hw.CPU.LoadAvg.Bar("CPU Avg", 0)
		// fatal(e)
		// s += ss + "\n\n"
		// ss, e = hw.MEM.Bar()
		// fatal(e)
		// s += ss + "\n\nCPU CORES:\n\n"
		// ss, e = hw.CPU.LoadBar()
		// fatal(e)
		// s += ss + "\n"

		fmt.Printf("\033[2J\033[1;1H\r%s\n", s)
	}
}
