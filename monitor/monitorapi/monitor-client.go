package monitorapi

//	"github.com/derekparker/delve/service/api"

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"context"
)

var (
	monitorConn *grpc.ClientConn
	monitorClt MonitorClient
)

//InitMonitor starts client
func InitMonitor(srv string){
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure()) // insecure
	var err error
	monitorConn, err = grpc.Dial(srv, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}else{
		log.Printf("Monitor Client Connected with %s",srv)
	}
	monitorClt = NewMonitorClient(monitorConn)
	log.Printf("Monitor Client is  %v",monitorClt)

}

func SendMes(mes *Mes){
	resp, err := monitorClt.SendReport(context.Background(),mes )

	if err != nil {
		log.Printf("Error in Sendmes %v",err)
	}else{
		if resp.Ok{
			log.Printf("Success! to send %v",mes)
		}
	}
}


func SendMessage(msgType string, chType int, src uint64, dst uint64, arg string){
	mes := &Mes{Msgtype:msgType, Chtype:int32(chType),Src: src,Dst: dst,Args:arg}
	resp, err := monitorClt.SendReport(context.Background(),mes )

	if err != nil {
		log.Printf("Error in Sendmes %v",err)
	}else{
		if resp.Ok{
			log.Printf("Success! to send %v",mes)
		}
	}
}

func (m *Mes)GetJson() string {
	s := fmt.Sprintf("{\"msgType\":\"%s\",\"chType\":%d,\"src\":%d,\"dst\":%d,\"arg\":\"%s\"}",
						m.Msgtype,m.Chtype, m.Src, m.Dst, m.Args)
	return s
}
