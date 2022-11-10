/*
Verify whether a replica/client has received matching messages from sufficient number of replicas
*/
package quorum

// Used for normal operation
var buffer BUFFER  //prepare certificate. Client uses it as reply checker.
var bufferc BUFFER //commit certificate. Client API uses it as reply checker.
var cer CERTIFICATE //used for vcbc only. Store the set of signatures

var n int
var f int
var quorum int
var squorum int
var half int


type Step int32

const (
	PP Step = 0
	CM Step = 1
)

/*
Clear in-memory data for view changes.
*/
func ClearCer() {
	bufferc.Init()
	buffer.Init()
	cer.Init()
}



func Add(id int64, hash string, msg []byte, step Step) {
	switch step {
	case PP:
		buffer.InsertValue(hash, id, msg, step)
	case CM:
		bufferc.InsertValue(hash, id, msg, step)
	}
}

func GetBuffercList(key string) []int64 {
	_, exist := bufferc.Buffer[key]
	if exist {
		s := bufferc.Buffer[key]
		return s.SetList()
	} else {
		return []int64{}
	}
}

func CheckQuorum(input string, step Step) bool {
	switch step {
	case PP:
		return buffer.GetLen(input) >= quorum
	case CM:
		return bufferc.GetLen(input) >= quorum
	}

	return false
}

func CheckCurNum(input string, step Step) int{
	switch step {
	case PP:
		return buffer.GetLen(input) 
	case CM:
		return bufferc.GetLen(input) 
	}
	return 0 
}

func CheckEqualQuorum(input string, step Step) bool {
	switch step {
	case PP:
		return buffer.GetLen(input) == quorum
	case CM:
		return bufferc.GetLen(input) == quorum
	}

	return false
}

func CheckSmallQuorum(input string,step Step) bool {
	switch step {
	case PP:
		return buffer.GetLen(input) >= squorum
	case CM:
		return bufferc.GetLen(input) >= squorum
	}

	return false
}

func CheckOverSmallQuorum(input string) bool {
	return bufferc.GetLen(input) >= squorum
}

func CheckEqualSmallQuorum(input string) bool {
	return bufferc.GetLen(input) == squorum
}

func ClearBuffer(input string, step Step) {
	switch step {
	case PP:
		buffer.Clear(input)
	case CM:
		bufferc.Clear(input)
	}
}

func ClearBufferPC(input string) {
	buffer.Clear(input)
	bufferc.Clear(input)
	cer.Clear(input)
}




func QuorumSize() int {
	return quorum
}

func HalfSize() int {
	return half
}

func SQuorumSize() int {
	return squorum
}

func FSize() int {
	return f
}

func NSize() int {
	return n
}



func SetQuorumSizes(num int) {
	n = num
	f = (n - 1) / 3
	quorum = (n + f + 1) / 2
	if (n+f+1)%2 > 0 {
		quorum += 1
	}
	squorum = f + 1
}

func CheckOverHalf(input string) bool {
	return bufferc.GetLen(input) >= half
}

func CheckHalf(input string) bool {
	return bufferc.GetLen(input) == half
}

func StartQuorum(num int) {

	n = num

	f = (n - 1) / 3
	quorum = (n + f + 1) / 2
	if (n+f+1)%2 > 0 {
		quorum += 1
	}
	squorum = f + 1

	bufferc.Init()
	buffer.Init()
	cer.Init()
}
