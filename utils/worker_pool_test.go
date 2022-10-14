package utils_test

import (
	"context"
	"sync"
	"testing"
	"time"

	. "github.com/galihsatriawan/eod/utils"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	testCases := map[string]struct {
		args         []int32
		expectResult []int32
		expectWorker int32
	}{
		"increment an array with one": {
			args: []int32{
				1, 2, 3,
			},
			expectResult: []int32{
				2, 3, 4,
			},
			expectWorker: 3,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			workerPool := NewWorkerPool(context.Background(), 3)
			workerPool.Run()
			mapWorkerId := map[int32]bool{}
			actualResult := tc.args
			var mutex sync.Mutex
			var wg sync.WaitGroup
			for i := 0; i < len(tc.args); i++ {
				wg.Add(1)
				tempIndex := i
				workerPool.AddTask(func(workerId int32) {
					defer wg.Done()
					mutex.Lock()
					mapWorkerId[workerId] = true
					actualResult[tempIndex]++
					mutex.Unlock()
					time.Sleep(1 * time.Second)
				})
			}
			wg.Wait()
			workerPool.Stop()
			assert.Equal(t, tc.expectResult, actualResult)
			assert.Equal(t, tc.expectWorker, int32(len(mapWorkerId)))
		})
	}
}
