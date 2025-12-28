package idx

import (
	"sync"
	"time"

	"github.com/go-xuan/utilx/stringx"
)

/*
首位：第一个bit作为符号位，正数为0。
时间戳：占用41bit，精确到毫秒。41位最好可以表示2^41-1毫秒，转化成单位年为69年。
机器号：占用10bit，最多可以容纳1024个节点。
序列号：占用12bit，每个节点每毫秒从0开始不断累加，最多可以累加到4095，一共可以产生4096个ID。
*/
const (
	workerBits     = uint(10)                   // 机器号位数，即2^10=1024个节点
	sequenceBits   = uint(12)                   // 序列号位数，即2^12=4096个ID
	workerShift    = sequenceBits               // 机器码偏移量
	timestampShift = workerBits + sequenceBits  // 时间戳偏移量
	workerMax      = int64(1<<workerBits - 1)   // 机器号最大值（即1023）
	sequenceMax    = int64(1<<sequenceBits - 1) // 序列号最大值（即4095）
	epoch          = int64(946656000000)        // 起始常量时间戳（毫秒）,此处选取的时间是2000-01-01 00:00:00
)

var flake *Flake

// SnowFlake 雪花ID生成器
func SnowFlake(id ...int64) *Flake {
	workerId := getWorkerId(id...)
	if flake == nil || flake.WorkerId != workerId {
		flake = newSnowflake(workerId)
	}
	return flake
}

// newSnowflake 创建新的雪花ID生成器
func newSnowflake(workerId int64) *Flake {
	return &Flake{
		Mutex:     new(sync.Mutex),
		WorkerId:  workerId,
		TimeStamp: 0,
		Sequence:  0,
	}
}

// 获取机器号
func getWorkerId(worker ...int64) int64 {
	if len(worker) > 0 && worker[0] != 0 {
		id := worker[0] % workerMax
		if id < 0 {
			return -id
		}
		return id
	}
	return 1
}

type Flake struct {
	*sync.Mutex       // 互斥锁
	WorkerId    int64 // 机器号,0~1023
	TimeStamp   int64 // 当前时间戳（毫秒）
	Sequence    int64 // 当前毫秒内的序列号
}

func (s *Flake) Value() int64 {
	s.Lock()
	defer s.Unlock()
	milli := time.Now().UnixMilli()
	if s.TimeStamp == milli {
		s.Sequence++
		if s.Sequence > sequenceMax {
			// 序列号超过最大值，等待到下一个毫秒
			for milli <= s.TimeStamp {
				milli = time.Now().UnixMilli()
			}
		}
	} else {
		s.Sequence = 0
		s.TimeStamp = milli
	}
	return (milli-epoch)<<timestampShift | (s.WorkerId << workerShift) | s.Sequence
}

func (s *Flake) String() string {
	return stringx.Int64(s.Value())
}
