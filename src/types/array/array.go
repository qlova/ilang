package list

import "github.com/qlova/ilang/src"

//var Type = ilang.NewType("array", "SHARE", "GRAB")
var Type = ilang.Array

func ScanStatement(ic *ilang.Compiler) {
	var name = ic.Scan(ilang.Name)
	var token = ic.Scan(0)
	
	//TODO CLEAN THIS UP!
	switch token {
		case "-":
			ic.Scan('-')
			var value = ic.Tmp("cut")
			ic.Assembly("PLACE ", name)
			ic.Assembly("POP ", value)
			
		case "+":
			if ic.Peek() == "+" {
					//List++ 
					//We append to the list the zero value of it's SubType.
					ic.Scan('+')
					
					ic.Assembly("PLACE ", name)
					ic.Assembly("PUT 0")
					
			} else if ic.Peek() == "=" {
				ic.Scan('=')
				
				var value = ic.ScanExpression()
				if !ic.ExpressionType.Equals(ilang.Number) {
					ic.RaiseError("Only numbers can be added to arrays!")
				}
				
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUT ", value)
				
			} else {
					ic.RaiseError("Expecting ++ or += got +", ic.Peek())
			}
			
		case "[":
			
			index := ic.ScanExpression()
			ic.Scan(']')
			
			token = ic.Scan(0)
			if token == "=" {
				
				value := ic.ScanExpression()
				if !ic.ExpressionType.Equals(ilang.Number) {
					ic.RaiseError("Only numbers can be added to arrays!")
				}
				
				ic.Assembly("PUSH ", index)
				ic.Assembly("PLACE ", name)
				ic.Assembly("SET ", value)
			
			} else {
				ic.RaiseError("Unexpected ", token, " expecting =")
			}
		default:
				ic.RaiseError("Unexpected ", token, " expecting -, [ or +")
	}
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	ic.Scan(']')
	return Type
}

func Shunt(ic *ilang.Compiler, name string) string {
	if ic.ExpressionType.Name == "array" {
		index := ic.ScanExpression()
		ic.Scan(']')
		if ic.ExpressionType != ilang.Number {
			ic.RaiseError("Index must be a number.")
		}
		
		var value = ic.Tmp("index")
		ic.Assembly("PUSH ", index)
		ic.Assembly("PLACE ", name)
		ic.Assembly("GET ", value)
		
		return value
	}
	return ""
}

var Number = Type

func init() {
	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol("[", ScanSymbol)
	ilang.RegisterShunt("[", Shunt)
	
	ilang.RegisterFunction("array", ilang.Method(Type, true, "PUSH 0\nMAKE\n"))
}

