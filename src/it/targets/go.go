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
	run := exec.Command("./"+path.Base(mainFile[:len(mainFile)-2])+".gob")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	return run.Run()	
}
func (Go) Export(mainFile string) error { 
	if runtime.GOOS == "linux" {

		return os.Rename(path.Base(mainFile[:len(mainFile)-2])+".gob", "../"+path.Base(mainFile[:len(mainFile)-2]))
		
	//TODO support exe on windows.
	} else {
		return errors.New("Cannot export on "+runtime.GOOS+ "systems!")
	}
}

func init() {
	RegisterTarget("go", Go{})
}
