package wuid

import (
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	epoch         int64  = 1672531200000
	timestampBits uint8  = 41
	machineIDBits uint8  = 10
	sequenceBits  uint8  = 12
	maxSequence   uint64 = (1 << sequenceBits) - 1
	shift         uint8  = machineIDBits + sequenceBits
)

var (
	machineID     uint64 = 1
	uniqueCounter uint64
	lastTimestamp int64
	sequence      uint64
)

func init() {
	now := time.Now().UnixMilli()
	uniqueCounter = ((uint64(now - epoch)) << shift) | (machineID << sequenceBits)
}

func tilNextMillis(t int64) int64 {
	ts := time.Now().UnixMilli()
	for ts <= t {
		ts = time.Now().UnixMilli()
	}
	return ts
}

type ID uint64

func New() ID {
	return ID(GenerateFastID())
}

func (id ID) String() string {
	return u64ToString(uint64(id))
}

func (id ID) Int64() int64 {
	return int64(id)
}

func (id ID) Uint64() uint64 {
	return uint64(id)
}

func GenerateFastID() uint64 {
	newID := atomic.AddUint64(&uniqueCounter, 1)
	return newID
}

func GenerateWithTimestamp() uint64 {
	current := uint64(time.Now().UnixMilli() - epoch)
	current &= (1 << timestampBits) - 1
	for {
		old := atomic.LoadUint64(&uniqueCounter)
		oldTimestamp := old >> shift
		var candidate uint64
		if current > oldTimestamp {
			candidate = (current << shift) | (machineID << sequenceBits)
		} else {
			candidate = old + 1
		}
		if atomic.CompareAndSwapUint64(&uniqueCounter, old, candidate) {
			return candidate
		}
	}
}

func u64ToString(n uint64) string {
	var buf [20]byte
	i := len(buf)
	for {
		i--
		buf[i] = '0' + byte(n%10)
		n /= 10
		if n == 0 {
			break
		}
	}
	return *(*string)(unsafe.Pointer(&struct {
		Data *byte
		Len  int
	}{&buf[i], len(buf) - i}))
}

func GenerateIDInt64() int64 {
	return int64(GenerateFastID())
}

func GenerateIDString() string {
	return u64ToString(GenerateFastID())
}

func main() {
	for i := 0; i < 10000; i++ {
		fmt.Printf("FastID %d: %s\n", i, GenerateIDString())
	}
}
