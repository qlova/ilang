package targets

var Game bool

type Target interface {
	Compile(filename string) error
	Run(filename string) error
	Export(filename string) error
}

var Targets = make(map[string]Target)

func RegisterTarget(ext string, target Target) {
	Targets[ext] = target
}

/*

type TargetProcessing struct {}

func (TargetProcessing) Compile(filename string) {}
func (TargetProcessing) Run(mainFile string) {
	os.Mkdir("Processing", 0777)
	os.Remove("Processing/Processing.pde")
	verify(os.Rename(path.Base(mainFile[:len(mainFile)-2]+".pde"), "Processing/Processing.pde"))

	dir, err := os.Getwd()
	verify(err)

	run := exec.Command("processing-java", "--sketch="+dir+"/Processing", "--run")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	verify(run.Run())
}
func (TargetProcessing) Export(filename string) {}
*/
