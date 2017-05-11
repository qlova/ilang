package main

import (
	"path"
	"os"
	"../"
)

func ic(input, directory string) {

	file, err := os.Open(input)
	if err != nil {
		return
	}
	
	ic := ilang.NewCompiler(file)
	
	if directory == "" {
		directory = path.Dir(input)
	}
	
	//Open the output file with the file type replaced to .u
	var filename = path.Base(input)[:len(path.Base(input))-2]+".u"
	
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
