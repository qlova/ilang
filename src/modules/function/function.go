package function

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/types/function"

var Flag = ilang.NewNamedFlag("Function")

func init() {
	ilang.RegisterToken([]string{"function"}, ScanFunction)
	ilang.RegisterToken([]string{"return"}, ScanReturn)
	ilang.RegisterExpression(FuncExpression)
	ilang.RegisterListener(Flag, FunctionEnd)
	ilang.RegisterShunt("(", ShuntFunctionCall)
}

func FuncExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	
	//Scan anonymous function.
	if token == "function" {
		ic.Scan('(')
		
		var name = ic.Tmp("anonymous")
		
		ic.DisableOutput = true
		//ic.Assembly("FUNCTION ", name)
		ic.GainScope()
		
		CreateFromArguments(name, ic)
		
		var copythisscope = ic.Scope[len(ic.Scope)-1]
		
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
		
		//Mock scan to load any needed requirements.
		ic.Insertion = append(ic.Insertion, plugin)
		
		for {	
			ic.ScanAndCompile()
			if !ic.GetFlag(Flag) {
				break
			}
		}
		
		ic.DisableOutput = false
		
		ic.SwapOutput()
		//Do it again properly.
		ic.Assembly("FUNCTION ", name)
		ic.GainScope()
		ic.Assembly(ic.DefinedFunctions[name].UnpackArguments)
		
		ic.Scope[len(ic.Scope)-1] = copythisscope
		
		ic.Insertion = append(ic.Insertion, plugin)
		
		for {	
			ic.ScanAndCompile()
			if !ic.GetFlag(Flag) {
				break
			}
		}
		
		ic.SwapOutput()
		
		var tmp = ic.Tmp("scope")
		ic.Assembly("SCOPE ", name)
		ic.Assembly("TAKE ", tmp)
		
		
		var t = function.Type
		t.Detail = new(ilang.UserType)
		t.Detail.Elements = ic.DefinedFunctions[name].Args
		
		ic.ExpressionType = t
		return tmp
	}
		
	if _, ok := ic.DefinedFunctions[token]; ok {
		if ic.Peek() == "(" {
			ic.ExpressionType = Flag
			
			return token
		}
	}
	return ""
}

func ShuntFunctionCall(ic *ilang.Compiler, name string) string {
	if ic.ExpressionType != Flag {
		return ""
	}
	
	var r = ic.ScanFunctionCall(name)
	ic.Scan(')')
	
	return ic.Shunt(r)
}

func ScanFunction(ic *ilang.Compiler) {
	ic.Header = false
	
	var name string = ic.Scan(ilang.Name)
	
	ic.Assembly("FUNCTION ", name)
	ic.Scan('(')
	ic.GainScope()
	
	CreateFromArguments(name, ic)
}

func FunctionEnd(ic *ilang.Compiler) {
	if ic.InOperatorFunction {
		ic.InOperatorFunction = false
		ic.Assembly("SHARE c")
	}
	ic.Assembly("RETURN")
}

func ScanReturn(ic *ilang.Compiler) {
	if !ic.CurrentFunction.Exists {
		ic.RaiseError("Cannot return, not in a function!")
	}
	
	if len(ic.CurrentFunction.Returns) == 1 {
		r := ic.ScanExpression()
		
		ic.SetVariable(r+".", ilang.Protected)
		
		if ic.CurrentFunction.Returns[0].Name == "thing" {
			ic.CurrentFunction.Returns[0] = ic.ExpressionType
		}
		
		/*if ic.CurrentFunction.Returns[0] == ilang.List {
			if !ic.ExpressionType.List && ic.ExpressionType != ilang.Array {
				ic.RaiseError("Cannot return '",ic.ExpressionType.Name,
				"', not a list!")
			}
			ic.CurrentFunction.Returns[0] = ic.ExpressionType
		}*/
		
		if !ic.ExpressionType.Equals(ic.CurrentFunction.Returns[0]) {
			ic.RaiseError("Cannot return '",ic.ExpressionType.GetComplexName(),
				"', expecting ",ic.CurrentFunction.Returns[0].GetComplexName())
		}
		
		
		ic.Assembly("%v %v", ic.ExpressionType.Push, r)
		
	}
	
	//Infer return type for this function, because I is clever and concise.
	if len(ic.CurrentFunction.Returns) == 0 {
		if ic.Peek() != "\n" {
			r := ic.ScanExpression()
			
			ic.SetVariable(r+".", ilang.Protected)
	
			ic.CurrentFunction.Returns = append(ic.CurrentFunction.Returns, ic.ExpressionType)
		
			ic.DefinedFunctions[ic.CurrentFunction.Name] = ic.CurrentFunction
		
			ic.Assembly("%v %v", ic.ExpressionType.Push, r)
		}
	}
	
	
	if len(ic.Scope) > 2 {
		//TODO garbage collection.
		ic.CollectGarbage()
		ic.Assembly("RETURN")
	}
}

func CreateFromArguments(name string, ic *ilang.Compiler) {
	var function ilang.Function
	function.Name = name
	
	//We need to reverse the POP's because of stack pain.
	if ic.Peek() != ")" {
		var toReverse []string
		for {
			//Identfy the type and add it to the function.
			var ArgumentType = ic.ScanSymbolicType()
		
			if ArgumentType == ilang.Variadic {
				function.Variadic = true
				ArgumentType = ilang.Array
			}
			function.Args = append(function.Args, ArgumentType)
		
			var name = ic.Scan(ilang.Name)
		
			ic.SetVariable(name, ArgumentType)
			ic.SetVariable(name+"_use", ilang.Used)

			toReverse = append(toReverse, ArgumentType.Pop+" "+name)
		
			token := ic.Scan(0)
		
			if token != "," {
				if token != ")" {
					ic.RaiseError()
				}
				break
			}
		}
		for i := len(toReverse)-1; i>=0; i-- {
			ic.Assembly(toReverse[i])
			function.UnpackArguments += toReverse[i]+"\n"
		}
	} else {
		ic.Scan(')')
	}
	
	
	token := ic.Scan(0)
	
	//Find out the return value.
	if token != "{" || (token == "{" && ic.Peek() == "}") {
	
		ic.NextToken = token
		var ReturnType = ic.ScanSymbolicType()
		
		function.Returns = append(function.Returns, ReturnType)
		
		
		if ReturnType == ilang.Number {
			ic.Scan(ilang.Name)
		}
		ic.Scan('{')
	}
	
	function.Exists = true
	function.Method = true
	
	ic.DefinedFunctions[name] = function
	
	ic.CurrentFunction = function
	
	ic.SetFlag(Flag)
}
