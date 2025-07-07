package statlog

import (
	"encoding/json"
	"fmt"
	"git.sophuwu.com/statlog/device"
	"time"
)

type HWInfo struct {
	Time int64      `json:"UnixMilli"`
	CPU  device.CPU `json:"CPU"`
	MEM  device.MEM `json:"MEM"`
}

func (hw *HWInfo) Update() {
	hw.Time = time.Now().UnixMilli()
	done := make(chan bool)
	go func() {
		hw.CPU.Update()
		done <- true
	}()
	go func() {
		hw.MEM.Update()
		done <- true
	}()
	<-done
	<-done
}
func (hw *HWInfo) String() string {
	return fmt.Sprintf("%s | %s\n", hw.MEM.String(), hw.CPU.String())
}
func (hw *HWInfo) JSON() (string, error) {
	b, e := json.MarshalIndent(hw, "", "  ")
	return string(b), e
}
