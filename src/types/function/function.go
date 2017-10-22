package function

import "github.com/qlova/ilang/src"

var Type = ilang.NewType("function", "RELAY", "TAKE")

func init() {
	ilang.RegisterDefault(ScanFuncStatement)
	ilang.RegisterSymbol("(", ScanFuncSymbol)
	ilang.RegisterExpression(FuncExpression)
}

func ScanFuncSymbol(ic *ilang.Compiler) ilang.Type {
	var t = Type
	
	var token = ic.Scan(0)
	if token != ")" {
		ic.NextToken = token
		t.Detail = new(ilang.UserType)
		
		var arg = ic.ScanSymbolicType()
		if arg == ilang.Undefined {
			ic.RaiseError()
		}
		if arg == ilang.Number {
			ic.Scan(0)
		}
		t.Detail.Elements = append(t.Detail.Elements, arg)
		
		for {
			var token = ic.Scan(0)
			if token == "," {
				
				var arg = ic.ScanSymbolicType()
				if arg == ilang.Undefined {
					ic.RaiseError()
				}
				if arg == ilang.Number {
					ic.Scan(0)
				}
				t.Detail.Elements = append(t.Detail.Elements, arg)
				
			} else if token == ")" {
				break
			} else {
				ic.RaiseError("Expecting , or )")
			}
		}
		
	}

	return t
}

func FuncExpression(ic *ilang.Compiler) string {
	token := ic.LastToken
	
	if _, ok := ic.DefinedFunctions[token]; ok {
	
		//TODO don't peek.	
		if ic.Peek() != "(" {
			ic.ExpressionType = Type
			var id = ic.Tmp("func")
			ic.Assembly("SCOPE ", token)
			ic.Assembly("TAKE ", id)
		
			return id
		}
		
	}
	return ""
}

/*
	Scan a function pipe statement.
		function()
		function = newfunction
*/
func ScanFuncStatement(ic *ilang.Compiler) bool {

	//TODO make names containing an underscore illegal.
	var name = ic.LastToken
	
	if ic.GetVariable(name).Name != "function" {
		return false
	}
	
	var t = ic.GetVariable(name)
	
	var token = ic.Scan(0)
	
	switch token {
		case "(":
			if t.Detail != nil && len(t.Detail.Elements) >  0 {
				for i := 0; i < len(t.Detail.Elements); i++ {
					var expr = ic.ScanExpression()
					if ic.ExpressionType != t.Detail.Elements[i] {
						ic.RaiseError("Type mismatch! argument ", i+1, " of ", name, " expects ", 
						t.Detail.Elements[i].GetComplexName(), ", instead got ", ic.ExpressionType.GetComplexName())
					}
					ic.Assembly(ic.ExpressionType.Push, " ", expr)
					if i < len(t.Detail.Elements)-1 {
						ic.Scan(',')
					}
				}
			}
			ic.Scan(')')
			ic.Assembly("EXE ", name)
		case "=":
			value := ic.ScanExpression()
			if ic.ExpressionType != Type {
				ic.RaiseError("Only ",Type.Name," values can be assigned to ",name,".")
			}
			ic.Assembly("PLACE ", value)
			ic.Assembly("RELOAD ", name)
		default:
			ic.ExpressionType = Type
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != ilang.Undefined {
				ic.RaiseError("blank expression!")
			}
	}
	return true
}
