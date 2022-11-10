

package communication

import (
	"sync"
	pb "proto/proto/communication"
)


type AddrConnMap struct {
	m map[string]pb.SendClient
	idmap map[string]string
	count map[string]int
	sync.RWMutex
}

func (s *AddrConnMap) Init() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[string]pb.SendClient)
	s.count = make(map[string]int)
	s.idmap = make(map[string]string)
}

func (s *AddrConnMap) Get(key string) (pb.SendClient, bool) {
	s.Lock()
	defer s.Unlock()
	_, exist := s.m[key]
	if exist {
		return s.m[key], true
	}
	var emptyConn pb.SendClient
	return emptyConn, false
}

func (s *AddrConnMap) GetID(key string) string {
	s.Lock()
	defer s.Unlock()
	val, exist := s.idmap[key]
	if exist {
		return val
	}
	return ""
}


func (s *AddrConnMap) GetCurCount(key string) int {
	s.Lock()
	defer s.Unlock()
	_, exist := s.count[key]
	if exist {
		return s.count[key]
	}
	return 0 
}


func (s *AddrConnMap) ResetCount(key string) {
	s.Lock()
	defer s.Unlock()
	s.count[key] = 0
}


func (s *AddrConnMap) GetAll() ([]pb.SendClient) {
	s.Lock()
	defer s.Unlock()
	var result []pb.SendClient 
	for _,v := range s.m{
		result = append(result,v)
	}
	return result
}

func (s *AddrConnMap) Insert(key string, value pb.SendClient) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (s *AddrConnMap) InsertID(key string, value string) {
	s.Lock()
	defer s.Unlock()
	s.idmap[key] = value
}

func (s *AddrConnMap) IncrementCount(key string) {
	s.Lock()
	defer s.Unlock()
	_,exist := s.count[key]
	if exist{
		s.count[key] = s.count[key] + 1
	}else{
		s.count[key] = 0 
	}
}
