package main

import "fmt"
import "github.com/google/go-github/github"
import "os"
import "time"
import "os/exec"
import "github.com/kardianos/osext"
import "path"
import "path/filepath"

func CheckForUpdate(uptodate time.Time) {
	client := github.NewClient(nil)
	
	repo, _, err := client.Repositories.Get("qlova", "ilang")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	//Check if we are already uptodate.
	if repo.UpdatedAt.Time.Sub(uptodate) <= 0 {
		exec.Command("touch", os.Args[0]).Start()
		return
	}
	
	fmt.Println("Update Available! Updating!")
	update := exec.Command(git, "pull")
	update.Dir, _ = osext.ExecutableFolder()
	update.Start()
	
	return
}

func verify(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	
	
	//Check if we should update IT and I.
	fpath, err := osext.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}
	info, err := os.Stat(fpath)
	if err != nil {
		 fmt.Println(err)
		return
	}
	
	//Make sure everything we need, is available.
	//Will also download missing components.
	SystemChecks()
	
	//If the executable hasn't been updated in an hour, check for an update.
	//(This will be scaled back before public use TODO)
	if time.Now().Sub(info.ModTime()) > time.Hour*24 {
		CheckForUpdate(info.ModTime())
	}
	
	//Command.
	if len(os.Args) > 1 {
		switch os.Args[1] {
			case "build", "run":
				//TODO clean this up along with (grep.go)
				filepath.Walk("./", func(name string, file os.FileInfo, err error) error {
					if !file.IsDir() && path.Ext(name) == ".i" {
						wg.Add(1)
						go grep(&wg, name)
					}
					return nil
				})
				wg.Wait()
				if mainFile == "" {
					fmt.Println("Could not find a 'software' block!")
					os.Exit(1)
				}
				
				//Compile.
				os.Mkdir(path.Dir(mainFile)+"/.it", 0700)
				compile := exec.Command(ic, "-o", path.Dir(mainFile)+"/.it", mainFile)
				compile.Stdout = os.Stdout
				compile.Stderr = os.Stderr
				verify(compile.Run())
				
				//Other languages.
				if len(os.Args) > 2 {
					compile = exec.Command(uct, os.Args[2], path.Base(mainFile[:len(mainFile)-2]+".u"))
					compile.Stdout = os.Stdout
					compile.Stderr = os.Stderr
					compile.Dir = path.Dir(mainFile)+"/.it/"
					verify(compile.Run())
				} else {
					compile = exec.Command(uct, "-go", path.Base(mainFile[:len(mainFile)-2]+".u"))
					compile.Stdout = os.Stdout
					compile.Stderr = os.Stderr
					compile.Dir = path.Dir(mainFile)+"/.it/"
					verify(compile.Run())
					compile = exec.Command(goc, "build", "-o",  "../"+path.Base(mainFile[:len(mainFile)-2]))
					compile.Stdout = os.Stdout
					compile.Stderr = os.Stderr
					compile.Dir = path.Dir(mainFile)+"/.it/"
					verify(compile.Run())
				
					if os.Args[1] == "run" {
						run := exec.Command("./"+path.Base(mainFile[:len(mainFile)-2]))
						run.Stdout = os.Stdout
						run.Stderr = os.Stderr
						run.Stdin = os.Stdin
						run.Run()
					}
				}
		}
	}
}
