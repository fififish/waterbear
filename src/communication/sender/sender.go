/*
Sender functions.
It implements all sending functions for replicas.
*/

package sender

import (
	"communication"
	"config"
	"context"
	"cryptolib"
	"fmt"
	"google.golang.org/grpc"
	"log"
	logging "logging"
	"message"
	pb "proto/proto/communication"
	"sync"
	"time"
	"utils"
)

var id int64
var err error

//var completed map[string]bool
var verbose bool

var wg sync.WaitGroup

var broadcastTimer int
var sleepTimerValue int
var reply []byte


var dialOpt []grpc.DialOption 
var connections communication.AddrConnMap




func BuildConnection(ctx context.Context, nid string, address string) bool {
	p := fmt.Sprintf("building a connection with %v", nid)
	logging.PrintLog(verbose, logging.NormalLog, p)

	/*if config.CommOption() == "TLS" {
		dialOpt = communication.GetDialOption()
	}*/
	conn, err := grpc.DialContext(ctx, address, dialOpt...)

	if err != nil {
		p := fmt.Sprintf("[Communication Sender Error] failed to bulid a connection with %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return false
	}
	c := pb.NewSendClient(conn)

	connections.Insert(address, c)
	connections.InsertID(address, nid)
	return true
}


func ByteSend(msg []byte, address string, msgType message.TypeOfMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(broadcastTimer)*time.Millisecond)
	defer cancel()

	if address == "" {
		return
	}
	nid := config.FetchReplicaID(address)
	c, built := connections.Get(address)
	existnid := connections.GetID(address)

	if !built || c == nil || nid != existnid{
		suc := BuildConnection(ctx, nid, address)
		if !suc {
			p := fmt.Sprintf("[Communication Sender Error] did not connect to node %s, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			
			communication.NotLive(nid)
			broadcastTimer = broadcastTimer * 2

			return
		}else{
			c,_ = connections.Get(address)
		}
	}

	switch msgType {
	case message.RBC_ALL:
		_, err = c.RBCSendByteMsg(ctx, &pb.RawMessage{Msg: msg})
		if err != nil {
			p := fmt.Sprintf("[Communication Sender Error] could not get reply from node %s when send ReplicaMsg, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			connections.Insert(address, nil)
			return
		}
	case message.ABA_ALL:
		_, err = c.ABASendByteMsg(ctx, &pb.RawMessage{Msg: msg})
		if err != nil {
			p := fmt.Sprintf("[Communication Sender Error] could not get reply from node %s when send ReplicaMsg, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			connections.Insert(address, nil)
			return
		}
	case message.PRF:
		_, err = c.PRFSendByteMsg(ctx, &pb.RawMessage{Msg: msg})
		if err != nil {
			p := fmt.Sprintf("[Communication Sender Error] could not get reply from node %s when send ReplicaMsg, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			connections.Insert(address, nil)
			return
		}
	case message.ECRBC_ALL:
		_, err = c.ECRBCSendByteMsg(ctx, &pb.RawMessage{Msg: msg})
		if err != nil {
			p := fmt.Sprintf("[Communication Sender Error] could not get reply from node %s when send ReplicaMsg, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			connections.Insert(address, nil)
			return
		}
	case message.CBC_ALL:
		_, err = c.CBCSendByteMsg(ctx, &pb.RawMessage{Msg: msg})
		if err != nil {
			p := fmt.Sprintf("[Communication Sender Error] could not get reply from node %s when send ReplicaMsg, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			connections.Insert(address, nil)
			return
		}
	case message.EVCBC_ALL:
		_, err = c.EVCBCSendByteMsg(ctx, &pb.RawMessage{Msg: msg})
		if err != nil {
			p := fmt.Sprintf("[Communication Sender Error] could not get reply from node %s when send ReplicaMsg, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			connections.Insert(address, nil)
			return
		}
	case message.MVBA_DISTRIBUTE:
		_, err = c.MVBASendByteMsg(ctx, &pb.RawMessage{Msg: msg})
		if err != nil {
			p := fmt.Sprintf("[Communication Sender Error] could not get reply from node %s when send ReplicaMsg, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			connections.Insert(address, nil)
			return
		}
	case message.RETRIEVE:
		_, err = c.RetrieveSendByteMsg(ctx, &pb.RawMessage{Msg: msg})
		if err != nil {
			p := fmt.Sprintf("[Communication Sender Error] could not get reply from node %s when send ReplicaMsg, set it to notlive: %v", nid, err)
			logging.PrintLog(true, logging.ErrorLog, p)
			communication.NotLive(nid)
			connections.Insert(address, nil)
			return
		}
	default:
		log.Fatalf("message type %v not supported", msgType)
	}
}

func MACBroadcast(msg []byte, mtype message.ProtocolType){
	
	nodes := FetchNodesFromConfig()

	for i := 0; i < len(nodes); i++ {
		nid := nodes[i]

		dest,_ := utils.StringToInt64(nid)
		request, err := message.SerializeWithMAC(id, dest, msg)
		if err != nil {
			logging.PrintLog(true, logging.ErrorLog, "[Sender Error] Not able to generate MAC")
			continue
		}
		
		if communication.IsNotLive(nid) {
			p := fmt.Sprintf("[Communication Sender] Replica %v is not live, don't send message to it", nid)
			logging.PrintLog(verbose, logging.NormalLog, p)
			continue
		}
		switch mtype{
		case message.RBC:
			go ByteSend(request, config.FetchAddress(nid), message.RBC_ALL)
		case message.ABA:
			go ByteSend(request, config.FetchAddress(nid), message.ABA_ALL)
		case message.CBC:
			go ByteSend(request, config.FetchAddress(nid), message.CBC_ALL)
		case message.EVCBC:
			go ByteSend(request, config.FetchAddress(nid), message.EVCBC_ALL)
		case message.MVBA:
			go ByteSend(request, config.FetchAddress(nid), message.MVBA_DISTRIBUTE)
		}
		
		
	}
}

func SendToNode(msg []byte, dest int64, mtype message.ProtocolType){
	
	nid := utils.Int64ToString(dest)

	request, err := message.SerializeWithMAC(id, dest, msg)
	if err != nil {
		logging.PrintLog(true, logging.ErrorLog, "[Sender Error] Not able to generate MAC")
		return 
	}
	
	if communication.IsNotLive(nid) {
		p := fmt.Sprintf("[Communication Sender] Replica %v is not live, don't send message to it", nid)
		logging.PrintLog(verbose, logging.NormalLog, p)
		return
	}
	switch mtype{
	case message.CBC:
		go ByteSend(request, config.FetchAddress(nid), message.CBC_ALL)
	case message.EVCBC:
		go ByteSend(request, config.FetchAddress(nid), message.EVCBC_ALL)
	}
		
}

/*
isDifferentData means whether the message in data is different, in other words, whether the messages sent to
different replicas are different.
 */
func MACBroadcastWithErasureCode(data [][]byte, mtype message.ProtocolType, isDifferentData bool){

	nodes := FetchNodesFromConfig()

	for i := 0; i < len(nodes); i++ {
		nid := nodes[i]

		dest,_ := utils.StringToInt64(nid)
		var request []byte
		var err error
		if isDifferentData{
			request, err = message.SerializeWithMAC(id, dest, data[i])
		}else {
			request, err = message.SerializeWithMAC(id, dest, data[0])
		}

		if err != nil {
			logging.PrintLog(true, logging.ErrorLog, "[Sender Error] Not able to generate MAC")
			continue
		}

		if communication.IsNotLive(nid) {
			p := fmt.Sprintf("[Communication Sender] Replica %v is not live, don't send message to it", nid)
			logging.PrintLog(verbose, logging.NormalLog, p)
			continue
		}
		switch mtype{
		case message.ECRBC:
			go ByteSend(request, config.FetchAddress(nid), message.ECRBC_ALL)
		case message.MVBA:
			go ByteSend(request, config.FetchAddress(nid), message.RETRIEVE)
		}

	}
}

func CoinBroadcast(nodeid int64, instanceid int, roundnum int, mtype message.TypeOfMessage){
	nodes := FetchNodesFromConfig()

	for i := 0; i < len(nodes); i++ {
		nid := nodes[i]
		
		shareinput := utils.IntToBytes(instanceid+roundnum)
		request := cryptolib.GenPRFShare(shareinput)
		
		m := message.ReplicaMessage{
			Source: 	 nodeid, 
			Instance:    instanceid, 
			Round:       roundnum,
			Payload: 	 request,
			Hash: 		 shareinput, //Used as C for coin messages
		}
		
		msgbyte,_ := m.Serialize()

		if communication.IsNotLive(nid) {
			p := fmt.Sprintf("[Communication Sender] Replica %v is not live, don't send message to it", nid)
			logging.PrintLog(verbose, logging.NormalLog, p)
			continue
		}
	
		go ByteSend(msgbyte, config.FetchAddress(nid), mtype)

		
	}
}

/*
Used for membership protocol to fetch list of nodes
Output
	[]string: a list of nodes (in the string type)
*/
func FetchNodesFromConfig() []string {
	return config.FetchNodes()
}

func StartSender(rid string) {
	log.Printf("Starting sender %v", rid)
	config.LoadConfig()
	verbose = config.FetchVerbose()

	id, err = utils.StringToInt64(rid) // string to int64
	if err != nil {
		p := fmt.Sprintf("[Communication Sender Error] Replica id %v is not valid. Double check the configuration file", id)
		logging.PrintLog(true, logging.ErrorLog, p)
		return
	}

	
	// Set up a connection to the server.

	dialOpt = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		//grpc.WithKeepaliveParams(kacp),
	}
	
	connections.Init()

	
	verbose = config.FetchVerbose()
	communication.StartConnectionManager()
	broadcastTimer = config.FetchBroadcastTimer()
	sleepTimerValue = config.FetchSleepTimer()
}

func SetId(newnid int64) {
	id = newnid
}
