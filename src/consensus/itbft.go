
package consensus

import (
	aba "aba/biasedit"
	"broadcast/ecrbc"
	"broadcast/rbc"
	"config"
	"fmt"
	"log"
	"logging"
	"quorum"
	"time"
	"utils"
)


func MonitorRBCStatus(e int){
	for {
		if epoch.Get() > e {
			p := fmt.Sprintf("[Consensus RBC] Current epoch %v is greater than the input epoch %v", epoch.Get(), e)
			logging.PrintLog(true, logging.NormalLog, p)
			return 
		}
	
		for i:=0; i<n; i++{
			instanceid := GetInstanceID(members[i])
			status := rbc.QueryStatus(instanceid)
			if !astatus.GetStatus(instanceid) && status{
				astatus.Insert(instanceid,true)
				go StartABA(instanceid, 1)
			}

			if astatus.GetCount() >= quorum.QuorumSize(){
				go StartOtherABAs()
			}

		}
		time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
	}
}

func MonitorECRBCStatus(e int){
	for {
		if epoch.Get() > e {
			p := fmt.Sprintf("[Consensus RBC] Current epoch %v is greater than the input epoch %v", epoch.Get(), e)
			logging.PrintLog(true, logging.NormalLog, p)
			return
		}

		for i:=0; i<n; i++{
			instanceid := GetInstanceID(members[i])
			status := ecrbc.QueryStatus(instanceid)
			if !astatus.GetStatus(instanceid) && status{
				astatus.Insert(instanceid,true)
				go StartABA(instanceid, 1)
			}

			if astatus.GetCount() >= quorum.QuorumSize(){
				go StartOtherABAs()
			}

		}
		time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
	}
}


func MonitorABAStatus(e int){
	for {
		if epoch.Get() > e {
			p := fmt.Sprintf("[Consensus ABA] Current epoch %v is greater than the input epoch %v", epoch.Get(), e)
			logging.PrintLog(true, logging.NormalLog, p)
			return 
		}
	
		for i:=0; i<n; i++{
			instanceid := GetInstanceID(members[i])
			status := aba.QueryStatus(instanceid)
			//if fstatus.GetStatus(instanceid) && status{
			//	p := fmt.Sprintf("[Consensus] Instance %v has been insert to fstatus %v",instanceid,fstatus)
			//	logging.PrintLog(true, logging.InfoLog, p)
			//}
			if !fstatus.GetStatus(instanceid) && status{
				//log.Printf("[%v] Instance has been decided!**************************************%v",instanceid,instanceid)
				fstatus.Insert(instanceid,true)
				go UpdateOutput(instanceid)
			}

			if fstatus.GetCount() == n{
				return 
			}

		}
		time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
	}
}

func UpdateOutputSet(instanceid int){
	for{
		v := rbc.QueryReq(instanceid)
		if v != nil{
			output.AddItem(v)
			break
		}else{
			time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
		}
	}
}

func UpdateOutput(instanceid int){
	p := fmt.Sprintf("[Consensus] Update Output for instance %v in epoch %v",instanceid,epoch.Get())
	logging.PrintLog(true, logging.NormalLog, p)
	value := aba.QueryValue(instanceid)

	if value == 0{
		outputCount.Increment()
	}else{
		outputSize.Increment()
		outputCount.Increment()
		go UpdateOutputSet(instanceid)
	}
	//p = fmt.Sprintf("[Consensus] outputCount %v for epoch %v",outputCount.Get(),epoch.Get())
	//logging.PrintLog(true, logging.InfoLog, p)
	//elock.Lock()
	if outputCount.Get() == n && curStatus.Get()!=READY{
		curStatus.Set(READY)
		//elock.Unlock()
		ExitEpoch()
		//InitITBFT()
		return
	}
	//elock.Unlock()
}



func StartABA(instanceid int, input int){
	if bstatus.GetStatus(instanceid){
		return 
	}
	bstatus.Insert(instanceid,true)
	//log.Printf("[%v] Starting ABA from zero with input %v in epoch %v", instanceid, input,epoch.Get())
	if config.MaliciousNode(){
		switch config.MaliciousMode() {
		case 0:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				//log.Printf("I'm a malicious node %v, start ABA %v with %v!",id,instanceid,0)
				p := fmt.Sprintf("[%v] I'm a malicious node %v, start ABA with %v!",id,instanceid,0)
				logging.PrintLog(verbose, logging.NormalLog, p)
				input = 0
			}
		case 1:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid > 2 * quorum.FSize(){
				//log.Printf("I'm a malicious node %v, start ABA %v with %v!",id,instanceid,0)
				p := fmt.Sprintf("[%v] I'm a malicious node %v, start ABA with %v!",id,instanceid,0)
				logging.PrintLog(verbose, logging.NormalLog, p)
				input = 0
			}
		case 3:
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				//log.Printf("I'm a malicious node %v, start ABA %v with %v!",id,instanceid,input ^ 1)
				p := fmt.Sprintf("[%v] I'm a malicious node %v, start ABA with %v!",id,instanceid,input^1)
				logging.PrintLog(verbose, logging.NormalLog, p)
				if input != 2{
					input = input ^ 1
				}

			}
		
		}

	}
	aba.StartABAFromRoundZero(instanceid, input)
}

func StartOtherABAs(){
	//log.Printf("Start other ABAs")
	if otherlock.Get() == 1{
		return 
	}
	//log.Printf("Start other ABAs")
	for i:=0; i<n; i++{
		instanceid := GetInstanceID(members[i])
		if !astatus.GetStatus(instanceid){
			//log.Printf("[%v] Start other ABAs for %v with 0",instanceid,instanceid)
			go StartABA(instanceid, 0)
		}
	}
	otherlock.Set(1)
}


func StartITBFT(data []byte){
	/*rbc.InitRBC(id,n,verbose)
	aba.InitABA(id,n,verbose,members,sleepTimerValue)
	aba.SetEpoch(epoch.Get())
	rbc.SetEpoch(epoch.Get())*/
	log.Println("start IT-BFT")
	InitITBFT()
	t1 = utils.MakeTimestamp()
	if rbcType == RBC{
		if config.MaliciousNode() && config.MaliciousMode() == 2{
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				log.Printf("I'm a malicious node %v, don't propose RBC!",id)
			}else {
				rbc.StartRBC(GetInstanceID(iid),data)
			}

		}else {
			rbc.StartRBC(GetInstanceID(iid),data)
		}

		go MonitorRBCStatus(epoch.Get())
	}else if rbcType == ECRBC{
		if config.MaliciousNode() && config.MaliciousMode() == 2{
			intid ,err := utils.Int64ToInt(id)
			if err !=nil {
				log.Fatal("Failed transform int64 to int", err)
			}
			if intid < quorum.FSize(){
				log.Printf("I'm a malicious node %v, don't propose RBC!",id)
			}else {
				ecrbc.StartECRBC(GetInstanceID(iid),data)
			}
		}else {
			ecrbc.StartECRBC(GetInstanceID(iid),data)
		}


		go MonitorECRBCStatus(epoch.Get())
	}

	go MonitorABAStatus(epoch.Get())
}

func InitITBFT(){
	//rbc.InitRBC(id,n,verbose)
	//aba.InitABA(id,n,verbose,members,sleepTimerValue)
	InitStatus(n)
	aba.SetEpoch(epoch.Get())
	if rbcType == RBC{
		rbc.SetEpoch(epoch.Get())
	}else if rbcType == ECRBC{
		rbc.SetEpoch(epoch.Get())
		ecrbc.SetEpoch(epoch.Get())
	}
}