package main

import (
	"flag"
	"path"
	"os"
	"github.com/qlova/ilang/src"
)

var directory string
func init() {
	flag.StringVar(&directory, "o", "", "directory to output to")
}

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		return
	}
	
	ic := ilang.NewCompiler(file)
	
	if directory == "" {
		directory = path.Dir(flag.Arg(0))
	}
	
	//Open the output file with the file type replaced to .u
	var filename = path.Base(flag.Arg(0))[:len(path.Base(flag.Arg(0)))-2]+".u"
	
	if output, err := os.Create(directory+"/"+filename); err != nil {
		ic.RaiseError("Could not create output file!", err.Error())
	} else {
		ic.Output = output
	}
	
	
	if lib, err := os.Create(directory+"/ilang.u"); err != nil {
		ic.RaiseError("Could not create output library file!", err.Error())
	} else {
		ic.Lib = lib	
	}

	ic.Compile()
}
