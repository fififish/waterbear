package message

import (
	"encoding/json"
	pb "proto/proto/communication"
)

type ClientRequest struct{
	Type pb.MessageType
	ID int64 
	OP []byte // Message payload. Opt for contract. 
	TS int64 // Timestamp
}


func (r *ClientRequest) Serialize() ([]byte,error){
	jsons, err := json.Marshal(r)
	if err != nil {
		return []byte(""), err
	}
	return jsons, nil
}

func DeserializeClientRequest(input []byte) ClientRequest{
	var clientRequest = new(ClientRequest)
	json.Unmarshal(input, &clientRequest)
	return *clientRequest
}