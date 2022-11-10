/*
Loading configuration parameters from configuration file.
The file path is hard-coded to ../etc/conf.json.
TODO:
1) The file path can be tailored to different paths according to the deployment environments.
2) Need to validate the format of the configuration parameters
*/

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	logging "logging"
	"os"
	"path"
	"utils"
	"strings"
)

var nodeIDs []string
var nodes map[string]string
var portMap map[string]string
var nodesReverse map[string]string
var verbose bool
var maxBatchSize int
var sleepTimer int
var clientTimer int
var broadcastTimer int
var evalMode int
var evalInterval int
var local bool
var maliciousNode bool
var maliciousMode int
var maliciousNID []int64

var tParameter int

var cryptoOpt int
var splitPorts bool
var logOpt int
var consensus int
var rbc int

type System struct {
	MaxBatchSize       int       `json:"maxBatchSize"`                // Max batch size for consensus
	SleepTimer         int       `json:"sleepTimer"`                  // Timer for the while loops to monitor the status of requests. Should be a small value
	ClientTimer        int       `json:"clientTimer"`                 // Timer for clients to monitor the responses and see whether the requests should be re-transmitted.
	BroadcastTimer     int       `json:"broadcastTimer"`              // Timer used for replicas to send gRPC messages to each other. Should be set to a value that is close to RTT
	TParameter		   int       `json:"tParameter"`                  //coin set 1 when round less than TParameter
	Verbose            bool      `json:"verbose"`                     // Whether log messages should be printed.
	EvalMode           int       `json:"evalMode"`                    // Evaluation mode.
	EvalInterval       int       `json:"evalInterval"`                // Interval for assessing throughput
	CryptoOpt          int       `json:"cryptoOpt"`                   // Crypto library option
	LogOpt             int       `json:"logOpt"`  
	Local              bool      `json:"local"`                       // Local or not
	MaliciousNode      bool      `json:"maliciousNode"`               // Simulate a simple malicious node
	MaliciousMode      int       `json:"maliciousMode"`               //
	MaliciousNID       string    `json:"maliciousNID"`                // Malicious node id
	SplitPorts         bool      `json:"splitPorts"`                  // Split ports for request handler and server
	Consensus          int       `json:"consensus`					  // Protocol
	RBCType			   int       `json:"RBCType"`                     //RBC
	Replicas           []Replica `json:"replicas"`                    // Replica information
}

type Replica struct {
	ID   string `json:"id"`   // ID of the node
	Host string `json:"host"` // IP address
	Port string `json:"port"` // Port number
}

func LoadConfig() bool {

	nodes = make(map[string]string)
	nodesReverse = make(map[string]string)
	portMap = make(map[string]string)
	nodeIDs = make([]string, 0)


	exepath, err := os.Executable()
	if err != nil {
		p := fmt.Sprintf("[Configuration Error]  Failed to get path for the executable")
		logging.PrintLog(true, logging.ErrorLog, p)
		os.Exit(1)
		return false
	}

	p1 := path.Dir(exepath)
	homepath := path.Dir(p1)
	//fmt.Println("homepath %s", homepath)

	defaultFileName := homepath + "/etc/conf.json"
	f, err := os.Open(defaultFileName)
	if err != nil {
		p := fmt.Sprintf("[Configuration Error]  Failed to open config file: %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		os.Exit(1)
		return false
	}
	defer f.Close()
	var system System
	byteValue, _ := ioutil.ReadAll(f)

	json.Unmarshal(byteValue, &system)
	for i := 0; i < len(system.Replicas); i++ {
		nodeIDs = append(nodeIDs, system.Replicas[i].ID)
		addr := system.Replicas[i].Host + ":" + system.Replicas[i].Port
		nodes[system.Replicas[i].ID] = addr
		nodesReverse[addr] = system.Replicas[i].ID
		portMap[system.Replicas[i].ID] = ":" + system.Replicas[i].Port
	}



	maxBatchSize = system.MaxBatchSize

	sleepTimer = system.SleepTimer
	clientTimer = system.ClientTimer
	broadcastTimer = system.BroadcastTimer

	tParameter = system.TParameter
	if tParameter > 0{
		log.Printf("********************************************************************************\n")
		log.Printf("** Note coin in RABA will be always 1 when round is equal or lesser than %d.\n",tParameter)
		log.Printf("********************************************************************************\n")
	}

	verbose = system.Verbose
	evalMode = system.EvalMode
	evalInterval = system.EvalInterval
	local = system.Local
	cryptoOpt = system.CryptoOpt
	maliciousNode = system.MaliciousNode
	maliciousMode = system.MaliciousMode
	splitPorts = system.SplitPorts
	logOpt = system.LogOpt

	//get malicious node id
	mnidList := system.MaliciousNID
	mlist := strings.Split(mnidList, ",")
	consensus = system.Consensus
	rbc = system.RBCType

	for i := 0; i < len(mlist); i++ {
		tmp, err := utils.StringToInt64(mlist[i])
		if err != nil {
			fmt.Printf("Incorrect list of malicious node. Please double check conf.json, %v, mlist %v", mlist[i], mlist)
			continue
		}
		maliciousNID = append(maliciousNID, tmp)
	}

	return true
}


func FetchLogOpt() int {
	return logOpt
}

func TParameter() int {
	return tParameter
}

func MaxBatchSize() int {
	return maxBatchSize
}

// IP address of a node
func FetchAddress(id string) string {
	return nodes[id]
}

func FetchPort(id string) string {
	return portMap[id]
}

// Get list of nodes
func FetchNodes() []string {
	return nodeIDs
}

// Total number of replicas
func FetchNumReplicas() int {
	return len(nodes)
}

// Get the id (string format) of a node given the address (ip and port number)
func FetchReplicaID(addr string) string {
	return nodesReverse[addr]
}

func FetchSleepTimer() int {
	return sleepTimer
}

func FetchClientTimer() int {
	return clientTimer
}

func FetchBroadcastTimer() int {
	return broadcastTimer
}

func FetchVerbose() bool {
	return verbose
}

func EvalMode() int {
	return evalMode
}

func CryptoOption() int {
	return cryptoOpt
}

func EvalInterval() int {
	return evalInterval
}

func Local() bool {
	return local
}

func MaliciousNode() bool {
	return maliciousNode
}

func MaliciousMode() int {
	return maliciousMode
}

func MaliciousNID(nid int64) bool {
	for i := 0; i < len(maliciousNID); i++ {
		if nid == maliciousNID[i] {
			return true
		}
	}
	return false
}

func SplitPorts() bool {
	return splitPorts
}

func Consensus() int{
	return consensus
}

func RBCType() int{
	return rbc
}