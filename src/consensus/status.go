package consensus

import (
	"fmt"
	"log"
	"utils"
	"logging"
	"sync"
)

var epoch utils.IntValue //epoch number
var curStatus CurStatus
var output utils.ByteSet //output set
var astatus utils.IntBoolMap//rbc status
var bstatus utils.IntBoolMap//aba status
var fstatus utils.IntBoolMap

var avalues utils.IntBoolMap //aba outputvalues
var otherlock utils.IntValue

var elock sync.Mutex

var outputCount utils.IntValue ////number of decide 0 or 1
var outputSize utils.IntValue  //number of decide 1

func InitStatus(n int){
	baseinstance =  0 //Set up baseinstance to avoid conflict
	output = *utils.NewByteSet()
	astatus.Init()
	avalues.Init()
	bstatus.Init()
	fstatus.Init()
	outputCount.Init()
	outputSize.Init()
	epoch.Increment()
	otherlock.Init()
	log.Printf("Starting epoch %v", epoch.Get())
	p := fmt.Sprintf("*************Starting epoch %v************", epoch.Get())
	logging.PrintLog(verbose, logging.NormalLog, p)
}