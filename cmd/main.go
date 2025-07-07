package main

import (
	"fmt"
	"git.sophuwu.com/statlog"
)

func main() {
	hw := &statlog.HWInfo{}
	for i := 0; i < 5; i++ {
		hw.Update()
		s, err := hw.MEM.Bar()
		if err != nil {
			fmt.Println("\nError generating memory bar:", err)
			break
		}
		fmt.Printf("\r%s", s)
	}
	fmt.Println()
}
