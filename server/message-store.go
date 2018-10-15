package main

import (
	"fmt"
	"sync"
)

// MessageStore saves all massage in memory
//   this should change into other distributed service.
/*
type MessageStore interface{
	AddMessage(msgType string, chType int, mid uint64, src uint64, dst uint64, arg string)
	getSrcId(mid uint64) uint64
}*/

type message struct{
	msgType string
	chType int
	mid uint64
	src uint64
	dst uint64
	arg string
}


// real struct for MessageStore
type MessageStore struct {
	store map[uint64]message
	mutex sync.RWMutex
}


// CreateMessageStore creates base dataset
func CreateLocalMessageStore()  *MessageStore {
	mst := &MessageStore{}
	mst.init()
	return mst
}

func (mst *MessageStore) init(){
//	fmt.Println("Initialize LocalStore")
	mst.store = make(map[uint64]message)
	mst.mutex = sync.RWMutex{}
	fmt.Println("Initialize LocalStore ",mst.store)
}

func (mst *MessageStore) AddMessage(msgType string, chType int, mid uint64, src uint64, dst uint64, arg string){

	mes := message{msgType,chType, mid, src, dst, arg}
//	fmt.Printf("AddMessage %v\n",mes)
//	fmt.Printf("ls.store %v %d \n",ls.store, mid)
	mst.mutex.Lock()
	mst.store[mid] = mes
	mst.mutex.Unlock()
//	fmt.Println("OK.")
}

func (mst *MessageStore) getSrcId(mid uint64) uint64{
	mst.mutex.RLock()
	mes, ok  := mst.store[mid]
	mst.mutex.RUnlock()
	if !ok {
		fmt.Println("Cant find message id Error!")
		return 0
	}
	return mes.src
}

