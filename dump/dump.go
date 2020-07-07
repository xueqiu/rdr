package dump

import (
	"encoding/json"
	"fmt"
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
	data = getData(filename, cnt)
	return data, nil
}

// ToCliWriter dump rdb file statistical information to STDOUT.
func ToCliWriter(cli *cli.Context) {
	if cli.NArg() < 1 {
		fmt.Fprintln(cli.App.ErrWriter, " requires at least 1 argument")
		return
	}

	// parse rdbfile
	fmt.Fprintln(cli.App.Writer, "[")
	nargs := cli.NArg()
	for i := 0; i < nargs; i++ {
		file := cli.Args().Get(i)
		decoder := decoder.NewDecoder()
		go Decode(cli, decoder, file)
		cnt := NewCounter()
		cnt.Count(decoder.Entries)
		filename := filepath.Base(file)
		data := getData(filename, cnt)
		data["MemoryUse"] = decoder.GetUsedMem()
		data["CTime"] = decoder.GetTimestamp()
		jsonBytes, _ := json.MarshalIndent(data, "", "    ")
		fmt.Fprint(cli.App.Writer, string(jsonBytes))
		if i == nargs-1 {
			fmt.Fprintln(cli.App.Writer)
		} else {
			fmt.Fprintln(cli.App.Writer, ",")
		}
	}
	fmt.Fprintln(cli.App.Writer, "]")
}

// Decode ...
func Decode(c *cli.Context, decoder *decoder.Decoder, filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintf(c.App.ErrWriter, "open rdbfile err: %v\n", err)
		close(decoder.Entries)
		return
	}
	err = rdb.Decode(f, decoder)
	if err != nil {
		fmt.Fprintf(c.App.ErrWriter, "decode rdbfile err: %v\n", err)
		close(decoder.Entries)
		return
	}
}

func getData(filename string, cnt *Counter) map[string]interface{} {
	data := make(map[string]interface{})
	data["CurrentInstance"] = filename
	data["LargestKeys"] = cnt.GetLargestEntries(100)

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
	return data
}
