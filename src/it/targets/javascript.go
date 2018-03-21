package targets

import (
	"net/http"
	"sync"
	"fmt"
	"path"
	"os"
	"runtime"
	"os/exec"
	"time"
)

// Unix epoch time
var epoch = time.Unix(0, 0).Format(time.RFC1123)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
// a router (or subrouter) from being cached by an upstream proxy and/or client.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//      Expires: Thu, 01 Jan 1970 00:00:00 UTC
//      Cache-Control: no-cache, private, max-age=0
//      X-Accel-Expires: 0
//      Pragma: no-cache (for HTTP/1.0 proxies/clients)
func NoCache(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		f(w, r)
	}

	return fn
}

func init() {
	RegisterTarget("js", Javascript{})
}

type Javascript struct {}

func (Javascript) Compile(filename string) error { return nil }
func (t Javascript) Run(mainFile string) error {
	if Game {
		var wg sync.WaitGroup
	
		http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
			wg.Add(1)
			fmt.Fprintf(w, `
<html>
<head>
<style>
	body {
		margin: 0;
	}
</style>
</head>
<body>
	<script src="stack.js"></script>
	<script src="game.js"></script>
</body>
</html>
`)
			wg.Done()
		})
		http.HandleFunc("/stack.js",  func(w http.ResponseWriter, r *http.Request) {
			wg.Add(1)
			 http.ServeFile(w, r, "./stack.js")
			 wg.Done()
		})
		http.HandleFunc("/game.js",  func (w http.ResponseWriter, r *http.Request) {
			wg.Add(1)
			 http.ServeFile(w, r, path.Base(mainFile[:len(mainFile)-2]+".js"))
			 wg.Done()
		})
		http.HandleFunc("/data/", NoCache(func (w http.ResponseWriter, r *http.Request) {
			wg.Add(1)
			 http.ServeFile(w, r, ".."+r.URL.Path)
			  wg.Done()
		}))
		
		go func() {
			err := http.ListenAndServe(":9090", nil) // set listen port
			if err != nil {
				fmt.Println("ListenAndServe: ", err)
				os.Exit(1)
			}
		}()
		
		var url = "http://localhost:9090"
		var cmd string
		var args []string

		switch runtime.GOOS {
		case "windows":
			cmd = "cmd"
			args = []string{"/c", "start"}
		case "darwin":
			cmd = "open"
		default: // "linux", "freebsd", "openbsd", "netbsd"
			cmd = "xdg-open"
		}
		args = append(args, url)
		if err := exec.Command(cmd, args...).Run(); err != nil {
			return err
		}
		
		<-time.After(time.Second*3)
		wg.Wait()
	} else {
		run := exec.Command("nodejs", path.Base(mainFile[:len(mainFile)-2]+".js"))
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		return run.Run()	
	}
	return nil
}
func (Javascript) Export(filename string) error { return nil }
