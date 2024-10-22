package targets

import "path"
import "os"
import "os/exec"
import "runtime"
import "errors"

type Go struct {}

func (Go) Compile(mainFile string) error { 
	compile := exec.Command("go", "build", "-tags", "example", "-o",  path.Base(mainFile[:len(mainFile)-2])+".gob")
	compile.Stdout = os.Stdout
	compile.Stderr = os.Stderr
	return compile.Run() 
}
func (Go) Run(mainFile string) error {
		
	var args = []string{}
	if len(os.Args) > 3 {
		args = append(args, os.Args[3:]...)
	}
	
	run := exec.Command("./"+path.Base(mainFile[:len(mainFile)-2])+".gob", args...)
	run.Stdout = os.Stdout
	run.Stdin = os.Stdin
	run.Stderr = os.Stderr
	return run.Run()	
}
func (Go) Export(mainFile string) error { 
	if runtime.GOOS == "linux" ||  runtime.GOOS == "darwin" {

		return os.Rename(path.Base(mainFile[:len(mainFile)-2])+".gob", "../"+path.Base(mainFile[:len(mainFile)-2]))
		
	} else if runtime.GOOS == "windows" {
		
		return os.Rename(path.Base(mainFile[:len(mainFile)-2])+".gob", "../"+path.Base(mainFile[:len(mainFile)-2])+".exe")
		
	} else {
		return errors.New("Cannot export on "+runtime.GOOS+ " systems!")
	}
}

func init() {
	RegisterTarget("go", Go{})
} 
