package function

import "github.com/qlova/ilang/src"

var Type = ilang.NewType("function", "RELAY", "TAKE")

func init() {
	ilang.RegisterStatement(Type, ScanFuncStatement)
	ilang.RegisterSymbol("(", ScanFuncSymbol)
	ilang.RegisterExpression(FuncExpression)
}

func ScanFuncSymbol(ic *ilang.Compiler) ilang.Type {
	ic.Scan(')')
	return Type
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
func ScanFuncStatement(ic *ilang.Compiler) {

	//TODO make names containing an underscore illegal.
	var name = ic.Scan(ilang.Name)
	var token = ic.Scan(0)
	
	switch token {
		case "(":
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
}
