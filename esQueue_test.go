package esQueueTest

import (
	"fmt"
	"github.com/yireyun/go-queue"
	"sync"
	"testing"
)
var goRoutineCnt = 1000
func checkPutAndGet(checkData []int){
	l := len(checkData)
	for i := 0; i < l; i++{
		if checkData[i] != i{
			fmt.Printf("\ncheckPutAndGet error!! lost data!!! i:%v checkData[i]:%v l:%v\n", i, checkData[i], l)
		}
	}
}
func BenchmarkEsQueueReadContention(b *testing.B) {
	var checkData = make([]int, b.N)
	q := queue.NewQueue(1024*1024)
	var wgGet sync.WaitGroup
	wgGet.Add(goRoutineCnt)
	var wgPut sync.WaitGroup
	wgPut.Add(1)
	b.ResetTimer()

	go func() {
		for i := 0; i < b.N; i++ {
			ok, _ := q.Put(i)
			for !ok {
				ok, _= q.Put(i)
			}
		}
		wgPut.Done()
	}()

	for i := 0; i < goRoutineCnt; i++ {
		go func() {
			for i := 0; i < b.N / goRoutineCnt; i++ {
				val, ok, _ := q.Get()
				for !ok {
					val, ok, _= q.Get()
				}
				v := val.(int)
				checkData[v] = v
			}
			wgGet.Done()
		}()
	}
	wgGet.Wait()
	wgPut.Wait()
	for q.Quantity() > 0{
		val, ok, _ := q.Get()
		for !ok {
			val, ok, _ = q.Get()
		}
		v := val.(int)
		checkData[v] = v
	}
	checkPutAndGet(checkData)
}
