package main

import "fmt"
import "github.com/google/go-github/github"
import "os"
import "time"
import "os/exec"
import "github.com/kardianos/osext"
import "path"
import "context"
import "bufio"
import "io/ioutil"
import "net/http"

func CheckForUpdate(uptodate time.Time) {
	ctx := context.Background()

	client := github.NewClient(nil)
	
	repo, _, err := client.Repositories.Get(ctx, "qlova", "ilang")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	//Check if we are already uptodate.
	//TODO windows compatibillity.
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

func cargo(mode string) {
	
	dir := os.Getenv("HOME")+"/.cargo/target"
	os.Mkdir(dir, 0700)
	
	env := os.Environ()
	env = append(env, fmt.Sprintf("CARGO_TARGET_DIR="+dir))

	compile := exec.Command("cargo", mode)
	compile.Stdout = os.Stdout
	compile.Stderr = os.Stderr
	compile.Env = env
	verify(compile.Run())
	
	verify(os.Rename(dir+"/debug/"+path.Base(mainFile[:len(mainFile)-2]), "./"+path.Base(mainFile[:len(mainFile)-2])))
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
		//CheckForUpdate(info.ModTime())
	}
	
	
	//Command.
	if len(os.Args) > 1 {
		os.Chdir(os.Args[1])
	}
	
	files, _ := ioutil.ReadDir("./")
    for _, f := range files {
            if path.Ext(f.Name()) == ".i" {
            	wg.Add(1)
				go grep(&wg, f.Name())
			}
    }
	wg.Wait()
	if mainFile == "" {
		fmt.Println("Could not find a 'software' block!")
		os.Exit(1)
	}
	
	//Compile.
	
	programdir, _ := os.Getwd()
	
	os.Mkdir(path.Dir(mainFile)+"/.it", 0700)
	ic(mainFile, path.Dir(mainFile)+"/.it")
	os.Chdir(programdir)
	
	
	
	//Other languages.
	if len(os.Args) > 2 {
		
		
		
		switch os.Args[2] {
			case "rs": //Rust needs to be handled differently.
				os.Chdir(".it")
				uct(os.Args[2], path.Base(mainFile[:len(mainFile)-2]+".u"))
				os.Mkdir("src", 0755)
				verify(os.Rename(path.Base(mainFile[:len(mainFile)-2]+".rs"), "src/main.rs"))
				verify(os.Rename("stack.rs","src/stack.rs"))
				f, err := os.Create("Cargo.toml")
				verify(err)
				f.Write([]byte(`[package]
name = "`+path.Base(mainFile[:len(mainFile)-2])+`"
version = "0.1.0"

[dependencies]
num = "0.1"
rand = "0.3"
`))
				f.Close()
				
				cargo("build")
				
				
			default:
				os.Chdir(".it")
				uct(os.Args[2], path.Base(mainFile[:len(mainFile)-2]+".u"))
		}
		
		if os.Args[2] == "js" && Game {
			
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
			})
			http.HandleFunc("/stack.js", func (w http.ResponseWriter, r *http.Request) {
				 http.ServeFile(w, r, "./stack.js")
			})
			http.HandleFunc("/game.js", func (w http.ResponseWriter, r *http.Request) {
				 http.ServeFile(w, r, path.Base(mainFile[:len(mainFile)-2]+".js"))
			})
			
			fmt.Println("Go to http://localhost:9090 to play your game!")
			
			go func() {
				err := http.ListenAndServe(":9090", nil) // set listen port
				if err != nil {
					fmt.Println("ListenAndServe: ", err)
				}
			}()
			fmt.Println("\nPress 'Enter' to stop...")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
		
		
	} else {
		os.Chdir(".it")
		uct("go", path.Base(mainFile[:len(mainFile)-2]+".u"))
		compile := exec.Command(goc, "build", "-tags", "example", "-o",  "../"+path.Base(mainFile[:len(mainFile)-2])+ext)
		compile.Stdout = os.Stdout
		compile.Stderr = os.Stderr
		verify(compile.Run())
	
		os.Chdir("../")
		//if os.Args[1] == "run" {
			run := exec.Command("./"+path.Base(mainFile[:len(mainFile)-2]))
			run.Stdout = os.Stdout
			run.Stderr = os.Stderr
			run.Stdin = os.Stdin
			run.Run()
		//}
		
		fmt.Println("\n[SOFTWARE EXIT]\nPress 'Enter' to close...")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	}
}
