package uuid

import (
	"fmt"
	"sync"
	"time"
)

const (
	twepoch      = int64(1417937700000) // 默认起始的时间戳 1449473700000 。计算时，减去这个值
	nodeIdBits   = uint(10)             //节点 所占位置
	sequenceBits = uint(14)             //自增ID 所占用位置

	/*
	 * 1 符号位  |  39 时间戳                                    | 10 区域节点     | 14 （毫秒内）自增ID
	 * 0        |  0000000 00000000 00000000 00000000 00000000 | 00000000 00   |  00000000 000000
	 *
	 */
	maxNodeId = -1 ^ (-1 << nodeIdBits) //节点 ID 最大范围

	nodeIdShift        = sequenceBits //左移次数
	timestampLeftShift = sequenceBits + nodeIdBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
	nodeMask           = -1 ^ (-1 << nodeIdBits)
	maxNextIdsNum      = 100 //单次获取ID的最大数量
)

var uuidWorker *idWorker

func getUUID() *idWorker {
	if uuidWorker == nil {
		uuidWorker = newIdWorker()
	}
	return uuidWorker
}

type idWorker struct {
	sequence      int64 //序号
	lastTimestamp int64 //最后时间戳
	twepoch       int64
	mutex         sync.Mutex
}

// newIdWorker new a snowflake id generator object.
func newIdWorker() *idWorker {
	idWorker := &idWorker{}

	idWorker.lastTimestamp = -1
	idWorker.sequence = 0
	idWorker.twepoch = twepoch
	idWorker.mutex = sync.Mutex{}
	//fmt.Sprintf("worker starting. timestamp left shift %d,worker id bits %d, sequence bits %d, workerid %d", timestampLeftShift, nodeIdBits, sequenceBits, NodeId)
	return idWorker
}
func (id *idWorker) nextid() (int64, error) {
	timestamp := timeGen()
	if timestamp < id.lastTimestamp {
		return 0, fmt.Errorf("Clock moved backwards.  Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp)
	}
	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence + 1) & sequenceMask
		if id.sequence == 0 {
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else {
		id.sequence = 0
	}
	sarly := time.Now().Unix() % nodeMask
	id.lastTimestamp = timestamp
	return ((timestamp - id.twepoch) << timestampLeftShift) | (sarly << nodeIdShift) | id.sequence, nil
}

// timeGen generate a unix millisecond.
func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// tilNextMillis spin wait till next millisecond.
func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

// NextId get a snowflake id.
func NextId() (int64, error) {
	getUUID().mutex.Lock()
	defer getUUID().mutex.Unlock()
	return getUUID().nextid()
}

// NextIds get snowflake ids.
func NextIds(num int) ([]int64, error) {
	if num > maxNextIdsNum || num < 0 {
		return nil, fmt.Errorf("NextIds num: %d error", num)
	}
	ids := make([]int64, num)
	getUUID().mutex.Lock()
	defer getUUID().mutex.Unlock()
	for i := 0; i < num; i++ {
		ids[i], _ = getUUID().nextid()
	}
	return ids, nil
}
