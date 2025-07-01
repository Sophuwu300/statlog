package types

import (
	"fmt"
	"git.sophuwu.com/statlog/types/units"
)

type Bytes uint64

func (b Bytes) Human() string {
	n := float64(b)
	var i int
	for i = units.B; i < units.GiB; i *= units.KiB {
		if n < units.KiB {
			break
		}
		n /= units.KiB
	}
	return fmt.Sprintf("%.2f %s", n, units.Labels[i])
}

func (b Bytes) HumanSI() string {
	n := float64(b)
	var i int
	for i = units.B; i < units.GB; i *= units.KB {
		if n < units.KB {
			break
		}
		n /= units.KB
	}
	return fmt.Sprintf("%.2f %s", n, units.Labels[i])
}

func (b Bytes) Value() uint64 {
	return uint64(b)
}

func (b Bytes) String() string {
	return b.Human()
}

type Percent float64

func (p *Percent) String() string {
	return fmt.Sprintf("%.1f%%", *p)
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
func (P *Percent) SetValue(v float64) {
	*P = Percent(v)
}
func (P *Percent) SetValueInt(v int) {
	*P = Percent(v)
}
