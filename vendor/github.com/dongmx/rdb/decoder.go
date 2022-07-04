// Package rdb implements parsing and encoding of the Redis RDB file format.
package rdb

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/dongmx/rdb/crc64"
	"github.com/juju/errors"
)

type Info struct {
	Encoding    string
	Idle        uint64
	Freq        int
	SizeOfValue int
	Zips        uint64
}

type StreamPendingEntry struct {
	ID            []byte
	DeliveryTime  uint64
	DeliveryCount uint64
}

type StreamConsumerPendingEntry struct {
	ID []byte
}

type StreamConsumerData struct {
	Name     []byte
	SeenTime uint64
	Pending  []*StreamConsumerPendingEntry
}

type StreamGroup struct {
	Name        []byte
	LastEntryId string
	Pending     []*StreamPendingEntry
	Consumers   []*StreamConsumerData
}

type StreamGroups []*StreamGroup

// A Decoder must be implemented to parse a RDB file.
type Decoder interface {
	// StartRDB is called when parsing of a valid RDB file starts.
	StartRDB(ver int)
	// StartDatabase is called when database n starts.
	// Once a database starts, another database will not start until EndDatabase is called.
	StartDatabase(n int)
	// AUX field
	Aux(key, value []byte)
	// ResizeDB hint
	ResizeDatabase(dbSize, expiresSize uint32)
	// Set is called once for each string key.
	Set(key, value []byte, expiry int64, info *Info)
	// StartHash is called at the beginning of a hash.
	// Hset will be called exactly length times before EndHash.
	StartHash(key []byte, length, expiry int64, info *Info)
	// Hset is called once for each field=value pair in a hash.
	Hset(key, field, value []byte)
	// EndHash is called when there are no more fields in a hash.
	EndHash(key []byte)
	// StartSet is called at the beginning of a set.
	// Sadd will be called exactly cardinality times before EndSet.
	StartSet(key []byte, cardinality, expiry int64, info *Info)
	// Sadd is called once for each member of a set.
	Sadd(key, member []byte)
	// EndSet is called when there are no more fields in a set.
	EndSet(key []byte)
	// StartStream is called at the beginning of a stream.
	// Xadd will be called exactly length times before EndStream.
	StartStream(key []byte, cardinality, expiry int64, info *Info)
	// Xadd is called once for each id in a stream.
	Xadd(key, id, listpack []byte)
	// EndHash is called when there are no more fields in a hash.
	EndStream(key []byte, items uint64, lastEntryID string, cgroupsData StreamGroups)
	// StartList is called at the beginning of a list.
	// Rpush will be called exactly length times before EndList.
	// If length of the list is not known, then length is -1
	StartList(key []byte, length, expiry int64, info *Info)
	// Rpush is called once for each value in a list.
	Rpush(key, value []byte)
	// EndList is called when there are no more values in a list.
	EndList(key []byte)
	// StartZSet is called at the beginning of a sorted set.
	// Zadd will be called exactly cardinality times before EndZSet.
	StartZSet(key []byte, cardinality, expiry int64, info *Info)
	// Zadd is called once for each member of a sorted set.
	Zadd(key []byte, score float64, member []byte)
	// EndZSet is called when there are no more members in a sorted set.
	EndZSet(key []byte)
	// EndDatabase is called at the end of a database.
	EndDatabase(n int)
	// EndRDB is called when parsing of the RDB file is complete.
	EndRDB()
}

// Decode parses a RDB file from r and calls the decode hooks on d.
func Decode(r io.Reader, d Decoder) error {
	decoder := &decode{d, make([]byte, 8), bufio.NewReader(r), 0, 0, nil, 0}
	return decoder.decode()
}

// DecodeDump a byte slice from the Redis DUMP command. The dump does not contain the
// database, key or expiry, so they must be included in the function call (but
// can be zero values).
func DecodeDump(dump []byte, db int, key []byte, expiry int64, d Decoder) error {
	err := verifyDump(dump)
	if err != nil {
		return errors.Trace(err)
	}

	decoder := &decode{d, make([]byte, 8), bytes.NewReader(dump[1:]), 0, 0, nil, 0}
	decoder.event.StartRDB(0)
	decoder.event.StartDatabase(db)

	err = decoder.readObject(key, ValueType(dump[0]), expiry)

	decoder.event.EndDatabase(db)
	decoder.event.EndRDB()
	return errors.Trace(err)
}

type byteReader interface {
	io.Reader
	io.ByteReader
}

type decode struct {
	event  Decoder
	intBuf []byte
	r      byteReader

	lruIdle uint64
	lfuFreq int

	info       *Info
	rdbVersion int
}

// ValueType of redis type
type ValueType byte

// type value
const (
	TypeString  ValueType = 0
	TypeList    ValueType = 1
	TypeSet     ValueType = 2
	TypeZSet    ValueType = 3
	TypeHash    ValueType = 4
	TypeZSet2   ValueType = 5
	TypeModule  ValueType = 6
	TypeModule2 ValueType = 7

	TypeHashZipmap      ValueType = 9
	TypeListZiplist     ValueType = 10
	TypeSetIntset       ValueType = 11
	TypeZSetZiplist     ValueType = 12
	TypeHashZiplist     ValueType = 13
	TypeListQuicklist   ValueType = 14
	TypeStreamListPacks ValueType = 15
)

const (
	rdbVersion  = 9
	rdb6bitLen  = 0
	rdb14bitLen = 1
	rdb32bitLen = 0x80
	rdb64bitLen = 0x81
	rdbEncVal   = 3
	rdbLenErr   = math.MaxUint64

	rdbOpCodeModuleAux = 247
	rdbOpCodeIdle      = 248
	rdbOpCodeFreq      = 249
	rdbOpCodeAux       = 250
	rdbOpCodeResizeDB  = 251
	rdbOpCodeExpiryMS  = 252
	rdbOpCodeExpiry    = 253
	rdbOpCodeSelectDB  = 254
	rdbOpCodeEOF       = 255

	rdbModuleOpCodeEOF    = 0
	rdbModuleOpCodeSint   = 1
	rdbModuleOpCodeUint   = 2
	rdbModuleOpCodeFloat  = 3
	rdbModuleOpCodeDouble = 4
	rdbModuleOpCodeString = 5

	rdbLoadNone  = 0
	rdbLoadEnc   = (1 << 0)
	rdbLoadPlain = (1 << 1)
	rdbLoadSds   = (1 << 2)

	rdbSaveNode        = 0
	rdbSaveAofPreamble = (1 << 0)

	rdbEncInt8  = 0
	rdbEncInt16 = 1
	rdbEncInt32 = 2
	rdbEncLZF   = 3

	rdbZiplist6bitlenString  = 0
	rdbZiplist14bitlenString = 1
	rdbZiplist32bitlenString = 2

	rdbZiplistInt16 = 0xc0
	rdbZiplistInt32 = 0xd0
	rdbZiplistInt64 = 0xe0
	rdbZiplistInt24 = 0xf0
	rdbZiplistInt8  = 0xfe
	rdbZiplistInt4  = 15

	rdbLpHdrSize           = 6
	rdbLpHdrNumeleUnknown  = math.MaxUint16
	rdbLpMaxIntEncodingLen = 0
	rdbLpMaxBacklenSize    = 5
	rdbLpMaxEntryBacklen   = 34359738367
	rdbLpEncodingInt       = 0
	rdbLpEncodingString    = 1

	rdbLpEncoding7BitUint     = 0
	rdbLpEncoding7BitUintMask = 0x80

	rdbLpEncoding6BitStr     = 0x80
	rdbLpEncoding6BitStrMask = 0xC0

	rdbLpEncoding13BitInt     = 0xC0
	rdbLpEncoding13BitIntMask = 0xE0

	rdbLpEncoding12BitStr     = 0xE0
	rdbLpEncoding12BitStrMask = 0xF0

	rdbLpEncoding16BitInt     = 0xF1
	rdbLpEncoding16BitIntMask = 0xFF

	rdbLpEncoding24BitInt     = 0xF2
	rdbLpEncoding24BitIntMask = 0xFF

	rdbLpEncoding32BitInt     = 0xF3
	rdbLpEncoding32BitIntMask = 0xFF

	rdbLpEncoding64BitInt     = 0xF4
	rdbLpEncoding64BitIntMask = 0xFF

	rdbLpEncoding32BitStr     = 0xF0
	rdbLpEncoding32BitStrMask = 0xFF

	rdbLpEOF = 0xFF
)

func (d *decode) decode() error {
	err := d.checkHeader()
	if err != nil {
		return errors.Trace(err)
	}
	d.event.StartRDB(d.rdbVersion)
	var db uint64
	var expiry int64
	//var lruClock int64
	firstDB := true
	for {
		d.lruIdle = 0
		d.lfuFreq = 0

		objType, err := d.r.ReadByte()
		if err != nil {
			return errors.Wrap(err, errors.New("readfailed"))
		}
		switch objType {
		case rdbOpCodeFreq:
			b, err := d.r.ReadByte()
			d.lfuFreq = int(b)
			if err != nil {
				return errors.Trace(err)
			}
		case rdbOpCodeIdle:
			idle, _, err := d.readLength()
			if err != nil {
				return errors.Trace(err)
			}
			d.lruIdle = uint64(idle)
		case rdbOpCodeAux:
			auxKey, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			auxVal, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			d.event.Aux(auxKey, auxVal)
		case rdbOpCodeResizeDB:
			dbSize, _, err := d.readLength()
			if err != nil {
				return errors.Trace(err)
			}
			expiresSize, _, err := d.readLength()
			if err != nil {
				return errors.Trace(err)
			}
			d.event.ResizeDatabase(uint32(dbSize), uint32(expiresSize))
		case rdbOpCodeExpiryMS:
			_, err := io.ReadFull(d.r, d.intBuf)
			if err != nil {
				return errors.Trace(err)
			}
			expiry = int64(binary.LittleEndian.Uint64(d.intBuf))
		case rdbOpCodeExpiry:
			_, err := io.ReadFull(d.r, d.intBuf[:4])
			if err != nil {
				return errors.Trace(err)
			}
			expiry = int64(binary.LittleEndian.Uint32(d.intBuf)) * 1000
		case rdbOpCodeSelectDB:
			if !firstDB {
				d.event.EndDatabase(int(db))
			}
			db, _, err = d.readLength()
			if err != nil {
				return errors.Trace(err)
			}
			d.event.StartDatabase(int(db))
		case rdbOpCodeEOF:
			d.event.EndDatabase(int(db))
			d.event.EndRDB()
			return nil
		case rdbOpCodeModuleAux:
			return errors.Errorf("unsupport module")
		default:
			key, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			err = d.readObject(key, ValueType(objType), expiry)
			if err != nil {
				return errors.Trace(err)
			}
			expiry = 0
		}
	}

}

func (d *decode) readObject(key []byte, typ ValueType, expiry int64) error {
	d.info = &Info{
		Idle: d.lruIdle,
		Freq: d.lfuFreq,
	}
	switch typ {
	case TypeString:
		value, err := d.readString()
		if err != nil {
			return errors.Trace(err)
		}
		d.info.Encoding = "string"
		d.event.Set(key, value, expiry, d.info)
	case TypeList:
		length, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}
		d.info.Encoding = "linkedlist"
		d.event.StartList(key, int64(length), expiry, d.info)
		for length > 0 {
			length--
			value, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			d.event.Rpush(key, value)
		}
		d.event.EndList(key)
	case TypeListQuicklist:
		length, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}
		d.info.Encoding = "quicklist"
		d.info.Zips = length
		d.event.StartList(key, int64(-1), expiry, d.info)
		for length > 0 {
			length--
			d.readZiplist(key, 0, false)
		}
		d.event.EndList(key)
	case TypeSet:
		cardinality, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}
		d.info.Encoding = "hashtable"
		d.event.StartSet(key, int64(cardinality), expiry, d.info)
		for cardinality > 0 {
			cardinality--
			member, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			d.event.Sadd(key, member)
		}
		d.event.EndSet(key)
	case TypeZSet2:
		fallthrough
	case TypeZSet:
		cardinality, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}
		d.info.Encoding = "skiplist"
		d.event.StartZSet(key, int64(cardinality), expiry, d.info)
		for cardinality > 0 {
			cardinality--
			member, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			var score float64
			if typ == TypeZSet2 {
				score, err = d.readBinaryFloat64()
				if err != nil {
					return errors.Trace(err)
				}
			} else {
				score, err = d.readFloat64()
				if err != nil {
					return errors.Trace(err)
				}
			}
			d.event.Zadd(key, score, member)
		}
		d.event.EndZSet(key)
	case TypeHash:
		length, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}
		d.info.Encoding = "hashtable"
		d.event.StartHash(key, int64(length), expiry, d.info)
		for length > 0 {
			length--
			field, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			value, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			d.event.Hset(key, field, value)
		}
		d.event.EndHash(key)
	case TypeHashZipmap:
		return errors.Trace(d.readZipmap(key, expiry))
	case TypeListZiplist:
		return errors.Trace(d.readZiplist(key, expiry, true))
	case TypeSetIntset:
		return errors.Trace(d.readIntset(key, expiry))
	case TypeZSetZiplist:
		return errors.Trace(d.readZiplistZset(key, expiry))
	case TypeHashZiplist:
		return errors.Trace(d.readZiplistHash(key, expiry))
	case TypeStreamListPacks:
		return errors.Trace(d.readStream(key, expiry))
	case TypeModule:
		fallthrough
	case TypeModule2:
		return d.readModule(key, expiry)
	default:
		return fmt.Errorf("rdb: unknown object type %d for key %s", typ, key)
	}
	return nil
}

func (d *decode) readModule(key []byte, expiry int64) error {
	moduleid, _, err := d.readLength()
	if err != nil {
		return errors.Trace(err)
	}
	return fmt.Errorf("Not supported load module %v", moduleid)
}

func (d *decode) readStream(key []byte, expiry int64) error {
	cardinality, _, err := d.readLength()
	if err != nil {
		return errors.Trace(err)
	}
	d.info.Encoding = "listpack"
	d.event.StartStream(key, int64(cardinality), expiry, d.info)
	for cardinality > 0 {
		cardinality--

		streamID, err := d.readString()
		if err != nil {
			return errors.Trace(err)
		}
		/*
			IDms := strconv.FormatUint(binary.BigEndian.Uint64(streamID[:8]), 10)
			IDseq := strconv.FormatUint(binary.BigEndian.Uint64(streamID[8:]), 10)
			fmt.Println(string(key))
			fmt.Println(IDms + "-" + IDseq)
		*/
		listPack, err := d.readString()
		if err != nil {
			return errors.Trace(err)
		}
		d.event.Xadd(key, streamID, listPack)
	}
	var items, lastIDms, lastIDseq uint64
	items, _, err = d.readLength()
	if err != nil {
		return errors.Trace(err)
	}
	lastIDms, _, err = d.readLength()
	if err != nil {
		return errors.Trace(err)
	}
	lastIDseq, _, err = d.readLength()
	if err != nil {
		return errors.Trace(err)
	}

	lastEntryID := fmt.Sprintf("%d-%d", lastIDms, lastIDseq)

	//TODO output consumer groups
	var groupsCount uint64
	groupsCount, _, err = d.readLength()
	if err != nil {
		return errors.Trace(err)
	}

	cgroupsData := make(StreamGroups, 0, groupsCount)
	for groupsCount > 0 {
		groupsCount--

		cgname, err := d.readString()
		if err != nil {
			return errors.Trace(err)
		}
		gIDms, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}
		gIDseq, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}

		lastCgEntryID := fmt.Sprintf("%d-%d", gIDms, gIDseq)

		pelSize, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}

		groupPendingEntries := make([]*StreamPendingEntry, 0, pelSize)
		for pelSize > 0 {
			pelSize--
			// d.readUint64()
			rawid := make([]byte, 16)
			n, err := io.ReadFull(d.r, rawid)
			if err != nil {
				return errors.Trace(err)
			}
			if n != 16 {
				return errors.Errorf("expected %d got %d", 16, n)
			}

			deliveryTime, err := d.readUint64()
			if err != nil {
				return errors.Trace(err)
			}
			deliveryCount, _, err := d.readLength()
			if err != nil {
				return errors.Trace(err)
			}

			groupPendingEntries = append(groupPendingEntries, &StreamPendingEntry{
				ID:            rawid,
				DeliveryTime:  deliveryTime,
				DeliveryCount: deliveryCount,
			})
		}

		consumersNum, _, err := d.readLength()
		if err != nil {
			return errors.Trace(err)
		}

		consumersData := make([]*StreamConsumerData, 0, consumersNum)
		for consumersNum > 0 {
			consumersNum--
			cname, err := d.readString()
			if err != nil {
				return errors.Trace(err)
			}
			seenTime, err := d.readUint64()
			if err != nil {
				return errors.Trace(err)
			}
			pelSize, _, err := d.readLength()
			if err != nil {
				return errors.Trace(err)
			}
			consumerPendingEntries := make([]*StreamConsumerPendingEntry, 0, pelSize)
			for pelSize > 0 {
				pelSize--
				rawid := make([]byte, 16)
				n, err := io.ReadFull(d.r, rawid)
				if err != nil {
					return errors.Trace(err)
				}
				if n != 16 {
					return errors.Errorf("expected %d got %d", 16, n)
				}

				consumerPendingEntries = append(consumerPendingEntries, &StreamConsumerPendingEntry{ID: rawid})
			}

			consumersData = append(consumersData, &StreamConsumerData{
				Name:     cname,
				SeenTime: seenTime,
				Pending:  consumerPendingEntries,
			})
		}

		cgroupsData = append(cgroupsData, &StreamGroup{
			Name:        cgname,
			LastEntryId: lastCgEntryID,
			Pending:     groupPendingEntries,
			Consumers:   consumersData,
		})
	}

	d.event.EndStream(key, items, lastEntryID, cgroupsData)
	return nil
}

func (d *decode) readZipmap(key []byte, expiry int64) error {
	var length int
	zipmap, err := d.readString()
	if err != nil {
		return errors.Trace(err)
	}
	buf := newSliceBuffer(zipmap)
	lenByte, err := buf.ReadByte()
	if err != nil {
		return errors.Trace(err)
	}
	if lenByte >= 254 { // we need to count the items manually
		length, err = countZipmapItems(buf)
		length /= 2
		if err != nil {
			return errors.Trace(err)
		}
	} else {
		length = int(lenByte)
	}
	d.info.Encoding = "zipmap"
	d.info.SizeOfValue = len(zipmap)
	d.event.StartHash(key, int64(length), expiry, d.info)
	for i := 0; i < length; i++ {
		field, err := readZipmapItem(buf, false)
		if err != nil {
			return errors.Trace(err)
		}
		value, err := readZipmapItem(buf, true)
		if err != nil {
			return errors.Trace(err)
		}
		d.event.Hset(key, field, value)
	}
	d.event.EndHash(key)
	return nil
}

func readZipmapItem(buf *sliceBuffer, readFree bool) ([]byte, error) {
	length, free, err := readZipmapItemLength(buf, readFree)
	if err != nil {
		return nil, err
	}
	if length == -1 {
		return nil, nil
	}
	value, err := buf.Slice(length)
	if err != nil {
		return nil, err
	}
	_, err = buf.Seek(int64(free), 1)
	return value, err
}

func countZipmapItems(buf *sliceBuffer) (int, error) {
	n := 0
	for {
		strLen, free, err := readZipmapItemLength(buf, n%2 != 0)
		if err != nil {
			return 0, err
		}
		if strLen == -1 {
			break
		}
		_, err = buf.Seek(int64(strLen)+int64(free), 1)
		if err != nil {
			return 0, err
		}
		n++
	}
	_, err := buf.Seek(0, 0)
	return n, err
}

func readZipmapItemLength(buf *sliceBuffer, readFree bool) (int, int, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	switch b {
	case 253:
		s, err := buf.Slice(5)
		if err != nil {
			return 0, 0, err
		}
		return int(binary.BigEndian.Uint32(s)), int(s[4]), nil
	case 254:
		return 0, 0, fmt.Errorf("rdb: invalid zipmap item length")
	case 255:
		return -1, 0, nil
	}
	var free byte
	if readFree {
		free, err = buf.ReadByte()
	}
	return int(b), int(free), err
}

func (d *decode) readListPack() error {
	listpack, err := d.readString()
	//fmt.Println(len(listpack))
	//fmt.Println(listpack)
	if err != nil {
		return errors.Trace(err)
	}
	buf := newSliceBuffer(listpack)
	buf.Slice(4) // total bytes
	numElements, _ := buf.Slice(2)
	num := int64(binary.LittleEndian.Uint16(numElements))
	if err != nil {
		return errors.Trace(err)
	}
	for {
		num--
		b, _ := buf.Slice(1)
		if b[0] == byte(rdbLpEOF) {
			fmt.Println("eof")
			break
		}
		//lpGet(b, buf)
	}
	return nil
}

func lpGet(b []byte, buf *sliceBuffer) {
	var val int64
	var uval, negstart, negmax uint64
	fmt.Println(b[0], lpEncodingIs7BitUint(b[0]))
	if lpEncodingIs7BitUint(b[0]) {
		fmt.Println("lpEncodingIs7BitUint")
		negstart = math.MaxUint64
		negmax = 0
		uval = uint64(b[0] & 0x7F)
	} else if lpEncodingIs6BitStr(b[0]) {
		fmt.Println("lpEncodingIs6BitStr")
		len := lpEncoding6BitStrLen(b)
		str, _ := buf.Slice(int(len))
		fmt.Print(string(str))
	} else if lpEncodingIs13BitInt(b[0]) {
		fmt.Println("lpEncodingIs13BitInt")
		tmp, _ := buf.Slice(1)
		b = append(b, tmp...)
		uval = (uint64(b[0]&0x1f) << 8) | uint64(b[1])
		negstart = uint64(1) << 12
		negmax = 8191
	} else if lpEncodingIs16BitInt(b[0]) {
		fmt.Println("lpEncodingIs16BitInt")
		tmp, _ := buf.Slice(2)
		b = append(b, tmp...)
		uval = uint64(b[1]) |
			uint64(b[2])<<8
		negstart = uint64(1) << 15
		negmax = math.MaxUint16
	} else if lpEncodingIs24BitInt(b[0]) {
		fmt.Println("lpEncodingIs24BitInt")
		tmp, _ := buf.Slice(3)
		b = append(b, tmp...)
		uval = uint64(b[1]) |
			uint64(b[2])<<8 |
			uint64(b[3])<<16
		negstart = uint64(1) << 23
		negmax = math.MaxUint32 >> 8
	} else if lpEncodingIs32BitInt(b[0]) {
		fmt.Println("lpEncodingIs32BitInt")
		tmp, _ := buf.Slice(4)
		b = append(b, tmp...)
		uval = uint64(b[1]) |
			uint64(b[2])<<8 |
			uint64(b[3])<<16 |
			uint64(b[4])<<24
		negstart = uint64(1) << 31
		negmax = math.MaxUint32
	} else if lpEncodingIs64BitInt(b[0]) {
		fmt.Println("lpEncodingIs64BitInt")
		tmp, _ := buf.Slice(8)
		b = append(b, tmp...)
		uval = uint64(b[1]) |
			uint64(b[2])<<8 |
			uint64(b[3])<<16 |
			uint64(b[4])<<24 |
			uint64(b[5])<<32 |
			uint64(b[6])<<40 |
			uint64(b[7])<<48 |
			uint64(b[8])<<56
		negstart = uint64(1) << 63
		negmax = math.MaxUint64
	} else if lpEncodingIs12BitStr(b[0]) {
		fmt.Println("lpEncodingIs12BitStr")
		tmp, _ := buf.Slice(1)
		b = append(b, tmp...)
		len := lpEncoding12BitStrLen(b)
		fmt.Println(len)
		str, _ := buf.Slice(int(len))
		fmt.Print(string(str))
	} else if lpEncodingIs32BitStr(b[0]) {
		fmt.Println("lpEncodingIs32BitStr")
		tmp, _ := buf.Slice(4)
		b = append(b, tmp...)
		len := lpEncoding32BitStrLen(b)
		str, _ := buf.Slice(int(len))
		fmt.Print(string(str))
	} else {
		fmt.Println("else")
		uval = uint64(12345678900000000) + uint64(b[0])
		negstart = math.MaxUint64
		negmax = 0
	}

	if uval >= negstart {
		uval = negmax - uval
		val = int64(uval)
		val = -val - 1
	} else {
		val = int64(uval)
	}

	fmt.Printf("%d\n", val)
}

func lpEncodingIs7BitUint(b byte) bool {
	return (((b) & rdbLpEncoding7BitUintMask) == rdbLpEncoding7BitUint)
}
func lpEncodingIs6BitStr(b byte) bool {
	return (((b) & rdbLpEncoding6BitStrMask) == rdbLpEncoding6BitStr)
}
func lpEncodingIs13BitInt(b byte) bool {
	return (((b) & rdbLpEncoding13BitIntMask) == rdbLpEncoding13BitInt)
}
func lpEncodingIs12BitStr(b byte) bool {
	return (((b) & rdbLpEncoding12BitStrMask) == rdbLpEncoding12BitStr)
}
func lpEncodingIs16BitInt(b byte) bool {
	return (((b) & rdbLpEncoding16BitIntMask) == rdbLpEncoding16BitInt)
}
func lpEncodingIs24BitInt(b byte) bool {
	return (((b) & rdbLpEncoding24BitIntMask) == rdbLpEncoding24BitInt)
}
func lpEncodingIs32BitInt(b byte) bool {
	return (((b) & rdbLpEncoding32BitIntMask) == rdbLpEncoding32BitInt)
}
func lpEncodingIs64BitInt(b byte) bool {
	return (((b) & rdbLpEncoding64BitIntMask) == rdbLpEncoding64BitInt)
}
func lpEncodingIs32BitStr(b byte) bool {
	return (((b) & rdbLpEncoding32BitStrMask) == rdbLpEncoding32BitStr)
}
func lpEncoding6BitStrLen(b []byte) uint32 {
	return uint32(b[0] & 0x3F)
}
func lpEncoding12BitStrLen(b []byte) uint32 {
	return (uint32((b)[0]&0xF) << 8) | uint32((b)[1])
}
func lpEncoding32BitStrLen(b []byte) uint32 {
	return (uint32(b[1]) << 0) |
		(uint32(b[2]) << 8) |
		(uint32(b[3]) << 16) |
		(uint32(b[4]) << 24)
}

func (d *decode) readZiplist(key []byte, expiry int64, addListEvents bool) error {
	ziplist, err := d.readString()
	if err != nil {
		return errors.Trace(err)
	}
	buf := newSliceBuffer(ziplist)
	length, err := readZiplistLength(buf)
	if err != nil {
		return errors.Trace(err)
	}
	if addListEvents {
		d.info.Encoding = "ziplist"
		d.info.SizeOfValue = len(ziplist)
		d.event.StartList(key, length, expiry, d.info)
	}
	for i := int64(0); i < length; i++ {
		entry, err := readZiplistEntry(buf)
		if err != nil {
			return errors.Trace(err)
		}
		d.event.Rpush(key, entry)
	}
	if addListEvents {
		d.event.EndList(key)
	}
	return nil
}

func (d *decode) readZiplistZset(key []byte, expiry int64) error {
	ziplist, err := d.readString()
	if err != nil {
		return errors.Trace(err)
	}
	buf := newSliceBuffer(ziplist)
	cardinality, err := readZiplistLength(buf)
	if err != nil {
		return errors.Trace(err)
	}
	cardinality /= 2
	d.info.Encoding = "ziplist"
	d.info.SizeOfValue = len(ziplist)
	d.event.StartZSet(key, cardinality, expiry, d.info)
	for i := int64(0); i < cardinality; i++ {
		member, err := readZiplistEntry(buf)
		if err != nil {
			return errors.Trace(err)
		}
		scoreBytes, err := readZiplistEntry(buf)
		if err != nil {
			return errors.Trace(err)
		}
		score, err := strconv.ParseFloat(string(scoreBytes), 64)
		if err != nil {
			return errors.Trace(err)
		}
		d.event.Zadd(key, score, member)
	}
	d.event.EndZSet(key)
	return nil
}

func (d *decode) readZiplistHash(key []byte, expiry int64) error {
	ziplist, err := d.readString()
	if err != nil {
		return errors.Trace(err)
	}
	buf := newSliceBuffer(ziplist)
	length, err := readZiplistLength(buf)
	if err != nil {
		return errors.Trace(err)
	}
	length /= 2
	d.info.Encoding = "ziplist"
	d.info.SizeOfValue = len(ziplist)
	d.event.StartHash(key, length, expiry, d.info)
	for i := int64(0); i < length; i++ {
		field, err := readZiplistEntry(buf)
		if err != nil {
			return errors.Trace(err)
		}
		value, err := readZiplistEntry(buf)
		if err != nil {
			return errors.Trace(err)
		}
		d.event.Hset(key, field, value)
	}
	d.event.EndHash(key)
	return nil
}

func readZiplistLength(buf *sliceBuffer) (int64, error) {
	buf.Seek(8, 0) // skip the zlbytes and zltail
	lenBytes, err := buf.Slice(2)
	if err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint16(lenBytes)), nil
}

func readZiplistEntry(buf *sliceBuffer) ([]byte, error) {
	prevLen, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if prevLen == 254 {
		buf.Seek(4, 1) // skip the 4-byte prevlen
	}

	header, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	switch {
	case header>>6 == rdbZiplist6bitlenString:
		return buf.Slice(int(header & 0x3f))
	case header>>6 == rdbZiplist14bitlenString:
		b, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}
		return buf.Slice((int(header&0x3f) << 8) | int(b))
	case header>>6 == rdbZiplist32bitlenString:
		lenBytes, err := buf.Slice(4)
		if err != nil {
			return nil, err
		}
		return buf.Slice(int(binary.BigEndian.Uint32(lenBytes)))
	case header == rdbZiplistInt16:
		intBytes, err := buf.Slice(2)
		if err != nil {
			return nil, err
		}
		return []byte(strconv.FormatInt(int64(int16(binary.LittleEndian.Uint16(intBytes))), 10)), nil
	case header == rdbZiplistInt32:
		intBytes, err := buf.Slice(4)
		if err != nil {
			return nil, err
		}
		return []byte(strconv.FormatInt(int64(int32(binary.LittleEndian.Uint32(intBytes))), 10)), nil
	case header == rdbZiplistInt64:
		intBytes, err := buf.Slice(8)
		if err != nil {
			return nil, err
		}
		return []byte(strconv.FormatInt(int64(binary.LittleEndian.Uint64(intBytes)), 10)), nil
	case header == rdbZiplistInt24:
		intBytes := make([]byte, 4)
		_, err := buf.Read(intBytes[1:])
		if err != nil {
			return nil, err
		}
		return []byte(strconv.FormatInt(int64(int32(binary.LittleEndian.Uint32(intBytes))>>8), 10)), nil
	case header == rdbZiplistInt8:
		b, err := buf.ReadByte()
		return []byte(strconv.FormatInt(int64(int8(b)), 10)), err
	case header>>4 == rdbZiplistInt4:
		return []byte(strconv.FormatInt(int64(header&0x0f)-1, 10)), nil
	}

	return nil, fmt.Errorf("rdb: unknown ziplist header byte: %d", header)
}

func (d *decode) readIntset(key []byte, expiry int64) error {
	intset, err := d.readString()
	if err != nil {
		return errors.Trace(err)
	}
	buf := newSliceBuffer(intset)
	intSizeBytes, err := buf.Slice(4)
	if err != nil {
		return errors.Trace(err)
	}
	intSize := binary.LittleEndian.Uint32(intSizeBytes)

	if intSize != 2 && intSize != 4 && intSize != 8 {
		return fmt.Errorf("rdb: unknown intset encoding: %d", intSize)
	}

	lenBytes, err := buf.Slice(4)
	if err != nil {
		return errors.Trace(err)
	}
	cardinality := binary.LittleEndian.Uint32(lenBytes)

	d.info.SizeOfValue = len(intset)
	d.info.Encoding = "intset"
	d.event.StartSet(key, int64(cardinality), expiry, d.info)
	for i := uint32(0); i < cardinality; i++ {
		intBytes, err := buf.Slice(int(intSize))
		if err != nil {
			return errors.Trace(err)
		}
		var intString string
		switch intSize {
		case 2:
			intString = strconv.FormatInt(int64(int16(binary.LittleEndian.Uint16(intBytes))), 10)
		case 4:
			intString = strconv.FormatInt(int64(int32(binary.LittleEndian.Uint32(intBytes))), 10)
		case 8:
			intString = strconv.FormatInt(int64(int64(binary.LittleEndian.Uint64(intBytes))), 10)
		}
		d.event.Sadd(key, []byte(intString))
	}
	d.event.EndSet(key)
	return nil
}

func (d *decode) checkHeader() error {
	header := make([]byte, 9)
	_, err := io.ReadFull(d.r, header)
	if err != nil {
		return errors.Trace(err)
	}

	if !bytes.Equal(header[:5], []byte("REDIS")) {
		return fmt.Errorf("rdb: invalid file format")
	}

	version, _ := strconv.ParseInt(string(header[5:]), 10, 64)
	if version < 1 || version > rdbVersion {
		return fmt.Errorf("rdb: invalid RDB version number %d", version)
	}

	d.rdbVersion = int(version)
	return nil
}

func (d *decode) readString() ([]byte, error) {
	length, encoded, err := d.readLength()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if encoded {
		switch length {
		case rdbEncInt8:
			i, err := d.readUint8()
			return []byte(strconv.FormatInt(int64(int8(i)), 10)), errors.Trace(err)
		case rdbEncInt16:
			i, err := d.readUint16()
			return []byte(strconv.FormatInt(int64(int16(i)), 10)), errors.Trace(err)
		case rdbEncInt32:
			i, err := d.readUint32()
			return []byte(strconv.FormatInt(int64(int32(i)), 10)), errors.Trace(err)
		case rdbEncLZF:
			clen, _, err := d.readLength()
			if err != nil {
				return nil, errors.Trace(err)
			}
			ulen, _, err := d.readLength()
			if err != nil {
				return nil, errors.Trace(err)
			}
			compressed := make([]byte, clen)
			_, err = io.ReadFull(d.r, compressed)
			if err != nil {
				return nil, errors.Trace(err)
			}
			decompressed := lzfDecompress(compressed, int(ulen))
			if len(decompressed) != int(ulen) {
				return nil, fmt.Errorf("decompressed string length %d didn't match expected length %d", len(decompressed), ulen)
			}
			return decompressed, nil
		default:
			return nil, errors.Errorf("Unknown RDB string encoding type %d", length)
		}
	}

	if length == rdbLenErr {
		return nil, nil
	}

	str := make([]byte, length)
	_, err = io.ReadFull(d.r, str)
	if err != nil {
		return str, errors.Wrap(err, errors.New("readfailed"))
	}
	return str, nil
}

func (d *decode) readUint8() (uint8, error) {
	b, err := d.r.ReadByte()
	if err != nil {
		return uint8(b), errors.Wrap(err, errors.New("readfailed"))
	}
	return uint8(b), nil
}

func (d *decode) readUint16() (uint16, error) {
	_, err := io.ReadFull(d.r, d.intBuf[:2])
	if err != nil {
		return 0, errors.Wrap(err, errors.New("readfailed"))
	}
	return binary.LittleEndian.Uint16(d.intBuf), nil
}

func (d *decode) readUint32() (uint32, error) {
	_, err := io.ReadFull(d.r, d.intBuf[:4])
	if err != nil {
		return 0, errors.Wrap(err, errors.New("readfailed"))
	}
	return binary.LittleEndian.Uint32(d.intBuf), nil
}

func (d *decode) readUint64() (uint64, error) {
	_, err := io.ReadFull(d.r, d.intBuf)
	if err != nil {
		return 0, errors.Wrap(err, errors.New("readfailed"))
	}
	return binary.LittleEndian.Uint64(d.intBuf), nil
}

func (d *decode) readUint32Big() (uint32, error) {
	_, err := io.ReadFull(d.r, d.intBuf[:4])
	if err != nil {
		return 0, errors.Wrap(err, errors.New("readfailed"))
	}
	return binary.BigEndian.Uint32(d.intBuf), nil
}

func (d *decode) readUint64Big() (uint64, error) {
	_, err := io.ReadFull(d.r, d.intBuf)
	if err != nil {
		return 0, errors.Wrap(err, errors.New("readfailed"))
	}
	return binary.BigEndian.Uint64(d.intBuf), nil
}

func (d *decode) readBinaryFloat64() (float64, error) {
	floatBytes := make([]byte, 8)
	_, err := io.ReadFull(d.r, floatBytes)
	if err != nil {
		return 0, errors.Wrap(err, errors.New("readfailed"))
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(floatBytes)), nil
}

// Doubles are saved as strings prefixed by an unsigned
// 8 bit integer specifying the length of the representation.
// This 8 bit integer has special values in order to specify the following
// conditions:
// 253: not a number
// 254: + inf
// 255: - inf
func (d *decode) readFloat64() (float64, error) {
	length, err := d.readUint8()
	if err != nil {
		return 0, err
	}
	switch length {
	case 253:
		return math.NaN(), nil
	case 254:
		return math.Inf(0), nil
	case 255:
		return math.Inf(-1), nil
	default:
		floatBytes := make([]byte, length)
		_, err := io.ReadFull(d.r, floatBytes)
		if err != nil {
			return 0, err
		}
		f, err := strconv.ParseFloat(string(floatBytes), 64)
		return f, err
	}
}

func (d *decode) readLength() (uint64, bool, error) {
	b, err := d.r.ReadByte()
	if err != nil {
		return 0, false, errors.Wrap(err, errors.New("readfailed"))
	}
	// The first two bits of the first byte are used to indicate the length encoding type
	typ := (b & 0xc0) >> 6
	switch typ {
	case rdb6bitLen:
		// When the first two bits are 00, the next 6 bits are the length.
		return uint64(b & 0x3f), false, nil
	case rdb14bitLen:
		// When the first two bits are 01, the next 14 bits are the length.
		bb, err := d.r.ReadByte()
		if err != nil {
			return 0, false, errors.Wrap(err, errors.New("readfailed"))
		}
		return (uint64(b&0x3f) << 8) | uint64(bb), false, nil

	case rdbEncVal:
		// When the first two bits are 11, the next object is encoded.
		// The next 6 bits indicate the encoding type.
		return uint64(b & 0x3f), true, nil

	default:
		switch b {
		case rdb32bitLen:
			bb, err := d.readUint32Big()
			if err != nil {
				return 0, false, err
			}
			return uint64(bb), false, nil
		case rdb64bitLen:
			bb, err := d.readUint64Big()
			if err != nil {
				return 0, false, err
			}
			return bb, false, nil
		default:
			return 0, false, errors.Errorf("Unknown length encoding %d in rdbLoadLen()", b)
		}
		// When the first two bits are 10, the next 6 bits are discarded.
		// The next 4 bytes are the length.
		// length, err := d.readUint32Big()
		// return uint64(length), false, err
	}

}

func verifyDump(d []byte) error {
	if len(d) < 10 {
		return fmt.Errorf("rdb: invalid dump length")
	}
	version := binary.LittleEndian.Uint16(d[len(d)-10:])
	if version > uint16(rdbVersion) {
		return fmt.Errorf("rdb: invalid version %d, expecting %d", version, rdbVersion)
	}

	if binary.LittleEndian.Uint64(d[len(d)-8:]) != crc64.Digest(d[:len(d)-8]) {
		return fmt.Errorf("rdb: invalid CRC checksum")
	}

	return nil
}

func lzfDecompress(in []byte, outlen int) []byte {
	out := make([]byte, outlen)
	for i, o := 0, 0; i < len(in); {
		ctrl := int(in[i])
		i++
		if ctrl < 32 {
			for x := 0; x <= ctrl; x++ {
				out[o] = in[i]
				i++
				o++
			}
		} else {
			length := ctrl >> 5
			if length == 7 {
				length = length + int(in[i])
				i++
			}
			ref := o - ((ctrl & 0x1f) << 8) - int(in[i]) - 1
			i++
			for x := 0; x <= length+1; x++ {
				out[o] = out[ref]
				ref++
				o++
			}
		}
	}
	return out
}
