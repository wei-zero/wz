package ticker

import (
	"bytes"
	"fmt"
	"time"
)

type Ticker struct {
	lastTime int64
	ticks    []tick
}

type tick struct {
	msg  string
	time int64
}

func (m *Ticker) Tick(msg string) {
	if m.ticks == nil {
		m.ticks = []tick{}
	}

	m.ticks = append(m.ticks, tick{
		msg:  msg,
		time: time.Now().UnixNano(),
	})
}

// TimeCostMs return time cost in ms
func (m *Ticker) TimeCostMs(index int) int64 {
	if index <= 0 {
		return 0
	}
	if index >= len(m.ticks) {
		return -1
	}
	return (m.ticks[index].time - m.ticks[index-1].time) / 1000000
}

func (m *Ticker) String() string {
	var buffer bytes.Buffer
	for i := 1; i < len(m.ticks); i++ {
		buffer.WriteString(fmt.Sprintf("%s: %dms ", m.ticks[i].msg, (m.ticks[i].time-m.ticks[i-1].time)/1000000))
	}
	return buffer.String()
}

func (m *Ticker) StringNano() string {
	var buffer bytes.Buffer
	for i := 1; i < len(m.ticks); i++ {
		buffer.WriteString(fmt.Sprintf("%s: %dns ", m.ticks[i].msg, (m.ticks[i].time - m.ticks[i-1].time)))
	}
	return buffer.String()
}
