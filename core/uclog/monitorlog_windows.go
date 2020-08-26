package uclog

import (
	"encoding/json"
	"fmt"
)

type Monitorlog struct {
	Fields map[string]interface{}
}

func NewMonitorlog() *Monitorlog {
	l := &Monitorlog{}
	l.Fields = make(map[string]interface{})
	return l
}

func (this *Monitorlog) AddField(name string, value interface{}) {
	this.Fields[name] = value
}

func (this *Monitorlog) RemoveField(name string) {
}

func (this *Monitorlog) Done() {
	b, _ := json.Marshal(this.Fields)
	fmt.Println("monitorlog: ", string(b))
}

func (this *Monitorlog) DoneNew(v interface{}) {
	b, _ := json.Marshal(v)
	fmt.Println("monitorlog: ", string(b))
}

func (this *Monitorlog) DoneWithPrefix(v interface{}, prefix string) {
	b, _ := json.Marshal(v)
	fmt.Println("monitorlog:%s %s", prefix, string(b))
}
