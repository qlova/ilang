package ilang

type InterfaceMethod struct {
	Name string
	Returns []Type
}

//This is an interface, for example:
/*
	interface Printable {
		text() ""
	}
*/
type Interface struct {
	Name string
	Methods []InterfaceMethod
}

func (i Interface) GetType() Type {
	var s = Something
	s.Interface = &i
	return s
}

func (ic *Compiler) CallInterfaceMethod(name string) string {
	var intf = ic.ExpressionType.Interface
	
	if _, ok := ic.DefinedFunctions[name+"_m_"+intf.Name]; !ok {
		ic.LoadFunction("text_m_Something")
		ic.DefinedFunctions[name+"_m_"+intf.Name] = Function{Exists:true}
		ic.Library(`
DATA `+name+`_string "`+name+`_m_"
FUNCTION `+name+"_m_"+intf.Name+`
	GRAB something
	SHARE something
	RUN text_m_Something
	GRAB method
	
	PUSH 1
	PLACE something
	GET garbage
	IF garbage
		PUSH 0
		GET address
		SUB garbage garbage 1
		IF garbage
			SUB garbage garbage 1
			IF garbage
				PUSH address
				HEAP
			ELSE
				PUSH address
				#HEAPIT
			END
		ELSE
			PUSH address
			HEAP
		END
	ELSE
		PUSH 0
		GET value
		PUSH value
	END
	
	JOIN method `+name+`_string method
	SHARE method
	EVAL
RETURN
		`)
	}
	
	return "RUN "+name+"_m_"+intf.Name
}

func (ic *Compiler) ScanInterface() {
	var name = ic.Scan(Name)
	ic.Scan('{')
	ic.Scan('\n')
	
	var intf = Interface{Name:name}
	
	for {
		var method = ic.Scan(Name)
		
		if method == "}" {
			break
		}
		
		ic.Scan('(')
		ic.Scan(')')
		var r = ic.ScanSymbolicType()
		intf.Methods = append(intf.Methods, InterfaceMethod{Name:method, Returns:[]Type{r}})
		
		var token = ic.Scan(0)
		if token == "}" {
			break
		} else if token != "," && token != "\n" {
			ic.RaiseError()
		}
	}
	
	ic.DefinedInterfaces[name] = intf
}
