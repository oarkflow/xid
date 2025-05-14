package wuid

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	epoch         int64  = 1672531200000
	machineIDBits uint8  = 10
	sequenceBits  uint8  = 12
	maxSequence   uint64 = (1 << sequenceBits) - 1
	shift                = machineIDBits + sequenceBits
)

var (
	machineID     uint64
	uniqueCounter uint64
	timestampBase uint64
	seqCounter    uint64
)

func init() {
	machineID = generateRandomMachineID()
	now := time.Now().UnixMilli()
	uniqueCounter = (uint64(now-epoch) << shift) | (machineID << sequenceBits)
	base := uniqueCounter
	atomic.StoreUint64(&timestampBase, base)
	atomic.StoreUint64(&seqCounter, 0)
}

func generateRandomMachineID() uint64 {
	var b [2]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Fatalf("failed to seed machineID: %v", err)
	}
	return uint64(binary.BigEndian.Uint16(b[:]) & ((1 << machineIDBits) - 1))
}

type ID uint64

func New(timestamp ...bool) ID {
	if len(timestamp) > 0 && timestamp[0] {
		return ID(GenerateWithTimestamp())
	}
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
