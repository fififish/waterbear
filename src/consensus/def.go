package consensus

import (
	"broadcast/ecrbc"
	"broadcast/rbc"
	"log"
	"sync"
	"utils"
	"config"
	"cryptolib"
	"communication/sender"
	"aba/biasedit"
	cobalt "aba/cobalt"
	wb "aba/waterbear"
)



type ConsensusType int 

const (
	ITBFT 				ConsensusType = 1
	BEATCobalt 			ConsensusType = 4
	WaterBearBiased     ConsensusType = 14
)

type RbcType int

const (
	RBC					RbcType = 0
	ECRBC				RbcType = 1
)


func StartProcessing(data []byte){
	switch consensus{
	case ITBFT:
		StartITBFT(data)
	case BEATCobalt:
		StartBCBFT(data)
	case WaterBearBiased:
		StartWaterBear(data,true)
	}
	
}

func GetInstanceID(input int) int{
	return input + n*epoch.Get()//baseinstance*epoch.Get()
}

func GetIndexFromInstanceID(input int, e int) int{
	return input-n*e
}

func GetInstanceIDsOfEpoch() []int{
	var output []int 
	for i:=0; i<len(members); i++{
		output = append(output, GetInstanceID(members[i]))
	}
	return output 
}


func StartHandler(rid string) {
	id, errs = utils.StringToInt64(rid)
	
	if errs != nil {
		log.Printf("[Error] Replica id %v is not valid. Double check the configuration file", id)
		return
	}
	iid, _ = utils.StringToInt(rid)

	config.LoadConfig()
	cryptolib.StartCrypto(id, config.CryptoOption())
	consensus = ConsensusType(config.Consensus())
	rbcType = RbcType(config.RBCType())

	
	n = config.FetchNumReplicas()
	curStatus.Init()
	epoch.Init()
	queue.Init()
	verbose = config.FetchVerbose()
	sleepTimerValue = config.FetchSleepTimer()
	
	nodes := config.FetchNodes()
	for i:=0; i<len(nodes); i++{
		nid,_ := utils.StringToInt(nodes[i])
		members = append(members,nid)
	}

	log.Printf("sleeptimer value %v", sleepTimerValue)
	switch consensus{
	case ITBFT:
		log.Printf("running IT-BFT")
		InitITBFT()

		if rbcType == RBC{
			log.Printf("running RBC")
			rbc.InitRBC(id,n,verbose)
		}else if rbcType == ECRBC {
			log.Printf("running ECRBC")
			rbc.InitRBC(id,n,verbose)
			ecrbc.InitECRBC(id,n,verbose)
		}

		biasedit.InitABA(id,n,verbose,members,sleepTimerValue)

	case BEATCobalt:
		log.Printf("running BEAT-Cobalt")
		InitBCBFT(false)

		if rbcType == RBC{
			log.Printf("running RBC")
			rbc.InitRBC(id,n,verbose)
		}else if rbcType == ECRBC {
			log.Printf("running ECRBC")
			//rbc.InitRBC(id,n,verbose)
			ecrbc.InitECRBC(id,n,verbose)
		}
		cobalt.InitABA(id,n,verbose,members,sleepTimerValue)
	case WaterBearBiased:
		log.Printf("running WaterBearBiased-BFT")
		InitWaterBearBFT(true)

		if rbcType == RBC{
			log.Printf("running RBC")
			rbc.InitRBC(id,n,verbose)
		}else if rbcType == ECRBC {
			log.Printf("running ECRBC")
			rbc.InitRBC(id,n,verbose)
			ecrbc.InitECRBC(id,n,verbose)
		}

		wb.InitABA(id,n,verbose,members,sleepTimerValue)
	default:
		log.Fatalf("Consensus type not supported")
	}

	sender.StartSender(rid)
	go RequestMonitor()
}


type QueueHead struct {
	Head string
	sync.RWMutex
}

func (c *QueueHead) Set(head string) {
	c.Lock()
	defer c.Unlock()
	c.Head = head
}

func (c *QueueHead) Get() string {
	c.RLock()
	defer c.RUnlock()
	return c.Head
}

type CurStatus struct{
	enum Status
	sync.RWMutex 
}

type Status int 

const (
	READY Status = 0
	PROCESSING Status = 1
)


func (c *CurStatus) Set(status Status) {
	c.Lock()
	defer c.Unlock()
	c.enum = status
}

func (c *CurStatus) Init() {
	c.Lock()
	defer c.Unlock()
	c.enum = READY
}

func (c *CurStatus) Get() Status {
	c.RLock()
	defer c.RUnlock()
	return c.enum
}

