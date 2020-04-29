package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rfyiamcool/gpool"
)

var (
	wg = sync.WaitGroup{}
)

func main() {
	gp, err := gpool.NewGPool(&gpool.Options{
		MaxWorker:   20,               // 最大的协程数
		MinWorker:   8,                // 最小的协程数
		JobBuffer:   6,                // 缓冲队列的大小
		IdleTimeout: 30 * time.Second, // 协程的空闲超时退出时间
	})

	if err != nil {
		panic(err.Error())
	}

	for index := 0; index < 200; index++ {
		wg.Add(1)
		idx := index
		gp.ProcessAsync(func() {
			resp, _ := http.Get("http://fueltank-1:82/ping")
			var result []byte
			resp.Body.Read(result)
			log.Println(idx, ".  result: ", string(result))
		})
	}

	wg.Done()
}
