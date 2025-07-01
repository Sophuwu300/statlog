package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.sophuwu.com/statlog/types"
	"os"
	"time"
)

type CPU struct {
	Load float64
	MHz  float64
	Temp int
}

func (cpu *CPU) loadTemp() {
	if cpu.Temp == -100 {
		return
	}
	var findCPUTypeNo = func(path string, fileName string, comp string) (string, error) {
		var b = make([]byte, 100)
		var i int = 0
		if fileName == "_label" {
			i = 1
		}
		var err error = nil
		for ; err == nil && string(b[:len(comp)]) != comp; i++ {
			b, err = os.ReadFile(path + string(i+48) + fileName)
		}
		if err != nil {
			return "", err
		}
		return string(i + 48 - 1), nil
	}
	Err := func(err error) bool {
		if err == nil {
			return false
		}
		cpu.Temp = -100
		return true
	}
	var (
		b         = make([]byte, 2)
		thrmPath  = "/sys/class/thermal/thermal_zone"
		hwMonPath = "/sys/class/hwmon/hwmon"
	)
	nStr, err := findCPUTypeNo(thrmPath, "/type", "x86_pkg_temp")
	if err != nil {
		nStr, err = findCPUTypeNo(hwMonPath, "/name", "k10temp")
		if Err(err) {
			return
		}
		nStrHW := nStr
		nStr, err = findCPUTypeNo(hwMonPath+nStrHW+"/temp", "_label", "Tdie")
		if Err(err) {
			return
		}
		b, err = os.ReadFile(hwMonPath + nStrHW + "/temp" + nStr + "_input")
	} else {
		b, err = os.ReadFile(thrmPath + nStr + "/temp")
	}
	if Err(err) {
		return
	}
	cpu.Temp = (int(b[0])-48)*10 + int(b[1]) - 48
}
func (c *CPU) loadMHz() {
	c.MHz = 0
	b, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return
	}
	var ns types.NumSeeker
	n := 0.0
	for _, v := range bytes.Split(b, []byte("\n")) {
		if bytes.HasPrefix(v, []byte("cpu MHz")) {
			ns.Init(v)
			c.MHz += ns.GetFloat()
			n++
		}
	}
	c.MHz /= n
}
func (c *CPU) loadUsage() {
	readStat := func(n *[4]float64) bool {

		b := make([]byte, 100)
		f, err := os.Open("/proc/stat")
		if err != nil {
			return true
		}
		_, err = f.Read(b)
		if err != nil {
			return true
		}
		f.Close()
		for i, j := 6, 0; j < 4; i++ {
			if b[i] >= 48 && b[i] <= 57 {
				n[j] = n[j]*10 + float64(b[i]) - 48
			} else if b[i] == ' ' {
				j++
			}
		}
		return false
	}
	var a, b [4]float64
	if readStat(&a) {
		return
	}
	time.Sleep(types.SampleDuration)
	if readStat(&b) {
		return
	}
	c.Load = ((b[0] + b[1] + b[2]) - (a[0] + a[1] + a[2])) / ((b[0] + b[1] + b[2] + b[3]) - (a[0] + a[1] + a[2] + a[3]))
}
func (c *CPU) update() {
	c.loadMHz()
	c.loadTemp()
	c.loadUsage()
}
func (c *CPU) GHzStr() string {
	return fmt.Sprintf("%.2f GHz", c.MHz/1000)
}
func (c *CPU) LoadStr() string {
	return fmt.Sprintf("%.1f %c", c.Load*100, '%')
}
func (c *CPU) TempStr() string {
	if c.Temp == -100 {
		return ""
	}
	return fmt.Sprintf("%d C", c.Temp)
}
func (c *CPU) String() string {
	return fmt.Sprintf("%6.6s %8.8s %5.5s", c.LoadStr(), c.GHzStr(), c.TempStr())
}
func (c *CPU) JSON() (string, error) {
	b, e := json.Marshal(c)
	return string(b), e
}
