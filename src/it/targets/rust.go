package targets

import "path"
import "os"
import "os/exec"
import "runtime"
import "errors"
import "fmt"

func verify(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Rust struct {}

func (Rust) Compile(mainFile string) error { 
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

	dir := os.Getenv("HOME")+"/.cargo/target"
	os.Mkdir(dir, 0700)
	
	env := os.Environ()
	env = append(env, fmt.Sprintf("CARGO_TARGET_DIR="+dir))

	compile := exec.Command("cargo", "build")
	compile.Stdout = os.Stdout
	compile.Stderr = os.Stderr
	compile.Env = env
	verify(compile.Run())
	
	verify(os.Rename(dir+"/debug/"+path.Base(mainFile[:len(mainFile)-2]), "./"+path.Base(mainFile[:len(mainFile)-2])+".rsb"))
	return nil
}
func (Rust) Run(mainFile string) error {
	run := exec.Command("./"+path.Base(mainFile[:len(mainFile)-2])+".rsb")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	return run.Run()	
}
func (Rust) Export(mainFile string) error { 
	if runtime.GOOS == "linux" {

		return os.Rename(path.Base(mainFile[:len(mainFile)-2])+".rsb", "../"+path.Base(mainFile[:len(mainFile)-2]))
		
	//TODO support exe on windows.
	} else {
		return errors.New("Cannot export on "+runtime.GOOS+ "systems!")
	}
}

func init() {
	RegisterTarget("rs", Rust{})
}
