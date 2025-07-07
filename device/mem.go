package device

import (
	"encoding/json"
	"fmt"
	"git.sophuwu.com/statlog/types"
	"os"
	"strings"
)

type MEM struct {
	Total   types.Bytes `json:"Total"`
	Free    types.Bytes `json:"Free"`
	Buffer  types.Bytes `json:"Buffer"`
	Used    types.Bytes `json:"Used"`
	Percent struct {
		Used types.Percent `json:"Used"`
		Buff types.Percent `json:"Buff"`
		Free types.Percent `json:"Free"`
	} `json:"Percent"`
}

func (m *MEM) Update() {
	b := make([]byte, 140)
	f, _ := os.Open("/proc/meminfo")
	_, err := f.Read(b)
	if err != nil {
		return
	}
	f.Close()
	var seeker types.NumSeeker
	seeker.Init(b)
	m.Total.FromKiB(seeker.GetNum())
	m.Free.FromKiB(seeker.GetNum())
	m.Used.FromKiB(seeker.GetNum())
	m.Used = m.Total - m.Used
	m.Buffer.FromKiB(seeker.GetNum() + seeker.GetNum())
	m.Percent.Used.CalcBytes(m.Used, m.Total)
	m.Percent.Buff.CalcBytes(m.Buffer, m.Total)
	m.Percent.Free.CalcBytes(m.Free, m.Total)
}

func (m *MEM) String() string {
	return fmt.Sprintf("%s (%s) USED  %s BUFF  %s FREE", m.Used, m.Percent.Used.Compact(), m.Buffer, m.Free)
}

func (m *MEM) Bar() (string, error) {
	w, _ := types.TermSize()
	if w < 40 {
		return "", fmt.Errorf("terminal too narrow")
	}

	s := "MEM: "
	// s += "[\033[1;35m" + m.Used.Human() + " USED\033[0m  "
	// s += "\033[1;36m" + m.Buffer.Human() + " BUFF\033[0m  "
	// s += "\033[1;32m" + m.Free.Human() + " FREE\033[0m]"
	s += "[\033[1;35mUSED\033[0m/"
	s += "\033[1;36mBUFF\033[0m/"
	s += "\033[1;32mFREE\033[0m] "

	l := len(strings.NewReplacer("\033[1;32m", "", "\033[1;36m", "", "\033[1;35m", "", "\033[0m", "").Replace(s))
	w -= l
	if w < 40 {
		return "", fmt.Errorf("terminal too narrow")
	}
	w -= 2

	s += "[\033[1;35m"
	i, ul := 0, int(m.Percent.Used.Value()*float64(w)/100)
	c := m.Percent.Used.Compact()
	s += c
	for i = len(c); i < w && i < ul; i++ {
		s += "|"
	}
	s += "\033[0m\033[1;36m"
	c = m.Percent.Buff.Compact()
	s += c
	i += len(c)
	for ul += int(m.Percent.Buff.Value() * float64(w) / 100); i < w && i < ul; i++ {
		s += "|"
	}
	s += "\033[0m\033[1;32m"
	c = m.Percent.Free.Compact()
	s += c
	i += len(c)
	for ; i < w; i++ {
		s += "|"
	}
	s += "\033[0m]"
	return s, nil
}

func (m *MEM) JSON() (string, error) {
	b, e := json.Marshal(m)
	return string(b), e
}
