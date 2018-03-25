package macro

import "github.com/qlova/ilang/src"

var DefinedMacros = make(map[string]Macro)

type Macro struct {
	Arguments []string
	ilang.Plugin
}

func init() {
	
	ilang.RegisterExpression(func(ic *ilang.Compiler) string {
		
		if p, ok := DefinedMacros[ic.LastToken]; ok {
			//Argument loop.
			ic.Scan('(')
			
			for i := 0; i < len(p.Arguments); i++ {
				ic.Aliases[p.Arguments[i]] = ic.Scan(0)
				
				if i < len(p.Arguments)-1 {
					ic.Scan(',')
				}
			}
			
			ic.Scan(')')
			
			ic.Insertion = append(ic.Insertion, p.Plugin)
			
			defer func() {
				for i := 0; i < len(p.Arguments); i++ {
					delete(ic.Aliases, p.Arguments[i])
				}
			}()

			return ic.ScanExpression()
		}
		return ""
	})
	
	ilang.RegisterDefault(func(ic *ilang.Compiler) bool {
		
		if p, ok := DefinedMacros[ic.LastToken]; ok {
		
			//Argument loop.
			ic.Scan('(')
			
			for i := 0; i < len(p.Arguments); i++ {
				ic.Aliases[p.Arguments[i]] = ic.Scan(0)
				
				if i < len(p.Arguments)-1 {
					ic.Scan(',')
				}
			}
			
			ic.Scan(')')
			
			ic.Insertion = append(ic.Insertion, p.Plugin)
			var current = len(ic.Insertion)
			
			for {
				ic.ScanAndCompile()
				if len(ic.Insertion) < current {
					break
				}
			}
			
			for i := 0; i < len(p.Arguments); i++ {
				delete(ic.Aliases, p.Arguments[i])
			}

			return true
		}
		return false
	})

	ilang.RegisterToken([]string{"macro"}, func(ic *ilang.Compiler) {
		var name = ic.Scan(ilang.Name)
		ic.Scan('(')
		
		var p Macro
		
		//Count the generic arguments.
		for {
			var argument = ic.Scan(ilang.Name)
			p.Arguments = append(p.Arguments, argument)
			
			var token = ic.Scan(0)
			if token == ")" {
				ic.Scan('{')
				break
			} else if token != "," {
				ic.RaiseError("Arguments should be seperated by commas and be untyped!")
			}
		}
		
		//Plugin Time!
		var plugin ilang.Plugin
		plugin.Line = ic.Scanner.Pos().Line

		var braces = 0
		for {
			var token = ic.Scan(0)
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
			plugin.Tokens = append(plugin.Tokens, token)
		}

		plugin.File = ic.Scanner.Pos().Filename
		
		p.Plugin = plugin
		
		DefinedMacros[name] = p
	})
}
