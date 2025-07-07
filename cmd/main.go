package main

import (
	"fmt"
	"git.sophuwu.com/statlog"
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
	var e error
	for bl {
		hw.Update()
		s, e = hw.CPU.LoadBar()
		fatal(e)
		fmt.Printf("\033[2J\033[1;1H\r%s\n", s)
	}
}
