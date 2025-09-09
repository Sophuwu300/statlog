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

var ErrTooNarrow = fmt.Errorf("terminal too narrow")

func Bar(title string, w, val, max int, vlbl string) (string, error) {
	if max <= 0 || val < 0 || val > max {
		return "", fmt.Errorf("invalid values: val=%d, max=%d", val, max)
	}
	if w == 0 {
		w, _ = TermSize()
	}
	s := title + " ["
	w = w - len(s) - 2 - len(vlbl)

	if w < 10 {
		return "", ErrTooNarrow
	}

	val = int(float64(val) * float64(w) / float64(max))

	i := 0
	for ; i < w && i < val; i++ {
		s += "|"
	}
	for ; i < w; i++ {
		s += " "
	}
	s += vlbl + "] "

	return s, nil
}

var gch = []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

func Graph(title string, max int) (func(w, h, val int) (string, error), error) {
	var vals []int
	if max <= 0 {
		return nil, fmt.Errorf("invalid max value: %d", max)
	}
	maxs := fmt.Sprintf("%d", max)
	return func(w, h, val int) (string, error) {
		if w == 0 {
			w, _ = TermSize()
		}
		if h == 0 {
			h = 2
		}
		if val < 0 || val > max {
			return "", fmt.Errorf("invalid value: %d, max: %d", val, max)
		}
		vals = append(vals, val)
		if len(vals) < w {
			ii := make([]int, w-len(vals))
			vals = append(ii, vals...)
		} else if len(vals) > 2*w {
			vals = vals[len(vals)-2*w:]
		}
		s := fmt.Sprintf("%s\033[%dC%s\033[%dB\033[%dD0\033[%dD", title, w-len(maxs)-len(title), maxs, h+1, 1, w)
		s += fmt.Sprintf("%s\033[%dD", title, len(title)+2)
		var j, i int
		vv := make([]int, 0)
		j = len(vals) - w
		if j < 0 {
			return "", fmt.Errorf("not enough data points: %d, required: %d", len(vals), w)
		}
		for ; j < len(vals); j++ {
			vv = append(vv, vals[j]*h*len(gch)/max)
		}
		var v int
		for i = 0; i < h; i++ {
			s += "\033[1A"
			for j, v = range vv {
				v -= i * len(gch)
				if v < len(gch) {
					if v < 0 {
						v = 0
					}
					s += gch[v]
				} else if v >= len(gch) {
					s += gch[len(gch)-1]
				} else {
					s += gch[v%len(gch)]
				}
			}
			s += fmt.Sprintf("\033[0m\033[%dD", j+2)
		}
		s += fmt.Sprintf("\033[%dB\033[0m\n", h+2)
		return s, nil
	}, nil
}
