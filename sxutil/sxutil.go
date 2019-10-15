package sxutil

// SMUtil.go is a helper utility package for Synergic Market

// Helper structures for Synergic Market

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/synerex/synerex_alpha/api/rideshare"
	"github.com/synerex/synerex_alpha/api/routing"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/participant"
	"github.com/synerex/synerex_alpha/api/simulation/route"

	"github.com/bwmarrin/snowflake"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/fleet"
	"github.com/synerex/synerex_alpha/api/ptransit"
	"github.com/synerex/synerex_alpha/nodeapi"
)

// IDType for all ID in Synergic Exchange
type IDType uint64

var (
	node       *snowflake.Node // package variable for keeping unique ID.
	nid        *nodeapi.NodeID
	nupd       *nodeapi.NodeUpdate
	numu       sync.RWMutex
	myNodeName string
	conn       *grpc.ClientConn
	clt        nodeapi.NodeClient
)

// DemandOpts is sender options for Demand
type DemandOpts struct {
	ID             uint64
	Target         uint64
	Name           string
	JSON           string
	Fleet          *fleet.Fleet
	RoutingService *routing.RoutingService
	RideShare      *rideshare.RideShare

	SetClockDemand       *clock.SetClockDemand
	ForwardClockDemand   *clock.ForwardClockDemand
	BackClockDemand      *clock.BackClockDemand
	GetAreaDemand        *area.GetAreaDemand
	GetAgentsDemand      *agent.GetAgentsDemand
	SetAgentsDemand      *agent.SetAgentsDemand
	GetParticipantDemand *participant.GetParticipantDemand
	SetParticipantDemand *participant.SetParticipantDemand
	GetRouteDemand       *route.GetRouteDemand
}

// SupplyOpts is sender options for Supply
type SupplyOpts struct {
	ID             uint64
	Target         uint64
	Name           string
	JSON           string
	Fleet          *fleet.Fleet
	PTService      *ptransit.PTService
	RoutingService *routing.RoutingService
	RideShare      *rideshare.RideShare

	SetClockSupply       *clock.SetClockSupply
	ForwardClockSupply   *clock.ForwardClockSupply
	BackClockSupply      *clock.BackClockSupply
	GetAreaSupply        *area.GetAreaSupply
	GetAgentsSupply      *agent.GetAgentsSupply
	SetAgentsSupply      *agent.SetAgentsSupply
	ForwardAgentsSupply  *agent.ForwardAgentsSupply
	GetParticipantSupply *participant.GetParticipantSupply
	SetParticipantSupply *participant.SetParticipantSupply
	GetRouteSupply       *route.GetRouteSupply
}

func init() {
	fmt.Println("Synergic Exchange Util init() is called!")
}

// InitNodeNum for initialize NodeNum again
func InitNodeNum(n int) {
	var err error
	node, err = snowflake.NewNode(int64(n))
	if err != nil {
		fmt.Println("Error in initializing snowflake:", err)
	} else {
		fmt.Println("Successfully Initialize node ", n)
	}
}

func GetNodeName(n int) string {
	ni, err := clt.QueryNode(context.Background(), &nodeapi.NodeID{NodeId: int32(n)})
	if err != nil {
		log.Printf("Error on QueryNode %v", err)
		return "Unknown"
	}
	return ni.NodeName
}

func SetNodeStatus(status int32, arg string) {
	numu.Lock()
	nupd.NodeStatus = status
	nupd.NodeArg = arg
	numu.Unlock()
}

func startKeepAlive() {
	for {
		//		fmt.Printf("KeepAlive %s %d\n",nupd.NodeStatus, nid.KeepaliveDuration)
		time.Sleep(time.Second * time.Duration(nid.KeepaliveDuration))
		if nid.Secret == 0 { // this means the node is disconnected
			break
		}
		numu.RLock()
		nupd.UpdateCount++
		clt.KeepAlive(context.Background(), nupd)
		numu.RUnlock()
	}
}

// RegisterNodeName is a function to register node name with node server address
func RegisterNodeName(nodesrv string, nm string, isServ bool) error { // register ID to server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure()) // insecure
	var err error
	conn, err = grpc.Dial(nodesrv, opts...)
	if err != nil {
		log.Printf("fail to dial: %v", err)
		return err
	}
	//	defer conn.Close()

	clt = nodeapi.NewNodeClient(conn)
	nif := nodeapi.NodeInfo{
		NodeName: nm,
		IsServer: isServ,
	}
	myNodeName = nm
	var ee error
	nid, ee = clt.RegisterNode(context.Background(), &nif)

	if ee != nil { // has error!
		log.Println("Error on get NodeID", ee)
		return ee
	} else {
		var nderr error
		node, nderr = snowflake.NewNode(int64(nid.NodeId))
		if nderr != nil {
			fmt.Println("Error in initializing snowflake:", err)
			return nderr
		} else {
			fmt.Println("Successfully Initialize node ", nid.NodeId)
		}
	}

	nupd = &nodeapi.NodeUpdate{
		NodeId:      nid.NodeId,
		Secret:      nid.Secret,
		UpdateCount: 0,
		NodeStatus:  0,
		NodeArg:     "",
	}
	// start keepalive goroutine
	go startKeepAlive()
	//	fmt.Println("KeepAlive started!")
	return nil
}

// UnRegisterNode de-registrate node id
func UnRegisterNode() {
	log.Println("UnRegister Node ", nid)
	resp, err := clt.UnRegisterNode(context.Background(), nid)
	nid.Secret = 0
	if err != nil || !resp.Ok {
		log.Print("Can't unregister", err, resp)
	}
}

// SMServiceClient Wrappter Structure for market client
type SMServiceClient struct {
	ClientID IDType
	MType    api.ChannelType
	Client   api.SynerexClient
	ArgJson  string
	MbusID   IDType
}

// NewSMServiceClient Creates wrapper structre SMServiceClient from SynerexClient
func NewSMServiceClient(clt api.SynerexClient, mtype api.ChannelType, argJson string) *SMServiceClient {
	s := &SMServiceClient{
		ClientID: IDType(node.Generate()),
		MType:    mtype,
		Client:   clt,
		ArgJson:  argJson,
	}
	return s
}

// GenerateIntID for generate uniquie ID
func GenerateIntID() uint64 {
	return uint64(node.Generate())
}

func (clt SMServiceClient) getChannel() *api.Channel {
	return &api.Channel{ClientId: uint64(clt.ClientID), Type: clt.MType, ArgJson: clt.ArgJson}
}

// IsSupplyTarget is a helper function to check target
func (clt *SMServiceClient) IsSupplyTarget(sp *api.Supply, idlist []uint64) bool {
	spid := sp.TargetId
	for _, id := range idlist {
		if id == spid {
			return true
		}
	}
	return false
}

// IsDemandTarget is a helper function to check target
func (clt *SMServiceClient) IsDemandTarget(dm *api.Demand, idlist []uint64) bool {
	dmid := dm.TargetId
	for _, id := range idlist {
		if id == dmid {
			return true
		}
	}
	return false
}

// ProposeSupply send proposal Supply message to server
func (clt *SMServiceClient) ProposeSupply(spo *SupplyOpts) uint64 {
	pid := GenerateIntID()
	sp := &api.Supply{
		Id:         pid,
		SenderId:   uint64(clt.ClientID),
		TargetId:   spo.Target,
		Type:       clt.MType,
		SupplyName: spo.Name,
		ArgJson:    spo.JSON,
	}
	//	log.Printf("spo Target %v", sp)

	switch clt.MType {
	case api.ChannelType_PARTICIPANT_SERVICE:
		if spo.GetParticipantSupply != nil {
			sp.WithGetParticipantSupply(spo.GetParticipantSupply)
		}
		if spo.SetParticipantSupply != nil {
			sp.WithSetParticipantSupply(spo.SetParticipantSupply)
		}

	case api.ChannelType_CLOCK_SERVICE:
		if spo.SetClockSupply != nil {
			sp.WithSetClockSupply(spo.SetClockSupply)
		}
		if spo.ForwardClockSupply != nil {
			sp.WithForwardClockSupply(spo.ForwardClockSupply)
		}
		if spo.BackClockSupply != nil {
			sp.WithBackClockSupply(spo.BackClockSupply)
		}

	case api.ChannelType_AREA_SERVICE:
		if spo.GetAreaSupply != nil {
			sp.WithGetAreaSupply(spo.GetAreaSupply)
		}

	case api.ChannelType_ROUTE_SERVICE:
		if spo.GetRouteSupply != nil {
			sp.WithGetRouteSupply(spo.GetRouteSupply)
		}

	case api.ChannelType_AGENT_SERVICE:
		if spo.GetAgentsSupply != nil {
			sp.WithGetAgentsSupply(spo.GetAgentsSupply)
		}
		if spo.SetAgentsSupply != nil {
			sp.WithSetAgentsSupply(spo.SetAgentsSupply)
		}
		if spo.ForwardAgentsSupply != nil {
			sp.WithForwardAgentsSupply(spo.ForwardAgentsSupply)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.ProposeSupply(ctx, sp)
	if err != nil {
		log.Printf("%v.ProposeSupply err %v, [%v]", clt, err, sp)
		return 0 // should check...
	}
	log.Println("ProposeSupply Response:", resp, ":PID ", pid)
	return pid
}

// SelectSupply send select message to server
func (clt *SMServiceClient) SelectSupply(sp *api.Supply) (uint64, error) {
	tgt := &api.Target{
		Id:       GenerateIntID(),
		SenderId: uint64(clt.ClientID),
		TargetId: sp.Id, /// Message Id of Supply (not SenderId),
		Type:     sp.Type,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.SelectSupply(ctx, tgt)
	if err != nil {
		log.Printf("%v.SelectSupply err %v", clt, err)
		return 0, err
	}
	log.Println("SelectSupply Response:", resp)
	// if mbus is OK, start mbus!
	clt.MbusID = IDType(resp.MbusId)
	if clt.MbusID != 0 {
		//		clt.SubscribeMbus()
	}

	return uint64(clt.MbusID), nil
}

// SelectDemand send select message to server
func (clt *SMServiceClient) SelectDemand(dm *api.Demand) error {
	tgt := &api.Target{
		Id:       GenerateIntID(),
		SenderId: uint64(clt.ClientID),
		TargetId: dm.Id,
		Type:     dm.Type,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.SelectDemand(ctx, tgt)
	if err != nil {
		log.Printf("%v.SelectDemand err %v", clt, err)
		return err
	}
	log.Println("SelectDemand Response:", resp)
	return nil
}

// SubscribeSupply  Wrapper function for SMServiceClient
func (clt *SMServiceClient) SubscribeSupply(ctx context.Context, spcb func(*SMServiceClient, *api.Supply), wg *sync.WaitGroup) error {
	ch := clt.getChannel()
	smc, err := clt.Client.SubscribeSupply(ctx, ch)
	log.Printf("Test3 %v", ch)
	wg.Done()
	if err != nil {
		log.Printf("%v SubscribeSupply Error %v", clt, err)
		return err
	}
	for {
		var sp *api.Supply
		sp, err = smc.Recv() // receive Demand
		if err != nil {
			if err == io.EOF {
				log.Print("End Supply subscribe OK")
			} else {
				log.Printf("%v SMServiceClient SubscribeSupply error [%v]", clt, err)
			}
			break
		}
		//		log.Println("Receive SS:", sp)
		// call Callback!
		spcb(clt, sp)
	}
	return err
}

// SubscribeDemand  Wrapper function for SMServiceClient
func (clt *SMServiceClient) SubscribeDemand(ctx context.Context, dmcb func(*SMServiceClient, *api.Demand), wg *sync.WaitGroup) error {
	ch := clt.getChannel()
	dmc, err := clt.Client.SubscribeDemand(ctx, ch)
	log.Printf("Test3 %v", ch)
	wg.Done()
	if err != nil {
		log.Printf("%v SubscribeDemand Error %v", clt, err)
		return err // sender should handle error...
	}
	for {
		var dm *api.Demand
		dm, err = dmc.Recv() // receive Demand
		if err != nil {
			if err == io.EOF {
				log.Print("End Demand subscribe OK")
			} else {
				log.Printf("%v SMServiceClient SubscribeDemand error [%v]", clt, err)
			}
			break
		}
		//		log.Println("Receive SD:",*dm)
		// call Callback!
		dmcb(clt, dm)
	}
	return err
}

// SubscribeMbus  Wrapper function for SMServiceClient
func (clt *SMServiceClient) SubscribeMbus(ctx context.Context, mbcb func(*SMServiceClient, *api.MbusMsg)) error {

	mb := &api.Mbus{
		ClientId: uint64(clt.ClientID),
		MbusId:   uint64(clt.MbusID),
	}

	smc, err := clt.Client.SubscribeMbus(ctx, mb)
	if err != nil {
		log.Printf("%v Synerex_SubscribeMbusClient Error %v", clt, err)
		return err // sender should handle error...
	}
	for {
		var mes *api.MbusMsg
		mes, err = smc.Recv() // receive Demand
		if err != nil {
			if err == io.EOF {
				log.Print("End Mbus subscribe OK")
			} else {
				log.Printf("%v SMServiceClient SubscribeMbus error %v", clt, err)
			}
			break
		}
		//		log.Printf("Receive Mbus Message %v", *mes)
		// call Callback!
		mbcb(clt, mes)
	}
	return err
}

func (clt *SMServiceClient) SendMsg(ctx context.Context, msg *api.MbusMsg) error {
	if clt.MbusID == 0 {
		return errors.New("No Mbus opened!")
	}
	msg.MsgId = GenerateIntID()
	msg.SenderId = uint64(clt.ClientID)
	msg.MbusId = uint64(clt.MbusID)
	_, err := clt.Client.SendMsg(ctx, msg)

	return err
}

func (clt *SMServiceClient) CloseMbus(ctx context.Context) error {
	if clt.MbusID == 0 {
		return errors.New("No Mbus opened!")
	}
	mbus := &api.Mbus{
		ClientId: uint64(clt.ClientID),
		MbusId:   uint64(clt.MbusID),
	}
	_, err := clt.Client.CloseMbus(ctx, mbus)
	if err == nil {
		clt.MbusID = 0
	}
	return err
}

// RegisterDemand sends Typed Demand to Server
func (clt *SMServiceClient) RegisterDemand(dmo *DemandOpts) uint64 {
	id := GenerateIntID()
	ts := ptypes.TimestampNow()
	dm := api.Demand{
		Id:         id,
		SenderId:   uint64(clt.ClientID),
		Type:       clt.MType,
		DemandName: dmo.Name,
		Ts:         ts,
		ArgJson:    dmo.JSON,
	}

	switch clt.MType {
	case api.ChannelType_PARTICIPANT_SERVICE:
		if dmo.GetParticipantDemand != nil {
			dm.WithGetParticipantDemand(dmo.GetParticipantDemand)
		}
		if dmo.SetParticipantDemand != nil {
			dm.WithSetParticipantDemand(dmo.SetParticipantDemand)
		}

	case api.ChannelType_CLOCK_SERVICE:
		if dmo.SetClockDemand != nil {
			dm.WithSetClockDemand(dmo.SetClockDemand)
		}
		if dmo.ForwardClockDemand != nil {
			dm.WithForwardClockDemand(dmo.ForwardClockDemand)
		}
		if dmo.BackClockDemand != nil {
			dm.WithBackClockDemand(dmo.BackClockDemand)
		}

	case api.ChannelType_AREA_SERVICE:
		if dmo.GetAreaDemand != nil {
			dm.WithGetAreaDemand(dmo.GetAreaDemand)
		}

	case api.ChannelType_ROUTE_SERVICE:
		if dmo.GetRouteDemand != nil {
			dm.WithGetRouteDemand(dmo.GetRouteDemand)
		}

	case api.ChannelType_AGENT_SERVICE:
		if dmo.GetAgentsDemand != nil {
			dm.WithGetAgentsDemand(dmo.GetAgentsDemand)
		}
		if dmo.SetAgentsDemand != nil {
			dm.WithSetAgentsDemand(dmo.SetAgentsDemand)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//	log.Printf("Now RegisterDemand with %v",dm)
	//	log.Printf("Arg %v",dm.ArgOneof)
	//	log.Printf("Arg %v",dm.GetArg_RideShare())

	_, err := clt.Client.RegisterDemand(ctx, &dm)

	//	resp, err := clt.Client.RegisterDemand(ctx, &dm)
	if err != nil {
		log.Printf("%v.RegisterDemand err %v", clt, err)
		return 0
	}
	//	log.Println(resp)
	dmo.ID = id // assign ID
	return id
}

// RegisterSupply sends Typed Supply to Server
func (clt *SMServiceClient) RegisterSupply(spo *SupplyOpts) uint64 {
	id := GenerateIntID()
	ts := ptypes.TimestampNow()
	sp := api.Supply{
		Id:         id,
		SenderId:   uint64(clt.ClientID),
		Type:       clt.MType,
		SupplyName: spo.Name,
		Ts:         ts,
		ArgJson:    spo.JSON,
	}

	switch clt.MType {
	case api.ChannelType_PARTICIPANT_SERVICE:
		if spo.GetParticipantSupply != nil {
			sp.WithGetParticipantSupply(spo.GetParticipantSupply)
		}
		if spo.SetParticipantSupply != nil {
			sp.WithSetParticipantSupply(spo.SetParticipantSupply)
		}

	case api.ChannelType_CLOCK_SERVICE:
		if spo.SetClockSupply != nil {
			sp.WithSetClockSupply(spo.SetClockSupply)
		}
		if spo.ForwardClockSupply != nil {
			sp.WithForwardClockSupply(spo.ForwardClockSupply)
		}
		if spo.BackClockSupply != nil {
			sp.WithBackClockSupply(spo.BackClockSupply)
		}

	case api.ChannelType_AREA_SERVICE:
		if spo.GetAreaSupply != nil {
			sp.WithGetAreaSupply(spo.GetAreaSupply)
		}

	case api.ChannelType_ROUTE_SERVICE:
		if spo.GetRouteSupply != nil {
			sp.WithGetRouteSupply(spo.GetRouteSupply)
		}

	case api.ChannelType_AGENT_SERVICE:
		if spo.GetAgentsSupply != nil {
			sp.WithGetAgentsSupply(spo.GetAgentsSupply)
		}
		if spo.SetAgentsSupply != nil {
			sp.WithSetAgentsSupply(spo.SetAgentsSupply)
		}
		if spo.ForwardAgentsSupply != nil {
			sp.WithForwardAgentsSupply(spo.ForwardAgentsSupply)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//	resp , err := clt.Client.RegisterSupply(ctx, &dm)

	_, err := clt.Client.RegisterSupply(ctx, &sp)
	if err != nil {
		log.Printf("Error for sending:RegisterSupply to  Synerex Server as %v ", err)
		return 0
	}
	//	log.Println("RegiterSupply:", smo, resp)
	spo.ID = id // assign ID
	return id
}

// Confirm sends confirm message to sender
func (clt *SMServiceClient) Confirm(id IDType) error {
	tg := &api.Target{
		Id:       GenerateIntID(),
		SenderId: uint64(clt.ClientID),
		TargetId: uint64(id),
		Type:     clt.MType,
		MbusId:   uint64(id),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.Confirm(ctx, tg)
	if err != nil {
		log.Printf("%v Confirm Failier %v", clt, err)
		return err
	}
	clt.MbusID = id
	log.Println("Confirm Success:", resp)
	return nil
}
