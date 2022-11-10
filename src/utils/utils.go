
package utils

import (
	"time"
	"strconv"
	"sync"
	"bytes"
	"encoding/binary"
)

type IntValue struct {
	m int
	sync.RWMutex
}

func (s *IntValue) Init() {
	s.Lock()
	defer s.Unlock()
	s.m = 0
}

func (s *IntValue) Get() int {
	s.Lock()
	defer s.Unlock()
	return s.m
}

func (s *IntValue) Increment() int {
	s.Lock()
	defer s.Unlock()
	s.m = s.m + 1
	return s.m
}

func (s *IntValue) Set(n int) {
	s.Lock()
	defer s.Unlock()
	s.m = n
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func Int64ToString(input int64) string {
	return strconv.FormatInt(input, 10)
}


func StringToInt(input string) (int, error) {
	return strconv.Atoi(input)
}

func IntToString(input int) string {
	return strconv.Itoa(input)
}

func StringToInt64(input string) (int64, error) {
	return strconv.ParseInt(input, 10, 64)
}

func StringToBytes(input string) []byte {
	return []byte(input)
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func BytesToString(input []byte) string {
	return string(input[:])
}

func IntToInt64(input int) int64 {
	return int64(input)
}

func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func Int64ToInt(input int64) (int, error) {
	tmp := strconv.FormatInt(input, 10)
	output, err := strconv.Atoi(tmp)
	return output, err
}

func SerializeBytes(input [][]byte) []byte{
	if len(input) == 0{
		return nil 
	}
	var output []byte
	for i:=0; i<len(input); i++{
		output = append(output, input[i] ...)
	}
	return output
}
