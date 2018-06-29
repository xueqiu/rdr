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

package main

import "github.com/dongmx/rdb/nopdecoder"
import "strconv"

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
	Entries  chan *Entry
	m        MemProfiler
	tmpStore map[string]*Entry
	nopdecoder.NopDecoder
}

// NewDecoder new a rdb decoder
func NewDecoder() *Decoder {
	return &Decoder{
		Entries:  make(chan *Entry, 1024),
		m:        MemProfiler{},
		tmpStore: map[string]*Entry{},
	}
}

// Set is called once for each string key.
func (d *Decoder) Set(key, value []byte, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += d.m.SizeofString(value)
	bytes += d.m.TopLevelObjOverhead()
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
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
func (d *Decoder) StartHash(key []byte, length, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.TopLevelObjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	bytes += d.m.HashtableOverhead(uint64(length))
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "hash",
		NumOfElem: uint64(length),
	}
	d.tmpStore[keyStr] = e
}

// Hset is called once for each field=value pair in a hash.
func (d *Decoder) Hset(key, field, value []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	lenOfElem := d.m.ElemLen(field) + d.m.ElemLen(value)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(field)
		e.LenOfLargestElem = lenOfElem
	}
	e.Bytes += d.m.SizeofString(field)
	e.Bytes += d.m.SizeofString(value)
	e.Bytes += d.m.HashtableEntryOverhead()
	e.Bytes += 2 * d.m.RobjOverhead()
	d.tmpStore[keyStr] = e
}

// EndHash is called when there are no more fields in a hash.
func (d *Decoder) EndHash(key []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	d.Entries <- e
	delete(d.tmpStore, keyStr)
}

// StartSet is called at the beginning of a set.
// Sadd will be called exactly cardinality times before EndSet.
func (d *Decoder) StartSet(key []byte, cardinality, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.TopLevelObjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	bytes += d.m.HashtableOverhead(uint64(cardinality))
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "set",
		NumOfElem: uint64(cardinality),
	}
	d.tmpStore[keyStr] = e
}

// Sadd is called once for each member of a set.
func (d *Decoder) Sadd(key, member []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	lenOfElem := d.m.ElemLen(member)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(member)
		e.LenOfLargestElem = lenOfElem
	}
	e.Bytes += d.m.SizeofString(member)
	e.Bytes += d.m.HashtableEntryOverhead()
	e.Bytes += d.m.RobjOverhead()
	d.tmpStore[keyStr] = e
}

// EndSet is called when there are no more fields in a set.
// Same as EndHash
func (d *Decoder) EndSet(key []byte) {
	d.EndHash(key)
}

// StartList is called at the beginning of a list.
// Rpush will be called exactly length times before EndList.
// If length of the list is not known, then length is -1
func (d *Decoder) StartList(key []byte, length, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.TopLevelObjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	bytes += d.m.LinkedListEntryOverhead() * uint64(length)
	bytes += d.m.LinkedlistOverhead()
	bytes += d.m.RobjOverhead() * uint64(length)
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "list",
		NumOfElem: uint64(length),
	}
	d.tmpStore[keyStr] = e
}

// Rpush is called once for each value in a list.
func (d *Decoder) Rpush(key, value []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	lenOfElem := d.m.ElemLen(value)
	if _, err := strconv.ParseInt(string(value), 10, 32); err == nil {
		e.Bytes += 4
	} else {
		e.Bytes += d.m.SizeofString(value)
	}
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(value)
		e.LenOfLargestElem = lenOfElem
	}
	d.tmpStore[keyStr] = e
}

// EndList is called when there are no more values in a list.
func (d *Decoder) EndList(key []byte) {
	d.EndHash(key)
}

// StartZSet is called at the beginning of a sorted set.
// Zadd will be called exactly cardinality times before EndZSet.
func (d *Decoder) StartZSet(key []byte, cardinality, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.TopLevelObjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	bytes += d.m.SkiplistOverhead(uint64(cardinality))
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "sortedset",
		NumOfElem: uint64(cardinality),
	}
	d.tmpStore[keyStr] = e
}

// Zadd is called once for each member of a sorted set.
func (d *Decoder) Zadd(key []byte, score float64, member []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	lenOfElem := d.m.ElemLen(member)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(member)
		e.LenOfLargestElem = lenOfElem
	}
	e.Bytes += 8 // sizeof(score)
	e.Bytes += d.m.SizeofString(member)
	e.Bytes += 2 * d.m.RobjOverhead()
	e.Bytes += d.m.SkiplistEntryOverhead()
	d.tmpStore[keyStr] = e
}

// EndZSet is called when there are no more members in a sorted set.
func (d *Decoder) EndZSet(key []byte) {
	d.EndHash(key)
}

// EndRDB is called when parsing of the RDB file is complete.
func (d *Decoder) EndRDB() {
	close(d.Entries)
}
