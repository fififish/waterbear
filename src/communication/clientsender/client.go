package clientsender

import (
	"communication"
	"config"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	logging "logging"
	pb "proto/proto/communication"
	"time"
	"sync"
	"utils"
	"cryptolib"
)

var clientTimer int
var verbose bool
var dialOpt []grpc.DialOption // tls dial option
var connections communication.AddrConnMap
var wg sync.WaitGroup
var id int64
var err error

func BuildConnection(ctx context.Context, nid string, address string) bool{
	p := fmt.Sprintf("[Client Sender] builidng a connection with %v", nid)
	logging.PrintLog(verbose, logging.NormalLog, p)

	
	conn, err := grpc.DialContext(ctx, address, dialOpt...)
	if err != nil {
		p := fmt.Sprintf("[Client Sender] failed to bulid a connection with %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return false
	}
	c := pb.NewSendClient(conn)

	connections.Insert(address, c)
	return true
}

func SendRequest(rtype pb.MessageType, t1 int64, op []byte, address string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(clientTimer)*time.Millisecond)
	defer cancel()

	nid := config.FetchReplicaID(address)
	//v, err := utils.StringToInt64(nid)

	if config.SplitPorts() {
		address = communication.UpdateAddress(address)
	}

	c,built := connections.Get(address)
	if !built || c == nil{
		suc := BuildConnection(ctx, nid, address)
		if !suc{

			p := fmt.Sprintf("[Client Sender Error] did not connect to node %s, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			clientTimer = clientTimer * 2

			//CatchSendRequestError(v, op, rtype, t1)
			wg.Done()
			return 
		}else{
			c,_ = connections.Get(address)
		}
	}

	var r *pb.RawMessage
	r, err = c.SendRequest(ctx, &pb.Request{Type: rtype, Request: op})

	if err != nil {
		p := fmt.Sprintf("[Client Sender Error] could not get reply from node %s, %v", nid, err)
		logging.PrintLog(true, logging.ErrorLog, p)
		//CatchSendRequestError(v, op, rtype, t1)
		connections.Insert(address, nil)
		wg.Done()
		return	
	}

	re := string(r.GetMsg())
	log.Printf("Got a reply %s", re)
	wg.Done()
}


func BroadcastRequest(rtype pb.MessageType, op []byte) {
	t1 := utils.MakeTimestamp()
	
	nodes := communication.FetchNodesFromConfig()
	
	for i := 0; i < len(nodes); i++ {
		nid := nodes[i]
		if communication.IsNotLive(nid) {
			p := fmt.Sprintf("[Client Sender] Replica %s is not live, not sending any message to it", nid)
			logging.PrintLog(verbose, logging.NormalLog, p)
			continue
		}
		wg.Add(1)
		p := fmt.Sprintf("[Client Sender] Send a %v Request to Replica %v", rtype, nid)
		logging.PrintLog(verbose, logging.NormalLog, p)
	
		go SendRequest(rtype, t1, op, config.FetchAddress(nid))
	
	}
	
	wg.Wait()
}


func StartClientSender(cid string, loadkey bool) {

	config.LoadConfig()
	verbose = config.FetchVerbose()

	id, err = utils.StringToInt64(cid) 
	if err != nil {
		log.Printf("[Client Sender Error] Client id %v is not valid. Double check the configuration file", id)
		return
	}


	cryptolib.StartCrypto(id, config.CryptoOption())

	communication.StartConnectionManager()

	connections.Init()

	clientTimer = config.FetchClientTimer()


	dialOpt = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

}
