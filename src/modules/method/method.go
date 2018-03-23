package method

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/modules/function"

var Flag = ilang.NewNamedFlag("Method")
var New = ilang.NewNamedFlag("New")

func Get(ic *ilang.Compiler, t ilang.Type, name string) *ilang.Function {
	
	if ic.FunctionExists(name+"_m_"+t.GetComplexName()) {
		
		var f = ic.DefinedFunctions[name+"_m_"+t.GetComplexName()]
		return &f
		
	} else if ic.FunctionExists(name+"_m_"+t.Super) {
		
		var f = ic.DefinedFunctions[name+"_m_"+t.Super]
		return &f
		
	} else {
		
		return nil
	}
	
}

func init() {
	ilang.RegisterToken([]string{"method"}, ScanMethod)
	ilang.RegisterListener(New, NewEnd)
	
	ilang.RegisterDefault(func(ic *ilang.Compiler) bool {
		token := ic.LastToken
		
		if ic.TypeExists(token) {
			if ic.DefinedTypes[token].Empty() && ic.Peek() == "." {
				
				var TYPE = ic.DefinedTypes[token]
				
				for {
					ic.Scan('.')
				
					var name = ic.Scan(ilang.Name)
					
					if TYPE.Detail != nil {
						if index, ok := TYPE.Detail.Table[name]; ok {
							TYPE = TYPE.Detail.Elements[index]
							continue
						}
					}
					
					ic.ExpressionType = function.Flag
					
					ic.Shunt(name+"_m_"+TYPE.GetComplexName())
					ic.Scan('(') //BUG I don't know why this has to be here.
					return true
				}
			}
		}
		
		if ic.GetFlag(Flag) && ic.LastDefinedType.Detail != nil {
			if _, ok := ic.LastDefinedType.Detail.Table[token]; ok {
				ic.NextToken = ic.LastDefinedType.Name
				ic.NextNextToken = "."
				ic.NextNextNextToken = token
				ic.ScanStatement()
				return true
			}
		}
		
		return false
	})
	
	ilang.RegisterExpression(func(ic *ilang.Compiler) string {
		token := ic.LastToken
		
		if ic.TypeExists(token) {
			if ic.DefinedTypes[token].Empty() && ic.Peek() == "." {
				
				var TYPE = ic.DefinedTypes[token]
				
				for {
					ic.Scan('.')
					var name = ic.Scan(ilang.Name)
					
					if TYPE.Detail != nil {
						if index, ok := TYPE.Detail.Table[name]; ok {
							TYPE = TYPE.Detail.Elements[index]
							continue
						}
					}
					
					ic.ExpressionType = function.Flag
					
					//Error checking.
					f, ok := ic.DefinedFunctions[name+"_m_"+TYPE.GetComplexName()]
					if !ok {
						ic.RaiseError("Method ",  name, " for ", TYPE.GetComplexName(), " does not exist!")
					}
					if len(f.Returns) == 0 {
						ic.RaiseError("Method ", TYPE.GetComplexName(), ".", name, "() does not return anything!")
					}
					
					var test = ic.Shunt(name+"_m_"+TYPE.GetComplexName())
					return test
				}
			}
		}
		
		if ic.GetFlag(Flag) {
			if ic.TypeExists(token) && ic.LastDefinedType.Super == token {
	 			ic.ExpressionType = ic.DefinedTypes[ic.LastDefinedType.Super]
				return ic.Shunt(ic.LastDefinedType.Name)
		 	}
		
			//Scope methods with multiple arguments inside the method.
			//eg. method Package.dosomething(22)
			// in a Package method, dosomething(22) should be local.
			if _, ok  := ic.DefinedFunctions[token+"_m_"+ic.LastDefinedType.Name]; ok {
				var f = token+"_m_"+ic.LastDefinedType.Name
				
				if !ic.LastDefinedType.Empty() {
					ic.Assembly(ic.LastDefinedType.Push," ", ic.LastDefinedType.Name)
				}
				
				ic.ExpressionType = function.Flag
				
				return ic.Shunt(f)
			}
		}
		
	 	return ""
	})
	
	ilang.RegisterVariable(func(ic *ilang.Compiler, name string) ilang.Type {
		//Allow table values to be indexed in a method.
		if ic.GetFlag(Flag) && ic.LastDefinedType.Detail != nil {
			if _, ok := ic.LastDefinedType.Detail.Table[name]; ok {
				ic.DisableOwnership = true
				var value = ic.IndexUserType(ic.LastDefinedType.Name, name)
				ic.DisableOwnership = false
				
				ic.AssembleVar(name, value)
				ic.SetVariable(name+"_use", ilang.Used)
				ic.SetVariable(name+".", ilang.Protected)
				return ic.ExpressionType
			}
		}
		return ilang.Undefined
	}) 
	
	ilang.RegisterShunt(".", ShuntMethodCall)
}

func Sync(ic *ilang.Compiler, variables ...string) {
	if ic.GetFlag(Flag) && ic.LastDefinedType.Detail != nil {
	
		for _, variable := range variables {
	
			if _, ok := ic.LastDefinedType.Detail.Table[variable]; ok {
				ic.SetUserType(ic.LastDefinedType.Name, variable, variable)
			}
		}
	}
}

func ShuntMethodCall(ic *ilang.Compiler, name string) string {
	var index = ic.Scan(ilang.Name)
	
	if _, ok  := ic.DefinedFunctions[index+"_m_"+ic.ExpressionType.GetComplexName()]; ok && ic.Peek() == "(" {
		var f = index+"_m_"+ic.ExpressionType.GetComplexName()
		ic.Assembly(ic.ExpressionType.Push," ", name)
		ic.ExpressionType = function.Flag
		return ic.Shunt(f)
	}
	
	ic.NextToken = index
	return ""
}

func NewEnd(ic *ilang.Compiler) {
	ic.Assembly("SHARE ", ic.LastDefinedType.Name)
	ic.LoseScope()
}

func Call(ic *ilang.Compiler, name string, t ilang.Type) {
	if ic.FunctionExists(name+"_m_"+t.GetComplexName()) {
		ic.Assembly(ic.RunFunction(name+"_m_"+t.GetComplexName()))
	} else if ic.FunctionExists(name+"_m_"+t.Super) {
		ic.Assembly(ic.RunFunction(name+"_m_"+t.Super))
	} else {
		ic.RaiseError("Method ", name, " for type ", t.GetComplexName(), " does not exist!")
	}
}

func ScanMethod(ic *ilang.Compiler) {
	ic.Header = false
	
	var name string = ic.Scan(ilang.Name)
	
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
	
	
		if MethodType.Name == "Graphics" && name == "new" {
			ic.NewGame = true
		}
		if MethodType.Name == "Graphics" && name == "draw" {
			ic.DrawGame = true
		}
		if MethodType.Name == "Graphics" && name == "update" {
			ic.UpdateGame = true
		}
	
		name += "_m_"+MethodType.Name
	
		ic.Assembly("FUNCTION ", name)
		ic.GainScope()

		if name == "new_m_"+MethodType.Name {
		
			ic.Assembly("PUSH ", len(MethodType.Detail.Elements))
			ic.Assembly("MAKE")
			ic.Assembly("GRAB ", MethodType.Name)
			
			ic.SetVariable(MethodType.Name, MethodType)
			ic.SetVariable(MethodType.Name+"_use", ilang.Used)
			ic.SetVariable(MethodType.Name+".", ilang.Protected)
		
		} else if  MethodType.Detail != nil && len(MethodType.Detail.Elements) > 0 {	
			ic.Assembly("%v %v", MethodType.Pop, MethodType.Name)
			ic.SetVariable(MethodType.Name, MethodType)
			ic.SetVariable(MethodType.Name+"_use", ilang.Used)
			ic.SetVariable(MethodType.Name+".", ilang.Protected)
		}
	
		function.CreateFromArguments(name, ic)
		ic.SetFlag(Flag)
		
		f = ic.DefinedFunctions[name]
		if name == "new_m_"+MethodType.Name {
			ic.GainScope()
			ic.SetFlag(New)
			f.Returns = []ilang.Type{MethodType}
			
			
			
		}
	
	
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
		
		
		name = ic.Scan(ilang.Name)
		name += "_m_"+t.GetComplexName()
		
		ic.Assembly("FUNCTION ", name)
		ic.GainScope()
		ic.Scan('(')
	
		function.CreateFromArguments(name, ic)
		
		if t.Detail == nil || len(t.Detail.Elements) > 0 {
			ic.Assembly("%v %v", t.Pop, t.Name)
			ic.SetVariable(t.Name, t)
			ic.SetVariable(t.Name+"_use", ilang.Used)
		}
		
		f = ic.DefinedFunctions[name]
		ic.SetFlag(Flag)
		
		f.Method = true
		ic.DefinedFunctions[name] = f
	
		ic.InsertPlugins(name)
	
	} else {
		var symbol = token
		var other = ic.Scan(ilang.Name)
		ic.Scan('{')
		
		if t, ok := ic.DefinedTypes[name]; ok {
			MethodType = t
		}
		
		var a = MethodType
		
		MethodType = ic.DefinedTypes[other]
		
		ic.LastDefinedType = MethodType
		
		var b = MethodType
		
		ilang.NewOperator(a, symbol, b, "SHARE %a\n SHARE %b\nRUN "+a.Name+"_"+symbol+"_"+b.Name+"\nGRAB %c", true)
		
		ic.Assembly("FUNCTION %s_%s_%s\n", a.Name, symbol, b.Name)
		ic.GainScope()
		ic.Assembly("GRAB b\nGRAB a\nARRAY c\n")
		for range a.Detail.Elements {
			ic.Assembly("PUT 0\n")
		}
		ic.InOperatorFunction = true
		
		ic.SetFlag(function.Flag)
	
		ic.SetVariable("c", a)
		ic.SetVariable("a", a)
		ic.SetVariable("b", b)
	}
}

