package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"wendal/howeyc/fsnotify"
	"wendal/gor"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
)

const (
	VER = "3.7.0"
)

var (
	http_addr   = flag.String("http", ":8080", "Http addr for Preview or Server")
	args        []string
	_compileVer = 0
	_watch_js   = `
<script type="text/javascript" src="http://lib.sinaapp.com/js/jquery/1.8.3/jquery.min.js"></script>
<script type="text/javascript">
	$(function() {
	var _gor_compile_ver = 0;
	function _sc(data) {
		if (parseInt(data) != _gor_compile_ver) {
			location.reload(true);
		} else {
			setTimeout(_gor_f5, 1000);
		}
	}
	function _gor_f5() {$.get("/_api/f5?ver="+_gor_compile_ver, "", _sc)};
	$.get("/_api/f5", "",  function(data) {_gor_compile_ver = parseInt(data);setTimeout(_gor_f5, 1000);});
	});
</script>
	`
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)//SetFlags设置标准记录器的输出标志。 标志位是Ldate、Ltime等。
	log.Println("gor ver " + VER)
}

func main() {
	flag.Parse()		//解析命令行
	args = flag.Args()		//Args返回非标志命令行参数。   return[]string
	if len(args) == 0 || len(args) > 4 {
		PrintUsage()		//输出帮助HELP
		os.Exit(1)		//Exit使当前程序以给定的状态代码退出。
								// 按照惯例，代码0表示成功，非零表示错误。 程序立即终止；未运行延迟函数。 为了便于移植，状态码应该在[0，125]范围内。
	}
	switch args[0] {
	default:
		PrintUsage()
		os.Exit(1)
	case "config":
		cnf, err := gor.ReadConfig(".")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("RuhohSpec: ", cnf["RuhohSpec"])
		buf, err := json.MarshalIndent(cnf, "", "  ") 		//MarshalIndent类似于Marshal，但应用Indent来格式化输出。
																			//输出中的每个JSON元素都将从一个以prefix开头的新行开始
																			//根据缩进嵌套，后跟一个或多个缩进副本。
		if err != nil {
			log.Fatal(err)
		}
		log.Println("global config\n", string(buf))
	case "new":
		fallthrough
	case "init":
		if len(args) == 1 {
			log.Fatalln(os.Args[0], "new", "<dir>")		//Fatalln等价于Println（）后跟操作系统退出(1).
		}
		CmdInit(args[1])
	case "posts":
		gor.ListPosts()
	case "payload":
		payload, err := gor.BuildPayload("./")
		if err != nil {
			log.Fatal(err)
		}
		buf, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(buf))
	case "compile":
		fallthrough
	case "build":
		fallthrough
	case "c":
		_compile()
	case "post":
		if len(args) == 1 {
			log.Fatal("gor post <title>")
		} else if len(args) == 2 {
			gor.CreateNewPost(args[1])
		} else {
			gor.CreateNewPostWithImgs(args[1], args[2])
		}
	case "addimg":
		if len(args) == 3 {
			gor.AddImgs(args[1], args[2], "")
		} else if len(args) == 4 {
			gor.AddImgs(args[1], args[2], args[3])
		} else {
			log.Fatal("gor post <title> <dir> or <date>")
		}
	case "http":
		_http()
	case "preview":
		gor.HTML_EXT = _watch_js
		_compile()
		go _http()
		_watch()
	case "pprof":
		f, _ := os.OpenFile("gor.pprof", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		for i := 0; i < 100; i++ {
			_compile()
		}
	case ".update.zip.go":
		d, _ := ioutil.ReadFile("gor-content.zip")
		_zip, _ := os.OpenFile("zip.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		header := `package main
const INIT_ZIP="`
		_zip.Write([]byte(header))
		encoder := base64.NewEncoder(base64.StdEncoding, _zip)
		encoder.Write(d)
		encoder.Close()
		_zip.Write([]byte(`"`))
		_zip.Sync()
		_zip.Close()
	}
}

func _http() {
	log.Println("Listen at ", *http_addr)
	sm := http.NewServeMux()
	sm.HandleFunc("/_api/f5", f5)
	sm.Handle("/", http.FileServer(http.Dir("compiled")))
	log.Println(http.ListenAndServe(*http_addr, sm))
}

func _compile() {
	err := gor.Compile()
	if err != nil {
		log.Fatal(err)
	}
}

func _watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("fsnotify fail, on-fly watch is disable")
		return
	}
	done := make(chan bool)
	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() {
					_compile()
					_compileVer += 1
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
	path := "posts"
	if len(args) == 3 {
		path = args[2]
	}
	log.Println("Start watching on ", path)
	err = watcher.Watch(path)
	if err != nil {
		log.Fatal(err)
	}

	<-done

	/* ... do stuff ... */
	watcher.Close()
}

// -----------------------------
// HTTP APIs

func f5(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%d", _compileVer)))
	return
}
