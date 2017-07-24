package v

import "github.com/qlova/ilang/src"

func init() {
	ilang.RegisterToken([]string{"var", "ver", "变量"}, ScanVar)
}

func ScanVar(ic *ilang.Compiler) {
	if len(ic.Scope) > 1 {
		name := ic.Scan(ilang.Name)
		token := ic.Scan(0)
		
		//Single assignment.
		if token == "=" {
			ic.AssembleVar(name, ic.ScanExpression())
		
		//Multiple assignment.
		} else if token == "," {
			
			var names = []string{name}
			for {
				token = ic.Scan(0)
				if token == "=" {
					break	
				}
				names = append(names, token)
			}
			
			for i, name := range names {
				ic.AssembleVar(name, ic.ScanExpression())
				if i < len(names)-1 { 
					ic.Scan(',')
				}
			}
		
		//Uninitialised variables are illegal.
		} else {
			ic.RaiseError("A variable should have a value assigned to it with an '=' sign.")
		}
	} else {
		ic.RaiseError("Global variables are not supported. Use a constant?")				
	}
}
