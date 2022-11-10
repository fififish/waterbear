package utils

import (
	"sync"
)

/***************[]byte set*******************/
type ByteSet struct {
	m map[string]bool
	sync.RWMutex
}

func NewByteSet() *ByteSet {
	return &ByteSet{
	  m: map[string]bool{},
	}
}
  
func (s *ByteSet) AddItem(item []byte) {
	s.Lock()
	defer s.Unlock()
	s.m[BytesToString(item)] = true
}
  
func (s *ByteSet) RemoveItem(item []byte) {
	s.Lock()
	s.Unlock()
	delete(s.m, BytesToString(item))
}

func (s *ByteSet) HasItem(item []byte) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[BytesToString(item)]
	return ok
}

func (s *ByteSet) Len() int {
	return len(s.SetList())
}

func (s *ByteSet) ClearSet() {
	_ = s.Lock
	defer s.Unlock()
	s.m = map[string]bool{}
}

func (s *ByteSet) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *ByteSet) SetList() [][]byte {
	s.RLock()
	defer s.RUnlock()
	var list [][]byte
	for item := range s.m {
		list = append(list, StringToBytes(item))
	}
	return list
}


/***************Int64 set*******************/
type Set struct {
	m map[int64]bool
	sync.RWMutex
}

func NewSet() *Set {
	return &Set{
	  m: map[int64]bool{},
	}
}
  
func (s *Set) AddItem(item int64) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}
  
func (s *Set) RemoveItem(item int64) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

func (s *Set) HasItem(item int64) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) Len() int {
	return len(s.SetList())
}

func (s *Set) ClearSet() {
	_ = s.Lock
	defer s.Unlock()
	s.m = map[int64]bool{}
}

func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *Set) SetList() []int64 {
	s.RLock()
	defer s.RUnlock()
	list := []int64{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}



/*************** Int64 set without lock *******************/
type Set_N struct {
	m map[int64]bool
}

func NewSet_N() *Set_N {
	return &Set_N{
		m: map[int64]bool{},
	}
}

func (s *Set_N) AddItem(item int64) {
	s.m[item] = true
}

func (s *Set_N) RemoveItem(item int64) {
	delete(s.m, item)
}

func (s *Set_N) HasItem(item int64) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set_N) Len() int {
	return len(s.SetList())
}

func (s *Set_N) ClearSet() {
	s.m = map[int64]bool{}
}

func (s *Set_N) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *Set_N) SetList() []int64 {
	list := []int64{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}







/***************Int set*******************/
type IntSet struct {
	m map[int]bool
	count map[int]int
	sync.RWMutex
}

func NewIntSet() *IntSet {
	return &IntSet{
	  m: map[int]bool{},
	  count: map[int]int{},
	}
}

func (s *IntSet) Init(){
	s.Lock()
	defer s.Unlock()
	s.m = make(map[int]bool)
	s.count = make(map[int]int)
}

func IntSetAddItem(set IntSet,item int){
	set.AddItem(item)
}
func (s *IntSet) AddItem(item int) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
	v,exist := s.count[item]
	if !exist{
		s.count[item] = 1
	}else{
		s.count[item] = v+1
	}
}
func (s *IntSet) RemoveItem(item int) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
	delete(s.count, item)
}

func (s *IntSet) IsTrue(item int) bool{
	s.Lock()
	defer s.Unlock()
	if s.m[item] == true{
		return true
	}
	return false
}

func (s *IntSet) GetCount(item int) int{
	s.Lock()
	defer s.Unlock()
	v,exist := s.count[item]
	if !exist{
		return 0 
	}
	return v
}

func (s *IntSet) Len() int {
	return len(s.IntSetList())
}

func (s *IntSet) IntSetList() []int {
	s.RLock()
	defer s.RUnlock()
	list := []int{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s *IntSet) SetValue(seqSet []int){
	for _,v := range seqSet{
		s.AddItem(v)
	}
}





/***************Int set without lock *******************/
type IntSet_N struct {
	m map[int]bool
	count map[int]int
}

func NewIntSet_N() *IntSet_N {
	return &IntSet_N{
		m: map[int]bool{},
		count: map[int]int{},
	}
}

func (s *IntSet_N) Init(){
	s.m = make(map[int]bool)
	s.count = make(map[int]int)
}

func IntSetAddItem_N(set IntSet_N,item int){
	set.AddItem(item)
}
func (s *IntSet_N) AddItem(item int) {
	s.m[item] = true
	v,exist := s.count[item]
	if !exist{
		s.count[item] = 1
	}else{
		s.count[item] = v+1
	}
}
func (s *IntSet_N) RemoveItem(item int) {
	delete(s.m, item)
	delete(s.count, item)
}

func (s *IntSet_N) IsTrue(item int) bool{
	if s.m[item] == true{
		return true
	}
	return false
}

func (s *IntSet_N) GetCount(item int) int{
	v,exist := s.count[item]
	if !exist{
		return 0
	}
	return v
}

func (s *IntSet_N) Len() int {
	return len(s.IntSetList())
}

func (s *IntSet_N) IntSetList() []int {
	list := []int{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s *IntSet_N) SetValue(seqSet []int){
	for _,v := range seqSet{
		s.AddItem(v)
	}
}