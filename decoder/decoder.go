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
	"fmt"
	"os"
	"strconv"

	"github.com/dongmx/rdb"
	"github.com/dongmx/rdb/nopdecoder"
)

// Entry is info of a redis recored
type Entry struct {
	Key                string
	Bytes              uint64
	Type               string
	NumOfElem          uint64
	LenOfLargestElem   uint64
	FieldOfLargestElem string
}

// Decoder decode rdb file
type Decoder struct {
	Entries chan *Entry
	m       MemProfiler

	usedMem int64
	ctime   int64
	count   int
	rdbVer int

	currentInfo  *rdb.Info
	currentEntry *Entry

	nopdecoder.NopDecoder
}

// NewDecoder new a rdb decoder
func NewDecoder() *Decoder {
	return &Decoder{
		Entries: make(chan *Entry, 1024),
		m:       MemProfiler{},
	}
}

func (d *Decoder) sendEntry() {
	d.Entries <- d.currentEntry
	d.currentEntry = nil
}

func (d *Decoder) GetTimestamp() int64 {
	return d.ctime
}

func (d *Decoder) GetUsedMem() int64 {
	return d.usedMem
}

func (d *Decoder) StartRDB(ver int) {
	d.rdbVer = ver
}

func (d *Decoder) Aux(key, value []byte) {
	switch string(key) {
	case "ctime":
		{
			n, err := strconv.ParseInt(string(value), 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ParseInt(", string(value), "):", err)
			}
			d.ctime = n
		}

	case "used-mem":
		{
			n, err := strconv.ParseInt(string(value), 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ParseInt(", string(value), "):", err)
			}
			d.usedMem = n
		}

	}
}

func (d *Decoder) StartStream(key []byte, cardinality, expiry int64, info *rdb.Info) {
	keyStr := string(key)
	bytes := d.m.TopLevelObjOverhead(key, expiry)
	bytes += d.m.StreamOverhead()
	bytes += d.m.SizeofStreamRadixTree(uint64(cardinality))

	d.currentInfo = info
	d.currentEntry = &Entry{
		Key:              keyStr,
		Bytes:            bytes,
		Type:             "stream",
		NumOfElem:        0,
		LenOfLargestElem: 0,
	}
}

func (d *Decoder) Xadd(key, id, listpack []byte) {
	e := d.currentEntry
	e.Bytes += d.m.mallocOverhead(uint64(len(listpack)))
}

func (d *Decoder) EndStream(key []byte, items uint64, lastEntryID string, cgroupsData rdb.StreamGroups) {
	e := d.currentEntry

	for _, cg := range cgroupsData {
		pendingLength := uint64(len(cg.Pending))
		e.Bytes += d.m.SizeofStreamRadixTree(pendingLength)
		e.Bytes += d.m.StreamNACK(pendingLength)

		for _, c := range cg.Consumers {
			e.Bytes += d.m.StreamConsumer(c.Name)
			pendingLength := uint64(len(cg.Pending))
			e.Bytes += d.m.SizeofStreamRadixTree(pendingLength)
		}
	}

	d.sendEntry()
}

// Set is called once for each string key.
func (d *Decoder) Set(key, value []byte, expiry int64, info *rdb.Info) {
	keyStr := string(key)
	bytes := d.m.TopLevelObjOverhead(key, expiry)
	bytes += d.m.SizeofString(value)

	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "string",
		NumOfElem: d.m.ElemLen(value),
	}
	d.Entries <- e
}

// StartHash is called at the beginning of a hash.
// Hset will be called exactly length times before EndHash.
func (d *Decoder) StartHash(key []byte, length, expiry int64, info *rdb.Info) {
	keyStr := string(key)

	bytes := d.m.TopLevelObjOverhead(key, expiry)

	if info.SizeOfValue > 0 {
		bytes += uint64(info.SizeOfValue)
	} else if info.Encoding == "hashtable" {
		bytes += d.m.HashtableOverhead(uint64(length))
	} else {
		panic(fmt.Sprintf("unexpected size(0) or encoding:%s", info.Encoding))
	}

	d.currentInfo = info
	d.currentEntry = &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "hash",
		NumOfElem: uint64(length),
	}
}

// Hset is called once for each field=value pair in a hash.
func (d *Decoder) Hset(key, field, value []byte) {
	e := d.currentEntry

	lenOfElem := d.m.ElemLen(field) + d.m.ElemLen(value)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(field)
		e.LenOfLargestElem = lenOfElem
	}

	if d.currentInfo.Encoding == "hashtable" {
		e.Bytes += d.m.SizeofString(field)
		e.Bytes += d.m.SizeofString(value)
		e.Bytes += d.m.HashtableEntryOverhead()

		if d.rdbVer < 8 {
			e.Bytes += 2 * d.m.RobjOverhead()
		}
	}
}

// EndHash is called when there are no more fields in a hash.
func (d *Decoder) EndHash(key []byte) {
	d.sendEntry()
}

// StartSet is called at the beginning of a set.
// Sadd will be called exactly cardinality times before EndSet.
func (d *Decoder) StartSet(key []byte, cardinality, expiry int64, info *rdb.Info) {
	d.StartHash(key, cardinality, expiry, info)
}

// Sadd is called once for each member of a set.
func (d *Decoder) Sadd(key, member []byte) {
	e := d.currentEntry
	lenOfElem := d.m.ElemLen(member)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(member)
		e.LenOfLargestElem = lenOfElem
	}

	if d.currentInfo.Encoding == "hashtable" {
		e.Bytes += d.m.SizeofString(member)
		e.Bytes += d.m.HashtableEntryOverhead()

		if d.rdbVer < 8 {
			e.Bytes += d.m.RobjOverhead()
		}
	}
}

// EndSet is called when there are no more fields in a set.
// Same as EndHash
func (d *Decoder) EndSet(key []byte) {
	d.sendEntry()
}

// StartList is called at the beginning of a list.
// Rpush will be called exactly length times before EndList.
// If length of the list is not known, then length is -1
func (d *Decoder) StartList(key []byte, length, expiry int64, info *rdb.Info) {
	keyStr := string(key)

	d.currentInfo = info
	bytes := d.m.TopLevelObjOverhead(key, expiry)

	//bug here length would be -1 if it is quicklist
	//bytes += d.m.RobjOverhead() * uint64(length)
	d.currentEntry = &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "list",
		NumOfElem: 0,
	}
}

// Rpush is called once for each value in a list.
func (d *Decoder) Rpush(key, value []byte) {
	//keyStr := string(key)
	e := d.currentEntry
	e.NumOfElem++

	switch d.currentInfo.Encoding {
	case "quicklist":
		e.Bytes += d.m.ZiplistEntryOverhead(value)

	case "ziplist":
		e.Bytes += d.m.ZiplistEntryOverhead(value)

	case "linkedlist":
		sizeInlist := uint64(0)
		if _, err := strconv.ParseInt(string(value), 10, 32); err != nil {
			sizeInlist = d.m.SizeofString(value)
		}

		e.Bytes += d.m.LinkedListEntryOverhead()
		e.Bytes += sizeInlist

		if d.rdbVer < 8 {
			e.Bytes += d.m.RobjOverhead()
		}

	default:
		panic(fmt.Sprintf("unknown encoding:%s", d.currentInfo.Encoding))
	}

	lenOfElem := d.m.ElemLen(value)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(value)
		e.LenOfLargestElem = lenOfElem
	}
}

// EndList is called when there are no more values in a list.
func (d *Decoder) EndList(key []byte) {
	e := d.currentEntry

	switch d.currentInfo.Encoding {
	case "quicklist":
		e.Bytes += d.m.QuicklistOverhead(d.currentInfo.Zips)
		e.Bytes += d.m.ZiplistHeaderOverhead() * d.currentInfo.Zips

	case "ziplist":
		e.Bytes += d.m.ZiplistHeaderOverhead()

	case "linkedlist":
		e.Bytes += d.m.LinkedlistOverhead()

	default:
		panic(fmt.Sprintf("unknown encoding:%s", d.currentInfo.Encoding))
	}

	d.sendEntry()
}

// StartZSet is called at the beginning of a sorted set.
// Zadd will be called exactly cardinality times before EndZSet.
func (d *Decoder) StartZSet(key []byte, cardinality, expiry int64, info *rdb.Info) {
	keyStr := string(key)

	bytes := d.m.TopLevelObjOverhead(key, expiry)
	d.currentInfo = info

	if info.SizeOfValue > 0 {
		bytes += uint64(info.SizeOfValue)
	} else if info.Encoding == "skiplist" {
		bytes += d.m.SkiplistOverhead(uint64(cardinality))
	} else {
		panic(fmt.Sprintf("unexpected size(0) or encoding:%s", info.Encoding))
	}

	d.currentEntry = &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "sortedset",
		NumOfElem: uint64(cardinality),
	}
}

// Zadd is called once for each member of a sorted set.
func (d *Decoder) Zadd(key []byte, score float64, member []byte) {
	e := d.currentEntry
	lenOfElem := d.m.ElemLen(member)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(member)
		e.LenOfLargestElem = lenOfElem
	}

	if d.currentInfo.Encoding == "skiplist" {
		e.Bytes += 8 // sizeof(score)
		e.Bytes += d.m.SizeofString(member)
		e.Bytes += d.m.SkiplistEntryOverhead()

		if d.rdbVer < 8 {
			e.Bytes += d.m.RobjOverhead()
		}
	}
}

// EndZSet is called when there are no more members in a sorted set.
func (d *Decoder) EndZSet(key []byte) {
	d.sendEntry()
}

// EndRDB is called when parsing of the RDB file is complete.
func (d *Decoder) EndRDB() {
	close(d.Entries)
}
