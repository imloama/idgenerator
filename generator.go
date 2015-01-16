//主键生成器
//代码内容主要来源于tweeter的snowflake，实现对应的go语言版本
package idgenerator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	minusOne           int64  = -1
	WorkerIdBits       uint64 = 5
	DataCenterIdBits   uint64 = 5
	MaxWorkerId        int64  = minusOne ^ (minusOne << WorkerIdBits)
	MaxDataCenterId    int64  = minusOne ^ (minusOne << DataCenterIdBits)
	SequenceBits       uint64 = 12
	WorkerIdShift      uint64 = SequenceBits
	DataCenterIdShift  uint64 = SequenceBits + WorkerIdBits
	TimestampLeftShift uint64 = SequenceBits + WorkerIdBits + DataCenterIdBits
	SequenceMask       int64  = minusOne ^ (minusOne << SequenceBits)
	twepoch            int64  = 1288834974657
)

var lock sync.Mutex

type IdGenerator interface {
	NextId() int64
}

type IdWorker struct {
	workerId      int64
	dataCenterId  int64
	sequence      int64
	lastTimestamp int64
}

//生成ID生成器实例
func NewIdWorker(workerId, dataCenterId int64) (*IdWorker, error) {
	if workerId > MaxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("worker Id can't be greater than %d or less than 0", workerId))
	}

	if dataCenterId > MaxDataCenterId || dataCenterId < 0 {
		return nil, errors.New(fmt.Sprintf("datacenter Id can't be greater than %d or less than 0", dataCenterId))
	}

	return &IdWorker{
		workerId:      workerId,
		dataCenterId:  dataCenterId,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

//服务器时间不正确，可能引发panic
func (self *IdWorker) NextId() int64 {
	lock.Lock()
	defer lock.Unlock()
	timestamp := timeGen()
	// fmt.Printf("timestamp:%d\n", timestamp)
	if timestamp < self.lastTimestamp {
		panic(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", self.lastTimestamp-timestamp))
	}
	if timestamp == self.lastTimestamp {
		self.sequence = (self.sequence + 1) & SequenceMask
		if self.sequence == 0 {
			timestamp = tilNextMillis(self.lastTimestamp)
		}
	} else {
		self.sequence = 0
	}
	self.lastTimestamp = timestamp
	id := ((timestamp - twepoch) << TimestampLeftShift) | (self.dataCenterId << DataCenterIdShift) | (self.workerId << WorkerIdShift) | self.sequence
	// fmt.Printf("id:%d,timestamp:%d,timestamp - twepoch:%d,dataCenterId:%d,DataCenterIdShift:%d,workerId:%d,WorkerIdShift:%d,sequence:%d\n", id, timestamp, (timestamp - twepoch), self.dataCenterId, DataCenterIdShift, self.workerId, WorkerIdShift, self.sequence)
	return id
}

//16进制的id,大写
func (self *IdWorker) HexId() string {
	id := self.NextId()
	str := strconv.FormatInt(id, 16)
	str = strings.ToUpper(str)
	return str
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp < lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

//返回ms 毫秒
func timeGen() int64 {
	return time.Now().UnixNano() / 1e6
}
