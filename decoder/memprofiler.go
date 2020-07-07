// Copyright 2017 XUEQIU.COM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package decoder

import (
	"math/rand"
	"sort"
	"strconv"
)

var (
	skiplistMaxLevel    = 32
	skiplistP           = 0.25
	redisSharedInterges = int64(10000)
	longSize            = uint64(8)
	pointerSize         = uint64(8)
	jemallocSizeClasses = []uint64{
		8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 448, 512, 640, 768, 896, 1024,
		1280, 1536, 1792, 2048, 2560, 3072, 3584, 4096, 5120, 6144, 7168, 8192, 10240, 12288, 14336, 16384, 20480, 24576,
		28672, 32768, 40960, 49152, 57344, 65536, 81920, 98304, 114688, 131072, 163840, 196608, 229376, 262144, 327680,
		393216, 458752, 524288, 655360, 786432, 917504, 1048576, 1310720, 1572864, 1835008, 2097152, 2621440, 3145728,
		3670016, 4194304, 5242880, 6291456, 7340032, 8388608, 10485760, 12582912, 14680064, 16777216, 20971520, 25165824,
		29360128, 33554432, 41943040, 50331648, 58720256, 67108864, 83886080, 100663296, 117440512, 134217728, 167772160,
		201326592, 234881024, 268435456, 335544320, 402653184, 469762048, 536870912, 671088640, 805306368, 939524096,
		1073741824, 1342177280, 1610612736, 1879048192, 2147483648, 2684354560, 3221225472, 3758096384, 4294967296,
		5368709120, 6442450944, 7516192768, 8589934592, 10737418240, 12884901888, 15032385536, 17179869184, 21474836480,
		25769803776, 30064771072, 34359738368, 42949672960, 51539607552, 60129542144, 68719476736, 85899345920,
		103079215104, 120259084288, 137438953472, 171798691840, 206158430208, 240518168576, 274877906944, 343597383680,
		412316860416, 481036337152, 549755813888, 687194767360, 824633720832, 962072674304, 1099511627776, 1374389534720,
		1649267441664, 1924145348608, 2199023255552, 2748779069440, 3298534883328, 3848290697216, 4398046511104,
		5497558138880, 6597069766656, 7696581394432, 8796093022208, 10995116277760, 13194139533312, 15393162788864,
		17592186044416, 21990232555520, 26388279066624, 30786325577728, 35184372088832, 43980465111040, 52776558133248,
		61572651155456, 70368744177664, 87960930222080, 105553116266496, 123145302310912, 140737488355328, 175921860444160,
		211106232532992, 246290604621824, 281474976710656, 351843720888320, 422212465065984, 492581209243648,
		562949953421312, 703687441776640, 844424930131968, 985162418487296, 1125899906842624, 1407374883553280,
		1688849860263936, 1970324836974592, 2251799813685248, 2814749767106560, 3377699720527872, 3940649673949184,
		4503599627370496, 5629499534213120, 6755399441055744, 7881299347898368, 9007199254740992, 11258999068426240,
		13510798882111488, 15762598695796736, 18014398509481984, 22517998136852480, 27021597764222976, 31525197391593472,
		36028797018963968, 45035996273704960, 54043195528445952, 63050394783186944, 72057594037927936, 90071992547409920,
		108086391056891904, 126100789566373888, 144115188075855872, 180143985094819840, 216172782113783808,
		252201579132747776, 288230376151711744, 360287970189639680, 432345564227567616, 504403158265495552,
		576460752303423488, 720575940379279360, 864691128455135232, 1008806316530991104, 1152921504606846976,
		1441151880758558720, 1729382256910270464, 2017612633061982208, 2305843009213693952, 2882303761517117440,
		3458764513820540928, 4035225266123964416, 4611686018427387904, 5764607523034234880, 6917529027641081856,
		8070450532247928832, 9223372036854775808, 11529215046068469760, 13835058055282163712, 16140901064495857664,
	}
)

// MemProfiler get memory use for all kinds of data stuct
type MemProfiler struct{}

// mallocOverhead used memory
func (m *MemProfiler) mallocOverhead(size uint64) uint64 {
	idx := sort.Search(len(jemallocSizeClasses),
		func(i int) bool { return jemallocSizeClasses[i] >= size })
	if idx < len(jemallocSizeClasses) {
		return jemallocSizeClasses[idx]
	}
	return size
}

// TopLevelObjOverhead get memory use of a top level object
// Each top level object is an entry in a dictionary, and so we have to include
// the overhead of a dictionary entry
func (m *MemProfiler) TopLevelObjOverhead(key []byte, expiry int64) uint64 {
	return m.HashtableEntryOverhead() + m.SizeofString(key) + m.RobjOverhead() + m.KeyExpiryOverhead(expiry)
}

// HashtableOverhead get memory use of a hashtable
// See  https://github.com/antirez/redis/blob/unstable/src/dict.h
// See the structures dict and dictht
// 2 * (3 unsigned longs + 1 pointer) + int + long + 2 pointers
//
// Additionally, see **table in dictht
// The length of the table is the next power of 2
// When the hashtable is rehashing, another instance of **table is created
// Due to the possibility of rehashing during loading, we calculate the worse
// case in which both tables are allocated, and so multiply
// the size of **table by 1.5
func (m *MemProfiler) HashtableOverhead(size uint64) uint64 {
	return 4 + 7*longSize + 4*pointerSize + nextPower(size)*pointerSize*3/2
}

func (m *MemProfiler) SizeofStreamRadixTree(numElements uint64) uint64 {
	numNodes := uint64(float64(numElements) * 2.5)
	return 16 * numElements + numNodes * 4 + numNodes * 30 * 8
}

func (m *MemProfiler) StreamOverhead() uint64 {
	return 2 * pointerSize + 8 + 16 + // stream struct
		pointerSize + 8 * 2 // rax struct
}

func (m *MemProfiler) StreamConsumer(name []byte) uint64 {
	return pointerSize * 2 + 8 + m.SizeofString(name)
}

func (m *MemProfiler) StreamCG() uint64 {
	return pointerSize * 2 + 16
}

func (m *MemProfiler) StreamNACK(length uint64) uint64 {
	return length * (pointerSize + 8 + 8)
}

// HashtableEntryOverhead get memory use of hashtable entry
// See  https://github.com/antirez/redis/blob/unstable/src/dict.h
// Each dictEntry has 3 pointers
// typedef struct dictEntry {
//     void *key;
//     union {
//         void *val;
//         uint64_t u64;
//         int64_t s64;
//         double d;
//     } v;
//     struct dictEntry *next;
// } dictEntry;
func (m *MemProfiler) HashtableEntryOverhead() uint64 {
	return 3 * pointerSize
}

// LinkedlistOverhead get memory use of a linked list
// See https://github.com/antirez/redis/blob/unstable/src/adlist.h
// A list has 5 pointers + an unsigned long
func (m *MemProfiler) LinkedlistOverhead() uint64 {
	return longSize + 5*pointerSize
}

// LinkedListEntryOverhead get memory use of a linked list entry
// See https://github.com/antirez/redis/blob/unstable/src/adlist.h
// A node has 3 pointers
func (m *MemProfiler) LinkedListEntryOverhead() uint64 {
	return 3 * pointerSize
}

// SkiplistOverhead get memory use of a skiplist
func (m *MemProfiler) SkiplistOverhead(size uint64) uint64 {
	return 2*pointerSize + m.HashtableOverhead(size) + (2*pointerSize + 16)
}

// SkiplistEntryOverhead get memory use of a skiplist entry
func (m *MemProfiler) SkiplistEntryOverhead() uint64 {
	return m.HashtableEntryOverhead() + 2*pointerSize + 8 + (pointerSize+8)*zsetRandLevel()
}

func (m *MemProfiler) QuicklistOverhead(size uint64) uint64 {
	quicklist := 2 * pointerSize + 8 + 2 * 4
	quickitem := 4 * pointerSize + 8 + 2 * 4
	return quicklist + size * quickitem
}

func (m *MemProfiler) ZiplistHeaderOverhead() uint64 {
	return 4 + 4 + 2 + 1
}

func (m *MemProfiler) ZiplistEntryOverhead(value []byte) uint64 {
	header := 0
	size := 0

	if n, err := strconv.ParseInt(string(value), 10, 64); err == nil {
		header = 1
		switch {
		case n < 12: size = 0
		case n < 256: size = 1
		case n < 65536: size = 2
		case n < 16777216: size = 3
		case n < 4294967296: size = 4
		default:
			size = 8
		}
	} else {
		size = len(value)
		if size <= 63 {
			header = 1
		} else if size <= 16383 {
			header = 2
		} else {
			header = 5

			if size >= 254 {
				header += 5
			}
		}
	}

	return uint64(header + size)
}

// KeyExpiryOverhead get memory useage of a key expiry
// Key expiry is stored in a hashtable, so we have to pay for the cost of a hashtable entry
// The timestamp itself is stored as an int64, which is a 8 bytes
func (m *MemProfiler) KeyExpiryOverhead(expiry int64) uint64 {
	//If there is no expiry, there isn't any overhead
	if expiry <= 0 {
		return 0
	}
	return m.HashtableEntryOverhead() + 8
}

// RobjOverhead get memory useage of a robj
// typedef struct redisobject {
//     unsigned type:4;
//     unsigned encoding:4;
//     unsigned lru:lru_bits; /* lru time (relative to server.lruclock) */
//     int refcount;
//     void *ptr;
// } robj;
const LRU_BITS = 24

func (m *MemProfiler) RobjOverhead() uint64 {
	return pointerSize + 4 + 4 + LRU_BITS + 4
}

// SizeofString get memory use of a string
// https://github.com/antirez/redis/blob/unstable/src/sds.h
func (m *MemProfiler) SizeofString(bytes []byte) uint64 {
	str := string(bytes)
	num, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		if num < redisSharedInterges && num > 0 {
			return 0
		}
		return 8
	}
	l := uint64(len(str))
	return m.mallocOverhead(l + 8 + 1)
}

// ElemLen get length of a element
func (m *MemProfiler) ElemLen(element []byte) uint64 {
	MaxInt64 := int64(1<<63 - 1)
	MinInt64 := int64(-1 << 63)
	if num, err := strconv.ParseInt(string(element), 10, 64); err == nil {
		if num < MinInt64 || num > MaxInt64 {
			return 16
		}
		return 8
	}
	return uint64(len(element))
}

func nextPower(size uint64) uint64 {
	power := uint64(1)
	for power <= size {
		power = power << 1
	}
	return power
}

func zsetRandLevel() uint64 {
	level := 1
	rint := rand.Intn(0xFFFF)
	for rint < int(0xFFFF*1/4) { //skiplistP
		level++
		rint = rand.Intn(0xFFFF)
	}
	if level < skiplistMaxLevel {
		return uint64(level)
	}
	return uint64(skiplistMaxLevel)
}
