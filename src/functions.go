package main

import "strings"

type Function struct {
	Exists bool
	Loaded bool
	Import string
	Inline bool
	Data string
	
	Method bool
	
	Variadic bool
	
	Returns []Type
	Args []Type
}

func (ic *Compiler) RunFunction(name string) string {
	if strings.Contains(name, "_m_Something") && ic.ExpressionType.Interface != nil {
		var sort = name[:len(name)-len("_m_Something")]
		return ic.CallInterfaceMethod(sort)
	}

	f, ok := ic.DefinedFunctions[name]
	if !ok {
		ic.RaiseError(name, " does not exist!")
	}
	
	ic.LoadFunction(name)
	
	if f.Import != "" {
		ic.LoadFunction(f.Import)
	}
	
	if f.Inline {
		return f.Data
	} else if ic.Fork {
		ic.Fork = false
		
		var returns string
		for _, v := range f.Args {
			returns += "\n"+v.Pop+" "+ic.Tmp("")
		}
		
		return "FORK "+name+returns
	} else {
		return "RUN "+name
	}
}

func (ic *Compiler) ScanFunctionCall(name string) string {
	f := ic.DefinedFunctions[name]
	
	//TODO allow variadic arguments along with normal arguments.
	if f.Variadic {
		id := ic.Tmp("variadic")
		
		ic.Assembly("ARRAY ", id)
		for {
			value := ic.ScanExpression()
			ic.Assembly("PLACE ", id)
			ic.Assembly("PUT ", value)
			
			token := ic.Scan(0)
			if token != "," {
				if token != ")" {
					ic.RaiseError()
				}
				ic.NextToken = ")"
				break
			}
		}
	
		ic.Assembly("SHARE ", id)
	} else if len(f.Args) > 0 {
	
		var defaultarguments bool
		for i := range f.Args {
			if defaultarguments {
				if f.Args[i] == Number {
					ic.Assembly("PUSH 0")
				} else {
					ic.RaiseError("Not enough arguments!", " Expected ", len(f.Args), " but got ", i+1)
				}
				continue
			}
			arg := ic.ScanExpression()
			
			if f.Args[i] != ic.ExpressionType {
				if f.Args[i] == User {
					f.Args[i] = ic.ExpressionType
				} else {
					ic.RaiseError("Type mismatch! Argument ",i+1," of '"+name+"()' expects ",
						f.Args[i].Name,", got ",ic.ExpressionType.Name) 
				}
			}
			
			ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
			
			token := ic.Scan(0)
			if token != "," && token != ")" {
				ic.RaiseError()
			}
			if token == ")" {
				if i+1 != len(f.Args) {
					defaultarguments = true
				}
				ic.NextToken = ")"
			}
		}
	} else {
		
		//Calls methods such as:
		//(Similar to function overloading!)
		/*
			read(Type) read(AnotherType)
		*/
		if f.Method && ic.Peek() != ")" {
			arg := ic.ScanExpression()
			
			//Hardcoded LEN optimisation.
			if name == "len" && ic.ExpressionType.Push == "SHARE" {
				ic.ExpressionType = Number
				return "#"+arg
			}
			
			InheritMethods:
			if _, ok := ic.DefinedFunctions[name+"_m_"+ic.ExpressionType.Name]; !ok {
				if ic.ExpressionType.Super != "" {
					ic.ExpressionType = ic.DefinedTypes[ic.ExpressionType.Super]
					goto InheritMethods
				}
				ic.RaiseError("Method ",name," for type ",ic.ExpressionType.Name, "does not exist!")
			}
			
			//Only pass the argument if it has a value, for example, the following type would not be passed:
			// type Blank {}
			if ic.ExpressionType.Detail == nil || len(ic.ExpressionType.Detail.Elements) > 0 {
				ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
			}
			f = ic.DefinedFunctions[name+"_m_"+ic.ExpressionType.Name]
			name = name+"_m_"+ic.ExpressionType.Name
		}
	}
	
	ic.Assembly(ic.RunFunction(name))
	
	if len(f.Returns) > 0 {
		id := ic.Tmp("result")
		
		var ReturnType = f.Returns[0]
		
		
		ic.Assembly("%v %v", ReturnType.Pop, id)
		ic.ExpressionType = ReturnType
		
		return id
	}	
	return ""
}

func (ic *Compiler) ScanFunction() {
	var name string = ic.Scan(Name)
	
	ic.Assembly("FUNCTION ", name)
	ic.Scan('(')
	ic.GainScope()
	
	ic.function(name)
}

func (ic *Compiler) ScanNew() {
	var sort = ic.Scan(Name)
	var name string = "new_m_"+sort
	
	if name == "new_m_Game" {
		ic.NewGame = true
	}
	
	ic.Assembly("FUNCTION ", name)
	ic.Scan(')')
	ic.Scan('{')
	ic.GainScope()
	
	ic.SetFlag(New)
	ic.SetFlag(InMethod)
	ic.SetFlag(InFunction)
	
	ic.Assembly("PUSH ", len(ic.DefinedTypes[sort].Detail.Elements))
	ic.Assembly("MAKE")
	ic.Assembly("GRAB ", ic.DefinedTypes[sort].Name)
	ic.SetVariable(ic.DefinedTypes[sort].Name, ic.DefinedTypes[sort])
	
	ic.DefinedFunctions[name] = Function{Exists:true, Returns:[]Type{ic.DefinedTypes[sort]}}
	
	ic.InsertPlugins(name)
}

func (ic *Compiler) ScanMethod() {
	var name string = ic.Scan(Name)
	
	f := ic.DefinedFunctions[name]
	f.Method = true
	ic.DefinedFunctions[name] = f
	
	/*if name == "new" {
		ic.Scan('(')
		ic.ScanNew()
		return
	}*/	
		
	var MethodType = ic.LastDefinedType
	
	var token = ic.Scan(0)
	if token == "(" {
		token = ic.Scan(0)
		if token != ")" {
			if t, ok := ic.DefinedTypes[token]; ok {
				MethodType = t
			} else {
				ic.NextToken = token
			}
		}
		
		ic.LastDefinedType = MethodType
	
	
		if MethodType.Name == "Game" && name == "new" {
			ic.NewGame = true
		}
		if MethodType.Name == "Game" && name == "draw" {
			ic.DrawGame = true
		}
		if MethodType.Name == "Game" && name == "update" {
			ic.UpdateGame = true
		}
	
		name += "_m_"+MethodType.Name
	
		ic.Assembly("FUNCTION ", name)
		ic.GainScope()

		if len(MethodType.Detail.Elements) > 0 {	
			ic.Assembly("%v %v", MethodType.Pop, MethodType.Name)
			ic.SetVariable(MethodType.Name, MethodType)
			ic.SetVariable(MethodType.Name+"_use", Used)
		}
	
		ic.function(name)
		f = ic.DefinedFunctions[name]
		if name == "new_m_"+MethodType.Name {
			ic.SetFlag(New)
			f.Returns = []Type{MethodType}
		}
		ic.SetFlag(InMethod)
	
	
		f.Method = true
		ic.DefinedFunctions[name] = f
	
		ic.InsertPlugins(name)
	
	//Functional methods.
	} else if token == "." {	
	
		if !ic.TypeExists(name) {
			ic.RaiseError("Undefined type: ", name)
		}
		
		
		t := ic.DefinedTypes[name]
		ic.LastDefinedType = t
		
		
		name = ic.Scan(Name)
		name += "_m_"+t.Name
		
		ic.Assembly("FUNCTION ", name)
		ic.GainScope()
		ic.Scan('(')
	
		ic.function(name)
		
		if len(t.Detail.Elements) > 0 {
			ic.Assembly("%v %v", t.Pop, t.Name)
			ic.SetVariable(t.Name, t)
			ic.SetVariable(t.Name+"_use", Used)
		}
		
		f = ic.DefinedFunctions[name]
		ic.SetFlag(InMethod)
		
		f.Method = true
		ic.DefinedFunctions[name] = f
	
		ic.InsertPlugins(name)
	
	} else {
		var symbol = token
		var other = ic.Scan(Name)
		ic.Scan('{')
		
		if t, ok := ic.DefinedTypes[name]; ok {
			MethodType = t
		}
		
		var a = MethodType
		
		MethodType = ic.DefinedTypes[other]
		
		ic.LastDefinedType = MethodType
		
		var b = MethodType
		
		NewOperator(a, symbol, b, "SHARE %a\n SHARE %b\nRUN "+a.Name+"_"+symbol+"_"+b.Name+"\nGRAB %c", true)
		
		ic.Assembly("FUNCTION %s_%s_%s\n", a.Name, symbol, b.Name)
		ic.GainScope()
		ic.Assembly("GRAB b\nGRAB a\nARRAY c\n")
		for range a.Detail.Elements {
			ic.Assembly("PUT 0\n")
		}
		ic.InOperatorFunction = true
		
		ic.SetFlag(InFunction)
	
		ic.SetVariable("c", a)
		ic.SetVariable("a", a)
		ic.SetVariable("b", b)
	}
}

func (ic *Compiler) function(name string) {
	var function Function
	
	//We need to reverse the POP's because of stack pain.
	if ic.Peek() != ")" {
		var toReverse []string
		for {
			//Identfy the type and add it to the function.
			var ArgumentType = ic.ScanSymbolicType()
		
			if ArgumentType == Variadic {
				function.Variadic = true
				ArgumentType = Array
			}
			function.Args = append(function.Args, ArgumentType)
		
			var name = ic.Scan(Name)
		
			ic.SetVariable(name, ArgumentType)
			ic.SetVariable(name+"_use", Used)
		
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
		
		if ReturnType == Number {
			ic.Scan(Name)
		}
		ic.Scan('{')
	}
	
	function.Exists = true
	function.Method = true
	
	ic.DefinedFunctions[name] = function
	
	ic.CurrentFunction = function
	
	ic.SetFlag(InFunction)
}
