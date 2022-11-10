
package biasedit

import (
	"config"
	"fmt"
	"log"
	"time"
	"math/rand"
	"utils"
	"message"
	"quorum"
	"logging"
	"sync"
	"communication/sender"
	"broadcast/rbc"
)


type ABAStatus int

const (
	STATUS_READY    	  ABAStatus = 0
	STATUS_AUX       	  ABAStatus = 1
	STATUS_CONF      	  ABAStatus = 2
	STATUS_DECIDED	   	  ABAStatus = 3
	STATUS_TERMINATE	  ABAStatus = 4
)

var baseinstance int
var round utils.IntIntMap //round number
var bvals utils.IntIntSetMap //set of sent bval values
var bin_values utils.IntIntSetMap //bvalsr in paper
var auxvals utils.IntIntSetMap // statistic value for main-vote
var auxnodes utils.IntIntInt64SetMap // statistic replicas for main-vote
var confvals utils.IntIntSetMap  // statistic value for final-vote, cvals in paper
var confnodes utils.IntIntInt64SetMap // statistic replicas for final-vote


var instancestatus utils.IntIntMap // status for each instance
var finalstatus utils.IntIntMap // status for each instance
var decidedround utils.IntIntMap
var decidedvalue utils.IntIntMap
var alock sync.Mutex
var block sync.Mutex
var clock sync.Mutex

var cachedMsg utils.IntIntBytesMapArr  //have index of replica, map GetIndex(instanceid (==replica id?)) to round num and msg
var astatus utils.IntIntSetMap
var bvalMap utils.IntIntDoubleSetMap  //map instance to its number of replicas, and value(0 or 1)

func HandleCachedMsg(instanceid int, r int){
	ro,_ := round.Get(instanceid)
	if ro != r {
		return 
	}
	stat, _ := finalstatus.Get(instanceid)
	if stat == int(STATUS_TERMINATE){
		return 
	}

	msgs := cachedMsg.GetAndClear(instanceid,r)

	if len(msgs) == 0 {
		return 
	}

	p:= fmt.Sprintf("[%v] Handling cached message for round %v, len(msgs) %v", instanceid, r, len(msgs))
	logging.PrintLog(verbose, logging.NormalLog, p)
	//log.Printf("=== Handle cache message")
	for i:=0; i<len(msgs); i++{
		m := message.DeserializeReplicaMessage(msgs[i])
		switch m.Mtype{
		case message.ABA_BVAL:
			go HandleBVAL(m)
		case message.ABA_AUX:
			go HandleAUX(m)
		case message.ABA_CONF:
			go ProcessConf(msgs[i])
		}
	}

}

func CacheMsg(instanceid int, roundnum int, msg []byte){
	cachedMsg.InsertValue(instanceid, roundnum, msg)
}

func HandleBVAL(m message.ReplicaMessage){
	p := fmt.Sprintf("[%v] Handling bval message, round %v, vote %v from %v", m.Instance, m.Round, m.Value, m.Source)
	logging.PrintLog(verbose, logging.NormalLog, p)
	r,_ := round.Get(m.Instance)
	
	HandleCachedMsg(m.Instance, r)
	if m.Round < r{
		return 
	}
	if m.Round > r{
		msgbyte,_ := m.Serialize()
		CacheMsg(m.Instance,m.Round,msgbyte)
		return 
	}
	

	bvalMap.Insert(m.Instance,m.Round , m.Value, m.Source)
	//upon receiving pre-vote(v) from f + 1 replicas, pre-vote(v) has not been sent
	if bvalMap.GetCount(m.Instance,m.Round ,m.Value) >= quorum.SQuorumSize(){
		SendBval(m.Instance, m.Value)
	}
	//upon receiving pre-voter(v) from 2f + 1 nodes, main-vote() has not been sent
	if bvalMap.GetCount(m.Instance, m.Round ,m.Value) >= quorum.QuorumSize(){
		ProceedToAux(m)
	}
	
}

func HandleAUX(m message.ReplicaMessage){
	p := fmt.Sprintf("[%v] Handling aux message, round %v, vote %v from %v", m.Instance, m.Round, m.Value, m.Source)
	logging.PrintLog(verbose, logging.NormalLog, p)
	r,_ := round.Get(m.Instance)
	
	HandleCachedMsg(m.Instance, r)
	
	if m.Round < r{
		return 
	}
	if m.Round > r{
		msgbyte,_ := m.Serialize()
		CacheMsg(m.Instance,m.Round,msgbyte)
		return
	}
	if !bin_values.Contains(m.Instance,m.Round,m.Value){
		msgbyte,_ := m.Serialize()
		CacheMsg(m.Instance,m.Round,msgbyte)
		return
	}
	

	auxvals.InsertValue(m.Instance,m.Round,m.Value)
	auxnodes.Insert(m.Instance,m.Round,m.Source)

	//upon receiving n−f main-vote()
	if auxnodes.GetLen(m.Instance,m.Round) >= quorum.QuorumSize(){
		ProceedToConf(m.Instance)
	}
	
}

//broadcast pre-vote()
func SendBval(instanceid int, value int){
	r,_ := round.Get(instanceid)

	if config.MaliciousNode(){
		switch config.MaliciousMode() {
		case 0:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send bval %v for round %v!", instanceid,id, 0,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				value = 0
			}
		case 1:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid > 2 * quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send bval %v for round %v!", instanceid,id, 0,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				value = 0
			}
		case 3:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send bval %v but not %v for round %v!", instanceid,id, value ^ 1,value,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				if value != 2{
					value = value ^ 1
				}

			}

		}

	}

	if bvals.Contains(instanceid,r ,value){
		return
	}

	bvals.InsertValue(instanceid, r,value)
	p:= fmt.Sprintf("[%v] Sending bval %v since it previously has not been sent.", instanceid, value)
	logging.PrintLog(verbose, logging.NormalLog, p)

	msg := message.ReplicaMessage{
		Mtype:		 message.ABA_BVAL, 
		Instance:    instanceid, 
		Source: 	 id, 
		Value: 		 value,
		Round: 		 r, 
		Epoch:		 epoch.Get(), 
	}

	msgbyte, err := msg.Serialize()
	if err !=nil{
		log.Fatalf("failed to serialize bval ABA message")
	}
	sender.MACBroadcast(msgbyte, message.ABA)
}
//broadcast main-vote(v)
func ProceedToAux(m message.ReplicaMessage){
	bin_values.InsertValue(m.Instance,m.Round,m.Value)
	alock.Lock() 
	r,_ := round.Get(m.Instance)
	stat, _ := instancestatus.Get(m.Instance)
	if stat >= int(STATUS_AUX){
		alock.Unlock() 
		return 
	}
	instancestatus.Insert(m.Instance,int(STATUS_AUX))
	alock.Unlock() 
	
	p := fmt.Sprintf("[%v] Sending aux %v for round %v", m.Instance, m.Value, r)
	logging.PrintLog(verbose, logging.NormalLog, p)

	msg := message.ReplicaMessage{
		Mtype:		 message.ABA_AUX, 
		Instance:    m.Instance, 
		Source: 	 id, 
		Value: 		 m.Value,
		Round: 		 r, 
		Epoch:		 epoch.Get(), 
	}

	if config.MaliciousNode(){
		switch config.MaliciousMode() {
		case 0:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, Sending aux %v for round %v!", m.Instance,id, 0,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				msg.Value = 0
			}
		case 1:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid > 2 * quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, Sending aux %v for round %v!", m.Instance,id, 0,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				msg.Value = 0
			}
		case 3:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, Sending aux %v but not %v for round %v!", m.Instance,id, msg.Value ^ 1,msg.Value,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				if msg.Value != 2{
					msg.Value = msg.Value ^ 1
				}

			}

		}

	}

	msgbyte, err := msg.Serialize()
	if err !=nil{
		log.Fatalf("failed to serialize bval ABA message")
	}
	sender.MACBroadcast(msgbyte, message.ABA)
}


//r-broadcast final-vote(v)
func ProceedToConf(instanceid int){

	block.Lock() 
	stat, _ := instancestatus.Get(instanceid)
	if stat >= int(STATUS_CONF){
		block.Unlock() 
		return 
	}
	instancestatus.Insert(instanceid,int(STATUS_CONF))
	block.Unlock() 

	r,_ := round.Get(instanceid)
	proposedvalue := 2 //used to represent \bot value
	//if vals only has one value, r-broadcast final-voter(ρ)
	if auxvals.GetLen(instanceid,r) == 1 {
		proposedvalue = auxvals.GetValue(instanceid,r)[0]
	}
	msg := message.ReplicaMessage{
		Mtype:		 message.ABA_CONF, 
		Instance:    instanceid, 
		Source: 	 id, 
		Value: 		 proposedvalue,
		Round: 		 r, 
		Epoch:		 epoch.Get(), 
	}


	if config.MaliciousNode(){
		switch config.MaliciousMode() {
		case 0:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, Sending conf %v for round %v!", msg.Instance,id, 0,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				msg.Value = 0
			}
		case 1:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid > 2 * quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, Sending conf %v for round %v!", msg.Instance,id, 0,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				msg.Value = 0
			}
		case 3:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, Sending conf %v but not %v for round %v!", msg.Instance,id, msg.Value ^ 1,msg.Value,r)
				logging.PrintLog(verbose, logging.NormalLog, p)
				if msg.Value != 2{
					msg.Value = msg.Value ^ 1
				}

			}

		}

	}

	msgbyte, err := msg.Serialize()
	if err !=nil{
		log.Fatalf("failed to serialize bval ABA message")
	}
	
	p := fmt.Sprintf("[%v] Conf send %v for round %v,epoch %v, rbc id %v", instanceid, proposedvalue, r,msg.Epoch ,GetRBCInstance(instanceid,iid,r))
	logging.PrintLog(verbose, logging.NormalLog, p)
	
	rbc.StartRBC(GetRBCInstance(instanceid,iid,r),msgbyte)
	go MonitorConfStatus(instanceid, epoch.Get())
}

func MonitorConfStatus(instanceid int, epochnum int){
	//p := fmt.Sprintf("[%v] Monitor Conf Status. epoch %v", instanceid,epochnum)
	//logging.PrintLog(verbose, logging.ErrorLog, p)
	for {
		
		stat, _ := instancestatus.Get(instanceid)
		//if epoch.Get() > epochnum || stat == int(STATUS_DECIDED){
		if stat == int(STATUS_DECIDED){
			p := fmt.Sprintf("[%v] stop monitor conf status.%v  %v  %v  ", instanceid,epoch.Get(),epochnum,stat)
			logging.PrintLog(verbose, logging.NormalLog, p)
			return 
		}

		r,_ := round.Get(instanceid)
		for i:=0; i<n; i++{
			rbcid := GetRBCInstance(instanceid,members[i],r)
			status := rbc.QueryStatus(rbcid)
			
			if !astatus.Contains(instanceid,r,members[i]) && status{
				p := fmt.Sprintf("[%v] has been delivered, process conf,RBCID %v,members[i] %v", instanceid,rbcid,members[i])
				logging.PrintLog(verbose, logging.NormalLog, p)
				astatus.InsertValue(instanceid,r,members[i])
				ProcessConf(rbc.QueryReq(rbcid))
			}

		}
		time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
	}
}

func ProcessConf(msg []byte){
	if msg == nil{
		//log.Printf("Get a nil message from RBC!")
		return
	}
	m := message.DeserializeReplicaMessage(msg)

	//clock.Lock()
	r,_ := round.Get(m.Instance)

	p := fmt.Sprintf("[%v] received conf %v for round %v from %v, curr %v", m.Instance, m.Value, m.Round, m.Source, r)
	logging.PrintLog(verbose, logging.NormalLog, p)

	if m.Round < r{
		p := fmt.Sprintf("[%v] A smaller round %v from %v, curr %v", m.Instance,m.Round, m.Source, r)
		logging.PrintLog(verbose, logging.NormalLog, p)
		//clock.Unlock()
		return 
	}
	if m.Round > r{
		p := fmt.Sprintf("[%v] cache conf %v for a greater round %v from %v, curr %v", m.Instance, m.Value, m.Round, m.Source, r)
		logging.PrintLog(verbose, logging.NormalLog, p)
		CacheMsg(m.Instance,m.Round,msg)
		//clock.Unlock()
		return 
	}
	

	if (m.Value!=2 && !bin_values.Contains(m.Instance,m.Round,m.Value)) || (m.Value==2 && (!bin_values.Contains(m.Instance,m.Round,0) || !bin_values.Contains(m.Instance,m.Round,1))){
		p := fmt.Sprintf("[%v] cache conf %v for round %v from %v, curr %v", m.Instance, m.Value, m.Round, m.Source, r)
		logging.PrintLog(verbose, logging.NormalLog, p)
		CacheMsg(m.Instance,m.Round,msg)
		//clock.Unlock()
		return
	}

	confvals.InsertValue(m.Instance,m.Round,m.Value)
	confnodes.Insert(m.Instance,m.Round,m.Source)

	
	p = fmt.Sprintf("[%v] handling conf %v for round %v from %v, curr %v", m.Instance, m.Value, m.Round, m.Source, r)
	logging.PrintLog(verbose, logging.NormalLog, p)


	
	dr,_ := decidedround.Get(m.Instance)
	stat, _ := finalstatus.Get(m.Instance)
	rstat, _ := instancestatus.Get(m.Instance)

	//p = fmt.Sprintf("[%v] rstat %v, dr %v, r %v, stat %v, m.Epoch %v, epoch.Get() %v.", m.Instance, rstat,dr,r,stat,m.Epoch,epoch.Get())
	//logging.PrintLog(verbose, logging.ErrorLog, p)

	if rstat==int(STATUS_CONF) && dr!=r && stat == int(STATUS_DECIDED) && m.Epoch == epoch.Get(){

		finalstatus.Insert(m.Instance,int(STATUS_TERMINATE))
		//log.Printf("[%v] terminate in round %v", m.Instance, r)
		p := fmt.Sprintf("[%v] terminate in round %v", m.Instance, r)
		logging.PrintLog(verbose, logging.NormalLog, p)

		//clock.Unlock()
		return 
	}
	if stat == int(STATUS_TERMINATE){
		//clock.Unlock()
		return 
	}
	var estr int
	//upon r-delivering n − f final-vote()
	if confnodes.GetLen(m.Instance,m.Round) >= quorum.QuorumSize(){
		if confvals.GetLen(m.Instance,m.Round) == 1 && m.Epoch == epoch.Get(){
			estr = confvals.GetValue(m.Instance,m.Round)[0]
			finalstatus.Insert(m.Instance,int(STATUS_DECIDED))
			decidedround.Insert(m.Instance,r)
			decidedvalue.Insert(m.Instance,estr)
			//log.Printf("[%v] decide %v in round %v", m.Instance, m.Value, r)
			p := fmt.Sprintf("[%v] decide %v in round %v", m.Instance, m.Value, r)
			logging.PrintLog(verbose, logging.NormalLog, p)

		}else{
			if confvals.GetCount(m.Instance,m.Round,0) >= quorum.SQuorumSize(){
				estr = 0 
			}else if confvals.GetCount(m.Instance,m.Round,1) >= quorum.SQuorumSize(){
				estr = 1
			}else{
				estr = GetCoin(r)
			}
		}
		p := fmt.Sprintf("[%v] enter the next round with %v,cur round %v", m.Instance, estr,r)
		logging.PrintLog(verbose, logging.NormalLog, p)

		
		InitParametersForInstance(m.Instance,r)
		//round.Increment(m.Instance)
		round.Set(m.Instance,r+1)
		StartABA(m.Instance,estr)
		instancestatus.Insert(m.Instance,int(STATUS_READY))
		
	}
	//clock.Unlock()
}



func GetCoin(roundnum int) int{
	if config.TParameter() > 0{
		log.Println("Get coin in TP")
		if roundnum <= config.TParameter(){
			return 1 //biased ba
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		return r.Intn(1)
	}

	if roundnum == 0{
		return 1 //biased ba
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(1)
}

