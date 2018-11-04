package sxutil

// SMUtil.go is a helper utility package for Synergic Market

// Helper structures for Synergic Market

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/fleet"
	"github.com/synerex/synerex_alpha/api/ptransit"
	"github.com/synerex/synerex_alpha/nodeapi"
)

// IDType for all ID in Synergic Market
type IDType uint64

var (
	node       *snowflake.Node // package variable for keeping unique ID.
	nid        *nodeapi.NodeID
	myNodeName string
	conn       *grpc.ClientConn
	clt        nodeapi.NodeClient
)

// DemandOpts is sender options for Demand
type DemandOpts struct {
	ID     uint64
	Target uint64
	Name   string
	JSON   string
}

// SupplyOpts is sender options for Supply
type SupplyOpts struct {
	ID        uint64
	Target    uint64
	Name      string
	JSON      string
	Fleet     *fleet.Fleet
	PTService *ptransit.PTService
}

func init() {
	fmt.Println("Synergic Market Util init() is called!")

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
	}
	return ni.NodeName
}

// RegisterNodeName is a function to register node name with node server address
func RegisterNodeName(nodesrv string, nm string, isServ bool) { // register ID to server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure()) // insecure
	var err error
	conn, err = grpc.Dial(nodesrv, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
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
		log.Fatalln("Error on get NodeID", ee)
	} else {

		var nderr error
		node, nderr = snowflake.NewNode(int64(nid.NodeId))
		if nderr != nil {
			fmt.Println("Error in initializing snowflake:", err)
		} else {
			fmt.Println("Successfully Initialize node ", nid.NodeId)
		}
	}
}

// UnRegisterNode de-registrate node id
func UnRegisterNode() {
	log.Println("UnRegister Node ", nid)
	resp, err := clt.UnRegisterNode(context.Background(), nid)
	if err != nil || !resp.Ok {
		log.Print("Can't unregister", err, resp)
	}
}

// SMServiceClient Wrappter Structure for market client
type SMServiceClient struct {
	ClientID IDType
	MType    api.MarketType
	Client   api.SMarketClient
	ArgJson  string
}

// NewSMServiceClient Creates wrapper structre SMServiceClient from SMarketClient
func NewSMServiceClient(clt api.SMarketClient, mtype api.MarketType, argJson string) *SMServiceClient {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.ProposeSupply(ctx, sp)
	if err != nil {
		log.Fatalf("%v.ProposeSupply err %v", clt, err)
	}
	log.Println("ProposeSupply Response:", resp)
	return pid
}

// SelectSupply send select message to server
func (clt *SMServiceClient) SelectSupply(sp *api.Supply) {
	tgt := &api.Target{
		Id:       GenerateIntID(),
		SenderId: uint64(sp.SenderId), // use senderId
		TargetId: sp.Id,
		Type:     sp.Type,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.SelectSupply(ctx, tgt)
	if err != nil {
		log.Fatalf("%v.SelectSupply err %v", clt, err)
	}
	log.Println("SelectSupply Response:", resp)
}

// SelectDemand send select message to server
func (clt *SMServiceClient) SelectDemand(dm *api.Demand) {
	tgt := &api.Target{
		Id:       GenerateIntID(),
		SenderId: uint64(dm.SenderId), // use senderId
		TargetId: dm.Id,
		Type:     dm.Type,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.SelectDemand(ctx, tgt)
	if err != nil {
		log.Fatalf("%v.SelectDemand err %v", clt, err)
	}
	log.Println("SelectDemand Response:", resp)
}

// SubscribeSupply  Wrapper function for SMServiceClient
func (clt *SMServiceClient) SubscribeSupply(ctx context.Context, spcb func(*SMServiceClient, *api.Supply)) {
	ch := clt.getChannel()
	smc, err := clt.Client.SubscribeSupply(ctx, ch)
	if err != nil {
		log.Fatalf("%v SubscribeSupply Error %v", clt, err)
	}
	for {
		sp, err := smc.Recv() // receive Demand
		if err != nil {
			if err == io.EOF {
				log.Print("End Supply subscribe OK")
			} else {
				log.Fatalf("%v SMServiceClient SubscribeSupply error %v", clt, err)
			}
			break
		}
		log.Printf("Receive Message %v", *sp)
		// call Callback!
		spcb(clt, sp)
	}
}

// SubscribeDemand  Wrapper function for SMServiceClient
func (clt *SMServiceClient) SubscribeDemand(ctx context.Context, dmcb func(*SMServiceClient, *api.Demand)) {
	ch := clt.getChannel()
	dmc, err := clt.Client.SubscribeDemand(ctx, ch)
	if err != nil {
		log.Fatalf("%v SubscribeDemand Error %v", clt, err)
	}
	for {
		dm, err := dmc.Recv() // receive Demand
		if err != nil {
			if err == io.EOF {
				log.Print("End Demand subscribe OK")
			} else {
				log.Fatalf("%v SMServiceClient SubscribeDemand error %v", clt, err)
			}
			break
		}
		log.Printf("Receive Message %v", *dm)
		// call Callback!
		dmcb(clt, dm)
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.RegisterDemand(ctx, &dm)
	if err != nil {
		log.Fatalf("%v.RegisterDemand err %v", clt, err)
	}
	log.Println(resp)
	dmo.ID = id // assign ID
	return id
}

// RegisterSupply sends Typed Supply to Server
func (clt *SMServiceClient) RegisterSupply(smo *SupplyOpts) uint64 {
	id := GenerateIntID()
	ts := ptypes.TimestampNow()
	dm := api.Supply{
		Id:         id,
		SenderId:   uint64(clt.ClientID),
		Type:       clt.MType,
		SupplyName: smo.Name,
		Ts:         ts,
		ArgJson:    smo.JSON,
	}

	switch clt.MType {
	case api.MarketType_RIDE_SHARE:
		sp := api.Supply_Arg_Fleet{
			smo.Fleet,
		}
		dm.ArgOneof = &sp
	case api.MarketType_PT_SERVICE:
		sp := api.Supply_Arg_PTService{
			smo.PTService,
		}
		dm.ArgOneof = &sp
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.RegisterSupply(ctx, &dm)
	if err != nil {
		log.Fatalf("Error for sending:RegisterSupply to  Synerex Server as %v", err)
	}
	log.Println("RegiterSupply:", smo, resp)
	smo.ID = id // assign ID
	return id
}

// Confirm sends confirm message to sender
func (clt *SMServiceClient) Confirm(id IDType) {
	tg := &api.Target{
		Id:       GenerateIntID(),
		SenderId: uint64(clt.ClientID),
		TargetId: uint64(id),
		Type:     clt.MType,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.Confirm(ctx, tg)
	if err != nil {
		log.Fatalf("%v Confirm Failier %v", clt, err)
	}
	log.Println("Confirm Success:", resp)
}
