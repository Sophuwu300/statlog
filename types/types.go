package types

import (
	"fmt"
	"golang.org/x/term"
)

type MHz int

func (m MHz) String() string {
	if m < 0 {
		return ""
	}
	return fmt.Sprintf("%4d MHz", m)
}

type Celsius int

func (c Celsius) String() string {
	if c <= -50 {
		return ""
	}
	return fmt.Sprintf("%3d C", c)
}

func TermSize() (w int, h int) {
	var e error
	w, h, e = term.GetSize(0)
	if e != nil {
		w = 80
		h = 24
	}
	return
}
