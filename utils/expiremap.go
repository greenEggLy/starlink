package utils

import (
	"sync"
	"time"
)

type val[V []string | string | chan []byte | chan string] struct {
	data        V
	expiredVime int64
}

const delChannelCap = 100

type ExpiredMap[K string, V []string | string | chan []byte | chan string] struct {
	m       map[K]*val[V]
	timeMap map[int64][]K
	lck     *sync.Mutex
	stop    chan struct{}
}

func NewExpiredMap[K string, V []string | string | chan []byte | chan string]() *ExpiredMap[K, V] {
	e := ExpiredMap[K, V]{
		m:       make(map[K]*val[V]),
		lck:     new(sync.Mutex),
		timeMap: make(map[int64][]K),
		stop:    make(chan struct{}),
	}
	go e.run(time.Now().Unix())
	return &e
}

type delMsg[K string, V []string | string | chan []byte | chan string] struct {
	keys []K
	t    int64
}

// background goroutine 主动删除过期的key
// 数据实际删除时间比应该删除的时间稍晚一些，这个误差会在查询的时候被解决。
func (e *ExpiredMap[K, V]) run(now int64) {
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()

	delCh := make(chan *delMsg[K, V], delChannelCap)
	go func() {
		for v := range delCh {
			e.multiDelete(v.keys, v.t)
		}
	}()
	for {
		select {
		case <-t.C:
			now++ //这里用now++的形式，直接用time.Now().Unix()可能会导致时间跳过1s，导致key未删除。
			e.lck.Lock()
			if keys, found := e.timeMap[now]; found {
				e.lck.Unlock()
				delCh <- &delMsg[K, V]{keys: keys, t: now}
			} else {
				e.lck.Unlock()
			}
		case <-e.stop:
			close(delCh)
			return
		}
	}
}

func (e *ExpiredMap[K, V]) Set(key K, value V, expireSeconds int64) {
	if expireSeconds <= 0 {
		return
	}
	e.lck.Lock()
	defer e.lck.Unlock()
	expiredTime := time.Now().Unix() + expireSeconds
	e.m[key] = &val[V]{
		data:        value,
		expiredVime: expiredTime,
	}
	e.timeMap[expiredTime] = append(e.timeMap[expiredTime], key) //过期时间作为key，放在map中
}

func (e *ExpiredMap[K, V]) Get(key K) (found bool, value V) {
	e.lck.Lock()
	defer e.lck.Unlock()
	if found = e.checkDeleteKey(key); !found {
		return
	}
	value = e.m[key].data
	return
}

func (e *ExpiredMap[K, V]) Delete(key K) {
	e.lck.Lock()
	delete(e.m, key)
	e.lck.Unlock()
}

func (e *ExpiredMap[K, V]) Remove(key K) {
	e.Delete(key)
}

func (e *ExpiredMap[K, V]) multiDelete(keys []K, t int64) {
	e.lck.Lock()
	defer e.lck.Unlock()
	delete(e.timeMap, t)
	for _, key := range keys {
		delete(e.m, key)
	}
}

func (e *ExpiredMap[K, V]) Length() int { //结果是不准确的，因为有未删除的key
	e.lck.Lock()
	defer e.lck.Unlock()
	return len(e.m)
}

func (e *ExpiredMap[K, V]) Size() int {
	return e.Length()
}

// 返回key的剩余生存时间 key不存在返回负数
func (e *ExpiredMap[K, V]) VVL(key K) int64 {
	e.lck.Lock()
	defer e.lck.Unlock()
	if !e.checkDeleteKey(key) {
		return -1
	}
	return e.m[key].expiredVime - time.Now().Unix()
}

func (e *ExpiredMap[K, V]) Clear() {
	e.lck.Lock()
	defer e.lck.Unlock()
	e.m = make(map[K]*val[V])
	e.timeMap = make(map[int64][]K)
}

func (e *ExpiredMap[K, V]) Close() { // todo 关闭后在使用怎么处理
	e.lck.Lock()
	defer e.lck.Unlock()
	e.stop <- struct{}{}
	//e.m = nil
	//e.timeMap = nil
}

func (e *ExpiredMap[K, V]) Stop() {
	e.Close()
}

func (e *ExpiredMap[K, V]) DoForEach(handler func(interface{}, interface{})) {
	e.lck.Lock()
	defer e.lck.Unlock()
	for k, v := range e.m {
		if !e.checkDeleteKey(k) {
			continue
		}
		handler(k, v)
	}
}

func (e *ExpiredMap[K, V]) DoForEachWithBreak(handler func(interface{}, interface{}) bool) {
	e.lck.Lock()
	defer e.lck.Unlock()
	for k, v := range e.m {
		if !e.checkDeleteKey(k) {
			continue
		}
		if handler(k, v) {
			break
		}
	}
}

func (e *ExpiredMap[K, V]) checkDeleteKey(key K) bool {
	if val, found := e.m[key]; found {
		if val.expiredVime <= time.Now().Unix() {
			delete(e.m, key)
			//delete(e.timeMap, val.expiredVime)
			return false
		}
		return true
	}
	return false
}

// get all data
func (e *ExpiredMap[K, V]) GetAll() []V {
	e.lck.Lock()
	defer e.lck.Unlock()
	m := []V{}
	for k, v := range e.m {
		if !e.checkDeleteKey(k) {
			continue
		}
		if v == nil {
			continue
		}
		m = append(m, v.data)
	}
	return m
}

func (e *ExpiredMap[K, V]) GetMapAll() map[K]V {
	e.lck.Lock()
	defer e.lck.Unlock()
	result := make(map[K]V)
	for k, v := range e.m {
		if !e.checkDeleteKey(k) {
			continue
		}
		if v == nil {
			continue
		}
		result[k] = v.data
	}
	return result
}

func (e *ExpiredMap[K, V]) GetAllKeys() []K {
	e.lck.Lock()
	defer e.lck.Unlock()
	m := []K{}
	for k, v := range e.m {
		if !e.checkDeleteKey(k) {
			continue
		}
		if v == nil {
			continue
		}
		m = append(m, k)
	}
	return m
}

func (e *ExpiredMap[K, V]) GetAndDelete(key K) (found bool, value V) {
	e.lck.Lock()
	defer e.lck.Unlock()

	if found = e.checkDeleteKey(key); !found {
		return
	}
	value = e.m[key].data
	delete(e.m, key)
	return
}
