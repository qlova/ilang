package collection

import "github.com/qlova/ilang/src"

var Type = ilang.NewType("collection", "PUSH", "PULL")

func init() {
	ilang.RegisterStatement(Type, ScanStatement)
	ilang.RegisterSymbol("<", ScanSymbol)
	ilang.RegisterExpression(ScanExpression)
	
	ilang.RegisterFunction("collection", ilang.Method(Type, true, "PUSH 1"))
	
	//Set Operations.
	ilang.NewOperator(Type, "=", Type, "VAR %c\nSEQ %c %a %b", true, ilang.Number)
	ilang.NewOperator(Type, "+", Type, "VAR %c\nMUL %c %a %b", false)
	ilang.NewOperator(Type, "-", Type, "VAR %c\nDIV %c %a %b", false)
	ilang.NewOperator(Type, "+=", Type, "MUL %a %a %b", false, ilang.Undefined)
	ilang.NewOperator(Type, "-=", Type, "DIV %a %a %b", false, ilang.Undefined)
	ilang.NewOperator(Type, "<=", Type, "VAR %c\nMOD %c %b %a\nDIV %c %c 0", true, ilang.Number)
	ilang.NewOperator(Type, ">", Type, "VAR %c\nMOD %c %b %a\n", true, ilang.Number)
	ilang.NewOperator(Type, ">=", Type, "VAR %c\nMOD %c %b %a\nVAR %t\nSEQ %t %a %b\nADD%c %c %t", true, ilang.Number)
	ilang.NewOperator(Type, "<", Type, "VAR %c\nMOD %c %b %a\nDIV %c %c 0\nVAR %t\nSNE %t %a %b\nMUL %c %c %t", true, ilang.Number)
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	ic.Scan('>')
	return Type
}

func ScanExpression(ic *ilang.Compiler) string {
	if ic.LastToken == "<" {
		var id = ic.Tmp("collection")

		ic.Assembly("VAR ", id)
		ic.Assembly("ADD ", id, " 1 0")

		for {
			var token = ic.Scan(0)
			if token == ">" {
				break
			}
	
			if prime, ok := ic.SetItems[token]; ok {
				ic.Assembly("MUL %s %s %v", id, id, prime)
			} else {
				ic.SetItems[token] = Primes[ic.SetItemCount]
				ic.Assembly("MUL %s %s %v", id, id, Primes[ic.SetItemCount])
		
				ic.SetItemCount++
			}
	
			token = ic.Scan(0)
			if token != "," {
				if token == ">" {
					break
				} else {
					ic.RaiseError("expecting >")
				}
			}
		}

		ic.ExpressionType = Type

		return id
	}
	return ""
}

/*
	Scan a pipe statement, eg.
		file("text to write")
		file = newfile
*/
func ScanStatement(ic *ilang.Compiler) {
	ic.ScanNumericStatement()
}
