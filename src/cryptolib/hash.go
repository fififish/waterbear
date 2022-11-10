
package cryptolib

import (
	"bytes"
	"crypto/sha256"
	"github.com/cbergoon/merkletree"
)

/*
Generate a hash value based on input
Input
	input: []byte type
*/
func GenHash(input []byte) []byte{
	if len(input) == 0 {
		return []byte("")
	}
	switch cryptoOption{
	default:
		h := sha256.New()
		h.Write(input)
		return h.Sum(nil)
	}
	return []byte("")
}

func GenInstanceHash(instanceid []byte, input []byte) []byte{
	if len(input) == 0 && len(instanceid) == 0 {
		return []byte("")
	}

	switch cryptoOption{
	default:
		h := sha256.New()
		h.Write(instanceid)
		h.Write(input)
		return h.Sum(nil)
	}
	return []byte("")
}

func GenABAInstanceHash(input1 []byte, input2 []byte, input3 []byte) []byte{
	if len(input1) == 0 && len(input2) == 0 && len(input3) == 0{
		return []byte("")
	}

	switch cryptoOption{
	default:
		h := sha256.New()
		h.Write(input1)
		h.Write(input2)
		h.Write(input3)
		return h.Sum(nil)
	}
	return []byte("")
}

type MTContent struct{
	x []byte 
}

func (t MTContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (t MTContent) Equals(other merkletree.Content) (bool, error) {
	return bytes.Compare(t.x,other.(MTContent).x)==0, nil
}

func ObtainMerkleNodeHash(input []byte) []byte{
	mn := MTContent{x: input}
	h,_ := mn.CalculateHash()
	return h
}

func GenMerkleTreeRoot(input [][]byte) []byte{
	var list []merkletree.Content 
	for i:=0; i<len(input); i++{
		list = append(list, MTContent{x: input[i]})
	}
	t,err := merkletree.NewTree(list)
	if err != nil{
		return nil 
	}
	mr := t.MerkleRoot()
	return mr 
}



func ObtainMerklePath(input [][]byte) ([][][]byte,[][]int64) {
	
	var result [][][]byte
	var indexresult [][]int64

	var list []merkletree.Content 
	for i:=0; i<len(input); i++{
		list = append(list, MTContent{x: input[i]})
	}
	t,err := merkletree.NewTree(list)
	if err != nil{
		return nil, nil 
	}

	for idx := 0; idx < len(input); idx ++{
		if idx > len(t.Leafs){
			return nil,nil
		}
		curNode := t.Leafs[idx]
		intermediate,index,_ := t.GetMerklePath(curNode.C)
		result = append(result, intermediate)
		indexresult = append(indexresult, index)
	}
	
	return result,indexresult
}
