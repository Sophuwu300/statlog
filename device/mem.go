package device

import (
	"encoding/json"
	"fmt"
	"git.sophuwu.com/statlog/types"
	"os"
)

type MEM struct {
	Total   types.Bytes
	Free    types.Bytes
	Buffer  types.Bytes
	Used    types.Bytes
	Percent struct {
		Used     types.Percent
		WithBuff types.Percent
	}
}

func (m *MEM) update() {
	b := make([]byte, 140)
	f, _ := os.Open("/proc/meminfo")
	_, err := f.Read(b)
	if err != nil {
		return
	}
	f.Close()
	var seeker types.NumSeeker
	seeker.Init(b)
	m.Total = types.Bytes(seeker.GetNum())
	m.Free = types.Bytes(seeker.GetNum())
	m.Used = m.Total - types.Bytes(seeker.GetNum())
	m.Buffer = types.Bytes(seeker.GetNum() + seeker.GetNum())
	m.Percent.Used.CalcBytes(m.Used, m.Total)
	m.Percent.WithBuff.CalcBytes(m.Used+m.Buffer, m.Total)
}

func (m *MEM) String() string {
	return fmt.Sprintf("%s USED  %s BUFF  %s FREE", types.Bytes(m.Used), m.Buffer, m.Free)
}

func (m *MEM) JSON() (string, error) {
	b, e := json.Marshal(m)
	return string(b), e
}
