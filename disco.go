package disco

import (
	"math/bits"
	"reflect"
	"unsafe"
)

const (
	_STATE     = 32
	_STATE64   = _STATE >> 3
	_STATEM    = _STATE - 1
	_HSTATE64M = (_STATE64 >> 1) - 1
	_STATE64M  = _STATE64 - 1
	_P         = 0xFFFFFFFFFFFFFFFF - 58
	_Q         = 13166748625691186689
)

type state struct {
	ds  [4]uint64
	ds8 []uint8
}

func newState() state {
	var s state
	s.ds8 = u64to8(s.ds[:])
	return s
}

func (s *state) mix(A int) {
	B := A + 1
	s.ds[A] *= _P
	s.ds[A] = bits.RotateLeft64(s.ds[A], 64-23)
	s.ds[A] *= _Q

	s.ds[B] ^= s.ds[A]

	s.ds[B] *= _P
	s.ds[B] = bits.RotateLeft64(s.ds[B], 64-23)
	s.ds[B] *= _Q
}

func (s *state) round(m64 []uint64, m8 []uint8, len int) {
	var index uint64
	var sindex int64
	var Len = uint64(len) >> 3
	var counter uint64 = 0xfaccadaccad09997
	var counter8 uint8 = 137

	for index = 0; index < Len; index++ {
		s.ds[sindex] += bits.RotateLeft64(m64[index]+index+counter+1, 64-23)
		counter += ^m64[index] + 1
		if sindex == _HSTATE64M {
			s.mix(0)
		} else if sindex == _STATE64M {
			s.mix(2)
			sindex = -1
		}
		sindex++
	}

	s.mix(1)

	Len = index << 3
	sindex = int64(index) & (_STATEM)

	for index = Len; index < uint64(len); index++ {
		s.ds8[sindex] += bits.RotateLeft8(m8[index]+uint8(index)+counter8+1, 8-(23%8))
		counter8 += ^m8[sindex] + 1
		s.mix(int(index) % _STATE64M)
		if sindex >= _STATEM {
			sindex = -1
		}
		sindex++
	}

	s.mix(0)
	s.mix(1)
	s.mix(2)
}

func u64to8(s []uint64) []uint8 {
	h := *(*reflect.SliceHeader)(unsafe.Pointer(&s))
	h.Len *= 8
	h.Cap *= 8
	buf := *(*[]uint8)(unsafe.Pointer(&h))
	return buf
}

func u64to32(s []uint64) []uint32 {
	h := *(*reflect.SliceHeader)(unsafe.Pointer(&s))
	h.Len *= 2
	h.Cap *= 2
	buf := *(*[]uint32)(unsafe.Pointer(&h))
	return buf
}

func u8to64(s []uint8) []uint64 {
	h := *(*reflect.SliceHeader)(unsafe.Pointer(&s))
	h.Len /= 8
	h.Cap /= 8
	buf := *(*[]uint64)(unsafe.Pointer(&h))
	return buf
}

func BEBB4185_64(key []uint8, seed uint64) uint64 {
	s := newState()

	var seed64 [2]uint64
	seed32 := u64to32(seed64[:])
	seed8 := u64to8(seed64[:])

	key64 := u8to64(key)

	// the cali number from the Matrix (1999)
	seed32[0] = 0xc5550690
	seed32[0] -= uint32(seed)
	seed32[1] = 1 + uint32(seed)
	seed32[2] = ^(1 - uint32(seed))
	seed32[3] = (1 + uint32(seed)) * 0xf00dacca

	// nothing up my sleeve
	s.ds[0] = 0x123456789abcdef0
	s.ds[1] = 0x0fedcba987654321
	s.ds[2] = 0xaccadacca80081e5
	s.ds[3] = 0xf00baaf00f00baaa

	s.round(key64, key, len(key))
	s.round(seed64[:], seed8, 16)
	s.round(s.ds[:], s.ds8, _STATE)

	var h [_STATE / 8]uint64

	h[0] = s.ds[2]
	h[1] = s.ds[3]

	h[0] += h[1]

	return h[0]
}
