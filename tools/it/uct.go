package main

import (
	"github.com/qlova/uct/assembler"

	
	_ "github.com/qlova/uct/targets/go"
	_ "github.com/qlova/uct/targets/u"
)

import "os"
import "path"
import "fmt"
import "io"

func uct(name, target string) {
	var asm = new(assembler.Assembler)
	
	file4, err := os.Open(path.Dir(os.Args[3])+"/.it/main.u") 
	if err != nil {
		fmt.Println(err.Error())
		return
	} 
	
	
	asm.Input = append(asm.Input, file4)
		
	asm.Linker = func(name string) io.ReadCloser {
		file, err := os.Open(path.Dir(os.Args[3])+"/.it/"+name) 
		if err != nil {
			fmt.Println([]byte(name))
			panic("Could not open "+name)
		}
		
		return file
	}
	
	asm.Output, err = os.Create(path.Dir(os.Args[3])+"/.it/main."+os.Args[2])
	if err != nil {
		fmt.Println("Could not open ", asm.Input, "! ", err.Error())
		return
	}
	
	err = asm.Assemble(path.Ext(path.Dir(os.Args[3])+"/.it/main."+os.Args[2])[1:])
	if err != nil {
		fmt.Println("Could not assemble ! ", err.Error())
		return
	}
	
	file4.Close()
}
