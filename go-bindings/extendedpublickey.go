package blschia

// #cgo LDFLAGS: -L../build -lchiabls -lstdc++ -lgmp
// #cgo CXXFLAGS: -std=c++14 -I../src -I../build/contrib/relic/include -I../contrib/relic/include
// #include <stdbool.h>
// #include <stdlib.h>
// #include "blschia.h"
import "C"
import (
	"encoding/hex"
	"runtime"
)

// ExtendedPublicKey represents a BIP-32 style extended public key
type ExtendedPublicKey struct {
	key C.CExtendedPublicKey
}

// ExtendedPublicKeyFromBytes parses a public key and chain code from bytes
func ExtendedPublicKeyFromBytes(data []byte) *ExtendedPublicKey {
	// Get a C pointer to bytes
	cBytesPtr := C.CBytes(data)
	defer C.free(cBytesPtr)

	var key ExtendedPublicKey
	key.key = C.CExtendedPublicKeyFromBytes(cBytesPtr)
	runtime.SetFinalizer(&key, func(p *ExtendedPublicKey) { p.Free() })

	return &key
}

// ExtendedPublicKeyFromString constructs a new extended public key from hex string
func ExtendedPublicKeyFromString(hexString string) (*ExtendedPublicKey, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	return ExtendedPublicKeyFromBytes(bytes), nil
}

// Free releases memory allocated by the key
func (key *ExtendedPublicKey) Free() {
	C.CExtendedPublicKeyFree(key.key)
	runtime.KeepAlive(key)
}

// Serialize returns the serialized byte representation of the
// ExtendedPublicKey object
func (key *ExtendedPublicKey) Serialize() []byte {
	ptr := C.CExtendedPublicKeySerialize(key.key)
	defer C.free(ptr)
	runtime.KeepAlive(key)
	return C.GoBytes(ptr, C.CExtendedPublicKeySizeBytes())
}

// GetPublicKey returns the public key for the given node
func (key *ExtendedPublicKey) GetPublicKey() *PublicKey {
	var pk PublicKey
	pk.pk = C.CExtendedPublicKeyGetPublicKey(key.key)

	runtime.SetFinalizer(&pk, func(p *PublicKey) { p.Free() })
	runtime.KeepAlive(key)

	return &pk
}

var childComparator uint32 = (1 << 31)

// PublicChild derives a child extended public key
func (key *ExtendedPublicKey) PublicChild(i uint32) *ExtendedPublicKey {
	// Hardened children have i >= 2^31. Non-hardened have i < 2^31
	if i >= childComparator {
		panic("cannot derive hardened children from public key")
	}
	if key.GetDepth() >= 255 {
		panic("cannot go further than 255 levels")
	}

	var child ExtendedPublicKey
	child.key = C.CExtendedPublicKeyPublicChild(key.key, C.uint(i))

	runtime.SetFinalizer(&child, func(p *ExtendedPublicKey) { p.Free() })
	runtime.KeepAlive(key)

	return &child
}

// GetVersion returns the version bytes
func (key *ExtendedPublicKey) GetVersion() uint32 {
	version := uint32(C.CExtendedPublicKeyGetVersion(key.key))
	runtime.KeepAlive(key)
	return version
}

// GetDepth returns the depth byte
func (key *ExtendedPublicKey) GetDepth() uint8 {
	depth := uint8(C.CExtendedPublicKeyGetDepth(key.key))
	runtime.KeepAlive(key)
	return depth
}

// GetParentFingerprint returns the parent fingerprint
func (key *ExtendedPublicKey) GetParentFingerprint() uint32 {
	parentFingerprint := uint32(C.CExtendedPublicKeyGetParentFingerprint(key.key))
	runtime.KeepAlive(key)
	return parentFingerprint
}

// GetChildNumber returns the child number
func (key *ExtendedPublicKey) GetChildNumber() uint32 {
	childNumber := uint32(C.CExtendedPublicKeyGetChildNumber(key.key))
	runtime.KeepAlive(key)
	return childNumber
}

// GetChainCode returns the ChainCode for the given node
func (key *ExtendedPublicKey) GetChainCode() *ChainCode {
	var cc ChainCode
	cc.cc = C.CExtendedPublicKeyGetChainCode(key.key)

	runtime.SetFinalizer(&cc, func(p *ChainCode) { p.Free() })
	runtime.KeepAlive(key)

	return &cc
}

// Equal tests if one ExtendedPublicKey object is equal to another
func (key *ExtendedPublicKey) Equal(other *ExtendedPublicKey) bool {
	isEqual := bool(C.CExtendedPublicKeyIsEqual(key.key, other.key))
	runtime.KeepAlive(key)
	runtime.KeepAlive(other)
	return isEqual
}
