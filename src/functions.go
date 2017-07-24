package ilang

import "strings"

type Function struct {
	Exists bool
	Loaded bool
	Import string
	
	Inline bool
	Data string
	Assemble func(*Compiler) string
	
	List bool //Does the method operate on a list?
	
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
		if strings.Contains(name, "collect_m_") {
			return "RUN "+name
		}
		if strings.Contains(name, "_flag_") {
			ic.RaiseError("Serious bug! Cannot create function from flag.")
		}
		ic.RaiseError(name, " does not exist!")
	}
	
	ic.LoadFunction(name)
	
	if f.Import != "" {
		ic.LoadFunction(f.Import)
	}
	
	if f.Inline {
		if f.Assemble != nil {
			return f.Assemble(ic)
		}
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
				
				//Hacky varidic lists!
				if i == len(f.Args)-1 && *f.Args[i].SubType == ic.ExpressionType {
					var tmp = ic.Tmp("varaidic")
					ic.Assembly("ARRAY ", tmp)
					ic.Assembly("PUT ", ic.GetPointerTo(arg))
					
					for {
						var token = ic.Scan(0)
						if token != "," {
							if token == ")" {
								ic.NextToken = ")"
								break
							}
							ic.RaiseError("Expecting , or )")
						}
						var value = ic.ScanExpression()
						
						if *f.Args[i].SubType != ic.ExpressionType {
							ic.RaiseError("Type mismatch! Variadic arguments of '"+name+"()' expect ",
								f.Args[i].SubType.Name,", got ",ic.ExpressionType.Name) 
						}
						
						ic.Assembly("PLACE ", tmp)
						ic.Assembly("PUT ", ic.GetPointerTo(value))
					}
					
					ic.Assembly("SHARE %v", tmp)
					
					break
					
				}
			
				ic.RaiseError("Type mismatch! Argument ",i+1," of '"+name+"()' expects ",
					f.Args[i].Name,", got ",ic.ExpressionType.Name) 
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
			/*if ic.ExpressionType == Text.MakeList() {
				ic.ExpressionType.Name = "textarray"
			}*/
			
			if _, ok := ic.DefinedFunctions[name+"_m_"+ic.ExpressionType.GetComplexName()]; !ok {
				if ic.ExpressionType.Super != "" {
					ic.ExpressionType = ic.DefinedTypes[ic.ExpressionType.Super]
					goto InheritMethods
				}
				ic.RaiseError("Method ",name," for type ",ic.ExpressionType.GetComplexName(), "does not exist!")
			}
			
			//Only pass the argument if it has a value, for example, the following type would not be passed:
			// type Blank {}
			if ic.ExpressionType.Detail == nil || !ic.ExpressionType.Empty() {
				ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
			}
			f = ic.DefinedFunctions[name+"_m_"+ic.ExpressionType.GetComplexName()]
			name = name+"_m_"+ic.ExpressionType.GetComplexName()
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
