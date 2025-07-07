package main

import (
	"fmt"
	"git.sophuwu.com/statlog"
)

func main() {
	hw := &statlog.HWInfo{}
	hw.Update()
	fmt.Println(hw)
}
