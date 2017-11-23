package main

import "fmt"
import "os"
import path "path/filepath"
import "io/ioutil"
import "github.com/qlova/ilang/src/it/targets"

type Gamer interface {
	SetGame()
}

/*func CheckForUpdate(uptodate time.Time) {
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
}*/

func verify(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	
	
	var TargetLanguage string
	var Mode string
	
	if len(os.Args) == 1 {
		fmt.Println("Usage: it run/build/export ext")
		return
	} 
	
	//Make sure everything we need, is available.
	//Will also download missing components.
	//fmt.Print("Checking System... ")
	SystemChecks()
	//fmt.Println("ok!")
	
	//Command.
	if len(os.Args) > 1 {
		os.Chdir(os.Args[1])
		
		var ext string
		//Set TargetLanguage to be the extension of the directory.
		if os.Args[1][len(os.Args[1])-1] == '/' {
			ext = path.Ext(os.Args[1][:len(os.Args[1])-1])
		} else {
			ext = path.Ext(os.Args[1])
		}
		if ext != "" && len(ext) > 1 {
			TargetLanguage = ext[1:]
		}
	}
	
	files, _ := ioutil.ReadDir("./")
    for _, f := range files {
            if path.Ext(f.Name()) == ".i" {
            	wg.Add(1)
				go grep(&wg, f.Name())
			}
    }
    
    //fmt.Print("Finding a file... ")
    
	wg.Wait()
	if mainFile == "" {
		fmt.Println("Could not find a 'software' block!")
		os.Exit(1)
	}
	//fmt.Println("Found!")
	
	//Compile.
	
	programdir, _ := os.Getwd()
	
	os.Mkdir(path.Dir(mainFile)+"/.it", 0700)
	ic(mainFile, path.Dir(mainFile)+"/.it")
	os.Chdir(programdir)
	
	//Other languages.
	if len(os.Args) > 2 {
		
		Mode = os.Args[1]
		TargetLanguage = os.Args[2]
		
	} else if TargetLanguage == "" {
		
		TargetLanguage = "go"	
		
	}
	
	//fmt.Print("Transpilling for ", TargetLanguage, "... "
	os.Chdir(".it")
	uct(TargetLanguage, path.Base(mainFile[:len(mainFile)-2]+".u"))
	
	if target, ok := targets.Targets[TargetLanguage]; ok {
		targets.Game = Game
	
		verify(target.Compile(mainFile))
		
		switch Mode {
			case "run":
				verify(target.Run(mainFile))
			case "export":
				verify(target.Export(mainFile))
			case "build":
				return
			default:
				fmt.Println("Unidentified command: ", Mode)
				return
		}	
		//Cleanup.
		verify(os.Chdir("../"))
		verify(os.RemoveAll(".it"))
		return
	} else if !ok {
		fmt.Println(TargetLanguage+" files are not a supported target!")
	}
	
}
	
