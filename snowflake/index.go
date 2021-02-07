package snowflake

import (
	"errors"
	"sync"
	"time"

	"github.com/xm-chentl/go-uuid"
)

const (
	workerIDBits     = uint64(5) // 10bit 工作机器ID中的 5bit workerID
	dataCenterIDBits = uint64(5) // 10bit 工作机器ID中的 5bit dataCenterID
	sequenceBits     = uint64(12)

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits) // 节点ID的最大值, 用于防止溢出
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft = uint8(22) // timeLeft = workerIDBits + sequenceBits // 时间戳向左偏移量
	dataLeft = uint8(17) // dataLeft = dataCenterIDBits + sequenceBits
	workLeft = uint8(12) // workLeft = sequenceBits // 节点IDx向左偏移量

	// 2020-05-20 08:00:00 +0800 CST
	twepoch = int64(1589923200000) // 常量时间戳(毫秒)
)

type workerImpl struct {
	mu           sync.Mutex // 加互斥锁，确保并发安全性
	dataCenterID int64      // 该节点的 数据中心ID
	workerID     int64      // 该节点的ID
	lastStamp    int64      // 记录上一次ID的时间戳
	sequence     int64      // 当前毫秒已经生成的ID序列号(从0 开始累加) 1毫秒内最多生成4096个ID
	twepoch      int64
}

func (w *workerImpl) Generate() (int64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	timeStamp := w.getCurrentMilliSeconds()
	if timeStamp < w.lastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}
	if w.lastStamp == timeStamp {
		w.sequence = (w.sequence + 1) & maxSequence
		if w.sequence == 0 {
			for timeStamp <= w.lastStamp {
				timeStamp = w.getCurrentMilliSeconds()
			}
		}
	} else {
		w.sequence = 0
	}
	w.lastStamp = timeStamp
	id := ((timeStamp - w.twepoch) << timeLeft) | (w.dataCenterID << dataLeft) | (w.workerID << workLeft) | w.sequence

	return int64(id), nil
}

func (w *workerImpl) getCurrentMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// New 分布式情况下,我们应通过外部配置文件或其他方式为每台机器分配独立的id
func New(dataCenterID, workerID int64) uuid.IUUID {
	return &workerImpl{
		dataCenterID: dataCenterID,
		workerID:     workerID,
		lastStamp:    0,
		sequence:     0,
		twepoch:      twepoch,
	}
}

// New 实例一个雪花算法实例
func NewCfg(dataCenterID, workerID int64, cfgTime time.Time) uuid.IUUID {
	return &workerImpl{
		dataCenterID: dataCenterID,
		workerID:     workerID,
		lastStamp:    0,
		sequence:     0,
		twepoch:      cfgTime.UnixNano(),
	}
}
