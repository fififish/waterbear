/*
Management of the connections between replicas and clients.
Instead of waiting for the timeout of a connection,
we take a proactive approach where a node is put in 'blacklist' if it cannot be connected.
The nodes that are in blacklist will be moved out of the list if they join the system.
*/

/*TODO:
The current version simply 'Blacklist' nodes that have connection error and puts nodes back if they join the system again.
We need a scheme to optimize this either via
1) enhancing the implementation to avoid slowdown
or
2) put node back to the list when a receiver is reachable again.
Current version puts a node back upon a join request (by the same node)

In the future version, we can integrate this with recovery module
*/

package communication

import (
	"config"
	"log"
	"utils"
	"strings"
)

var connection utils.StringBoolMap
var connectionMap utils.StringIntMap
var maxLimit = 3

/*
Get port number for server api
*/
func GetPortNumber(portNum string) string {
	tmpPort, _ := utils.StringToInt(portNum[1:])
	pn := ":" + utils.IntToString(tmpPort+1000)
	return pn
}

func UpdateAddress(address string) string {
	alist := strings.Split(address, ":")
	pn := GetPortNumber(":" + alist[1])
	newaddress := alist[0] + pn
	return newaddress
}


func IsNotLive(key string) bool {
	result, exist := connection.Get(key)
	if exist {
		return result
	}
	return false
}

/*
Set node to not alive
Input
	key: string format, a node id
*/
func NotLive(key string) {
	v, exist := connectionMap.Get(key)
	if !exist {
		connectionMap.Insert(key, 0)
	} else {
		connectionMap.Insert(key, v+1)
	}
	
	if v > maxLimit {
		connection.Insert(key, true)
	}

}

func SetLive(key string) {
	connectionMap.Insert(key, 0)
	connection.Insert(key, false)
}

func FetchNodesFromConfig() []string{
	return config.FetchNodes()
}

/*
Start connection manager
*/
func StartConnectionManager() {
	log.Printf("starting connection manager")
	connection.Init()
	connectionMap.Init()
}
