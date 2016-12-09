package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type threadPool struct {
	synG    sync.WaitGroup
	syn     sync.Mutex
	workNum int
	job     chan func() error
	done    int
}

func (t *threadPool) addJob(f func() error) {
	t.job <- f
}

func (t *threadPool) run() {
	for i := 0; i < t.workNum; i++ {
		t.synG.Add(1)
		go func(i int) {
			log.Printf("%d goroutine start....\n", i)
			for {
				work, ok := <-t.job
				if !ok {
					log.Printf("%d goroutine stop\n", i)
					break
				}
				err := work()
				if err != nil {
					log.Printf("%d goroutine occour error: %s\n", i, err.Error())
				}
				t.syn.Lock()
				t.done++
				t.syn.Unlock()
				log.Printf("%d goroutine have done one job....\n", i)
			}
			t.synG.Done()
		}(i)
	}
}

func (t *threadPool) stop() {
	close(t.job)
}

func job() error {
	syn.Lock()
	fmt.Println(count)
	count++
	syn.Unlock()
	return nil
}

var count int
var syn sync.Mutex

func main() {
	tp := threadPool{
		workNum: 3,
		job:     make(chan func() error),
	}
	tp.run()
	//add job
	for i := 0; i < 100; i++ {
		tp.addJob(job)
	}

	//check jobs num
	for tp.done != count {
		time.Sleep(time.Second * 2)
	}

	tp.stop()
	tp.synG.Wait()

}
