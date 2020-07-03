package main

import (
	"encoding/gob"
	"flag"
	tracer "github.com/dp1993132/niceTrace/tracer"
	"log"
	"net"
	"time"
)

var lisAddr string

func main(){
	flag.StringVar(&lisAddr,"listen",":9005","监听地址")
	flag.Parse()

	lis,err:=net.Listen("tcp4",lisAddr)
	tra := tracer.NewTracer(10000, func(res map[string]time.Duration) {
		for k,v:= range res {
			log.Printf("[%s] %f ms",k,v.Milliseconds())
		}
	})
	if err !=nil {
		log.Println(err)
		return
	}

	go func() {
		for  {
			conn,err:=lis.Accept()
			if err != nil {
				continue
			}
			go func(conn net.Conn) {
				dc := gob.NewDecoder(conn)
				for {
					var act tracer.Action
					dc.Decode(&act)
					tra.PushAction(&act)
				}
			}(conn)
		}
	}()

	log.Println("tracer 启动 ",lisAddr)
	for {
		<-time.After(time.Second)
		tra.SendCallBackAction()
	}
}
