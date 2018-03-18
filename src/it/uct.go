package main

import uc "github.com/qlova/uct/src"
import (
	"os"
	"path/filepath"
)

func uct(language, file string) {
	//Set the assembler.
	for _, asm := range uc.Assemblers {
		if asm.Ext == language {
			uc.SetAssembler(asm)
			
			var extension = filepath.Ext(file)
			var err error
			out, err := os.Create(file[0:len(file)-len(extension)]+"."+asm.Ext)
			if err != nil {
				uc.Output = os.Stdout
			} else {
				os.Chmod(file[0:len(file)-len(extension)]+"."+asm.Ext, 0755)
				uc.Output = out
				defer out.Close()
			}
		}
	}
	
	if !uc.AssemblerReady() {
		os.Stderr.Write([]byte("Please provide an assembler!"))
		os.Exit(1)
	}
	
	{
		

		uc.SetFileName(file)
	}

	//Write any necessary headers.
	uc.Output.Write(uc.Header())
	
	//Reset the assembler's instruction count. 
	err := uc.Assemble(file)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()+"\n"))
		os.Exit(1)
	}
	//fmt.Println(aliases)
	
	//Write any necessary footers.
	uc.Output.Write(uc.Footer())
}
