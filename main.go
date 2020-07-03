package main

var lisAddr string

func main(){
	//flag.StringVar(&lisAddr,"listen",":9005","监听地址")
	//flag.Parse()
	//
	//lis,err:=net.Listen("tcp4",lisAddr)
	//tra := tracer.NewTracer(10000, func(res *tracer.Res) {
	//	//for k,v:= range res {
	//	//	log.Printf("[%s] %f ms",k,v.Milliseconds())
	//	//}
	//})
	//if err !=nil {
	//	log.Println(err)
	//	return
	//}
	//
	//go func() {
	//	for  {
	//		conn,err:=lis.Accept()
	//		if err != nil {
	//			continue
	//		}
	//		go func(conn net.Conn) {
	//			dc := gob.NewDecoder(conn)
	//			for {
	//				var act tracer.Action
	//				dc.Decode(&act)
	//				tra.PushAction(&act)
	//			}
	//		}(conn)
	//	}
	//}()
	//
	//log.Println("tracer 启动 ",lisAddr)
	//for {
	//	<-time.After(time.Second)
	//	tra.SendCallBackAction()
	//}
}
