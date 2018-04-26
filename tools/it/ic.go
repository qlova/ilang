package main

import 	"github.com/qlova/uct/compiler"
import	"github.com/qlova/ilang/syntax"

import "os"
import "fmt"
import "path"

func ic(name string) {
	var c = new(compiler.Compiler)
	c.StdErr = append(c.StdErr, os.Stdout)
	c.SetSyntax(i.Syntax())
	
	file1, err := os.Open(name) 
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	
	c.AddInput(file1)
	c.Language = compiler.English
	
	os.Mkdir(".it", 0755)
	
	file2, err := os.Create(path.Dir(name)+"/.it/main.u") 
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	
	c.Output = file2
	
	file3, err := os.Create(path.Dir(name)+"/.it/header.u") 
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	
	c.Header = file3
	
	c.Link("header.u")
	c.Compile()
	if c.Errors {
		os.Exit(1)
	}
	
	file1.Close()
	file2.Close()
	file3.Close()
}
