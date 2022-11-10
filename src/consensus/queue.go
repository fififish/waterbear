/*
This file implements the buffer of messages and client requests (as queue).
*/

package consensus


import (
	"sync"
	"log"
	"bytes"
	"config"
	"cryptolib"
	pb "proto/proto/communication"
	"utils"
)

type Queue struct{
	Q []pb.RawMessage
	R []int64
	sync.RWMutex
}

func (q *Queue) Init(){
	q.Q =[]pb.RawMessage{}
}

func (q *Queue) Length() int{
	return len(q.Q)
}

func (q *Queue) Append(item []byte){
	q.Lock()
	defer q.Unlock()
	q.Q = append(q.Q, pb.RawMessage{Msg: item})
}

func (q *Queue) AppendBatch(items [][]byte) {
	q.Lock()
	defer q.Unlock()
	for i:=0;i<len(items);i++{
		q.Q = append(q.Q,pb.RawMessage{Msg: items[i]})
	}
}

func (q *Queue) Grab() []pb.RawMessage{
	q.Lock()
	defer q.Unlock()
	return q.Q
}


func (q *Queue) GrabWithMaxLen() []pb.RawMessage{
	q.Lock()
	defer q.Unlock()
	if len(q.Q) > config.MaxBatchSize(){
		return q.Q[:config.MaxBatchSize()]
	}
	return q.Q
}

func (q *Queue) GrabWtihMaxLenAndClear() []pb.RawMessage{
	q.Lock()
	defer q.Unlock()
	if len(q.Q) > config.MaxBatchSize(){
		ret := q.Q[:config.MaxBatchSize()]
		q.Q = q.Q[config.MaxBatchSize():]
		return ret
	}
	ret := q.Q
	q.Q = []pb.RawMessage{}
	return ret
}


func (q *Queue) GrabQLen() int{
	q.Lock()
	defer q.Unlock()
	return len(q.Q)
}

func (q *Queue) GrabFirst() (pb.RawMessage,bool){
	q.Lock()
	defer q.Unlock()
	if len(q.Q) == 0{
		var empty pb.RawMessage
		return empty,false
	}
	return q.Q[0],true
}

func (q *Queue) FetchFirst() (string,bool){
	q.Lock()
	defer q.Unlock()
	if (len(q.Q))==0{
		return "",false
	}
	return utils.BytesToString(cryptolib.GenHash(q.Q[0].GetMsg())),true
}

func (q *Queue) Contains(item pb.RawMessage) (int,bool){
	q.RLock()
	defer q.RUnlock()
	if (len(q.Q)==0){
		return 0,false
	}
	for i:=0; i<len(q.Q); i++{
		if (bytes.Compare(item.GetMsg(),q.Q[i].GetMsg())==0){
			return i,true
		}
	}
	return 0,false
}

func (q *Queue) RemoveFirst(){
	q.Lock()
	defer q.Unlock()
	if len(q.Q) == 0 {
		return 
	}
	q.Q = q.Q[1:]
}

func (q *Queue) Remove(hash string, msgs []pb.RawMessage){
	q.Lock()
	defer q.Unlock()

	if (len(q.Q)==0){
		return
	}

	var tmp []pb.RawMessage
	for i:=0; i<len(q.Q); i++{
		exist := false
		for j:=0; j<len(msgs); j++{
			item := msgs[j]
			if (bytes.Compare(item.GetMsg(),q.Q[i].GetMsg())==0){
				exist = true
			}
		}
		if (!exist){
			tmp = append(tmp, q.Q[i])
		}
	}
	q.Q = tmp
}

func (q *Queue) RemoveItem(hash []byte){
	q.Lock()
	defer q.Unlock()

	if (len(q.Q)==0){
		return
	}

	var tmp []pb.RawMessage
	for i:=0; i<len(q.Q); i++{
		exist := false
		if (bytes.Compare(hash,cryptolib.GenHash(q.Q[i].GetMsg()))==0){
			exist = true
		}
		if (!exist){
			tmp = append(tmp, q.Q[i])
		}
	}
	q.Q = tmp
}

func (q *Queue) IsEmpty() bool{
	return len(q.Q)==0
}

func (q *Queue) Clear(){
	q.Lock()
	defer q.Unlock()
	q.Q = []pb.RawMessage{}
}

func (q *Queue) ClearFraction(size int){
	q.Lock()
	defer q.Unlock()
	if len(q.Q) > size{
		q.Q = q.Q[size:]
	}else{
		q.Q = []pb.RawMessage{}
	}
}

func (q *Queue) PrintQueue (){
	for i:=0; i< len(q.Q); i++ {
		log.Println("Number %d: %s", i, q.Q[i].GetMsg())
	}
}