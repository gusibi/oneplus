package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

const TEMPLATE = `<html>
<head>
   <link rel="stylesheet" href="http://media.gusibi.mobi/highlight/static/styles/atom-one-dark.css">
   <script src="http://media.gusibi.mobi/highlight/static/highlight.site.pack.js"></script>
   <script>hljs.initHighlightingOnLoad();</script>
   <style type="text/css">
        body {
            /* background-image: url(https://www.google.com);
            margin: 0; */
            font-size: 16px;
        }
        
        .hljs {
            padding: 8px;
        }
        
        .hljs span {
            padding: 0px;
            line-height: 20px;
        }
        
        pre,
        code {
            margin: 0;
            padding-top: 0;
        }
    </style>
</head>
<body style="width: 640px;">
<pre>
<code class="{{.Language}}">{{.Code}}</code>
</pre>
</body>
</html>`

const TEMPLATE_HELLO = `
<!doctype html>
<html>
<head>
<link rel="stylesheet" href="http://media.gusibi.mobi/highlight/static/styles/atom-one-dark.css">
<script src="http://media.gusibi.mobi/highlight/static/highlight.site.pack.js"></script>
<script>hljs.initHighlightingOnLoad();</script>
</head>
<body style="width: 560px;">
<pre>
    <code class="c">
def maxSubArrayLen(nums, k):

    result, acc = 0, 0
    dic = {0: -1}

    for i in xrange(len(nums)):
        acc += nums[i]
        if acc not in dic:
            dic[acc] = i
        if acc - k in dic:
            result = max(result, i - dic[acc-k])
        print(dic)

    return result


if __name__ == "__main__":
    print([1, 2, 3, 4], 10)
    maxSubArrayLen([1, 2, 3, 4], 10)
    print([1, -2, -3, 4], -1)
    maxSubArrayLen([1, -2, -3, 4], -1)
</code>
</pre>
<img src="https://www.google.com" alt="">
</body>
</html>`

func EnvGet(env, default_value string) string {
	value := os.Getenv(env)
	if value == "" {
		if default_value != "" {
			return default_value
		}
	}
	return value
}

func MD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	md5 := fmt.Sprintf("%x", h.Sum(nil))
	return md5
}

type CodeBody struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func EchoHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(w, req.URL.Path)
}

func CreateHtml(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Panic(err)
	}
	var code CodeBody
	err = json.Unmarshal(body, &code)
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	id := MD5(code.Code)
	codeModel := Code{ID: id, Code: code.Code, Language: code.Language}
	c := codeModel.Get(id)
	js, err := json.Marshal(codeModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if c != nil {
		log.Printf("already exsit: %s\n", id)
		w.Write(js)
	} else {
		codeModel.Create()
		w.Write(js)
	}
}
func CodeRender(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Panic(err)
	}
	var code CodeBody
	err = json.Unmarshal(body, &code)
	if err != nil {
		log.Panic(err)
	}
	t, _ := template.New("hello").Parse(TEMPLATE)
	t.Execute(w, &code)
}

func RenderCode(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	codeModel := Code{}
	code := codeModel.Get(id)
	if code == nil {
		http.Error(w, "not found", http.StatusNotFound)
	}
	t, _ := template.New("hello").Parse(TEMPLATE)
	t.Execute(w, &code)
}

func Sleep(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	seconds, err := strconv.Atoi(ps.ByName("seconds"))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
	}
	if seconds > 2 {
		seconds = 2
	} else if seconds < 0 {
		seconds = 1
	}
	log.Println("sleep seconds: ", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
	io.WriteString(w, "wake up")
}
