package main

import (
	"fmt"
	//"math/rand"
	"testing"
	"time"
	. "github.com/dp1993132/niceTrace/tracer"
)

func TestNewTracer(t *testing.T) {
	tcer := NewTracer(3000, func(res map[string]time.Duration) {
		for k,v:=range res{
			fmt.Println(k,v)
		}
	})
	go func() {
		for {
			go func() {
				act:=tcer.BeginAction("test")
				<-time.After(time.Second)
				tcer.EndAction(act)
			}()
			<-time.After(time.Second)
		}
	}()
	for {
		<-time.After(time.Second * 5)
		tcer.SendCallBackAction()
	}
}

func TestSendMsg(t *testing.T)  {
	net.tc
}
