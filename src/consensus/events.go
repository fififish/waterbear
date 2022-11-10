package consensus

import (
	"encoding/json"
	"log"
	"message"
	"quorum"
	"time"
	"utils"
)



var verbose bool   //verbose level
var id int64  //id of server
var iid int  //id in type int, start a RBC using it to instanceid
var errs error
var queue Queue     // cached client requests
var queueHead QueueHead // hash of the request that is in the fist place of the queue
var sleepTimerValue int // sleeptimer for the while loop that continues to monitor the queue or the request status
var consensus ConsensusType
var rbcType	RbcType
var n int 
var members []int
var t1 int64
var baseinstance int

var batchSize int
var requestSize int

func ExitEpoch(){
	t2 := utils.MakeTimestamp()
	if (t2-t1) == 0{
		log.Printf("Latancy is zero!")
		return
	}
	if outputSize.Get() == 0{
		log.Printf("Finish zero instacne!")
		return
	}
	log.Printf("*****epoch ends with  output size %v, latency %v ms, throughput %d, quorum tps %d", outputSize.Get(), t2-t1,int64(outputSize.Get()*batchSize*1000)/(t2-t1), int64(quorum.QuorumSize()*batchSize*1000)/(t2-t1))

}

func CaptureRBCLat(){
	t3 := utils.MakeTimestamp()
	if (t3-t1) == 0{
		log.Printf("Latancy is zero!")
		return
	}
	log.Printf("*****RBC phase ends with %v ms", t3-t1)

}

func CaptureLastRBCLat(){
	t3 := utils.MakeTimestamp()
	if (t3-t1) == 0{
		log.Printf("Latancy is zero!")
		return
	}
	log.Printf("*****Final RBC phase ends with %v ms", t3-t1)

}

func RequestMonitor() {
	for {
		if curStatus.Get()==READY && !queue.IsEmpty() { 
			curStatus.Set(PROCESSING)

			batch := queue.GrabWtihMaxLenAndClear()
			rops := message.RawOPS{
				OPS : batch,
			}
			data,err := rops.Serialize()
			if err != nil{
				continue 
			}
			StartProcessing(data)
		} else {
			time.Sleep(time.Duration(sleepTimerValue) * time.Millisecond)
		}
	}
}


func HandleRequest(request []byte, hash string) {
	//log.Printf("Handling request")
	//rawMessage := message.DeserializeMessageWithSignature(request)
	//m := message.DeserializeClientRequest(rawMessage.Msg)

	/*if !cryptolib.VerifySig(m.ID, rawMessage.Msg, rawMessage.Sig) {
		log.Printf("[Authentication Error] The signature of client request has not been verified.")
		return
	}*/
	//log.Printf("Receive len %v op %v\n",len(request),m.OP)
	batchSize = 1
	requestSize = len(request)
	queue.Append(request)
}


func HandleBatchRequest(requests []byte) {
	requestArr := DeserializeRequests(requests)
	//var hashes []string
	Len := len(requestArr)
	log.Printf("Handling batch requests with len %v\n",Len)
	//for i:=0;i<Len;i++{
	//	hashes = append(hashes,string(cryptolib.GenHash(requestArr[i])))
	//}
	//for i:=0;i<Len;i++{
	//	HandleRequest(requestArr[i],hashes[i])
	//}
	/*for i:=0;i<Len;i++{
		rawMessage := message.DeserializeMessageWithSignature(requestArr[i])
		m := message.DeserializeClientRequest(rawMessage.Msg)

		if !cryptolib.VerifySig(m.ID, rawMessage.Msg, rawMessage.Sig) {
			log.Printf("[Authentication Error] The signature of client logout request has not been verified.")
			return
		}
	}*/
	batchSize = Len
	requestSize = len(requestArr[0])
	queue.AppendBatch(requestArr)
}

func DeserializeRequests(input []byte)[][]byte{
	var requestArr [][]byte
	json.Unmarshal(input,&requestArr)
	return requestArr
}
