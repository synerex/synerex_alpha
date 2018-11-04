package main

//go:generate protoc -I ../api --go_out=paths=source_relative:../api common/common.proto
//go:generate protoc -I ../api --go_out=paths=source_relative:../api adservice/adservice.proto
//go:generate protoc -I ../api  --go_out=paths=source_relative:../api fleet/fleet.proto
//go:generate protoc -I ../api  --go_out=paths=source_relative:../api library/library.proto
//go:generate protoc -I ../api  --go_out=paths=source_relative:../api rideshare/rideshare.proto
//go:generate protoc -I ../api  --go_out=paths=source_relative:../api ptransit/ptransit.proto

//go:generate protoc -I ../api -I .. --go_out=plugins=grpc:../api smarket.proto

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"path"
	"sync"
	"time"

	"strconv"

	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/monitor/monitorapi"
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	port    = flag.Int("port", 10000, "The Synerex Server Listening Port")
	nodesrv = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	monitor = flag.String("monitor", "127.0.0.1:9998", "Monitor Server")
)

type synerexServerInfo struct {
	demandChans  [api.MarketType_END][]chan *api.Demand // create slices for each MarketType(each slice contains channels)
	supplyChans  [api.MarketType_END][]chan *api.Supply
	demandMap    map[sxutil.IDType]chan *api.Demand // map from IDtype to Demand channel
	supplyMap    map[sxutil.IDType]chan *api.Supply // map from IDtype to Supply channel
	waitConfirms map[sxutil.IDType]chan *api.Target // confirm maps
	mu           sync.RWMutex
	messageStore *MessageStore // message store
}

func init() {
	////	sxutil.InitNodeNum(0)
}

// Implementation of each Protocol API
func (s *synerexServerInfo) RegisterDemand(c context.Context, dm *api.Demand) (r *api.Response, e error) {
	// send demand for desired channels
	s.mu.RLock()
	chs := s.demandChans[dm.GetType()]
	for i := range chs {
		ch := chs[i]
		ch <- dm
	}
	s.mu.RUnlock()
	r = &api.Response{Ok: true, Err: ""}
	return r, nil
}

func (s *synerexServerInfo) RegisterSupply(c context.Context, sp *api.Supply) (r *api.Response, e error) {
	s.mu.RLock()
	chs := s.supplyChans[sp.GetType()]
	for i := range chs {
		ch := chs[i]
		ch <- sp
	}
	s.mu.RUnlock()
	r = &api.Response{Ok: true, Err: ""}
	return r, nil
}
func (s *synerexServerInfo) ProposeDemand(c context.Context, dm *api.Demand) (r *api.Response, e error) {
	s.mu.RLock()
	chs := s.demandChans[dm.GetType()]
	for i := range chs {
		ch := chs[i]
		ch <- dm
	}
	s.mu.RUnlock()
	r = &api.Response{Ok: true, Err: ""}
	return r, nil
}
func (s *synerexServerInfo) ProposeSupply(c context.Context, sp *api.Supply) (r *api.Response, e error) {
	s.mu.RLock()
	chs := s.supplyChans[sp.GetType()]
	for i := range chs {
		ch := chs[i]
		ch <- sp
	}
	s.mu.RUnlock()
	r = &api.Response{Ok: true, Err: ""}
	return r, nil
}
func (s *synerexServerInfo) ReserveSupply(c context.Context, tg *api.Target) (r *api.ConfirmResponse, e error) {
	//	chs := s.demandChans[tg.GetType()]
	//	dm := &api.Demand{}

	r = &api.ConfirmResponse{Ok: true, Err: ""}
	return r, nil
}

func (s *synerexServerInfo) SelectSupply(c context.Context, tg *api.Target) (r *api.ConfirmResponse, e error) {
	ch, ok := s.demandMap[sxutil.IDType(tg.GetSenderId())]
	if !ok {
		r = &api.ConfirmResponse{Ok: false, Err: "Can't find demand receiver from SelectSupply"}
		e = errors.New("Cant find channel in SelectSupply")
		return
	}
	id := sxutil.GenerateIntID()
	dm := &api.Demand{
		Id:       id, // generate ID from market server
		SenderId: tg.SenderId,
		TargetId: tg.TargetId,
		Type:     tg.Type,
	}
	ch <- dm // send select message

	tch := make(chan *api.Target)
	s.waitConfirms[sxutil.IDType(id)] = tch
	tb := <-tch // got confirm!
	if tb.TargetId == id {
		r = &api.ConfirmResponse{Ok: true, Err: ""}
		return r, nil
	} else {
		r = &api.ConfirmResponse{Ok: false, Err: "should not happen"}
		return r, errors.New("Should not happen")
	}
	// TODO: should check response from target .. umm.
}

func (s *synerexServerInfo) ReserveDemand(c context.Context, tg *api.Target) (r *api.ConfirmResponse, e error) {
	r = &api.ConfirmResponse{Ok: true, Err: ""}
	return r, nil
}
func (s *synerexServerInfo) SelectDemand(c context.Context, tg *api.Target) (r *api.ConfirmResponse, e error) {
	// select!

	r = &api.ConfirmResponse{Ok: true, Err: ""}
	return r, nil
}

func (s *synerexServerInfo) Confirm(c context.Context, tg *api.Target) (r *api.Response, e error) {
	// check waitConfirms
	ch, ok := s.waitConfirms[sxutil.IDType(tg.TargetId)]
	if !ok {
		r = &api.Response{Ok: false, Err: "Can't find channel"}
		return r, errors.New("can't find channels for Confirm")
	}
	ch <- tg // send OK
	r = &api.Response{Ok: true, Err: ""}
	return r, nil
}

// go routine which wait demand channel and sending demands to each providers.
func demandServerFunc(ch chan *api.Demand, stream api.SMarket_SubscribeDemandServer) {
	for {
		select {
		case sp := <-ch:
			err := stream.Send(sp)
			if err != nil {
				//				log.Printf("Error in DemandServer Error %v", err)
				return
			}
		}
	}
}

// remove channel from slice

func removeDemandChannelFromSlice(sl []chan *api.Demand, c chan *api.Demand) []chan *api.Demand {
	for i, ch := range sl {
		if ch == c {
			return append(sl[:i], sl[i+1:]...)
		}
	}
	log.Printf("Cant find channel %v in removeChannel", c)
	return nil
}

func removeSupplyChannelFromSlice(sl []chan *api.Supply, c chan *api.Supply) []chan *api.Supply {
	for i, ch := range sl {
		if ch == c {
			return append(sl[:i], sl[i+1:]...)
		}
	}
	log.Printf("Cant find channel %v in removeChannel", c)
	return nil
}

// SubscribeDemand is called form client to subscribe channel
func (s *synerexServerInfo) SubscribeDemand(ch *api.Channel, stream api.SMarket_SubscribeDemandServer) error {
	// TODO: we can check the duplication of node id here! (especially 1024 snowflake node ID)
	if _, ok := s.demandMap[sxutil.IDType(ch.GetClientId())]; ok { // check the availability of duplicated client ID
		return errors.New("duplicated ClientID")
	}
	//	log.Printf("SubscribeDemand %v",ch)
	//	log.Printf("SubscribeDemand stream  %v",stream)

	// It is better to logging here.
	//	monitorapi.SendMes(&monitorapi.Mes{Message:"Subscribe Demand", Args: fmt.Sprintf("Type:%d,From: %x  %s",ch.Type,ch.ClientId, ch.ArgJson )})
	monitorapi.SendMessage("SubscribeDemand", int(ch.Type), ch.ClientId, 0, ch.ArgJson)

	subCh := make(chan *api.Demand, 10)
	// We should think about thread safe coding.
	tp := ch.GetType()
	idt := sxutil.IDType(ch.GetClientId())
	s.mu.Lock()
	s.demandChans[tp] = append(s.demandChans[tp], subCh)
	s.demandMap[idt] = subCh // mapping from clientID to channel
	s.mu.Unlock()
	demandServerFunc(subCh, stream) // infinite go routine?
	// if this returns, stream might be closed.
	// we should remove channel
	s.mu.Lock()
	delete(s.demandMap, idt) // remove map from idt
	s.demandChans[tp] = removeDemandChannelFromSlice(s.demandChans[tp], subCh)
	log.Printf("Remove Demand Stream Channel %v", ch)
	s.mu.Unlock()
	return nil
}

func supplyServerFunc(ch chan *api.Supply, stream api.SMarket_SubscribeSupplyServer) {
	for {
		select {
		case sp := <-ch:
			err := stream.Send(sp)
			if err != nil {
				//				log.Printf("Error SupplyServer Error %v", err)
				return
			}
		}
	}
}

func (s *synerexServerInfo) SubscribeSupply(ch *api.Channel, stream api.SMarket_SubscribeSupplyServer) error {
	subCh := make(chan *api.Supply, 10)
	tp := ch.GetType()

	//	monitorapi.SendMes(&monitorapi.Mes{Message:"Subscribe Supply", Args: fmt.Sprintf("Type:%d, From: %x %s",ch.Type,ch.ClientId,ch.ArgJson )})
	monitorapi.SendMessage("SubscribeSupply", int(ch.Type), ch.ClientId, 0, ch.ArgJson)

	s.mu.Lock()
	s.supplyChans[tp] = append(s.supplyChans[tp], subCh)
	s.mu.Unlock()
	supplyServerFunc(subCh, stream)
	// this supply stream may closed. so take care.

	s.mu.Lock()
	s.supplyChans[tp] = removeSupplyChannelFromSlice(s.supplyChans[tp], subCh)
	log.Printf("Remove Supply Stream Channel %v", ch)
	s.mu.Unlock()
	return nil
}

func newServerInfo() *synerexServerInfo {
	var ms synerexServerInfo
	s := &ms
	s.demandMap = make(map[sxutil.IDType]chan *api.Demand)
	s.supplyMap = make(map[sxutil.IDType]chan *api.Supply)
	s.waitConfirms = make(map[sxutil.IDType]chan *api.Target)

	s.messageStore = CreateLocalMessageStore()

	return s
}

var (
	NodeBits uint8 = 10
	StepBits uint8 = 12

	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask  int64 = nodeMax << StepBits
	nodeShift uint8 = StepBits
	nodeMap         = make(map[int]string)
)

func idToNode(id uint64) string {
	nodeNum := int(int64(id) & nodeMask >> nodeShift)
	var ok bool
	var str string
	if str, ok = nodeMap[nodeNum]; !ok {
		str = sxutil.GetNodeName(nodeNum)
	}
	return str + ":" + strconv.Itoa(nodeNum)
}

func unaryServerInterceptor(logger *logrus.Logger, s *synerexServerInfo) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		var args string
		var msgType int
		var srcId, tgtId, mid uint64
		method := path.Base(info.FullMethod)
		switch method {
		// Demand
		case "RegisterDemand", "ProposeDemand":
			dm := req.(*api.Demand)
			msgType = int(dm.Type)
			srcId = dm.SenderId
			tgtId = dm.TargetId
			mid = dm.Id
			args = "Type:" + strconv.Itoa(int(dm.Type)) + ":" + strconv.FormatUint(dm.Id, 16) + ":" + idToNode(dm.SenderId) + "->" + strconv.FormatUint(dm.TargetId, 16)
			// Supply
		case "RegisterSupply", "ProposeSupply":
			sp := req.(*api.Supply)
			msgType = int(sp.Type)
			srcId = sp.SenderId
			tgtId = sp.TargetId
			mid = sp.Id
			args = "Type:" + strconv.Itoa(int(sp.Type)) + ":" + strconv.FormatUint(sp.Id, 16) + ":" + idToNode(sp.SenderId) + "->" + strconv.FormatUint(sp.TargetId, 16)
			// Target
		case "SelectSupply", "Confirm", "SelectDemand":
			tg := req.(*api.Target)
			msgType = int(tg.Type)
			mid = tg.Id
			srcId = tg.SenderId
			tgtId = tg.TargetId
			args = "Type:" + strconv.Itoa(int(tg.Type)) + ":" + strconv.FormatUint(tg.Id, 16) + ":" + idToNode(tg.Id) + "->" + strconv.FormatUint(tg.TargetId, 16)
		}

		//		monitorapi.SendMes(&monitorapi.Mes{Message:method+":"+args, Args:""})

		dstId := s.messageStore.getSrcId(tgtId)
		monitorapi.SendMessage(method, msgType, srcId, dstId, args)

		// register for messageStore
		s.messageStore.AddMessage(method, msgType, mid, srcId, dstId, args)

		// Obtain log using defer
		defer func(begin time.Time) {
			// Obtain method name from info
			method := path.Base(info.FullMethod)
			took := time.Since(begin)
			fields := logrus.Fields{
				"method": method,
				"took":   took,
			}
			if err != nil {
				fields["error"] = err
				logger.WithFields(fields).Error("Failed")
			} else {
				logger.WithFields(fields).Info("Succeeded")
			}
		}(time.Now())

		// handler = RPC method
		reply, hErr := handler(ctx, req)
		if hErr != nil {
			err = hErr
		}

		return reply, err
	}
}

// Stream Interceptor
func streamServerInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var err error
		//		var args string
		log.Printf("streamserver intercept...")
		method := path.Base(info.FullMethod)
		switch method {
		case "SubscribeDemand":
		case "SubscribeSupply":
		}
		//		monitorapi.SendMes(&monitorapi.Mes{Message:method, Args:args})

		defer func(begin time.Time) {
			// Obtain method name from info
			method := path.Base(info.FullMethod)
			took := time.Since(begin)
			fields := logrus.Fields{
				"method": method,
				"took":   took,
			}
			if err != nil {
				fields["error"] = err
				logger.WithFields(fields).Error("Failed")
			} else {
				logger.WithFields(fields).Info("Succeeded")
			}
		}(time.Now())

		// handler = RPC method
		if hErr := handler(srv, stream); err != nil {
			err = hErr
		}
		log.Printf("streamserver intercept..end .")
		return err
	}
}

func prepareGrpcServer(s *synerexServerInfo, opts ...grpc.ServerOption) *grpc.Server {
	gcServer := grpc.NewServer(opts...)
	api.RegisterSMarketServer(gcServer, s)
	return gcServer
}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "SynerexServer", true)

	monitorapi.InitMonitor(*monitor)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	logger := logrus.New()

	s := newServerInfo()
	opts = append(opts, grpc.UnaryInterceptor(unaryServerInterceptor(logger, s)))

	// for more precise monitoring , we do not use StreamIntercepter.
	//	opts = append(opts, grpc.StreamInterceptor(streamServerInterceptor(logger)))

	grpcServer := prepareGrpcServer(s, opts...)
	log.Printf("Start Synergic Exchange Server, connection waiting at port :%d ...", *port)
	grpcServer.Serve(lis)

}
