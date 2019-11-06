package DataStructures

import (
	"container/list"
)

type Queue struct {
	In  chan<- interface{}
	Out <-chan interface{}
}

func NewQueue(maxBufferSize int) *Queue  {
	in := make(chan interface{}, maxBufferSize)
	out := make(chan interface{})
	go func() {
		queue := list.New()
		frontValue := func() interface{} {
			if queue.Len() == 0 {
				return nil
			}
			return queue.Front().Value
		}
		outChannel := func() chan interface{} {
			if queue.Len() == 0 {
				return nil
			}
			return out
		}
		for queue.Len() > 0 || in != nil {
			select {
				case v, ok := <-in:
					if ok {
						queue.PushBack(v)
					} else {
						in = nil
					}
				case outChannel() <- frontValue():
					queue.Remove(queue.Front())
			}
		}
		close(out)
	}()
	return &Queue{
		In:  in,
		Out: out,
	}
}