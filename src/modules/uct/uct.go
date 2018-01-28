/*
	This package exposes uct directly to the programmer and allows them to inline uct assemby directly.
*/

package uct

import "github.com/qlova/ilang/src"
import "strings"

func init() {
	ilang.RegisterToken([]string{"."}, ScanAssembly)
}

func ScanAssembly(ic *ilang.Compiler) {
	cmd := ic.Scan(ilang.Name)
	asm := strings.ToUpper(cmd)

	var data bool
	if cmd == "data" {
		data = true
	}

	//Are we in a block of code?
	var block = false

	var peeking = ic.Scan(0) 
	if peeking  == "{" {
		block = true
		ic.Scan('\n')
		asm = ""
	} else {
		ic.NextToken = peeking 
	}

	//Do some magic so that we can use variables in inline assembly.
	//Keep track of braces so we can have blocks of code.
	var braces = 0
	var first = true
	var second = false
	var last = ""
	for {
		var token = ic.Scan(0)
		if token == "\\" && ic.Peek() == "t" {
			ic.Scan(0)
			asm += "\t"
			continue
		}
		if strings.ContainsAny(token, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			ic.GetVariable(token)
			if cmd == "grab" {
				ic.SetVariable(token, ilang.Array)
				ic.SetVariable(token+"_use", ilang.Used)
			}
			if cmd == "pull" || cmd == "pop" || cmd == "get" {
				ic.SetVariable(token, ilang.Number)
				ic.SetVariable(token+"_use", ilang.Used)
			}
		}
		if data {
			ic.SetVariable(token, ilang.Text)
			data = false
		}
		if token == "\n" {
			if block {
				asm = strings.ToUpper(cmd)+" "+asm
			}
			if ic.Header {
				ic.Library(asm)
			} else {
				ic.Assembly(asm)
			}
			if !block {
				break
			} else {
				asm = ""
			}
		} else {
			if asm == "" {
				asm = token
				
			}  else if first || (token[0] == '"') || 
				(strings.ContainsAny(last, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_") &&
				 strings.ContainsAny(token, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")) {
				asm += " "+token
				
			}else if strings.ContainsAny(token, "+-/*().=[]<>{}:;!@#$%^&*") {
				asm += token
				
			} else {
				asm += token
			}
		}
	
		if block {
			if token == "}"  {
		 		if braces == 0 {
					break
				} else {
					braces--
				}
			}
			if token == "{" {
				braces++
			}
		}
		
		if first {
			first = false
			second = true
		} else if second {
			second = false
		}
		
		last = token
	}
}
