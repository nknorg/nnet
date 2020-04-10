package chord

import (
	"bytes"
	"log"
	"math/big"
	"sync"
)

var (
	two            = big.NewInt(2)
	powerCache     = make(map[uint32]*big.Int, 1)
	powerCacheLock sync.RWMutex
)

func getPower(exp uint32) *big.Int {
	powerCacheLock.RLock()
	power := powerCache[exp]
	powerCacheLock.RUnlock()

	if power == nil {
		power = &big.Int{}
		power.Exp(two, big.NewInt(int64(exp)), nil)
		powerCacheLock.Lock()
		powerCache[exp] = power
		powerCacheLock.Unlock()
	}

	return power
}

// CompareID returns -1, 0, 1 if id1 <, =, > id2 respectively
func CompareID(id1, id2 []byte) int {
	l1, l2 := len(id1), len(id2)
	if l1 > l2 {
		return -CompareID(id2, id1)
	}
	if l1 < l2 {
		tmp := make([]byte, l2)
		copy(tmp[l2-l1:], id1)
		id1 = tmp
	}
	return bytes.Compare(id1, id2)
}

// BigIntToID converts a big.Int with m bit rate to Chord ID bytes
func BigIntToID(i big.Int, m uint32) []byte {
	b := i.Bytes()
	lb := uint32(len(b))
	lID := m / 8
	if lb < lID {
		id := make([]byte, lID)
		copy(id[lID-lb:], b)
		return id
	}
	if lb > lID {
		log.Println("[WARNING] Big integer has more bytes than ID.")
	}
	return b
}

// IDToBigInt converts a Chord ID bytes to big.Int
func IDToBigInt(id []byte) big.Int {
	idInt := big.Int{}
	idInt.SetBytes(id)
	return idInt
}

// Between checks if a key is between two ID exclusively
func Between(id1, id2, key []byte) bool {
	if CompareID(id1, id2) > 0 {
		return CompareID(id1, key) < 0 || CompareID(id2, key) > 0
	}
	return CompareID(id1, key) < 0 && CompareID(id2, key) > 0
}

// BetweenLeftIncl checks if a key is between two ID, left inclusive
func BetweenLeftIncl(id1, id2, key []byte) bool {
	if CompareID(id1, id2) > 0 {
		return CompareID(id1, key) <= 0 || CompareID(id2, key) > 0
	}
	return CompareID(id1, key) <= 0 && CompareID(id2, key) > 0
}

// BetweenRightIncl checks if a key is between two ID, right inclusive
func BetweenRightIncl(id1, id2, key []byte) bool {
	if CompareID(id1, id2) > 0 {
		return CompareID(id1, key) < 0 || CompareID(id2, key) >= 0
	}
	return CompareID(id1, key) < 0 && CompareID(id2, key) >= 0
}

// BetweenIncl checks if a key is between two ID, both inclusive
func BetweenIncl(id1, id2, key []byte) bool {
	if CompareID(id1, id2) > 0 {
		return CompareID(id1, key) <= 0 || CompareID(id2, key) >= 0
	}
	return CompareID(id1, key) <= 0 && CompareID(id2, key) >= 0
}

// Offset computes (id + offset) % (2^m)
func Offset(id []byte, offset *big.Int, m uint32) []byte {
	idInt := IDToBigInt(id)

	sum := big.Int{}
	sum.Add(&idInt, offset)

	idInt.Mod(&sum, getPower(m))

	return BigIntToID(idInt, m)
}

// PowerOffset computes (id + 2^exp) % (2^m)
func PowerOffset(id []byte, exp uint32, m uint32) []byte {
	offset := big.Int{}
	offset.Exp(two, big.NewInt(int64(exp)), nil)
	return Offset(id, &offset, m)
}

// PrevID computes (id - 1) % (2^m)
func PrevID(id []byte, m uint32) []byte {
	return Offset(id, big.NewInt(-1), m)
}

// NextID computes (id + 1) % (2^m)
func NextID(id []byte, m uint32) []byte {
	return Offset(id, big.NewInt(1), m)
}

// Distance gcomputes the forward distance from a to b modulus a ring size
func Distance(a, b []byte, m uint32) *big.Int {
	// Convert to int
	var aInt, bInt big.Int
	(&aInt).SetBytes(a)
	(&bInt).SetBytes(b)

	// Compute the distances
	var dist big.Int
	(&dist).Sub(&bInt, &aInt)

	// Distance modulus ring size
	(&dist).Mod(&dist, getPower(m))
	return &dist
}
