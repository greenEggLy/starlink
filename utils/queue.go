package utils

import (
	"sync"
	"time"

	pb "starlink/pb"
)

type Queue struct {
	Content []*pb.LLAPosition
	Timeout int // timeout为0为无限延时， 小于0为不延时， 大于0为延时timeout秒
	MaxSize int // 队列容量, 小于或等于0为不限量，不限量时延时无效，大于0且到达上限时则开始延时
}

var lock = sync.Mutex{}

// 超过设定延时时间后， 元素会被抛弃
func (q *Queue) Put(msg *pb.LLAPosition) {
	lock.Lock()
	closeSingle := make(chan bool)
	succesSingle := make(chan bool)
	go func(close chan bool, success chan bool) {
		var t1 *time.Timer
		t1 = time.NewTimer(time.Second * time.Duration(q.Timeout))
		for {
			select {
			case <-t1.C:
				if q.Timeout == 0 {
					t1 = time.NewTimer(time.Second * time.Duration(q.Timeout))
					continue
				} else {
					success <- true
					return
				}
			default:
				if q.MaxSize != 0 && q.MaxSize == len(q.Content) {
					q.Content = q.Content[1:]
					q.Content = append(q.Content, msg)
				} else {
					q.Content = append(q.Content, msg)
					success <- true
					return
				}
			}
		}
	}(closeSingle, succesSingle)

	for {
		<-succesSingle
		lock.Unlock()
		return
	}
}

// 超过延时时间时会返回空字符串
func (q *Queue) Get() *pb.LLAPosition {
	lock.Lock()
	closeSingle := make(chan bool)
	succesSingle := make(chan *pb.LLAPosition)
	go func(close chan bool, output chan *pb.LLAPosition) {
		var t1 *time.Timer
		t1 = time.NewTimer(time.Second * time.Duration(q.Timeout))
		for {
			select {
			case <-t1.C:
				if q.Timeout == 0 {
					t1 = time.NewTimer(time.Second * time.Duration(q.Timeout))
					continue
				} else {
					output <- nil
					return
				}
			default:
				if len(q.Content) > 0 {
					msg := q.Content[0]
					q.Content = q.Content[1:]
					output <- msg
					return
				}
			}
		}
	}(closeSingle, succesSingle)

	for {
		res := <-succesSingle
		lock.Unlock()
		return res
	}
}

func (q *Queue) Size() int {
	return len(q.Content)
}

func NewQueue(timeout int, maxsize int) *Queue {
	return &Queue{
		Content: []*pb.LLAPosition{},
		Timeout: timeout,
		MaxSize: maxsize,
	}
}

// func main() {
// 	q := Queue{
// 		content: []string{},
// 		Timeout: 3,
// 		MaxSize: 4,
// 	}
// 	go func() {
// 		for i := 0; i < 50; i++ {
// 			st := "hello" + strconv.Itoa(i)
// 			q.put(st)
// 			fmt.Println(st)
// 		}
// 	}()
// 	for i := 0; i < 50; i++ {
// 		fmt.Println("read: " + q.get())
// 	}
// }
