package DataStructures

import (
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	maxSize := 100
	queue := NewQueue(maxSize)
	var result []int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for v := range queue.Out {
			result = append(result, v.(int))
		}
		wg.Done()
	}()

	var expected []int
	for i := 0 ; i < maxSize ; i ++ {
		queue.In <- i
		expected = append(expected, i)
	}
	close(queue.In)
	wg.Wait()

	for i := 0 ; i < maxSize ; i ++ {
		if expected[i] != result[i] {
			t.Errorf("Expected: %v\n Found %v\n", expected, result)
		}
	}

}