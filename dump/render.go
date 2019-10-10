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

package dump

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	for key := range counters.Items() {
		http.Redirect(w, r, "/instance/"+key.(string), http.StatusFound)
		return
	}
}

func rdbReveal(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// deep copy  tplCommonData into data
	data := map[string]interface{}{}
	for key, val := range tplCommonData {
		data[key] = val
	}

	path := p.ByName("path")

	c := counters.Get(path)
	if c == nil {
		return
	}
	counter := c.(*Counter)

	data["CurrentInstance"] = path
	data["LargestKeys"] = counter.GetLargestEntries(100)

	largestKeyPrefixesByType := map[string][]*PrefixEntry{}
	for _, entry := range counter.GetLargestKeyPrefixes() {
		// mem use less than 1M, and list is long enough, not necessary to add
		if entry.Bytes < 1000*1000 && len(largestKeyPrefixesByType[entry.Type]) > 50 {
			continue
		}
		largestKeyPrefixesByType[entry.Type] = append(largestKeyPrefixesByType[entry.Type], entry)
	}
	data["LargestKeyPrefixes"] = largestKeyPrefixesByType

	data["TypeBytes"] = counter.typeBytes
	data["TypeNum"] = counter.typeNum
	totleNum := uint64(0)
	for _, v := range counter.typeNum {
		totleNum += v
	}
	totleBytes := uint64(0)
	for _, v := range counter.typeBytes {
		totleBytes += v
	}
	data["TotleNum"] = totleNum
	data["TotleBytes"] = totleBytes

	lenLevelCount := map[string][]*PrefixEntry{}
	for _, entry := range counter.GetLenLevelCount() {
		lenLevelCount[entry.Type] = append(lenLevelCount[entry.Type], entry)
	}
	data["LenLevelCount"] = lenLevelCount
	ServeHTML(w, "base.html", "revel.html", data)
}
