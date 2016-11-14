package main

func (t Type) IsSomething() Type {
	if t.Name == "Something" {
		return t
	} else {
		return Undefined
	}
}

func (ic *Compiler) IndexSomething(name string, cast string) string {
	switch cast {
		case "number", "letter":
			var test = ic.Tmp("test")
			ic.Assembly("PUSH 2")
			ic.Assembly("PLACE ", name)
			ic.Assembly("GET ", test)
			ic.Assembly("SEQ %v %v %v", test, test, string2type[cast].Int)
			ic.Assembly("IF ", test)
			
			var num = ic.Tmp("number")
			ic.Assembly("PUSH 0")
			ic.Assembly("GET ", test)
			ic.Assembly("PUSH ", test)
			
			ic.Assembly("ELSE")
			ic.Assembly("ERROR 404")
			ic.Assembly("PUSH 0")
			ic.Assembly("END")
			ic.Assembly("PULL ", num)
			ic.ExpressionType = Number
			return num
		case "text", "array":
			var test = ic.Tmp("test")
			ic.Assembly("PUSH 2")
			ic.Assembly("PLACE ", name)
			ic.Assembly("GET ", test)
			ic.Assembly("SEQ %v %v %v", test, test, string2type[cast].Int)
			ic.Assembly("IF ", test)
			
			var address = ic.Tmp("address")
			ic.Assembly("PUSH 0")
			ic.Assembly("PLACE ", name)
			ic.Assembly("GET ", address)
			ic.Assembly("PUSH ", address)
			ic.Assembly("HEAP")
			var txt = ic.Tmp("txt")
			
			ic.Assembly("ELSE")
			ic.Assembly("ERROR 404")
			ic.Assembly("ARRAY ", txt)
			ic.Assembly("SHARE ", txt)
			ic.Assembly("END")
			
			ic.Assembly("GRAB ", txt)
			ic.ExpressionType = Text
			return txt
			
		default:
			ic.RaiseError("Cannot cast something to ", cast)
	}
	return ""
}

func (ic *Compiler) AssignSomething(name string, value string) {
	var intf = ic.GetVariable(name).Interface
	if intf != nil {
		var originalname = ic.ExpressionType.Name
		for _, method := range intf.Methods {
			InheritMethodForInterface:
			if f, ok := ic.DefinedFunctions[method.Name+"_m_"+ic.ExpressionType.Name]; !ok {
				if ic.ExpressionType.Super != "" {
					ic.ExpressionType = ic.DefinedTypes[ic.ExpressionType.Super]
					
					ic.Library(`
FUNCTION `+method.Name+"_m_"+originalname+`
	RUN `+method.Name+"_m_"+ic.ExpressionType.Name+`
RETURN
					`)
					
					goto InheritMethodForInterface
				}
				ic.RaiseError("Invalid assignment, value of type ", originalname, " does not implement the method ", method.Name)
			} else {
				ic.LoadFunction(method.Name+"_m_"+ic.ExpressionType.Name)
				if f.Inline {
					ic.Library(`
FUNCTION `+method.Name+"_m_"+ic.ExpressionType.Name+`
`+f.Data+`
RETURN
					`)
				}
			}
		}
	}
	switch ic.ExpressionType {
		case Number, Letter:
			var tmp = ic.Tmp("number")
			ic.Assembly("ARRAY ", tmp)
			ic.Assembly("PUT ", value)
			ic.Assembly("PUT 0")
			ic.Assembly("PUT ", ic.ExpressionType.Int)
			for _, v := range ic.ExpressionType.Name {
				ic.Assembly("PUT ", byte(v))
			}
			ic.Assembly("SHARE ", name)
			ic.Assembly("RUN collect_m_Something")
			ic.Assembly("PLACE ", tmp)
			ic.Assembly("RENAME ", name)
		case Text, Array, ic.ExpressionType.IsUser():
			var tmp = ic.Tmp("text")
			ic.Assembly("ARRAY ", tmp)
			ic.Assembly("SHARE ", value)
			ic.Assembly("PUSH 0")
			ic.Assembly("HEAP")
			var address = ic.Tmp("address")
			ic.Assembly("PULL ", address)
			ic.Assembly("PUT ", address)
			ic.Assembly("PUT 1")
			ic.Assembly("PUT ", ic.ExpressionType.Int)
			for _, v := range ic.ExpressionType.Name {
				ic.Assembly("PUT ", byte(v))
			}
			ic.Assembly("SHARE ", name)
			ic.Assembly("RUN collect_m_Something")
			ic.Assembly("PLACE ", tmp)
			ic.Assembly("RENAME ", name)
		default:
		ic.RaiseError(ic.ExpressionType.Name, " is not a something")
	}
}
