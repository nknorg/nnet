package chord

import (
	"bytes"
	"log"
	"math/big"
)

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

func IDToBigInt(id []byte) big.Int {
	idInt := big.Int{}
	idInt.SetBytes(id)
	return idInt
}

// Checks if a key is STRICTLY between two ID's exclusively
func Between(id1, id2, key []byte) bool {
	if CompareID(id1, id2) > 0 {
		return CompareID(id1, key) < 0 || CompareID(id2, key) > 0
	}
	return CompareID(id1, key) < 0 && CompareID(id2, key) > 0
}

// Checks if a key is between two ID's, left inclusive
func BetweenLeftIncl(id1, id2, key []byte) bool {
	if CompareID(id1, id2) > 0 {
		return CompareID(id1, key) <= 0 || CompareID(id2, key) > 0
	}
	return CompareID(id1, key) <= 0 && CompareID(id2, key) > 0
}

// Checks if a key is between two ID's, right inclusive
func BetweenRightIncl(id1, id2, key []byte) bool {
	if CompareID(id1, id2) > 0 {
		return CompareID(id1, key) < 0 || CompareID(id2, key) >= 0
	}
	return CompareID(id1, key) < 0 && CompareID(id2, key) >= 0
}

// Checks if a key is between two ID's, both inclusive
func BetweenIncl(id1, id2, key []byte) bool {
	if CompareID(id1, id2) > 0 {
		return CompareID(id1, key) <= 0 || CompareID(id2, key) >= 0
	}
	return CompareID(id1, key) <= 0 && CompareID(id2, key) >= 0
}

// Computes (id + 2^exp) % (2^mod)
func PowerOffsetBigInt(id []byte, exp uint32, mod uint32) big.Int {
	off := make([]byte, len(id))
	copy(off, id)

	idInt := IDToBigInt(id)

	two := big.NewInt(2)
	offset := big.Int{}
	offset.Exp(two, big.NewInt(int64(exp)), nil)

	sum := big.Int{}
	sum.Add(&idInt, &offset)

	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(mod)), nil)

	idInt.Mod(&sum, &ceil)

	return idInt
}

// Computes (id + 2^exp) % (2^mod)
func PowerOffset(id []byte, exp uint32, mod uint32) []byte {
	resInt := PowerOffsetBigInt(id, exp, mod)
	return BigIntToID(resInt, mod)
}

// Computes (id - 1) % (2^mod)
func PrevID(id []byte, mod uint32) []byte {
	idInt := IDToBigInt(id)
	prev := big.Int{}
	prev.Sub(&idInt, big.NewInt(1))
	return BigIntToID(prev, mod)
}

// Computes (id + 1) % (2^mod)
func NextID(id []byte, mod uint32) []byte {
	idInt := IDToBigInt(id)
	next := big.Int{}
	next.Add(&idInt, big.NewInt(1))
	return BigIntToID(next, mod)
}

// Computes the forward distance from a to b modulus a ring size
func Distance(a, b []byte, bits uint32) *big.Int {
	// Get the ring size
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	// Convert to int
	var aInt, bInt big.Int
	(&aInt).SetBytes(a)
	(&bInt).SetBytes(b)

	// Compute the distances
	var dist big.Int
	(&dist).Sub(&bInt, &aInt)

	// Distance modulus ring size
	(&dist).Mod(&dist, &ring)
	return &dist
}
