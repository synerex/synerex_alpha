package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/google/gops/agent"
	"log"
	"math/rand"
	"net"
	"sort"
	"sync"
	"time"

	nodepb "github.com/synerex/synerex_alpha/nodeapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

//go:generate protoc -I ../nodeapi --go_out=paths=source_relative,plugins=grpc:../nodeapi ../nodeapi/nodeid.proto

// NodeID Server for  keep all node ID
//    node ID = 0-1023. (less than 10 is for server)
// When we use sxutil, we need to support nodenum

// Function
//   register node (in the future authentication..)

// shuold use only at here

// MaxNodeNum  Max node Number
const MaxNodeNum = 1024

// MaxServerID  Max Market Server Node ID (Small number ID is for smarket server)
const MaxServerID = 10
const DefaultDuration int32 = 10   // need keepalive for each 10 sec.
const MaxDurationCount = 3   // duration count.

type eachNodeInfo struct {
	nodeName string
	secret   uint64
	address  string
	lastAlive time.Time
	count 	int32
	status 	int32
	arg		string
	duration int32     // duration for checking next time
}

type srvNodeInfo struct {
	nodeMap map[int32]*eachNodeInfo // map from nodeID to eachNodeInfo
}

var (
	port     = flag.Int("port", 9990, "NodeID Server Listening Port")
	srvInfo  srvNodeInfo
	lastNode int32 = MaxServerID // start ID from MAX_SERVER_ID to MAX_NODE_NUM
	lastPrint time.Time
	nmmu sync.RWMutex
)

func init() {
	log.Println("Starting Node ID Server..")
	rand.Seed(time.Now().UnixNano())
	s := &srvInfo
	s.nodeMap = make(map[int32]*eachNodeInfo)
	lastPrint = time.Now()
	go keepNodes(s)
}

// find unused ID from map.
func getNextNodeID(sv bool) int32 {
	var n int32
	if sv {
		n = 0
	} else {
		n = lastNode
	}
	nmmu.RLock()
	for {
		_, ok := srvInfo.nodeMap[n]
		if !ok {
			break
		}
		if sv {
			n = (n + 1) % MaxServerID
		} else {
			n = (n-9)%(MaxNodeNum-MaxServerID) + MaxServerID
		}
		if n == lastNode || n == 0 { // loop
			nmmu.RUnlock()
			return -1 // all id is full...
		}
	}
	nmmu.RUnlock()
	if !sv {
		lastNode = n
	}
	return n
}

func keepNodes(s *srvNodeInfo){
	for {
		time.Sleep(time.Second * time.Duration(DefaultDuration))
		killNodes := make([]int32,0)
		nmmu.Lock()
		for k, eni := range s.nodeMap {
			sub := time.Now().Sub(eni.lastAlive) / time.Second
			if sub > time.Duration(MaxDurationCount*DefaultDuration){
				killNodes = append(killNodes,k)
			}
		}
		for _, k := range killNodes {
			delete(s.nodeMap,k)
		}
		nmmu.Unlock()
	}
}


// display all node info
func (s *srvNodeInfo) listNodes() {
	nmmu.RLock()
	nk := make([]int32, len(s.nodeMap))
	i :=0
	for k := range s.nodeMap {
		nk[i] = k
		i++
	}
	sort.Slice(nk, func(i,j int) bool {return nk[i] < nk[j]})
	for i := range nk {
		eni := s.nodeMap[nk[i]]
		sub := time.Now().Sub(eni.lastAlive)/time.Second
		log.Printf("ID: %4d %20s %14s %3d %2d:%3d %s\n", nk[i], eni.nodeName, eni.address, int(sub), eni.count, eni.status, eni.arg)
	}
	nmmu.RUnlock()
}

func (s *srvNodeInfo) RegisterNode(cx context.Context, ni *nodepb.NodeInfo) (nid *nodepb.NodeID, e error) {
	// registration
	n := getNextNodeID(ni.IsServer)
	if n == -1 { // no extra node ID...
		e = errors.New("No extra nodeID")
		return nil, e
	}

	r := rand.Uint64() // secret for this node
	pr, ok := peer.FromContext(cx)
	var ipaddr string
	if ok {
		ipaddr = pr.Addr.String()
	} else {
		ipaddr = "0.0.0.0"
	}
	eni := eachNodeInfo{
		nodeName: ni.NodeName,
		secret:   r,
		address:  ipaddr,
		lastAlive: time.Now(),
		duration: DefaultDuration,
	}
	log.Println("Node Connection from :", ipaddr, ",", ni.NodeName)
	s.nodeMap[n] = &eni
	log.Println("------------------------------------------------------")
	s.listNodes()
	log.Println("------------------------------------------------------")
	nid = &nodepb.NodeID{NodeId: n, Secret: r, KeepaliveDuration: eni.duration}
	return nid, nil
}

func (s *srvNodeInfo) QueryNode(cx context.Context, nid *nodepb.NodeID) (ni *nodepb.NodeInfo, e error) {
	n := nid.NodeId
	eni, ok := s.nodeMap[n]
	if !ok {
		return nil, errors.New("Unregistered NodeID")
	}
	ni = &nodepb.NodeInfo{NodeName: eni.nodeName}
	return ni, nil
}

func (s *srvNodeInfo) KeepAlive(ctx context.Context, nu *nodepb.NodeUpdate) (nr *nodepb.Response, e error) {
	nid := nu.NodeId
	r := nu.Secret
	ni,ok := s.nodeMap[nid]
	if !ok  {
		fmt.Printf("Can't find node... It's killed %d", nid)
		return &nodepb.Response{Ok: false, Err: "Killed at Nodeserv"}, e
	}
	if r != ni.secret {
		e = errors.New("Secret Failed")
		return &nodepb.Response{Ok: false, Err: "Secret Failed"}, e
	}
	ni.lastAlive = time.Now()
	ni.count = nu.UpdateCount
	ni.status = nu.NodeStatus
	ni.arg = nu.NodeArg

	if ni.lastAlive.Sub(lastPrint) > time.Second *time.Duration(DefaultDuration/2) {
		log.Println("---KeepAlive------------------------------------------")
		s.listNodes()
		log.Println("------------------------------------------------------")
	}

	return &nodepb.Response{Ok: true, Err: ""}, nil
}

func (s *srvNodeInfo) UnRegisterNode(cx context.Context, nid *nodepb.NodeID) (nr *nodepb.Response, e error) {
	r := nid.Secret
	n := nid.NodeId
	ni,ok := s.nodeMap[n]
	if !ok  {
		fmt.Printf("Can't find node... It's killed")
		return &nodepb.Response{Ok: false, Err: "Killed at Nodeserv"}, e
	}

	if r != ni.secret { // secret failed
		e = errors.New("Secret Failed")
		log.Println("Invalid unregister")
		return &nodepb.Response{Ok: false, Err: "Secret Failed"}, e
	}

	log.Println("----------- Delete Node -----------", n, s.nodeMap[n].nodeName)
	delete(s.nodeMap, n)
	s.listNodes()
	log.Println("------------------------------------------------------")

	return &nodepb.Response{Ok: true, Err: ""}, nil
}

func prepareGrpcServer(opts ...grpc.ServerOption) *grpc.Server {
	nodeServer := grpc.NewServer(opts...)
	nodepb.RegisterNodeServer(nodeServer, &srvInfo)
	return nodeServer
}

func main() {
	if gerr := agent.Listen(agent.Options{}); gerr != nil{
		log.Fatal(gerr)
	}

	flag.Parse()
//	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	nodeServer := prepareGrpcServer(opts...)
	log.Printf("Start waiting Node Server at port :%d ...", *port)
	nodeServer.Serve(lis)
}
