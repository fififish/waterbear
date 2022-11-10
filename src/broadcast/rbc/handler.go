

package rbc

import (
	"communication/sender"
	"cryptolib"
	"log"
	"message"
	"quorum"
	"utils"
)

var id int64 
var n int 
var verbose bool
var epoch utils.IntValue

func StartRBC(instanceid int, input []byte){
	//log.Printf("Starting RBC %v for epoch %v\n", instanceid, epoch.Get())
	//p := fmt.Sprintf("[%v] Starting RBC for epoch %v", instanceid, epoch.Get())
	//logging.PrintLog(verbose, logging.NormalLog, p)

	msg := message.ReplicaMessage{
		Mtype:		 message.RBC_SEND, 
		Instance:    instanceid, 
		Source: 	 id, 
		TS: 		 utils.MakeTimestamp(), 
		Payload: 	 input,
		Epoch:		 epoch.Get(),
	}

	msgbyte, err := msg.Serialize()
	if err !=nil{
		log.Fatalf("failed to serialize RBC message")
	}
	sender.MACBroadcast(msgbyte, message.RBC)
}

func HandleRBCMsg(inputMsg []byte) {

	tmp := message.DeserializeMessageWithSignature(inputMsg)
	input := tmp.Msg
	content := message.DeserializeReplicaMessage(input)
	mtype := content.Mtype

	if !cryptolib.VerifyMAC(content.Source, tmp.Msg, tmp.Sig) {
		log.Printf("[Authentication Error] The signature of rbc message has not been verified.")
		return
	}

	//log.Printf("handling message from %v, type %v", source, mtype)
	switch mtype {
	case message.RBC_SEND:
		HandleSend(content)
	case message.RBC_ECHO:
		HandleEcho(content)
	case message.RBC_READY:
		HandleReady(content)
	default:
		log.Printf("not supported")
	}

}


func SetEpoch(e int){
	epoch.Set(e)
}

func InitRBC(thisid int64, numNodes int, ver bool){
	id = thisid
	n = numNodes
	verbose = ver
	quorum.StartQuorum(n)
	//log.Printf("ini rstatus %v",rstatus.GetAll())
	rstatus.Init()
	instancestatus.Init()
	cachestatus.Init()
	receivedReq.Init()
	received.Init()
	epoch.Init()
}

func ClearRBCStatus(instanceid int){
	rstatus.Delete(instanceid)
	instancestatus.Delete(instanceid)
	cachestatus.Delete(instanceid)
	receivedReq.Delete(instanceid)
}