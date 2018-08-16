/*	
 *  Native syntax.
 * 
 * 	Run code on a specified target. For example, in order to run Go code:
 * 	
 * 	.go System.out.println("Hello World");
 * 
 */

package native

import "github.com/qlova/uct/compiler"
import "strings"

var Statement = compiler.Statement {
	
	//Native code begins with a dot.
	Name: compiler.NoTranslation("."),

	OnScan: func(c *compiler.Compiler) {
		
		//This can be a UCT command or a language extension.
		Command := c.Scan()
		Assembly := ""

		var ThisIsADataCommand bool
		if Command == "data" {
			ThisIsADataCommand = true
		}

		//Are we in a block of code?
		var ThisIsABlock = false

		//Do some magic so that we can use variables in inline assembly.
		//Keep track of braces so we can have blocks of code.
		var NumberOfBraces = 0
		var FirstIteration = true
		var SecondIteration = false
		
		var CheckInBlock = true
		
		var LastToken = ""
		for {
			
			c.Scanners[len(c.Scanners)-1].Scan()
			var token = c.Scanners[len(c.Scanners)-1].TokenText()
			
			if CheckInBlock && token == "{" {
				ThisIsABlock = true
				c.Expecting("\n")
				CheckInBlock = false
				continue
			}
			
			if token == "\\" && c.Peek() == "t" {
				c.Scan()
				Assembly += "\t"
				continue
			} else if token == "\\" {
				Assembly += c.Scan()
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
			if ThisIsADataCommand {
				//ic.SetVariable(token, ilang.Text)
				ThisIsADataCommand = false
			}
			if token == "\n" {

				c.Native(Command, strings.TrimSpace(Assembly))
				
				if !ThisIsABlock {
					break
				} else {
					Assembly = ""
				}
			} else {
				if Assembly == "" {
					Assembly = token
					
				}  else if FirstIteration || (token[0] == '"') || 
					(strings.ContainsAny(LastToken, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789") &&
					strings.ContainsAny(token, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789")) {
					Assembly += " "+token
					
				}else if strings.ContainsAny(token, "+-/*().=[]<>{}:;!@#$%^&*") {
					Assembly += token
					
				} else {
					Assembly += token
				}
			}
		
			if ThisIsABlock {
				if token == "}" {
					if NumberOfBraces == 0 {
						break
					} else {
						NumberOfBraces--
					}
				}
				if token == "{" {
					NumberOfBraces++
				}
			}
			
			if FirstIteration {
				FirstIteration = false
				SecondIteration = true
			} else if SecondIteration {
				SecondIteration = false
			}
			
			LastToken = token
		}
	},
}
