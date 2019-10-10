package dump

import (
	"fmt"
	"net/http"
	"path/filepath"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/cli"
	"github.com/xueqiu/rdr/decoder"
	"github.com/xueqiu/rdr/static"
)

var counters = NewSafeMap()

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
	for _, file := range c.Args() {
		decoder := decoder.NewDecoder()
		go Decode(c, decoder, file)
		counter := NewCounter()
		counter.Count(decoder.Entries)
		filename := filepath.Base(file)
		counters.Set(filename, counter)
		instances = append(instances, filename)
		fmt.Fprintf(c.App.Writer, "parse %v  done\n", file)
	}

	// init html template
	// init common data in template
	InitHTMLTmpl()
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
