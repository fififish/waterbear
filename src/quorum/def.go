/*
Definitions for quorum class
*/

package quorum

import (
	"sync"
	"utils"
	"message"
)

/*
IntBuffer
	Map int to a set (int64)
	Used to check (given a sequence number) whether a replica receives matching messages from a quorum of replicas
V
	Map int to a list of MessageWithSignature
*/
type INTBUFFER struct {
	IntBuffer map[int]utils.Set
	V         map[int][]message.MessageWithSignature
	sync.RWMutex
}

/*
Buffer
	map string to a set (int64)
	Used to check (given a hash) whether a replica receives matching messages form a quorum of replicas
*/
type BUFFER struct {
	Buffer map[string]utils.Set
	sync.RWMutex
}

type CERTIFICATE struct {
	Certificate map[string][][]byte
	Identities map[string][]int64
	sync.RWMutex
}

func (p *CERTIFICATE) Init() {
	p.Lock()
	defer p.Unlock()
	p.Certificate = make(map[string][][]byte)
	p.Identities = make(map[string][]int64)
}

func (p *CERTIFICATE) Insert(key string, nid int64, msg []byte) {
	p.Lock()
	defer p.Unlock()
	certi, _ := p.Certificate[key]
	certi = append(certi, msg)
	p.Certificate[key] = certi
	identities, _ := p.Identities[key]
	identities = append(identities, nid)
	p.Identities[key] = identities 
}

func (p *CERTIFICATE) Clear(key string) {
	p.Lock()
	defer p.Unlock()
	
	delete(p.Certificate,key)
	delete(p.Identities,key)
}

func FetchCer(key string) []byte{
	result, exist := cer.Certificate[key]
	result2, exist2 := cer.Identities[key]
	if !exist || !exist2{
		return nil 
	}
	
	h:= utils.StringToBytes(key)
	msg := message.Signatures{
		Hash:	h, 
		Sigs: result,
		IDs: result2,
	}

	msgbyte, err := msg.Serialize()
	if err !=nil{
		return nil
	}
	return msgbyte
}

/*
Initialize INTBUFFER
*/
func (b *INTBUFFER) Init(n int) {
	b.IntBuffer = make(map[int]utils.Set)
	b.V = make(map[int][]message.MessageWithSignature)
}

/*
Add a value to INTBUFFER
Input
	input: integer, sequence number (int type)
	id: id of the node (int64 type)
	msg: MessageWithSignature, deserialized message
*/
func (b *INTBUFFER) InsertValue(input int, id int64, msg message.MessageWithSignature) {
	b.Lock()
	defer b.Unlock()
	_, exist := b.IntBuffer[input]
	if exist {
		s := b.IntBuffer[input]
		len1 := s.Len()
		s.AddItem(id)
		b.IntBuffer[input] = s
		len2 := s.Len()
		if len2 > len1 {
			b.V[input] = append(b.V[input], msg)
		}
	} else {
		s := *utils.NewSet()
		s.AddItem(id)
		b.IntBuffer[input] = s
		b.V[input] = append(b.V[input], msg)
	}
}

func (b *INTBUFFER) SettValue(input int, msg []message.MessageWithSignature, ids utils.Set) {
	b.Lock()
	defer b.Unlock()
	b.IntBuffer[input] = ids
	b.V[input] = msg

}

/*
Add value to V only. Used for view changes.
Input
	input: sequence number (int type)
	value: a list of messages
*/
func (b *INTBUFFER) InsertV(input int, value []message.MessageWithSignature) {
	b.Lock()
	defer b.Unlock()
	b.V[input] = value
}

/*
Get the number of messages fro IntBuffer
Input
	input: sequence number
Output
	int: number of messages (ids of the nodes)
*/
func (b *INTBUFFER) GetLen(input int) int {
	b.RLock()
	defer b.RUnlock()
	_, exist := b.IntBuffer[input]
	if exist {
		s := b.IntBuffer[input]
		return s.Len()
	}
	return 0
}

/*
Reset INTBUFFER given an input
Input
	input: int type, usually a sequence number
*/
func (b *INTBUFFER) Clear(input int) {
	b.Lock()
	defer b.Unlock()
	_, exist := b.IntBuffer[input]
	if exist {
		delete(b.IntBuffer, input)
	}

	_, exist1 := b.V[input]
	if exist1 {
		delete(b.V, input)
	}
}


/*
Initialize BUFFER
*/
func (b *BUFFER) Init() {
	b.Buffer = make(map[string]utils.Set)
}

/*
Add a value to BUFFER
Input
	input: hash of the message (key of the Buffer, string type)
	msg: message that needs to be added ([]byte type)
	nid: replica id (int64 type)
	step: PP (prepare) or CM (commit)
	seq: sequence number (int type)
*/
func (b *BUFFER) InsertValue(input string, nid int64, msg []byte, step Step) {
	b.Lock()
	defer b.Unlock()
	_, exist := b.Buffer[input]
	if exist {
		s := b.Buffer[input]
		len := s.Len()
		if len == 0 {
			s = *utils.NewSet()
		}
		s.AddItem(nid)
		b.Buffer[input] = s 
		
		if msg == nil{
			return 
		}
		len2 := s.Len()
		if len2 > len{
			cer.Insert(input,nid,msg)
		}
		
	} else {
		s := *utils.NewSet()
		s.AddItem(nid)
		b.Buffer[input] = s
		if msg == nil{
			return 
		}
		cer.Insert(input,nid,msg)
	}
}

/*
Get the number of messages in BUFFER
Input
	input: usually a hash (string type)
Output
	int: number of messages
*/
func (b *BUFFER) GetLen(input string) int {
	b.RLock()
	defer b.RUnlock()
	_, exist := b.Buffer[input]
	if exist {
		s := b.Buffer[input]
		return s.Len()
	}
	return 0
}

/*
Reset BUFFER given an input
Input
	input: usually a hash (string type)
*/
func (b *BUFFER) Clear(input string) {
	b.Lock()
	defer b.Unlock()
	_, exist := b.Buffer[input]
	if exist {
		delete(b.Buffer, input)
	}
}
