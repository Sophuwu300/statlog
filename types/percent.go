package types

import "fmt"

type Percent float64

func (p Percent) String() string {
	if p < 0 {
		return ""
	}
	return fmt.Sprintf("%3d %%", int(p))
}

func (p Percent) Compact() string {
	if p < 0 {
		return ""
	}
	return fmt.Sprintf("%.0f%%", p)
}

func (p *Percent) Value() float64 {
	return float64(*p)
}
func (P *Percent) Clamp() {
	if *P < 0 {
		*P = 0
	} else if *P > 100 {
		*P = 100
	}
}
func (P *Percent) CalcBytes(Vn, Vmax Bytes) {
	P.Calc(float64(Vn), float64(Vmax))
}
func (P *Percent) Calc(Vn, Vmax float64) {
	if Vmax == 0 {
		*P = 0
		return
	}
	*P = Percent(100 * Vn / Vmax)
}
func (P *Percent) FromFloat(v float64) {
	*P = Percent(v * 100)
}
func (P *Percent) SetValue(v float64) {
	*P = Percent(v)
}
func (P *Percent) SetValueInt(v int) {
	*P = Percent(v)
}

func (P *Percent) Bar(title string, w int) (s string, err error) {
	if w == 0 {
		w, _ = TermSize()
	}
	c := P.Compact()
	s = title + " ["
	w = w - len(s) - 2 - len(c)

	if w < 10 {
		return "", fmt.Errorf("terminal too narrow")
	}

	v := int(P.Value() * float64(w) / 100)
	i := 0
	for i = 0; i < v; i++ {
		s += "|"
	}
	for ; i < w; i++ {
		s += " "
	}
	s += c + "] "
	return s, nil
}
