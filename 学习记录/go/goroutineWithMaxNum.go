package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 模拟一个长任务
// 为了观察完成时间，执行时间设置为随机
func doWork() {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

func main() {
	maxNbConcurrentGoroutines := 5
	nbJobs := 20

	// 控制任务等待的信号量
	concurrentGoroutines := make(chan bool, maxNbConcurrentGoroutines)
	// 先填满信号，之后取到一个执行一个go routine，取不到的就等待，从而控制最大数量
	for i := 0; i < maxNbConcurrentGoroutines; i++ {
		concurrentGoroutines <- true
	}

	// 控制 main 退出
	var wg sync.WaitGroup

	// 执行每个任务
	for i := 1; i <= nbJobs; i++ {
		fmt.Printf("ID: %v: waiting to launch!\n", i)
		wg.Add(1)
		// 尝试获取等待信号量，拿得到就跑，拿不到就等
		<-concurrentGoroutines
		fmt.Printf("ID: %v: it's my turn!\n", i)
		go func(id int) {
			defer func() {
				fmt.Printf("ID: %v: all done!\n", id)
				concurrentGoroutines <- true
				wg.Done()
			}()
			doWork()
		}(i)
	}

	// 等待所有程序完成，主进程退出
	wg.Wait()
}
