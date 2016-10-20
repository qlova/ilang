package main

import (
	"flag"
	"path"
	"os"
)

//TODO, this needs to be a seperate package.

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		return
	}
	
	ic := NewCompiler(file)
	
	//Open the output file with the file type replaced to .u
	if output, err := os.Create(flag.Arg(0)[:len(flag.Arg(0))-2]+".u"); err != nil {
		ic.RaiseError("Could not create output file!", err.Error())
	} else {
		ic.Output = output
	}
	
	if lib, err := os.Create(path.Dir(flag.Arg(0))+"/ilang.u"); err != nil {
		ic.RaiseError("Could not create output library file!", err.Error())
	} else {
		ic.Lib = lib	
	}

	ic.Compile()
}
