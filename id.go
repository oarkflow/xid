package xid

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	Epoch     int64 = 1288834974657
	NodeBits  uint8 = 10
	StepBits  uint8 = 12
	mu        sync.Mutex
	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
)

const (
	maxJSONSize     = 22
	encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"
	encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

var (
	decodeBase32Map  [256]byte
	decodeBase58Map  [256]byte
	ErrInvalidBase58 = errors.New("invalid base58")
	ErrInvalidBase32 = errors.New("invalid base32")
	node             *Node
)

type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid snowflake ID %q", string(j.original))
}

func init() {
	for i := range decodeBase58Map {
		decodeBase58Map[i] = 0xFF
	}
	for i := range encodeBase58Map {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}
	for i := range decodeBase32Map {
		decodeBase32Map[i] = 0xFF
	}
	for i := range encodeBase32Map {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
	node, _ = NewNode(1)
}

type Node struct {
	epoch     time.Time
	time      int64
	node      int64
	step      int64
	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

type ID int64

func NewNode(nodeID int64) (*Node, error) {
	mu.Lock()
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits
	mu.Unlock()
	n := &Node{
		node:      nodeID,
		nodeMax:   -1 ^ (-1 << NodeBits),
		nodeMask:  (-1 ^ (-1 << NodeBits)) << StepBits,
		stepMask:  -1 ^ (-1 << StepBits),
		timeShift: NodeBits + StepBits,
		nodeShift: StepBits,
	}
	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	cur := time.Now()
	n.epoch = cur.Add(time.Unix(Epoch/1000, (Epoch%1000)*1e6).Sub(cur))
	return n, nil
}

func New() ID {
	return node.New()
}

func (n *Node) New() ID {
	for {
		now := time.Since(n.epoch).Milliseconds()
		old := atomic.LoadInt64(&n.time)
		if now > old {
			if atomic.CompareAndSwapInt64(&n.time, old, now) {
				atomic.StoreInt64(&n.step, 0)
				return ID((now << n.timeShift) | (n.node << n.nodeShift))
			}
			continue
		} else if now == old {
			step := atomic.AddInt64(&n.step, 1) - 1
			if step <= n.stepMask {
				return ID((now << n.timeShift) | (n.node << n.nodeShift) | step)
			}
			continue
		} else {
			continue
		}
	}
}

func (f ID) Int64() int64 {
	return int64(f)
}

func (f ID) String() string {
	var buf [20]byte
	i := len(buf)
	v := int64(f)
	for v >= 10 {
		i--
		q := v / 10
		buf[i] = byte('0' + v - q*10)
		v = q
	}
	i--
	buf[i] = byte('0' + v)
	return string(buf[i:])
}

func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

func (f ID) Base32() string {
	if f < 32 {
		return string(encodeBase32Map[f])
	}
	var buf [12]byte
	i := len(buf)
	v := f
	for v >= 32 {
		i--
		buf[i] = encodeBase32Map[v%32]
		v /= 32
	}
	i--
	buf[i] = encodeBase32Map[v]
	return string(buf[i:])
}

func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

func (f ID) Base58() string {
	if f < 58 {
		return string(encodeBase58Map[f])
	}
	var buf [11]byte
	i := len(buf)
	v := f
	for v >= 58 {
		i--
		buf[i] = encodeBase58Map[v%58]
		v /= 58
	}
	i--
	buf[i] = encodeBase58Map[v]
	return string(buf[i:])
}

func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

func (f ID) Bytes() []byte {
	var buf [20]byte
	i := len(buf)
	v := int64(f)
	for v >= 10 {
		i--
		q := v / 10
		buf[i] = byte('0' + v - q*10)
		v = q
	}
	i--
	buf[i] = byte('0' + v)
	return buf[i:]
}

func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + Epoch
}

func (f ID) Node() int64 {
	return (int64(f) & nodeMask) >> nodeShift
}

func (f ID) Step() int64 {
	return int64(f) & stepMask
}

func (f ID) MarshalJSON() ([]byte, error) {
	var buf [maxJSONSize]byte
	buf[0] = '"'
	n := len(buf) - 2
	v := int64(f)
	for v >= 10 {
		q := v / 10
		buf[n] = byte('0' + v - q*10)
		n--
		v = q
	}
	buf[n] = byte('0' + v)
	copy(buf[1:], buf[n:])
	buf[len(buf)-1] = '"'
	return buf[:len(buf)-(n-1)], nil
}

func (f *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}
	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}
	*f = ID(i)
	return nil
}
