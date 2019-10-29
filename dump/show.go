package dump

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/cli"
	"github.com/xueqiu/rdr/decoder"
	"github.com/xueqiu/rdr/static"
)

var counters = NewSafeMap()

func listPathFiles(pathname string) []string {
	var filenames []string
	fi, err := os.Lstat(pathname) // For read access.
	if err != nil {
		return filenames
	}
	if fi.IsDir() {
		files, err := ioutil.ReadDir(pathname)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			name := path.Join(pathname, f.Name())
			filenames = append(filenames, name)
		}
	} else {
		filenames = append(filenames, pathname)
	}
	return filenames
}

// Show parse rdbfile(s) and show statistical information by html
func Show(c *cli.Context) {
	if c.NArg() < 1 {
		fmt.Fprintln(c.App.ErrWriter, "show requires at least 1 argument")
		cli.ShowCommandHelp(c, "show")
		return
	}

	// parse rdbfile
	fmt.Fprintln(c.App.Writer, "start parsing...")
	instances := []string{}
	InitHTMLTmpl()
	go func() {
		for {
			for _, pathname := range c.Args() {
				for _, v := range listPathFiles(pathname) {
					filename := filepath.Base(v)

					if !counters.Check(filename) {
						decoder := decoder.NewDecoder()
						fmt.Fprintf(c.App.Writer, "start to parse %v \n", filename)
						go Decode(c, decoder, v)
						counter := NewCounter()
						counter.Count(decoder.Entries)
						counters.Set(filename, counter)
						fmt.Fprintf(c.App.Writer, "parse %v  done\n", filename)

						instances = append(instances, filename)
						// init html template
						// init common data in template
						tplCommonData["Instances"] = instances
					}
				}

			}
			time.Sleep(5 * time.Second)
		}
	}()

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
