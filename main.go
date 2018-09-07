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

import (
	"net/http"
	"os"

	"github.com/urfave/cli"

	"fmt"

	"path/filepath"

	"github.com/dongmx/rdb"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/julienschmidt/httprouter"
	"github.com/xueqiu/rdr/static"
	"encoding/json"
)

//go:generate go-bindata -prefix "static/" -o=static/static.go -pkg=static -ignore static.go static/...
//go:generate go-bindata -prefix "views/" -o=views/views.go -pkg=views -ignore views.go views/...

func decode(c *cli.Context, decoder *Decoder, filepath string) {
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

var counters = NewSafeMap()

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

// dump rdb file statistical information to STDOUT.
func dump(cli *cli.Context) {
	if cli.NArg() < 1 {
		fmt.Fprintln(cli.App.ErrWriter, " requires at least 1 argument")
		return
	}

	// parse rdbfile
	fmt.Fprintln(cli.App.Writer, "[")
	nargs := cli.NArg()
	for i := 0; i < nargs; i++ {
		file := cli.Args().Get(i)
		decoder := NewDecoder()
		go decode(cli, decoder, file)
		cnt := NewCounter()
		cnt.Count(decoder.Entries)
		filename := filepath.Base(file)
		data := getData(filename, cnt)
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

// show parse rdbfile(s) and show statistical information by html
func show(c *cli.Context) {
	if c.NArg() < 1 {
		fmt.Fprintln(c.App.ErrWriter, "show requires at least 1 argument")
		cli.ShowCommandHelp(c, "show")
		return
	}

	// parse rdbfile
	fmt.Fprintln(c.App.Writer, "start parsing...")
	instances := []string{}
	for _, file := range c.Args() {
		decoder := NewDecoder()
		go decode(c, decoder, file)
		counter := NewCounter()
		counter.Count(decoder.Entries)
		filename := filepath.Base(file)
		counters.Set(filename, counter)
		instances = append(instances, filename)
		fmt.Fprintf(c.App.Writer, "parse %v  done\n", file)
	}

	// init html template
	// init common data in template
	initHTMLTmpl()
	tplCommonData["Instances"] = instances

	// start http server
	staticFS := assetfs.AssetFS{
		Asset:     static.Asset,
		AssetDir:  static.AssetDir,
		AssetInfo: static.AssetInfo,
	}
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", &staticFS)
	router.GET("/", index)
	router.GET("/instance/:path", rdbReveal)
	fmt.Fprintln(c.App.Writer, "parsing finished, please access http://{$IP}:"+c.String("port"))
	listenErr := http.ListenAndServe(":"+c.String("port"), router)
	if listenErr != nil {
		fmt.Fprintf(c.App.ErrWriter, "Listen port err: %v\n", listenErr)
	}
}

// keys is function for command `keys`
// output all keys in rdbfile(s) get from args
func keys(c *cli.Context) {
	if c.NArg() < 1 {
		fmt.Fprintln(c.App.ErrWriter, "keys requires at least 1 argument")
		cli.ShowCommandHelp(c, "keys")
		return
	}
	for _, filepath := range c.Args() {
		decoder := NewDecoder()
		go decode(c, decoder, filepath)
		for e := range decoder.Entries {
			fmt.Fprintf(c.App.Writer, "%v\n", e.Key)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "rdr"
	app.Usage = "a tool to parse redis rdbfile"
	app.Version = "v0.0.1"
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	app.Commands = []cli.Command{
		cli.Command{
		    Name: "dump",
		    Usage: "dump statistical information of rdbfile to STDOUT",
		    ArgsUsage: "FILE1 [FILE2] [FILE3]...",
		    Action: dump,
		},
		cli.Command{
			Name:      "show",
			Usage:     "show statistical information of rdbfile by webpage",
			ArgsUsage: "FILE1 [FILE2] [FILE3]...",
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "port, p",
					Value: 8080,
					Usage: "Port for rdr to listen",
				},
			},
			Action: show,
		},
		cli.Command{
			Name:      "keys",
			Usage:     "get all keys from rdbfile",
			ArgsUsage: "FILE1 [FILE2] [FILE3]...",
			Action:    keys,
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.ErrWriter, "command %q can not be found.\n", command)
		cli.ShowAppHelp(c)
	}
	app.Run(os.Args)
}

