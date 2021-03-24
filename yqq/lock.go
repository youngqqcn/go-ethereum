package main

import (
	"fmt"
	"sync"
)

type Store struct {
	lock sync.RWMutex
	datas map[string]interface{}
}

func newStore() *Store{
	st := Store{
		sync.RWMutex{},
		make(map[string]interface{}, 1),
	}
	return &st
}

func (st *Store)set(key string, value interface{}) {
	st.lock.Lock()
	defer  st.lock.Unlock()
	st.datas[ key ] = value
}

func (st *Store)get(key string) interface{} {
	st.lock.RLock()
	defer st.lock.RUnlock()
	return st.datas[key]
}


func main()  {
	st := newStore()
	ch := make(chan bool)
	go func() {
		st.set("good", 111)
		ch <- true
	}()

	select {
	case <-ch:
		value := st.get("good")
		fmt.Println("value:", value)
	}
}