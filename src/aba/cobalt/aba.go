
package cobalt

import (
	"config"
	"time"
	"fmt"
	"log"
	"utils"
	"message"
	"quorum"
	"logging"
	"sync"
	"communication/sender"
	"aba/coin"
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
var round utils.IntIntMap    //round number
var bvals utils.IntIntSetMap //set of sent bval values
var bin_values utils.IntIntSetMap
var auxvals utils.IntIntSetMap
var auxnodes utils.IntIntInt64SetMap
var confvals utils.IntIntSetMap
var confnodes utils.IntIntInt64SetMap


var instancestatus utils.IntIntMap // status for each instance
var finalstatus utils.IntIntMap // status for each instance
var decidedround utils.IntIntMap
var decidedvalue utils.IntIntMap
var alock sync.Mutex
var block sync.Mutex
var clock sync.Mutex
var dlock sync.Mutex
var cachedMsg utils.IntIntBytesMapArr
var astatus utils.IntSetMap
var bvalMap utils.IntIntDoubleSetMap

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
	
	for i:=0; i<len(msgs); i++{
		m := message.DeserializeReplicaMessage(msgs[i])
		switch m.Mtype{
		case message.ABA_BVAL:
			go HandleBVAL(m)
		case message.ABA_AUX:
			go HandleAUX(m)
		case message.ABA_CONF:
			go HandleCONF(m)
		}
	}

}

func CacheMsg(instanceid int, roundnum int, msg []byte){
	cachedMsg.InsertValue(instanceid, roundnum, msg)
}

func HandleBVAL(m message.ReplicaMessage){
	
	r,_ := round.Get(m.Instance)
	p := fmt.Sprintf("[%v] Handling bval message, round %v, vote %v from %v, curR %v", m.Instance, m.Round, m.Value, m.Source,r)
	logging.PrintLog(verbose, logging.NormalLog, p)

	HandleCachedMsg(m.Instance, r)

	if m.Round < r{
		return 
	}
	if m.Round > r{
		msgbyte,_ := m.Serialize()
		CacheMsg(m.Instance,m.Round,msgbyte)
		return 
	}
	

	bvalMap.Insert(m.Instance,m.Round,m.Value, m.Source)
	
	if bvalMap.GetCount(m.Instance,m.Round,m.Value) >= quorum.SQuorumSize(){
		SendBval(m.Instance, m.Round, m.Value)
	}

	if bvalMap.GetCount(m.Instance,m.Round, m.Value) >= quorum.QuorumSize(){
		ProceedToAux(m)
	}
	
}

func HandleAUX(m message.ReplicaMessage){
	r,_ := round.Get(m.Instance)
	p := fmt.Sprintf("[%v] Handling aux message, round %v, vote %v from %v, curR %v", m.Instance, m.Round, m.Value, m.Source, r)
	logging.PrintLog(verbose, logging.NormalLog, p)

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

	
	if auxnodes.GetLen(m.Instance,m.Round) >= quorum.QuorumSize(){
		ProceedToConf(m.Instance)
	}
	
}

func SendBval(instanceid int, roundnum int, value int){


	//r,_ := round.Get(instanceid)

	if config.MaliciousNode(){
		switch config.MaliciousMode() {
		case 0:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send bval %v!", instanceid,id, 0)
				logging.PrintLog(verbose, logging.NormalLog, p)
				value = 0
			}
		case 1:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid > 2 * quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send bval %v!", instanceid,id, 0)
				logging.PrintLog(verbose, logging.NormalLog, p)
				value = 0
			}
		case 3:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send bval %v but not %v!", instanceid,id, value ^ 1,value)
				logging.PrintLog(verbose, logging.NormalLog, p)
				if value != 2{
					value = value ^ 1
				}

			}

		}

	}

	if bvals.Contains(instanceid, roundnum, value){
		return
	}

	bvals.InsertValue(instanceid, roundnum,value)
	p:= fmt.Sprintf("[%v] Sending bval %v since it previously has not", instanceid, value)
	logging.PrintLog(verbose, logging.NormalLog, p)

	msg := message.ReplicaMessage{
		Mtype:		 message.ABA_BVAL, 
		Instance:    instanceid, 
		Source: 	 id, 
		Value: 		 value,
		Round: 		 roundnum, 
	}

	msgbyte, err := msg.Serialize()
	if err !=nil{
		log.Fatalf("failed to serialize bval ABA message")
	}
	sender.MACBroadcast(msgbyte, message.ABA)
}

func ProceedToAux(m message.ReplicaMessage){
	p := fmt.Sprintf("[%v] round %v, inserting %v to binv", m.Instance, m.Round, m.Value)
	logging.PrintLog(verbose, logging.NormalLog, p)

	bin_values.InsertValue(m.Instance,m.Round,m.Value)
	
	//alock.Lock() 
	//r,_ := round.Get(m.Instance)
	stat, _ := instancestatus.Get(m.Instance)
	if stat >= int(STATUS_AUX){
	//	alock.Unlock() 
		return 
	}
	instancestatus.Insert(m.Instance,int(STATUS_AUX))
	//alock.Unlock() 
	
	p = fmt.Sprintf("[%v] Sending aux %v for round %v", m.Instance, m.Value, m.Round)
	logging.PrintLog(verbose, logging.NormalLog, p)

	msg := message.ReplicaMessage{
		Mtype:		 message.ABA_AUX, 
		Instance:    m.Instance, 
		Source: 	 id, 
		Value: 		 m.Value,
		Round: 		 m.Round, 
	}

	if config.MaliciousNode(){
		switch config.MaliciousMode() {
		case 0:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send aux %v for round %v!", m.Instance,id, 0,m.Round)
				logging.PrintLog(verbose, logging.NormalLog, p)
				msg.Value = 0
			}
		case 1:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid > 2 * quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send aux %v for round %v!", m.Instance,id, 0,m.Round)
				logging.PrintLog(verbose, logging.NormalLog, p)
				msg.Value = 0
			}
		case 3:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send aux %v but not %v for round %v!", m.Instance,id, msg.Value ^ 1,msg.Value,m.Round)
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

func ProceedToConf(instanceid int){

	//block.Lock() 
	stat, _ := instancestatus.Get(instanceid)
	if stat >= int(STATUS_CONF){
		//block.Unlock() 
		return 
	}
	instancestatus.Insert(instanceid,int(STATUS_CONF))
	//block.Unlock() 

	r,_ := round.Get(instanceid)

	proposedvalue := 2 //used to represent \bot value
	if auxvals.GetLen(instanceid,r) == 1 {
		proposedvalue = auxvals.GetValue(instanceid,r)[0]
	}
	msg := message.ReplicaMessage{
		Mtype:		 message.ABA_CONF, 
		Instance:    instanceid, 
		Source: 	 id, 
		Value: 		 proposedvalue,
		Round: 		 r, 
	}

	if config.MaliciousNode(){
		switch config.MaliciousMode() {
		case 0:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send conf %v for round %v!", msg.Instance,id, 0,msg.Round)
				logging.PrintLog(verbose, logging.NormalLog, p)
				msg.Value = 0
			}
		case 1:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid > 2 * quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send conf %v for round %v!", msg.Instance,id, 0,msg.Round)
				logging.PrintLog(verbose, logging.NormalLog, p)
				msg.Value = 0
			}
		case 3:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				p := fmt.Sprintf("[%v] I'm a malicious node %v, send conf %v but not %v for round %v!", msg.Instance,id, msg.Value ^ 1,msg.Value,msg.Round)
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
	
	p := fmt.Sprintf("[%v] Conf send %v for round %v", instanceid, proposedvalue, r)
	logging.PrintLog(verbose, logging.NormalLog, p)

	sender.MACBroadcast(msgbyte, message.ABA)
}

func HandleCONF(m message.ReplicaMessage){
	r,_ := round.Get(m.Instance)
	p := fmt.Sprintf("[%v] Handling conf message, round %v, vote %v from %v, curr %v", m.Instance, m.Round, m.Value, m.Source, r)
	logging.PrintLog(verbose, logging.NormalLog, p)

	//clock.Lock()
	//defer clock.Unlock()

	
	//go HandleCachedMsg(m.Instance, r)

	//clock.Lock()
	
	if m.Round < r{
	    //clock.Unlock()
		return 
	}
	if m.Round > r{
		msg,_ := m.Serialize()
		CacheMsg(m.Instance,m.Round,msg)
		//clock.Unlock()
		return 
	}

	if (m.Value!=2 && !bin_values.Contains(m.Instance,m.Round,m.Value)) || (m.Value==2 && (!bin_values.Contains(m.Instance,m.Round,0) || !bin_values.Contains(m.Instance,m.Round,1))){
		p := fmt.Sprintf("[%v] cache conf message, round %v, vote %v from %v, curr %v, %v", m.Instance, m.Round, m.Value, m.Source, r, bin_values.GetValue(m.Instance,m.Round))
		logging.PrintLog(verbose, logging.NormalLog, p)
		msg,_ := m.Serialize()
		CacheMsg(m.Instance,m.Round,msg)
		//clock.Unlock()
		return
	}

	go HandleCachedMsg(m.Instance, r)
	
	confvals.InsertValue(m.Instance,m.Round,m.Value)
	confnodes.Insert(m.Instance,m.Round,m.Source)
	
	
	//clock.Lock()
	dr,_ := decidedround.Get(m.Instance)
	stat, _ := finalstatus.Get(m.Instance)
	//rstat, _ := instancestatus.Get(m.Instance)
	
	var estr int
	
	r2,_ := round.Get(m.Instance)
	if r2 != r {
		p := fmt.Sprintf("[%v] r2 round %v, round %v, return", m.Instance, r2, r)
		logging.PrintLog(verbose, logging.NormalLog, p)
		return 
	}
	if confnodes.GetLen(m.Instance,m.Round) >= quorum.QuorumSize(){
		
		p := fmt.Sprintf("[%v] round %v, ready to get common coin", m.Instance, r)
		logging.PrintLog(verbose, logging.NormalLog, p)

		coinval := make(chan int)
		go GetCoin(m.Instance, r, coinval)
		c := <- coinval
		
		p = fmt.Sprintf("[%v] coin value for round %v is %v, mround %v, conf %v, r2 %v", m.Instance, r, c, m.Round, confvals.GetValue(m.Instance,m.Round), r2)
		logging.PrintLog(verbose, logging.NormalLog, p)
		
		//rstat==int(STATUS_CONF) &&
		if dr!=r2 && stat == int(STATUS_DECIDED){
			v,e := decidedvalue.Get(m.Instance)
			//log.Printf("[%v] round %v, v %v, c %v", m.Instance, r, v, c)
			if e && v==c{ //((v==1&&c==1)|| v==0){
				finalstatus.Insert(m.Instance,int(STATUS_TERMINATE))
				//log.Printf("[%v] ******************************** terminate in round %v", m.Instance, r)
				p := fmt.Sprintf("[%v] terminate in round %v", m.Instance, r)
				logging.PrintLog(verbose, logging.NormalLog, p)
		
				//clock.Unlock()
				return 
			}
		}
		if stat == int(STATUS_TERMINATE){
			//clock.Unlock()
			return 
		}
		
		
	
		//clock.Lock()
		//defer clock.Unlock() 
		if stat < int(STATUS_DECIDED) {
			l,tv := confvals.GetLenAndVal(m.Instance,r)
			//if confvals.GetLen(m.Instance) == 1{
			//	estr = confvals.Get(m.Instance)[0]
			p := fmt.Sprintf("[%v] get conf value in round %v: %v %v", m.Instance,r,l,tv)
			logging.PrintLog(verbose, logging.NormalLog, p)
			if l == 1{
				estr = tv
				if estr == c{
					finalstatus.Insert(m.Instance,int(STATUS_DECIDED))
					decidedround.Insert(m.Instance,r)
					decidedvalue.Insert(m.Instance,estr)
					//log.Printf("[%v] decide %v in round %v****************************************\n%v\n", m.Instance, m.Value, r,finalstatus.GetAll())
					p := fmt.Sprintf("[%v] decide %v in round %v", m.Instance, m.Value, r)
					logging.PrintLog(verbose, logging.NormalLog, p)
				}
			}else{
				estr = c
			}
		}else{
			v,_ := decidedvalue.Get(m.Instance)
			estr = v
		}
		
		//clock.Lock()
		//defer clock.Unlock()
		r2,_ = round.Get(m.Instance)
		if r2 != r {
			return 
		}
		//rstat, _ := instancestatus.Get(m.Instance)
		//clock.Unlock()
		//log.Printf("[%v] enter round %v with %v", m.Instance, r+1, estr)

		//if rstat == int(STATUS_CONF){
			p = fmt.Sprintf("[%v] enter round %v with %v", m.Instance, r+1, estr)
			logging.PrintLog(verbose, logging.NormalLog, p)

			//coin.ClearCoin(m.Instance)
			InitParametersForInstance(m.Instance,r)
			//round.Increment(m.Instance)
			round.Set(m.Instance,r+1)
			StartABA(m.Instance,r+1,estr)
			instancestatus.Insert(m.Instance,int(STATUS_READY))
		//}else{
		//	p = fmt.Sprintf("[%v] not entering round %v rstat %v", m.Instance, r, rstat)
		//	logging.PrintLog(verbose, logging.NormalLog, p)
		//}
		
		
	}//else{
	//	clock.Unlock()
	//}
	
}



func GetCoin(instanceid int, roundnum int, r chan int){
	if cointype && roundnum == 0{
		r <- 1
		return 
	}
	
	waitTime := 1

	//coin.ClearCoin(instanceid)
	coin.GenCoin(id, instanceid, roundnum)
	for{
		v := coin.QueryCoin(instanceid,roundnum)
		if v>=0{
			r <- v 
			
			return 
		}
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
		waitTime = waitTime * 2
	}

}

