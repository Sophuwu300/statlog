package types

import (
	"fmt"
	"git.sophuwu.com/statlog/types/units"
)

type Bytes uint64

func (b Bytes) Human() string {
	n := float64(b)
	var i int
	for i = units.B; n >= units.KiB && i < units.GiB; i *= units.KiB {
		n /= units.KiB
	}

	return fmt.Sprintf("%.2f %s", n, units.Labels[i])
}

func (b Bytes) HumanSI() string {
	n := float64(b)
	var i int
	for i = units.B; n >= units.KB && i < units.GB; i *= units.KB {
		n /= units.KB
	}
	return fmt.Sprintf("%.2f %s", n, units.Labels[i])
}

func (b *Bytes) FromKiB(v uint64) {
	*b = Bytes(v * units.KiB)
}

func (b Bytes) Value() uint64 {
	return uint64(b)
}

func (b Bytes) String() string {
	return b.Human()
}

type Percent float64

func (p Percent) String() string {
	if p < 0 {
		return ""
	}
	return fmt.Sprintf("%3.0f %%", p)
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
