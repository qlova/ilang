package main

import "os"
import "os/exec"
import "path"

import "fmt"
import "net/http"
import "runtime"
import "sync"

type Target interface {
	Compile(filename string)
	Run(filename string)
	Export(filename string)
}

var Targets = make(map[string]Target)

func init() {
	Targets["lua"] = TargetLua{}
	Targets["py"] = TargetPython{}
	Targets["js"] = TargetJavascript{}
	Targets["pde"] = TargetProcessing{}
}

type TargetLua struct {}

func (TargetLua) Compile(filename string) {}
func (TargetLua) Run(filename string) {
	if Game {
		os.Remove("main.lua")
		verify(os.Rename(path.Base(mainFile[:len(mainFile)-2]+".lua"), "main.lua"))

		run := exec.Command("love", ".")
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		verify(run.Run())
	} else {
		run := exec.Command("lua5.1", path.Base(mainFile[:len(mainFile)-2]+".lua"))
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		verify(run.Run())	
	}
}
func (TargetLua) Export(filename string) {}

type TargetPython struct {}

func (TargetPython) Compile(filename string) {}
func (TargetPython) Run(filename string) {
	run := exec.Command("python3", path.Base(mainFile[:len(mainFile)-2]+".py"))
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	verify(run.Run())	
}
func (TargetPython) Export(filename string) {}

type TargetJavascript struct {}

func (TargetJavascript) Compile(filename string) {}
func (TargetJavascript) Run(filename string) {
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
		verify(exec.Command(cmd, args...).Run())
		
		wg.Wait()
	} else {
		run := exec.Command("nodejs", path.Base(mainFile[:len(mainFile)-2]+".js"))
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		verify(run.Run())	
	}
}
func (TargetJavascript) Export(filename string) {}

type TargetProcessing struct {}

func (TargetProcessing) Compile(filename string) {}
func (TargetProcessing) Run(filename string) {
	os.Mkdir("Processing", 0777)
	os.Remove("Processing/Processing.pde")
	verify(os.Rename(path.Base(mainFile[:len(mainFile)-2]+".pde"), "Processing/Processing.pde"))

	dir, err := os.Getwd()
	verify(err)

	run := exec.Command("processing-java", "--sketch="+dir+"/Processing", "--run")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	verify(run.Run())
}
func (TargetProcessing) Export(filename string) {}
