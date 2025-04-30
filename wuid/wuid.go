package wuid

import (
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
	timestampBase uint64
	seqCounter    uint64
)

func init() {
	now := time.Now().UnixMilli()
	uniqueCounter = ((uint64(now - epoch)) << shift) | (machineID << sequenceBits)
	base := (uint64(now-epoch) << shift) | (machineID << sequenceBits)
	atomic.StoreUint64(&timestampBase, base)
	atomic.StoreUint64(&seqCounter, 0)
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
	return atomic.AddUint64(&uniqueCounter, 1)
}

func GenerateWithTimestamp() uint64 {
	now := uint64(time.Now().UnixMilli() - epoch)
	base := atomic.LoadUint64(&timestampBase)
	if now > (base >> shift) {
		newBase := (now << shift) | (machineID << sequenceBits)
		atomic.StoreUint64(&timestampBase, newBase)
		base = newBase
	}
	seq := atomic.AddUint64(&seqCounter, 1) & maxSequence
	return (base & ^maxSequence) | seq
}

//go:inline
//go:nosplit
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
	sh := struct {
		Data *byte
		Len  int
	}{&buf[i], len(buf) - i}
	return *(*string)(unsafe.Pointer(&sh))
}

func GenerateIDInt64() int64 {
	return int64(GenerateFastID())
}

func GenerateIDString() string {
	return u64ToString(GenerateFastID())
}
