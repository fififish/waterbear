

package message

import (
	"cryptolib"
	"encoding/json"
	pb "proto/proto/communication"
)

type ReplicaMessage struct {
	Mtype           TypeOfMessage
	Instance        int
	Source          int64
	Hash            []byte
	TS              int64
	Payload 		[]byte //message payload
	Value           int
	Maj             int
	Round           int
	Epoch			int
	MTRoot			[]byte
	MTBranch		[][]byte
	MTIndex			[]int64
}


type MessageWithSignature struct {
	Msg     []byte
	Sig     []byte
}

type RawOPS struct{
	OPS     []pb.RawMessage
}

type CBCMessage struct{
	Value  map[int][]byte
	RawData [][]byte
	MerkleBranch [][][]byte 
	MerkleIndexes [][]int64
}

type Signatures struct {
	Hash   []byte
	Sigs   [][]byte
	IDs    []int64
}

func (r *Signatures) Serialize() ([]byte, error) {
	jsons, err := json.Marshal(r)
	if err != nil {
		return []byte(""), err
	}
	return jsons, nil
}

func DeserializeSignatures(input []byte) ([]byte, [][]byte,[]int64) {
	var sigs = new(Signatures)
	json.Unmarshal(input, &sigs)
	return sigs.Hash, sigs.Sigs, sigs.IDs
}

func (r *CBCMessage) Serialize() ([]byte, error) {
	jsons, err := json.Marshal(r)
	if err != nil {
		return []byte(""), err
	}
	return jsons, nil
}

func DeserializeCBCMessage(input []byte) CBCMessage {
	var cbcMessage = new(CBCMessage)
	json.Unmarshal(input, &cbcMessage)
	return *cbcMessage
}


func (r *RawOPS) Serialize() ([]byte, error) {
	jsons, err := json.Marshal(r)
	if err != nil {
		return []byte(""), err
	}
	return jsons, nil
}

/*
Deserialize []byte to MessageWithSignature
*/
func DeserializeRawOPS(input []byte) RawOPS {
	var rawOPS = new(RawOPS)
	json.Unmarshal(input, &rawOPS)
	return *rawOPS
}

/*
Serialize MessageWithSignature
*/
func (r *MessageWithSignature) Serialize() ([]byte, error) {
	jsons, err := json.Marshal(r)
	if err != nil {
		return []byte(""), err
	}
	return jsons, nil
}

/*
Deserialize []byte to MessageWithSignature
*/
func DeserializeMessageWithSignature(input []byte) MessageWithSignature {
	var messageWithSignature = new(MessageWithSignature)
	json.Unmarshal(input, &messageWithSignature)
	return *messageWithSignature
}

func (r *ReplicaMessage) Serialize() ([]byte, error) {
	jsons, err := json.Marshal(r)
	if err != nil {
		return []byte(""), err
	}
	return jsons, nil
}

func DeserializeReplicaMessage(input []byte) ReplicaMessage {
	var replicaMessage = new(ReplicaMessage)
	json.Unmarshal(input, &replicaMessage)
	return *replicaMessage
}

func SerializeWithSignature(id int64, msg []byte) ([]byte, error) {
	request := MessageWithSignature{
		Msg:     msg,
		Sig:     cryptolib.GenSig(id, msg),
	}

	requestSer, err := request.Serialize()
	if err != nil {
		return []byte(""), err
	}
	return requestSer, err
}

func SerializeWithMAC(id int64, dest int64, msg []byte) ([]byte, error) {
	request := MessageWithSignature{
		Msg:     msg,
		Sig:     cryptolib.GenMAC(id, msg),
	}

	requestSer, err := request.Serialize()
	if err != nil {
		return []byte(""), err
	}
	return requestSer, err
}