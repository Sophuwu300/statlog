package types

import (
	"fmt"
	"math"
)

type NumSeeker struct {
	i int
	b []byte
}

func (n *NumSeeker) Init(b []byte) {
	n.b = b
	n.i = 0
}
func (n *NumSeeker) End() bool {
	return n.i >= len(n.b)
}
func (n *NumSeeker) Seek() {
	n.i++
}
func (n *NumSeeker) GetByte() uint8 {
	return uint8(n.b[n.i])
}
func (n *NumSeeker) GetOffset(off int) (uint8, error) {
	off += n.i
	if off < 0 || off >= len(n.b) {
		return 0, fmt.Errorf("index out of range: %d", n)
	}
	return uint8(n.b[off]), nil
}

func (n *NumSeeker) IsNum() bool {
	e := n.GetByte()
	return e >= 48 && e <= 57
}
func (n *NumSeeker) SeekToNum() {
	for ; !n.End(); n.Seek() {
		if n.IsNum() {
			break
		}
	}
}
func (n *NumSeeker) GetNum() uint64 {
	var num uint64 = 0
	for n.SeekToNum(); !n.End(); n.Seek() {
		if n.IsNum() {
			num = num*uint64(10) + uint64(uint8(n.GetByte())-48)
		} else {
			return num
		}
	}
	return num
}
func (n *NumSeeker) GetFloat() float64 {
	n.SeekToNum()
	if n.End() {
		return 0
	}
	var sign float64 = 1
	if b, _ := n.GetOffset(-1); b == '-' {
		sign = -1
		n.Seek()
	}
	ipart := float64(n.GetNum())
	if n.GetByte() == '.' {
		fpart := float64(n.GetNum())
		for i, j := 0.0, math.Log10(fpart); i < j; i++ {
			fpart *= 0.1
		}
		ipart += fpart
	}
	return ipart * sign
}

func (n *NumSeeker) GetNums() []uint64 {
	nums := make([]uint64, 0)
	for ; !n.End(); n.Seek() {
		nums = append(nums, n.GetNum())
	}
	return nums
}
