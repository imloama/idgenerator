package idgenerator

import (
	"fmt"
	"runtime"
	"testing"
)

func Test_NextId(t *testing.T) {
	runtime.GOMAXPROCS(2)
	idWorker, err := NewIdWorker(0, 0)
	if err != nil {
		t.Error("newidworker failed")
	}

	wchan := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			id1 := idWorker.NextId()
			t.Log(fmt.Sprintf("id1:%d", id1))
		}
		wchan <- true
	}()
	go func() {
		for i := 0; i < 100; i++ {
			id2 := idWorker.NextId()
			t.Log(fmt.Sprintf("id2:%d", id2))
		}
		wchan <- true
	}()
	c := 0
	for {
		b := <-wchan
		fmt.Printf("wchan:%v\n", b)
		if b {
			c++
		}
		if c >= 2 {
			break
		}

	}

}

func Test_HexId(t *testing.T) {
	runtime.GOMAXPROCS(2)
	idWorker, err := NewIdWorker(0, 0)
	if err != nil {
		t.Error("newidworker failed")
	}
	for i := 0; i < 10; i++ {
		str := idWorker.HexId()
		t.Log(str)
	}
}

func Test_x(t *testing.T) {
	/*
		    id:-8542765543168409600,timestamp - twepoch:132493773170,dataCenterId:0,DataCent
		erIdShift:24,workerId:0,WorkerIdShift:12,sequence:0
		id:4922997342669373440,timestamp - twepoch:132493886546,dataCenterId:0,DataCente
		rIdShift:24,workerId:0,WorkerIdShift:12,sequence:0
		((timestamp - twepoch) << TimestampLeftShift) | (self.dataCenterId << DataCenterIdShift) | (self.workerId << WorkerIdShift) | self.sequence
	*/

	var t1 int64 = 132493773170
	var i1 int64 = t1 << 48
	var i2 int64 = 1 << 24
	var i3 int64 = 0 << 12
	var t2 int64 = 0
	i := i1 | i2 | i3 | t2
	fmt.Printf("result:%d,i1:%d,i2:%d,i3:%d", i, i1, i2, i3)

}
