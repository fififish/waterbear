/*
Transaction types for client requests
*/

package client

import (
	pb "proto/proto/communication"
)

/*
Get message type by integer. This maps to the defined types in communication.proto
*/
var TypeOfTx = map[int]pb.MessageType{
	0: pb.MessageType_WRITE,
}

var TypeTx = map[string]pb.MessageType{
	"write": pb.MessageType_WRITE,
}
var TypeTx_int = map[string]int{
	"write": 0,
}