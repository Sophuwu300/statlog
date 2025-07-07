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
