package statlog

import (
	"encoding/json"
	"fmt"
	"git.sophuwu.com/statlog/device"
	"strings"
	"time"
)

type HWInfo struct {
	Time int64      `json:"UnixMilli"`
	CPU  device.CPU `json:"CPU"`
	MEM  device.MEM `json:"MEM"`
}

func (this *HWInfo) Update() {
	done := make(chan bool)
	go func() {
		this.CPU.update()
		done <- true
	}()
	go func() {
		this.MEM.update()
		done <- true
	}()
	<-done
	<-done
	this.Time = time.Now().UnixMilli()
}
func (i HWInfo) String() string {
	if CONFIG.Json {
		b, _ := json.Marshal(i)
		return string(b)
	}
	s := i.MEM.String()
	if SI[CONFIG.Unit] == '%' {
		s = i.MEM.Percent().String()
	}
	return fmt.Sprintf("MEM: %s | CPU: %s", s, i.CPU.String())
}

// func main() {
// 	var hw HWInfo
// 	for do := true; do; do = CONFIG.Repeat {
// 		hw.Update()
// 		fmt.Println(hw)
// 	}
// }
