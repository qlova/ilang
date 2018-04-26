package main

import "os"
import "fmt"
import "path"
import "github.com/qlova/ilang/tools/it/targets"

func main() {
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: it run/build/export/help ext [file.i]")
		return
	} 
	
	if os.Args[1] == "help" {
		os.Args = os.Args[1:]
		help()
		return
	}
	
	if len(os.Args) < 3 {
		os.Args = append(os.Args, "go")
	}
	
	if len(os.Args) < 4 {
		var dir, _ = os.Getwd()
		os.Args = append(os.Args, dir+"/"+path.Base(dir)+".i")
	}
	
	var TargetLanguage = os.Args[2]
	var Mode = os.Args[1]
	var File = os.Args[3]
	
	ic(File)
	uct(File, TargetLanguage)
	
	os.Chdir(".it")
	
	if target, ok := targets.Targets[TargetLanguage]; ok {
		
		if err := target.Compile(File); err != nil {
			fmt.Println(err)
			return
		}
		
		switch Mode {
			case "run":
				if err := target.Run(File); err != nil {
					fmt.Println(err)
					return
				}
				
				os.Chdir("../")
				os.RemoveAll(".it")
			case "export":
				if err := target.Export(File); err != nil {
					fmt.Println(err)
					return
				}
			case "build":
				return
				
			default:
				fmt.Println("Unidentified command: ", Mode)
				return
		}
		
	} else if !ok {
		fmt.Println(TargetLanguage+" files are not a supported target!")
	}
	
}
