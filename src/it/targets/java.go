package targets

import "path"
import "os"
import "os/exec"

type Java struct {}

func (Java) Compile(mainFile string) error { 
	compile := exec.Command("javac", "-Xlint:unchecked", path.Base(mainFile[:len(mainFile)-2])+".java")
	compile.Stdout = os.Stdout
	compile.Stderr = os.Stderr
	return compile.Run() 
}
func (Java) Run(mainFile string) error {
	
	var args = []string{path.Base(mainFile[:len(mainFile)-2])}
	if len(os.Args) > 3 {
		args = append(args, os.Args[3:]...)
	}
	
	run := exec.Command("java", args...)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	return run.Run()	
}
func (Java) Export(mainFile string) error { 

	f, err := os.Create("./Manifest.txt")
	verify(err)
	f.Write([]byte(`Main-Class: `+path.Base(mainFile[:len(mainFile)-2])+"\n"))
	f.Close()
	
	run := exec.Command("jar", "cfm", "../HelloWorld.jar", "Manifest.txt", path.Base(mainFile[:len(mainFile)-2])+".class",
	"Stack$Loader.class", "Stack$Pipe.class", "Stack$ArrayArray.class", "Stack$Number.class",
	"Stack$Pipe$Pipeable.class", "Stack$Array.class", "Stack$Opener.class",
	"Stack.class", "Stack$PipeArray.class")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	
	return run.Run()
}

func init() {
	RegisterTarget("java", Java{})
}
