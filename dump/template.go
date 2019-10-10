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
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/xueqiu/rdr/views"
)

var (
	tmpl          *template.Template
	tplCommonData = map[string]interface{}{}
	tplFuncMap    = make(template.FuncMap)
	isFirst       = true
)

func InitHTMLTmpl() {
	// init function maps
	tplFuncMap["isFirst"] = func() bool { res := isFirst; isFirst = false; return res }
	tplFuncMap["clearFirst"] = func() bool { isFirst = true; return isFirst }
	tplFuncMap["hash"] = func(str string) string { return fmt.Sprintf("%x", md5.Sum([]byte(str))) }
	tplFuncMap["humanizeBytes"] = humanize.Bytes
	tplFuncMap["humanizeComma"] = func(i uint64) string { return humanize.Comma(int64(i)) }

	// init views html template
	for _, name := range views.AssetNames() {
		// just support *.html
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		content, tmplErr := views.Asset(name)
		if tmplErr != nil {
			log.Printf("|ERROR|asset %v err %v", name, tmplErr)
			continue
		}

		if tmpl == nil {
			tmpl, tmplErr = template.New(name).Funcs(tplFuncMap).Parse(string(content))
		} else {
			tmpl, tmplErr = tmpl.New(name).Funcs(tplFuncMap).Parse(string(content))
		}

		if tmplErr != nil {
			log.Printf("|ERROR|parse template err %v", tmplErr)
		}
	}
}

// ServeHTML generate and write html to client
func ServeHTML(w http.ResponseWriter, layout string, content string, data map[string]interface{}) {
	var buf bytes.Buffer
	bodyTmplErr := tmpl.ExecuteTemplate(&buf, content, data)
	if bodyTmplErr != nil {
		log.Printf("|ERROR|ServeHTML bodyTmplErr ERROR %v", bodyTmplErr)
	}
	bodyHTML := template.HTML(buf.String())
	if len(data) == 0 {
		data = map[string]interface{}{}
	}
	data["LayoutContent"] = bodyHTML
	tmplErr := tmpl.ExecuteTemplate(w,
		layout,
		data)

	if tmplErr != nil {
		log.Printf("|ERROR|ServeHTML LayoutTmplErr ERROR %v", tmplErr)
	}
}
