package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.sophuwu.com/statlog/types"
	"os"
	"strings"
	"time"
)

type CPU struct {
	LoadAvg  types.Percent   `json:"LoadAvg"`
	LoadCore []types.Percent `json:"LoadCore"`
	MHzAvg   types.MHz       `json:"MHzAvg"`
	MHzCore  []types.MHz     `json:"MHzCore"`
	Temp     types.Celsius   `json:"Temp"`
}

func (c *CPU) LoadBar() (s string, err error) {
	w, _ := types.TermSize()
	l := len(c.LoadCore)
	ss := ""
	div := 3
	if w < 100 {
		div = 2
	} else if w > 200 {
		div = 4
	}
	w = (w - 2) / div
	for i := 0; i < l; i++ {
		if ss, err = c.LoadCore[i].Bar(fmt.Sprintf("%2d", i), w); err != nil {
			return "", err
		}
		s += ss
		if i%div == div-1 {
			s += "\n"
		} else {
			s += " "
		}
	}
	if s[len(s)-1] == ' ' || s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return s, nil
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
	cpu.Temp = (types.Celsius(b[0])-48)*10 + types.Celsius(b[1]) - 48
}
func (c *CPU) loadMHz() {
	c.MHzAvg = 0
	b, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return
	}
	var ns types.NumSeeker
	c.MHzCore = make([]types.MHz, 0)
	tot, n := 0.0, 0.0
	for _, v := range bytes.Split(b, []byte("\n")) {
		if bytes.HasPrefix(v, []byte("cpu MHz")) {
			ns.Init(v)
			n = ns.GetFloat()
			tot += n
			c.MHzCore = append(c.MHzCore, types.MHz(n))
		}
	}
	if len(c.MHzCore) == 0 {
		c.MHzAvg = 0
		c.MHzCore = nil
		return
	}
	c.MHzAvg = types.MHz(tot / float64(len(c.MHzCore)))
}

type coreLoad struct {
	v [2][2]float64
}

func (c *coreLoad) readStat(t, i *int, n *float64) {
	*n += c.v[*t][0]
	if *i == 3 {
		c.v[*t][1] = *n
		return
	}
	c.v[*t][0] = *n
}

func (c *coreLoad) percent() float64 {
	if c.v[0][1] == 0 {
		return 0.0
	}
	return (c.v[1][0] - c.v[0][0]) / (c.v[1][1] - c.v[0][1])
}

func (c *CPU) loadUse() {
	var cl []coreLoad
	loadVals := func(t int) error {
		i, j, k, n := 0, 0, 0, 0.0
		var el coreLoad
		elget := func() {
			el = coreLoad{[2][2]float64{{0.0, 0.0}, {0.0, 0.0}}}
		}
		elsv := func() {
			cl = append(cl, el)
		}
		if t >= 1 {
			t = 1
			elget = func() {
				el = cl[k]
			}
			elsv = func() {
				cl[k] = el
				k++
			}
		}

		b, err := os.ReadFile("/proc/stat")
		if err != nil {
			return err
		}
		for i = 0; i < len(b)-3 && (t == 0 || k < len(cl)); i++ {
			if !(b[i] == 'c' && b[i+1] == 'p' && b[i+2] == 'u') {
				break
			}
			for i < len(b) {
				if b[i] >= '0' && b[i] <= '9' {
					break
				}
				i++
			}
			elget()
			n = 0.0
			for j = 0; j < 4 && i < len(b); i++ {
				if b[i] >= 48 && b[i] <= 57 {
					n = n*10 + float64(b[i]) - 48
				} else if b[i] == ' ' {
					el.readStat(&t, &j, &n)
					n = 0.0
					j++
				} else {
					return fmt.Errorf("unexpected byte %c at index %d", b[i], i)
				}
			}
			if j < 3 {
				return fmt.Errorf("unexpected end of line at index %d", i)
			}
			elsv()
			for i < len(b) {
				if b[i] == '\n' {
					break
				}
				i++
			}
		}
		return nil
	}
	err := loadVals(0)
	if err != nil || len(cl) == 0 {
		goto onErr
	}
	time.Sleep(types.SampleDuration)
	err = loadVals(1)
	if err != nil || len(cl) == 0 {
		goto onErr
	}
	c.LoadAvg.FromFloat(cl[0].percent())
	cl = cl[1:]
	if len(cl) == 0 {
		goto onEmpty
	}
	c.LoadCore = make([]types.Percent, len(cl))
	for i, j := range cl {
		c.LoadCore[i].FromFloat(j.percent())
	}
	return
onErr:
	c.LoadAvg = -1.0
onEmpty:
	c.LoadCore = nil
}

func (c *CPU) Update() {
	c.loadMHz()
	c.loadTemp()
	c.loadUse()
}

func numSize(n int) int {
	if n == 0 {
		return 1
	}
	i := 0
	if n < 0 {
		n *= -1
		i++
	}
	for n > 0 {
		i++
		n /= 10
	}
	return i
}
func (c *CPU) LoadCoreStr() (s string) {
	if c.LoadCore == nil || len(c.LoadCore) == 0 {
		return s
	}
	var i int
	fmtstr := "cpu %." + fmt.Sprintf("%d", numSize(len(c.LoadCore))) + "d: %s"
	fn := func() {
		s += fmt.Sprintf(fmtstr, i, c.LoadCore[i])
	}
	if c.MHzCore != nil && len(c.LoadCore) == len(c.MHzCore) {
		fmtstr += "  %s"
		fn = func() {
			s += fmt.Sprintf(fmtstr, i, c.LoadCore[i], c.MHzCore[i])
		}
	}
	for i = 0; i < len(c.LoadCore); i++ {
		fn()
		if i%3 == 2 {
			s += "\n"
		} else {
			s += "  |  "
		}
	}
	if strings.HasSuffix(s, "\n") {
		s = s[:len(s)-1]
	} else if strings.HasSuffix(s, "  |  ") {
		s = s[:len(s)-5]
	}
	return
}

func (c CPU) String() string {
	return fmt.Sprintf("%s  %s  %s", c.LoadAvg, c.MHzAvg, c.Temp)
}
func (c *CPU) JSON() (string, error) {
	b, e := json.Marshal(c)
	return string(b), e
}
