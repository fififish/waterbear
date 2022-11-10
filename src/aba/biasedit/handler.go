
package biasedit

import (
	"broadcast/ecrbc"
	"config"
	"fmt"
	"log"
	"message"
	"quorum"
	"utils"
	"logging"
	"cryptolib"
	"communication/sender"
	"broadcast/rbc"
)

var id int64 
var iid int //id in type int
var n int 
var verbose bool
var members []int
var sleepTimerValue int
var mapMembers map[int]int  //map replica id to its index
var epoch utils.IntValue

func QueryStatus(instanceid int) bool{
	v,exist := finalstatus.Get(instanceid)
	return exist && v>=int(STATUS_DECIDED)
}
 //query decided value for instanceid
func QueryValue(instanceid int) int{
	v,exist:= decidedvalue.Get(instanceid)
	if !exist{
		return -1
	}
	return v
}


func StartABAFromRoundZero(instanceid int, input int){

	r,_ := round.Get(instanceid)

	if r > 0 {
		p := fmt.Sprintf("[%v] Round %v > 0 with %v", instanceid, r, input)
		logging.PrintLog(verbose, logging.NormalLog, p)
		return 
	}

	p := fmt.Sprintf("[%v] Starting ABA round %v with value %v", instanceid, r, input)
	logging.PrintLog(verbose, logging.NormalLog, p)

	HandleCachedMsg(instanceid, r)

	bvals.InsertValue(instanceid,r ,input)
	
	msg := message.ReplicaMessage{
		Mtype:		 message.ABA_BVAL, 
		Instance:    instanceid, 
		Source: 	 id, 
		Value: 		 input,
		Round: 		 r, 
		Epoch:		 epoch.Get(), 
	}

	msgbyte, err := msg.Serialize()
	if err !=nil{
		log.Fatalf("failed to serialize ABA message")
	}
	sender.MACBroadcast(msgbyte, message.ABA)
	
	if input == 1{
		ProceedToAux(msg)
		auxvals.InsertValue(instanceid,r,1)
		ProceedToConf(instanceid)
	}
	
}

func StartABA(instanceid int, input int){

	r,_ := round.Get(instanceid)

	p := fmt.Sprintf("[%v] Starting ABA round %v with value %v", instanceid, r, input)
	logging.PrintLog(verbose, logging.NormalLog, p)

	HandleCachedMsg(instanceid, r)


	if config.MaliciousNode(){
		switch config.MaliciousMode() {
		case 0:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, start ABA %v for round %v!", instanceid,id, 0,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				input = 0
			}
		case 1:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid > 2 * quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, start ABA %v for round %v!", instanceid,id, 0,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				input = 0
			}
		case 3:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, start ABA %v but not %v for round %v!", instanceid,id, input ^ 1,input,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				if input != 2{
					input = input ^ 1
				}

			}

		}

	}

	bvals.InsertValue(instanceid, r,input)
	
	msg := message.ReplicaMessage{
		Mtype:		 message.ABA_BVAL, 
		Instance:    instanceid, 
		Source: 	 id, 
		Value: 		 input,
		Round: 		 r, 
		Epoch:		 epoch.Get(), 
	}

	msgbyte, err := msg.Serialize()
	if err !=nil{
		log.Fatalf("failed to serialize ABA message")
	}
	sender.MACBroadcast(msgbyte, message.ABA)

	/*if r == 0 && input == 1{
		ProceedToAux(msg)
		auxvals.Insert(m.Instance,1)
		ProceedToConf(instanceid)
	}*/
	
}

func HandleABAMsg(inputMsg []byte) {

	tmp := message.DeserializeMessageWithSignature(inputMsg)
	input := tmp.Msg
	content := message.DeserializeReplicaMessage(input)
	mtype := content.Mtype

	if !cryptolib.VerifyMAC(content.Source, tmp.Msg, tmp.Sig) {
		log.Printf("[Authentication Error] The signature of aba message has not been verified.")
		return
	}

	//log.Printf("handling message from %v, type %v", source, mtype)
	p := fmt.Sprintf("handling message from %v, type %v", content.Source, mtype)
	logging.PrintLog(verbose, logging.NormalLog, p)
	switch mtype {
	case message.ABA_BVAL:
		HandleBVAL(content)
	case message.ABA_AUX:
		HandleAUX(content)
	default:
		log.Printf("not supported")
	}
}



func SetEpoch(e int){
	epoch.Set(e)
}

func InitABA(thisid int64, numNodes int, ver bool, mem []int, st int){
	id = thisid
	iid,_ = utils.Int64ToInt(id)
	n = numNodes
	verbose = ver
	quorum.StartQuorum(n)
	members = mem
	sleepTimerValue = st
	epoch.Init()

	round.Init()
	//initialize round numbers to 0 for all instances
	mapMembers = make(map[int]int)
	for i:=0; i<len(members); i++{
		round.Insert(members[i],0)
		mapMembers[members[i]] = i
	}

	InitParameters()

	instancestatus.Init()
	finalstatus.Init()
	decidedround.Init()
	decidedvalue.Init()
	

	astatus.Init()
	baseinstance = 1000 //hard-code to 1000 to avoid conflicts

}

//Init independently for testing ABA only 
func InitIndependentABA(id int64, n int, verbose bool){
	rbc.InitRBC(id,n,verbose) //Comment it 
}

func InitParameters(){
	cachedMsg.Init(n)
	bvals.Init()
	bin_values.Init()
	auxvals.Init()
	auxnodes.Init()
	confvals.Init()
	confnodes.Init()
	bvalMap.Init()
}

func GetRBCInstance(instanceid int, nid int, r int) int{
	return baseinstance*(instanceid+1) + n*r + nid
}


func InitParametersForInstance(instanceid int, r int){
	//bvals.Delete(instanceid)
	//bin_values.Delete(instanceid)
	//auxvals.Delete(instanceid)
	//auxnodes.Delete(instanceid)
	//confvals.Delete(instanceid)
	//confnodes.Delete(instanceid)
	//bvalMap.Delete(instanceid)
	//astatus.Delete(instanceid)
	
	for i:=0; i<n; i++{
		rbcid := GetRBCInstance(instanceid, members[i], r)
		rbc.ClearRBCStatus(rbcid)
		ecrbc.ClearECRBCStatus(rbcid)
	}
	
}

func GetIndex(instanceid int) int{
	return mapMembers[instanceid]
}


