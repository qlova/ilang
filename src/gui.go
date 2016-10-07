package main

import (
	"strings"
)

// gui main { <button> </button> }
func (ic *Compiler) ScanGui() {
	var token = ic.Scan(0)
	
	var name = token
	if token == "{" {
		name = "main"
		ic.GUIMainExists = true
	} else {
		ic.Scan('{')
	}
	
	var asm = "DATA gui_"+name+" \""
	ic.SetVariable("gui_"+name, Text)
	
	var design string
	var jsbraces = 0
	for {
		var token = ic.Scan(0)
		if token == "}"  {
	 		if jsbraces == 0 {
				break
			} else {
				jsbraces--
			}
		}
		if token == "{" {
			jsbraces++
		}
		if token != "\n" {
			if strings.Contains(token, "\"") {
				design += strings.Replace(token, "\"", "\\\"", -1)
			} else {
				design += token
			}
			if ic.Peek() == " " {
				design += " "		
			}
		}
	}
	
	asm += design+"\""
	ic.Assembly(asm)
	ic.GUIExists = true
}
