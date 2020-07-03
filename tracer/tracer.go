package tracer

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Res struct {
	D time.Duration
	C int64
	S time.Duration
	M map[string]uint64
}

type CallBackFunc func(res map[string]*Res)

type Record struct {
	time.Duration
	count int64
	statusMap map[string]uint64
}

func (t *Record)Add(duration time.Duration,status string){
	t.count++
	t.Duration+=duration

	if status != "" {
		defer func() {
			recover()
			t.statusMap = make(map[string]uint64)
		}()

		if n,ok:=t.statusMap[status];ok {
			t.statusMap[status] = n + 1
		} else {
			t.statusMap[status] = 1
		}
	}
}

func (t *Record)Get()*Res{
	//t.Lock()
	//defer t.Unlock()

	var res =  &Res{
		S: t.Duration,
		C: t.count,
	}

	t.Duration = t.Duration/time.Duration(t.count)
	t.count = 1
	res.D = t.Duration

	if t.statusMap != nil {
		res.M = t.statusMap
		t.statusMap = make(map[string]uint64)
	}
	return res
}

type Action struct {
	Start time.Time
	End time.Time
	Name string
	Status string
}

type Tracer struct {
	tm      map[string]*Record
	actions chan *Action
	cbFunc  CallBackFunc
	actp    *sync.Pool
}

func NewTracer(cap int,cb CallBackFunc)*Tracer{
	t := &Tracer{
		make(map[string]*Record),
		make(chan *Action,cap),
		cb,
		&sync.Pool{New: func() interface{} {
			return &Action{}
		},
		},
	}

	go t.run()
	return t
}

func (t *Tracer)BeginAction(name string)*Action{
	act := t.actp.Get().(*Action)
	act.Start = time.Now()
	act.Name = name
	act.End = time.Now()
	act.Status = ""

	return act
}
func (t *Tracer)EndAction(act *Action)bool{
	act.End = time.Now()

	return t.PushAction(act)
}

func (t *Tracer)run()  {
	for{
		act:= <-t.actions
		switch  act.Name {
		case "_cb":
			t.callBack()
		default:
			record,ok:=t.tm[act.Name]
			if !ok {
				record = &Record{}
				t.tm[act.Name] = record
			}
			record.Add(act.End.Sub(act.Start),act.Status)

			t.actp.Put(act)
		}
	}
}

func (t *Tracer)callBack()  {
	var res = make(map[string]*Res)
	for k,v := range t.tm {
		res[k]=v.Get()
	}
	t.cbFunc(res)
}

func (t *Tracer)SendCallBackAction(){
	t.actions <-&Action{Name:"_cb"}
}

func (t *Tracer)PushAction(action *Action) bool {
	select {
	case t.actions <-action:
		return true
	default:
		return false
	}
}

var logger *log.Logger

var DefaultTracer *Tracer
var CBFunc = func(res map[string]*Res) {
	for k,v:= range res {
		buf := bytes.NewBufferString("")
		fmt.Fprintf(buf,"[%s] 平均耗时 %d ms 执行次数 %d 总耗时 %f s ",k,v.D.Milliseconds(),v.C,v.S.Seconds())
		if v.M != nil {
			fmt.Fprintf(buf," |状态:")
			for mk,mv:=range v.M {
				fmt.Fprintf(buf,"(%s:%d)",mk,mv)
			}
		}
		logger.Println(buf.String())
	}
	logger.Println("-------------------------------------------------------------------------")
}

func init()  {
	fl,err := os.OpenFile("tracer.log",os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0644)
	if err != nil {
		return
	}
	logger = log.New(fl,"[trace]",log.Ldate|log.Ltime)
	DefaultTracer = NewTracer(10000, CBFunc)
}
