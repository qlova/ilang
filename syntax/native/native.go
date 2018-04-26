package native

import "github.com/qlova/uct/compiler"
import "strings"

var Statement = compiler.Statement {
	Name: compiler.NoTranslation("."),

	OnScan: func(c *compiler.Compiler) {
		cmd := c.Scan()
		asm := ""

		var data bool
		if cmd == "data" {
			data = true
		}

		//Are we in a block of code?
		var block = false

		var peeking = c.Scan() 
		if peeking  == "{" {
			block = true
			c.Expecting("\n")
			asm = ""
		} else {
			c.NextToken = peeking 
		}

		//Do some magic so that we can use variables in inline assembly.
		//Keep track of braces so we can have blocks of code.
		var braces = 0
		var first = true
		var second = false
		var last = ""
		for {
			c.Scanners[len(c.Scanners)-1].Scan()
			var token = c.Scanners[len(c.Scanners)-1].TokenText()
			
			if token == "\\" && c.Peek() == "t" {
				c.Scan()
				asm += "\t"
				continue
			} else if token == "\\" {
				asm += c.Scan()
				continue
			}
			if strings.ContainsAny(token, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
				//c.GetVariable(token)
				/*if cmd == "grab" {
					ic.SetVariable(token, ilang.Array)
					ic.SetVariable(token+"_use", ilang.Used)
				}
				if cmd == "pull" || cmd == "pop" || cmd == "get" {
					ic.SetVariable(token, ilang.Number)
					ic.SetVariable(token+"_use", ilang.Used)
				}*/
			}
			if data {
				//ic.SetVariable(token, ilang.Text)
				data = false
			}
			if token == "\n" {
				
				c.Native(cmd, asm)
				
				if !block {
					break
				} else {
					asm = ""
				}
			} else {
				if asm == "" {
					asm = token
					
				}  else if first || (token[0] == '"') || 
					(strings.ContainsAny(last, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789") &&
					strings.ContainsAny(token, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789")) {
					asm += " "+token
					
				}else if strings.ContainsAny(token, "+-/*().=[]<>{}:;!@#$%^&*") {
					asm += token
					
				} else {
					asm += token
				}
			}
		
			if block {
				if token == "}" {
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
	},
}
