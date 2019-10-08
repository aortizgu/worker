package worker

import (
	"log"
	"sync"
	"testing"
)

func Test_worker(t *testing.T) {
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	w := NewWorker()
	w.Start(10, func() bool {
		log.Println("Hola")
		waitgroup.Done()
		return true
	})
	waitgroup.Wait()
	w.Stop()

	waitgroup.Add(2)
	w.Start(1000, func() bool {
		log.Println("Hola2")
		waitgroup.Done()
		return true
	})
	waitgroup.Wait()
	w.Stop()
}
