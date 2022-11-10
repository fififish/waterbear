package coin 

import(
	"fmt"
	"logging"
	"log"
	"utils"
	"sync"
	"cryptolib"
	"message"
	"quorum"
	"communication/sender"
)


var curCoin utils.IntIntMapArr
var n int
var coinnodes utils.IntIntInt64SetMap
var coinShares utils.IntIntBytesMapArr
var alock sync.Mutex
var round utils.IntIntMap //round number
var cachedMsg utils.IntIntBytesMapArr
var mapMembers map[int]int
var members []int
var verbose bool
var coinmode bool 

func QueryCoin(instanceid int, roundnum int) int{
	//go HandleCachedMsg(instanceid)
	return curCoin.GetValue(instanceid,roundnum)
}

func GetIndex(instanceid int) int{
	return mapMembers[instanceid]
}


func GenCoin(nodeid int64, instanceid int, roundnum int){
	round.Insert(instanceid, roundnum)
	go HandleCachedMsg(instanceid, roundnum)
	sender.CoinBroadcast(nodeid, instanceid, roundnum, message.PRF)
}

func CacheMsg(instanceid int, roundnum int, msg []byte){
	cachedMsg.InsertValue(instanceid, roundnum, msg)
}

func HandleCachedMsg(instanceid int, r int){
	//r,_ := round.Get(instanceid)
	
	msgs := cachedMsg.GetAndClear(instanceid,r)

	if len(msgs) == 0 {
		return 
	}


	p:= fmt.Sprintf("[%v] Handling cached coin message for round %v, len(msgs) %v", instanceid, r, len(msgs))
	logging.PrintLog(verbose, logging.NormalLog, p)
	
	for i:=0; i<len(msgs); i++{
		HandleCoinMsg(msgs[i])
	}

}

func HandleCoinMsg(rawMsg []byte){
	m := message.DeserializeReplicaMessage(rawMsg)
	//log.Printf("message source %v, %v", m.Source, m.Instance)
	if !cryptolib.VerifyShare(m.Hash,m.Source,m.Payload){
		log.Printf("ERROR! Cannot verify share from node %v", m.Source)
		return 
	}
	
	//alock.Lock()
	r,_ := round.Get(m.Instance)

	p:= fmt.Sprintf("[%v] Handling coin message for round %v from node %v, curR %v", m.Instance, m.Round, m.Source, r)
	logging.PrintLog(verbose, logging.NormalLog, p)

	go HandleCachedMsg(m.Instance, r)

	if m.Round < r{
		//alock.Unlock()
		return 
	}
	if m.Round > r{
		//alock.Unlock()
		p:= fmt.Sprintf("[%v] cache coin message for round %v from node %v, curR %v", m.Instance, m.Round, m.Source, r)
		logging.PrintLog(verbose, logging.NormalLog, p)
		CacheMsg(m.Instance,m.Round,rawMsg)
		return 
	}

	//
	
	len1 := coinnodes.GetLen(m.Instance,m.Round)
	coinnodes.Insert(m.Instance,m.Round,m.Source)
	len2 := coinnodes.GetLen(m.Instance,m.Round)
	if len2 > len1{
		coinShares.InsertValueAndInt(m.Instance,m.Round,m.Payload,m.Source)
	}
	
	idarr,shares := coinShares.GetAllValue(m.Instance,m.Round)
	//get := curCoin.GetValue(m.Instance,r)
	//if get==-1 && len2 >= quorum.SQuorumSize() && len(idarr) >= quorum.SQuorumSize(){
	if len2 >= quorum.SQuorumSize() && len(idarr) >= quorum.SQuorumSize(){
		//log.Printf("[%v]----round %v, ready to combine", m.Instance, r)
		p:= fmt.Sprintf("[%v]----round %v, ready to combine", m.Instance, r)
		logging.PrintLog(verbose, logging.NormalLog, p)
		//if len(idarr) < quorum.SQuorumSize() || len(shares) < quorum.SQuorumSize(){
			//alock.Unlock()
		//	return 
		//}
		//idarr := coinnodes.Get(m.Instance)
		prfval := cryptolib.CombineShares(idarr,shares)
		h := cryptolib.GenHash(prfval)
		var coin int 
		if coinmode{
			coin = int(h[1])&1
		}else{
			coin = utils.BytesToInt(h)
		}
		
		//log.Printf("[Coin %v], round %v, value %v, len shares %v, idarr %v, prfval %s, h[0] %v", m.Instance, m.Round, coin, len(shares), idarr, prfval, int(h[0]))
		curCoin.InsertValue(m.Instance,m.Round,coin) 
		p = fmt.Sprintf("[%v]----round %v, coin value %v", m.Instance, r, coin)
		logging.PrintLog(verbose, logging.NormalLog, p)
	}
	//alock.Unlock()
}

func SetMode(cm bool){
	coinmode = cm 
}

func InitCoin(num int,id int64, k int, mem []int,ver bool){
	n = num
	verbose = ver
	coinmode = true 
	curCoin.Init()
	coinnodes.Init()
	coinShares.Init(n)
	cachedMsg.Init(n)
	round.Init()
	cryptolib.StartThSig(id,k)
	members = mem
	mapMembers = make(map[int]int)
	for i:=0; i<len(members); i++{
		round.Insert(members[i],0)
		mapMembers[members[i]] = i
	}
}

func ClearCoin(instanceid int){
	//coinShares.Delete(instanceid)
	//coinnodes.Delete(instanceid)
}