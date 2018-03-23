package ilang

import "strings"
import "fmt"

type Function struct {
	Name string

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
	
	UnpackArguments string
}

var ExportedInlineFunctions = make(map[string]bool)

func (ic *Compiler) ExportInlineFunction(name string, f Function) {
	if ExportedInlineFunctions[name] {
		return
	}
	
	ic.Library("FUNCTION ", name)
	ic.Library(f.Data)
	ic.Library("RETURN")
	
	ExportedInlineFunctions[name] = true
}

func (ic *Compiler) FunctionExists(name string) bool {
	_, ok := ic.DefinedFunctions[name]
	return ok
}

func (ic *Compiler) LoadFunction(name string) {
	
	//Maybe the function needs to be generated?
	for _, builder := range FunctionBuilders {
		var f = builder(name)
		if f != nil {
			ic.DefinedFunctions[name] = *f
			return
		}
	}
	
	f, ok := ic.DefinedFunctions[name]
	if !ok {
		ic.RaiseError(name, " does not exist!")
	}
	if !f.Inline && !f.Loaded {
		fmt.Fprintf(ic.Lib, f.Data)
		f.Loaded = true
		ic.DefinedFunctions[name] = f
	}
	if f.Import != "" {
		ic.LoadFunction(f.Import)
	}
}

func (ic *Compiler) RunFunction(name string) string {
	f, ok := ic.DefinedFunctions[name]
	if !ok {
		if strings.Contains(name, "collect_m_") {
			return "RUN "+name
		}
		if strings.Contains(name, "_flag_") {
			panic("Serious bug! Cannot create function from flag.")
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
	
	if len(f.Args) > 0 {
		
		for i := range f.Args {
			arg := ic.ScanExpression()
			
			
			if ! f.Args[i].Equals(ic.ExpressionType) {
			
				//Try converting the argument.
				if ic.CanCast(ic.ExpressionType, f.Args[i]) {
					
					ic.Assembly(ic.Cast(arg, ic.ExpressionType, f.Args[i]))
					
				} else {
					
					if ic.ExpressionType == Undefined {
						ic.RaiseError(ic.LastToken, " is undefined!")
					}
				
					//Hacky varidic lists!
					if i == len(f.Args)-1 && f.Args[i].SubType != nil && *f.Args[i].SubType == ic.ExpressionType {
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
						f.Args[i].GetComplexName(),", got ",ic.ExpressionType.GetComplexName()) 
				}

			}
			
			ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
			
			if i < len(f.Args)-1 {
				token := ic.Scan(0)
				if token != "," {
					ic.RaiseError()
				}
			}
		}
		
	} else {
		
		//Calls methods such as:
		//(Similar to function overloading!)
		/*
			read(Type) read(AnotherType)
		*/
		if f.Method && ic.Peek() != ")" {
			ic.DisableOwnership = true
				arg := ic.ScanExpression()
			ic.DisableOwnership = false
			
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
