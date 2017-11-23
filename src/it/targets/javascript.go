package targets

import (
	"net/http"
	"sync"
	"fmt"
	"path"
	"os"
	"runtime"
	"os/exec"
)

func init() {
	RegisterTarget("js", Javascript{})
}

type Javascript struct {}

func (Javascript) Compile(filename string) error { return nil }
func (t Javascript) Run(mainFile string) error {
	if Game {
		var wg sync.WaitGroup
		wg.Add(4)
	
		http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
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
		http.HandleFunc("/stack.js", func (w http.ResponseWriter, r *http.Request) {
			 http.ServeFile(w, r, "./stack.js")
			 wg.Done()
		})
		http.HandleFunc("/game.js", func (w http.ResponseWriter, r *http.Request) {
			 http.ServeFile(w, r, path.Base(mainFile[:len(mainFile)-2]+".js"))
			 wg.Done()
		})
		http.HandleFunc("/data/", func (w http.ResponseWriter, r *http.Request) {
			 http.ServeFile(w, r, ".."+r.URL.Path)
			 wg.Done()
		})
		
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
