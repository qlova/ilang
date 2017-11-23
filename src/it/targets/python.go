package targets

import "path"
import "os"
import "os/exec"
import "runtime"
import "errors"

type Python struct {}

func (Python) Compile(filename string) error { return nil }
func (Python) Run(mainFile string) error {
	run := exec.Command("python3", path.Base(mainFile[:len(mainFile)-2]+".py"))
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
