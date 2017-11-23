package targets

import "path"
import "os"
import "os/exec"
import "errors"

type Lua struct {}

func (Lua) Compile(filename string) error { return nil }
func (Lua) Run(mainFile string) error {
	run := exec.Command("lua", path.Base(mainFile[:len(mainFile)-2]+".lua"))
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	return run.Run()	
}
func (Lua) Export(mainFile string) error { 
	return errors.New("Cannot export lua files!")
}

func init() {
	RegisterTarget("lua", Lua{})
}
