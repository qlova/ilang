package generic

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/modules/function"
import "fmt"

var DefinedGenerics = make(map[string]Generic)

type Generic struct {
	Arguments []string
	Functions map[string]ilang.Function
	ilang.Plugin
}

func init() {
	ilang.RegisterExpression(func(ic *ilang.Compiler) string {
		if g, ok := DefinedGenerics[ic.LastToken]; ok {
			//TODO reuse generated generics.
			
			var generic_name = ic.LastToken
			var name = ic.LastToken+"_g"+fmt.Sprint(len(g.Functions))
			
			var types = make([]ilang.Type, 0, len(g.Arguments))
			
			var f ilang.Function
			
			//Generate the generic LOL.
			
			//Argument loop.
			ic.Scan('(')
			
			var nameReverse []string
			var typeReverse []ilang.Type
			for i := 0; i < len(g.Arguments); i++ {
				var arg = ic.ScanExpression()
				ic.Assembly(ic.ExpressionType.Push, " ", arg)
				types = append(types, ic.ExpressionType)
				
				f.Args = append(f.Args, ic.ExpressionType)
				
				nameReverse = append(nameReverse, g.Arguments[i])
				typeReverse = append(typeReverse, ic.ExpressionType)
				
				if i < len(g.Arguments)-1 {
					ic.Scan(',')
				}
			}
			ic.Scan(')')
			
			for fname, fg := range g.Functions {
				var match bool = false
				for i, arg := range fg.Args {
					if arg == f.Args[i] {
						match = true
					}
					if arg != f.Args[i] {
						match = false
						break
					}
				}
				
				//Generic already exists!
				if match {
					if len(ic.DefinedFunctions[fname].Returns) > 0 {
						ic.ExpressionType = ic.DefinedFunctions[fname].Returns[0]
					} else {
						ic.RaiseError(generic_name, " does not return any values! Cannot be used within an expression!")
					}
		
					ic.Assembly("RUN ", fname)
		
					var r = ic.Tmp("generic_result")
					ic.Assembly(ic.ExpressionType.Pop, " ", r)
		
					return r
				}
			}
			
			ic.SwapOutput()
			
			ic.Assembly("FUNCTION ", name)
			ic.GainScope()
			
			for i := len(nameReverse)-1; i>=0; i-- {
				ic.Assembly(typeReverse[i].Pop, " ", nameReverse[i])
				ic.SetVariable(nameReverse[i], typeReverse[i])
				ic.SetVariable(nameReverse[i]+"_use", ilang.Used)
			}

			ic.Insertion = append(ic.Insertion, g.Plugin)
	
			//Function stuff
			f.Name = name	
			f.Exists = true
			f.Method = true
	
			ic.DefinedFunctions[name] = f
	
			ic.CurrentFunction = f
	
			ic.SetFlag(function.Flag)
			
			g.Functions[name] = f 
			
			for {	
				ic.ScanAndCompile()
				if !ic.GetFlag(function.Flag) {
					break
				}
			}
			
			ic.SwapOutput()
			
			if len(ic.DefinedFunctions[name].Returns) > 0 {
				ic.ExpressionType = ic.DefinedFunctions[name].Returns[0]
			} else {
				ic.RaiseError(generic_name, " does not return any values! Cannot be used within an expression!")
			}
			
			ic.Assembly("RUN ", name)
			
			var r = ic.Tmp("generic_result")
			ic.Assembly(ic.ExpressionType.Pop, " ", r)
			
			ic.Scan(0)
			
			return r
		}
		return ""
	})
	
	ilang.RegisterToken([]string{"generic"}, func(ic *ilang.Compiler) {
		var name = ic.Scan(ilang.Name)
		ic.Scan('(')
		
		var g Generic
		
		//Count the generic arguments.
		for {
			var argument = ic.Scan(ilang.Name)
			g.Arguments = append(g.Arguments, argument)
			
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
		 			plugin.Tokens = append(plugin.Tokens, token)
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
		
		g.Plugin = plugin
		g.Functions = make(map[string]ilang.Function)
		
		DefinedGenerics[name] = g
	})
}
