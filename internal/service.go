package internal

import (
	"context"
	"log"
	"sync"

	"github.com/galihsatriawan/eod/utils"
)

type Service interface {
	EndOfDay(context.Context) error
}
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
func (s *service) EndOfDay(ctx context.Context) error {
	beforeEodList, err := s.repo.GetBeforeEod(ctx)
	if err != nil {
		log.Println("[Service][EndOfDay] failed get before eod", err)
		return err
	}
	var afterEodList AfterEodList
	// copy data to after eod
	for _, beforeEod := range beforeEodList {
		afterEod := NewAfterEod(*beforeEod)
		afterEodList = append(afterEodList, afterEod)
	}

	var wg sync.WaitGroup
	mutex := sync.Mutex{}
	workerPoolProcessThree := utils.NewWorkerPool(ctx, 8)
	workerPoolProcessThree.Run()
	for index := 0; index < 100; index++ {
		wg.Add(1)
		tempIndex := index
		go workerPoolProcessThree.AddTask(func(workerId int32) {
			defer wg.Done()
			mutex.Lock()
			afterEodList[tempIndex].Balanced += 10
			afterEodList[tempIndex].ProcessThree = func() *int64 {
				workerId64 := int64(workerId)
				return &workerId64
			}()
			mutex.Unlock()
		})
	}
	wg.Wait()
	workerPoolProcessThree.Stop()

	workerPoolProcessTwoA := utils.NewWorkerPool(ctx, 10)
	workerPoolProcessTwoA.Run()

	workerPoolProcessTwoB := utils.NewWorkerPool(ctx, 10)
	workerPoolProcessTwoB.Run()

	for _, afterEod := range afterEodList {
		wg.Add(1)
		tempAfterEod := afterEod
		if afterEod.Balanced > 150 {
			go workerPoolProcessTwoB.AddTask(func(workerId int32) {
				defer wg.Done()
				mutex.Lock()
				tempAfterEod.Balanced += 25
				tempAfterEod.ProcessTwoB = func() *int64 {
					workerId64 := int64(workerId)
					return &workerId64
				}()
				mutex.Unlock()
			})
		}

		if afterEod.Balanced <= 150 {
			go workerPoolProcessTwoA.AddTask(func(workerId int32) {
				defer wg.Done()
				mutex.Lock()
				tempAfterEod.FreeTransfer = 5
				tempAfterEod.ProcessTwoA = func() *int64 {
					workerId64 := int64(workerId)
					return &workerId64
				}()
				mutex.Unlock()
			})
		}

	}
	wg.Wait()
	workerPoolProcessTwoA.Stop()
	workerPoolProcessTwoB.Stop()

	workerPoolProcessOne := utils.NewWorkerPool(ctx, 10)
	workerPoolProcessOne.Run()
	for _, afterEod := range afterEodList {
		wg.Add(1)
		tempAfterEod := afterEod
		go workerPoolProcessOne.AddTask(func(workerId int32) {
			defer wg.Done()
			mutex.Lock()
			tempAfterEod.AverageBalanced = (tempAfterEod.Balanced + tempAfterEod.PreviousBalanced) / 2
			tempAfterEod.ProcessOne = int64(workerId)
			mutex.Unlock()
		})
	}
	wg.Wait()
	workerPoolProcessOne.Stop()
	err = s.repo.UpdateAfterEod(ctx, afterEodList)
	if err != nil {
		log.Println("[Service][EndOfDay] update after eod failed", err)
	}
	return err
}
