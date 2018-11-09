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
	limit []uint64    // for storing message history
	limit_pt int      // for message index
	limit_max int     // for max number of stored message

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
	mst.limit_max = 1000  // todo: check size.
	mst.limit_pt = 0
	mst.limit = make([]uint64,mst.limit_max)
	fmt.Println("Initialize LocalStore ",mst.store)
}


//todo: This is not efficient store. So we need to fix it.
func (mst *MessageStore) AddMessage(msgType string, chType int, mid uint64, src uint64, dst uint64, arg string){

	mes := message{msgType,chType, mid, src, dst, arg}
//	fmt.Printf("AddMessage %v\n",mes)
//	fmt.Printf("ls.store %v %d \n",ls.store, mid)
	mst.mutex.Lock()
	if mst.limit[mst.limit_pt] != 0 { // ring buffer, delete last one.
		delete(mst.store, mst.limit[mst.limit_pt])
	}
	mst.store[mid] = mes
	mst.limit[mst.limit_pt] = mid
	mst.limit_pt = (mst.limit_pt+1)%mst.limit_max
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

