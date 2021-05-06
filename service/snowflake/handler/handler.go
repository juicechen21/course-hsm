package handler

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits uint8 = 10 // 10位的机器标识，10位的长度最多支持部署1024个节点 2的10次方
	numberBits uint8 = 12 // 12位的计数序列号，每个节点每毫秒产生4096个ID序号 2的12次方
	workerMax int64 = -1 ^ (-1 << workerBits) // 1023 2的10次方减一
	numberMax int64 = -1 ^ (-1 << numberBits) // 4095 2的12次方减一
	timeShift uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	startTime int64 = 1100000000000 // 如果在程序跑了一段时间修改epoch这个值 可能会导致生成相同的ID
)

type Worker struct {
	mu sync.Mutex //互斥锁是一种常用的控制共享资源访问的方法，它能够保证同时只有一个goroutine可以访问共享资源
	timestamp int64
	workerId int64
	number int64
}

func NewWorker(workerId int64)(*Worker,error)  {
	if workerId < 0 || workerId > workerMax {
		return nil,errors.New("机器ID超出设定的数量！")
	}
	// 生成一个新的节点
	return &Worker{
		timestamp: 0,
		workerId: workerId,
		number: 0,
	},nil
}

func (w *Worker)GetId() int64 {
	w.mu.Lock() // 加锁
	defer w.mu.Unlock() //程序运行完成后解锁
	now := time.Now().UnixNano() / 1e6 // 获取当前13位时间戳
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	}else{
		w.number = 0
		w.timestamp = now
	}
	id := int64((now-startTime) << timeShift | (w.workerId << workerShift) | w.number)
	return id
}


