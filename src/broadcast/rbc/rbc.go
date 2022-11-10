
package rbc

import (
	"log"
	"fmt"
	"sync"
	"logging"
	"message"
	"utils"
	"cryptolib"
	"quorum"
	"communication/sender"
)

type RBCStatus int

const (
	STATUS_IDLE      RBCStatus = 0
	STATUS_SEND      RBCStatus = 1
	STATUS_ECHO      RBCStatus = 2
	STATUS_READY     RBCStatus = 3
)


var rstatus utils.IntBoolMap//broadcast status,only has value when  RBC Deliver
var instancestatus utils.IntIntMap // status for each instance, used in RBC
var cachestatus utils.IntIntMap // status for each instance
var receivedReq utils.IntByteMap  //req is serialized RawOPS or replica msg
var received utils.IntSet
var elock sync.Mutex
var rlock sync.Mutex

//check whether the instance has been deliver in RBC
func QueryStatus(instanceid int) bool{
	v,exist := rstatus.Get(instanceid)
	return v && exist
}

func QueryStatusCount() int{
	return rstatus.GetCount()
}

func QueryReq(instanceid int) []byte{
	v,exist := receivedReq.Get(instanceid)
	if !exist{
		return nil
	}
	return v
}

func HandleSend(m message.ReplicaMessage){
	result,exist := rstatus.Get(m.Instance)
	if exist && result{
		return 
	}

	p := fmt.Sprintf("[%v] Handling send message from node %v", m.Instance, m.Source)
	logging.PrintLog(verbose, logging.NormalLog, p)
	instancestatus.Insert(m.Instance,int(STATUS_SEND))

	msg := m 
	msg.Source = id 
	msg.Mtype = message.RBC_ECHO
	msg.Hash = cryptolib.GenInstanceHash(utils.IntToBytes(m.Instance),m.Payload)
	if !received.IsTrue(m.Instance){
		receivedReq.Insert(m.Instance, m.Payload)
		received.AddItem(m.Instance)
	}


	msgbyte, err := msg.Serialize()
	if err !=nil{
		log.Fatalf("failed to serialize echo message")
	}
	sender.MACBroadcast(msgbyte, message.RBC)

	v,exist := cachestatus.Get(m.Instance)
	if exist && v >= int(STATUS_ECHO){
		SendReady(m)
	}
	if exist && v == int(STATUS_READY){
		Deliver(m)
	}
}

func HandleEcho(m message.ReplicaMessage){
	result,exist := rstatus.Get(m.Instance)
	if exist && result{
		return 
	}


	p := fmt.Sprintf("[%v] Handling echo message from node %v", m.Instance, m.Source)
	logging.PrintLog(verbose, logging.NormalLog, p)

	hash := utils.BytesToString(m.Hash)
	quorum.Add(m.Source, hash, nil, quorum.PP)

	if quorum.CheckQuorum(hash, quorum.PP) {
		if !received.IsTrue(m.Instance){
			receivedReq.Insert(m.Instance, m.Payload)
			received.AddItem(m.Instance)
		}
		SendReady(m)
	}
}

func SendReady(m message.ReplicaMessage){
	elock.Lock() 
	stat, _ := instancestatus.Get(m.Instance)

	if stat == int(STATUS_SEND){
		instancestatus.Insert(m.Instance,int(STATUS_ECHO))
		elock.Unlock()
		p := fmt.Sprintf("Sending ready for instance id %v", m.Instance)
		logging.PrintLog(verbose, logging.NormalLog, p)

		msg := m 
		msg.Source = id 
		msg.Mtype = message.RBC_READY 
		msgbyte, err := msg.Serialize()
		if err !=nil{
			log.Fatalf("failed to serialize ready message")
		}
		sender.MACBroadcast(msgbyte, message.RBC)
	}else{
		v,exist := cachestatus.Get(m.Instance)
		elock.Unlock()
		if exist && v == int(STATUS_READY) {
			instancestatus.Insert(m.Instance,int(STATUS_ECHO))
			Deliver(m)
		}else{
			cachestatus.Insert(m.Instance, int(STATUS_ECHO))
		}
	}
}

func HandleReady(m message.ReplicaMessage){
	result,exist := rstatus.Get(m.Instance)
	if exist && result{
		return 
	}

	p := fmt.Sprintf("[%v] Handling ready message from node %v", m.Instance, m.Source)
	logging.PrintLog(verbose, logging.NormalLog, p)
	

	hash := utils.BytesToString(m.Hash)
	quorum.Add(m.Source, hash, nil, quorum.CM)

	if quorum.CheckEqualSmallQuorum(hash){
		if !received.IsTrue(m.Instance){
			receivedReq.Insert(m.Instance, m.Payload)
			received.AddItem(m.Instance)
		}
		SendReady(m)
	}

	if quorum.CheckQuorum(hash, quorum.CM) {
		Deliver(m)
	}
}

func Deliver(m message.ReplicaMessage){
	rlock.Lock() 
	stat, _ := instancestatus.Get(m.Instance)
	
	if stat == int(STATUS_ECHO){
		if !received.IsTrue(m.Instance){
			receivedReq.Insert(m.Instance, m.Payload)
			received.AddItem(m.Instance)
		}
		instancestatus.Insert(m.Instance,int(STATUS_READY))
		rlock.Unlock()
		
		p := fmt.Sprintf("[%v] RBC Deliver the request epoch %v, curEpoch %v", m.Instance, m.Epoch, epoch.Get())
		logging.PrintLog(verbose, logging.NormalLog, p)

		//if epoch.Get() == m.Epoch{
			rstatus.Insert(m.Instance,true)
		//if m.Instance<100{
		//	log.Printf("insert %v rstatus: %v",m.Instance,rstatus.GetAll())
		//}

		//}
		

	}else{
		rlock.Unlock()
		cachestatus.Insert(m.Instance, int(STATUS_READY))
	}
}