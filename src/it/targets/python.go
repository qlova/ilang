package targets

import "path"
import "os"
import "os/exec"
import "runtime"
import "errors"

type Python struct {}

func (Python) Compile(filename string) error { return nil }
func (Python) Run(mainFile string) error {
	var command = "python3"
	if runtime.GOOS == "windows" {
		command = "python"
	}
	
	os.Chdir("../")
	defer os.Chdir("./.it")
	
	var args = []string{"./.it/"+path.Base(mainFile[:len(mainFile)-2]+".py")}
	if len(os.Args) > 3 {
		args = append(args, os.Args[3:]...)
	}
	
	run := exec.Command(command, args...)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	return run.Run()	
}
func (Python) Export(mainFile string) error { 
	if runtime.GOOS == "linux" {
	
		//Export as a py file disguised as a zip, python knows what to do with these files.
		os.Rename(path.Base(mainFile[:len(mainFile)-2]+".py"), "__main__.py")
		run := exec.Command("zip", "../"+path.Base(mainFile[:len(mainFile)-2]+".py"), "__main__.py", "stack.py")
		return run.Run() 
		
	//TODO support py2exe on windows.
	} else {
		return errors.New("Cannot export on "+runtime.GOOS+ "systems!")
	}
}

func init() {
	RegisterTarget("py", Python{})
}
