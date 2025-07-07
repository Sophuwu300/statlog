package device

import (
	"encoding/json"
	"fmt"
	"git.sophuwu.com/statlog/types"
	"os"
)

type MEM struct {
	Total   types.Bytes `json:"Total"`
	Free    types.Bytes `json:"Free"`
	Buffer  types.Bytes `json:"Buffer"`
	Used    types.Bytes `json:"Used"`
	Percent struct {
		Used     types.Percent `json:"Used"`
		WithBuff types.Percent `json:"WithBuff"`
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
	m.Percent.WithBuff.CalcBytes(m.Used+m.Buffer, m.Total)
}

func (m *MEM) String() string {
	return fmt.Sprintf("%s (%s) USED  %s BUFF  %s FREE", types.Bytes(m.Used), m.Percent.Used.String(), m.Buffer, m.Free)
}

func (m *MEM) JSON() (string, error) {
	b, e := json.Marshal(m)
	return string(b), e
}
