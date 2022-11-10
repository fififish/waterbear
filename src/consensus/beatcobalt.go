
package consensus

import (
	aba "aba/cobalt"
	"broadcast/ecrbc"
	"broadcast/rbc"
	"config"
	"log"
	"quorum"
	"time"
	"utils"
)


func BCMonitorRBCStatus(e int){
	for {
		if epoch.Get() > e {
			return 
		}
	
		for i:=0; i<n; i++{
			if astatus.GetCount()==n {
				return 
			}
			instanceid := GetInstanceID(members[i])
			//instanceid := members[i]
			status := rbc.QueryStatus(instanceid)
			
			if !astatus.GetStatus(instanceid) && status{
				astatus.Insert(instanceid,true)
				//CaptureRBCLat()
				go BCStartABA(instanceid, 1)
			}

			switch consensus {
			case BEATCobalt:
				value := aba.QueryValue(instanceid)
				if value == 1{
					avalues.Insert(instanceid,true)
				}
				if avalues.GetCount() >= quorum.QuorumSize(){
					//CaptureLastRBCLat()
					go BCStartOtherABAs()
				}
			}
			

		}
		time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
	}
}


func BCMonitorECRBCStatus(e int){
	for {
		if epoch.Get() > e {
			return
		}

		
		for i:=0; i<n; i++{
			if astatus.GetCount()==n {
				return 
			}

			instanceid := GetInstanceID(members[i])
			//instanceid := members[i]
			status := ecrbc.QueryStatus(instanceid)

			if !astatus.GetStatus(instanceid) && status{
				astatus.Insert(instanceid,true)
				//CaptureRBCLat()
				go BCStartABA(instanceid, 1)
			}

			switch consensus {
			case BEATCobalt:
				value := aba.QueryValue(instanceid)
				if value == 1{
					avalues.Insert(instanceid,true)
				}
				if avalues.GetCount() >= quorum.QuorumSize(){
					//CaptureLastRBCLat()
					go BCStartOtherABAs()
				}
			}


		}
		time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
	}
}

func BCMonitorABAStatus(e int){
	for {
		if epoch.Get() > e {
			return 
		}
	
		for i:=0; i<n; i++{
			instanceid := GetInstanceID(members[i])
			//instanceid := members[i]
			status := aba.QueryStatus(instanceid)
			
			if !fstatus.GetStatus(instanceid) && status{
				fstatus.Insert(instanceid,true)
				//log.Printf("instance %v completed", instanceid)
				go BCUpdateOutput(instanceid)
			}

			if fstatus.GetCount() == n{
				return 
			}

		}
		time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
	}
}

func BCUpdateOutput(instanceid int){
	//log.Println("++++++++++++++++++++finish ",instanceid)
	value := aba.QueryValue(instanceid)

	if value == 0{
		outputCount.Increment()
	}else{
		outputCount.Increment()
		outputSize.Increment()
		for{
			var v []byte
			if rbcType == RBC{
				v = rbc.QueryReq(instanceid)
			}else if rbcType == ECRBC{
				v = ecrbc.QueryReq(instanceid)
			}

			//if v != nil{
			output.AddItem(v)
			break 
			//}
			//	break
			//}else{
			//	time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
			//}
		}
	}

	
	if curStatus.Get()!=READY && outputCount.Get() == n{
		curStatus.Set(READY)
		ExitEpoch()
		/*switch consensus{
		case BEATCobalt:
			InitBCBFT(false)
		case BiasedBEATCobalt:
			InitBCBFT(true)
		}*/
		return 
	}
}



func BCStartABA(instanceid int, input int){
	if bstatus.GetStatus(instanceid){
		return 
	}
	bstatus.Insert(instanceid,true)
	switch consensus{
	case BEATCobalt:
		aba.StartABA(instanceid, 0, input)
	default:
		log.Fatalf("This script only supports BEAT-Cobalt and biased BEAT-Cobalt")
	}
	
}

func BCStartOtherABAs(){
	//log.Printf("Start other ABAs")
	if otherlock.Get() == 1{
		return 
	}
	for i:=0; i<n; i++{
		instanceid := GetInstanceID(members[i])
		//instanceid := members[i]
		if !astatus.GetStatus(instanceid){
			go BCStartABA(instanceid, 0)
		}
	}
	otherlock.Set(1)
}

func StartBCBFT(data []byte){
	switch consensus{
	case BEATCobalt:
		InitBCBFT(false)
	}
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
		//rbc.StartRBC(iid,data)
		go BCMonitorRBCStatus(epoch.Get())
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
		//ecrbc.StartECRBC(iid,data)
		go BCMonitorECRBCStatus(epoch.Get())
	}

	go BCMonitorABAStatus(epoch.Get())
}

func InitBCBFT(ct bool){
	//rbc.InitRBC(id,n,verbose)
	//aba.InitABA(id,n,InitStatus,members,sleepTimerValue)

	aba.InitCoinType(ct)
	InitStatus(n)

}