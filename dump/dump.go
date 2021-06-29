package dump

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dongmx/rdb"
	"github.com/urfave/cli"
	"github.com/xueqiu/rdr/decoder"
)

// Dump rdb file statistical information
func Dump(path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	decoder := decoder.NewDecoder()
	go func() {
		f, err := os.Open(path)
		defer close(decoder.Entries)
		if err != nil {
			fmt.Printf("open rdbfile err: %v\n", err)
			return
		}
		err = rdb.Decode(f, decoder)
		if err != nil {
			fmt.Printf("decode rdbfile err: %v\n", err)
			return
		}
	}()
	cnt := NewCounter()
	cnt.Count(decoder.Entries)
	filename := filepath.Base(path)
	data = GetData(filename, cnt, 100, 0)
	return data, nil
}

// ToCliWriter dump rdb file statistical information to STDOUT.
func ToCliWriter(cli *cli.Context, topN, sizeThreshold int) {
	ToWriter(cli.Args(), cli.App.Writer, topN, sizeThreshold)
}

// ToWriter .
func ToWriter(files []string, writer io.Writer, topN, sizeThreshold int) {
	if len(files) < 1 {
		fmt.Fprintln(writer, " requires at least 1 argument")
		return
	}

	// parse rdbfile
	fmt.Fprintln(writer, "[")
	for i, file := range files {
		decoder := decoder.NewDecoder()
		go Decode(decoder, file, writer)
		cnt := NewCounter()
		cnt.Count(decoder.Entries)
		filename := filepath.Base(file)
		data := GetData(filename, cnt, topN, sizeThreshold)
		data["MemoryUse"] = decoder.GetUsedMem()
		data["CTime"] = decoder.GetTimestamp()
		jsonBytes, _ := json.MarshalIndent(data, "", "    ")
		fmt.Fprint(writer, string(jsonBytes))
		if i == len(files)-1 {
			fmt.Fprintln(writer)
		} else {
			fmt.Fprintln(writer, ",")
		}
	}
	fmt.Fprintln(writer, "]")
}

// Decode ...
func Decode(decoder *decoder.Decoder, filepath string, writer io.Writer) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintf(writer, "open rdbfile err: %v\n", err)
		close(decoder.Entries)
		return
	}
	err = rdb.Decode(f, decoder)
	if err != nil {
		fmt.Fprintf(writer, "decode rdbfile err: %v\n", err)
		close(decoder.Entries)
		return
	}
}

func GetData(filename string, cnt *Counter, topN, sizeThreshold int) map[string]interface{} {
	data := make(map[string]interface{})
	data["CurrentInstance"] = filename
	data["LargestKeys"] = cnt.GetLargestEntries(topN, sizeThreshold)

	largestKeyPrefixesByType := map[string][]*PrefixEntry{}
	for _, entry := range cnt.GetLargestKeyPrefixes() {
		// if mem usage is less than 1M, and the list is long enough, then it's unnecessary to add it.
		if entry.Bytes < 1000*1000 && len(largestKeyPrefixesByType[entry.Type]) > 50 {
			continue
		}
		largestKeyPrefixesByType[entry.Type] = append(largestKeyPrefixesByType[entry.Type], entry)
	}
	data["LargestKeyPrefixes"] = largestKeyPrefixesByType

	data["TypeBytes"] = cnt.typeBytes
	data["TypeNum"] = cnt.typeNum
	totalNum := uint64(0)
	for _, v := range cnt.typeNum {
		totalNum += v
	}
	totalBytes := uint64(0)
	for _, v := range cnt.typeBytes {
		totalBytes += v
	}
	data["TotleNum"] = totalNum
	data["TotleBytes"] = totalBytes

	lenLevelCount := map[string][]*PrefixEntry{}
	for _, entry := range cnt.GetLenLevelCount() {
		lenLevelCount[entry.Type] = append(lenLevelCount[entry.Type], entry)
	}
	data["LenLevelCount"] = lenLevelCount

	var slotBytesHeap slotHeap
	for slot, length := range cnt.slotBytes {
		heap.Push(&slotBytesHeap, &SlotEntry{
			Slot: slot, Size: length,
		})
	}

	var slotSizeHeap slotHeap
	for slot, size := range cnt.slotNum {
		heap.Push(&slotSizeHeap, &SlotEntry{
			Slot: slot, Size: size,
		})
	}

	slotBytes := make(slotHeap, 0, topN)
	slotNums := make(slotHeap, 0, topN)

	for i := 0; i < topN; i++ {
		continueFlag := false
		if slotBytesHeap.Len() > 0 {
			continueFlag = true
			slotBytes = append(slotBytes, heap.Pop(&slotBytesHeap).(*SlotEntry))
		}
		if slotSizeHeap.Len() > 0 {
			continueFlag = true
			slotNums = append(slotNums, heap.Pop(&slotSizeHeap).(*SlotEntry))
		}

		if !continueFlag {
			break
		}
	}

	data["SlotBytes"] = slotBytes
	data["SlotNums"] = slotNums

	return data
}
