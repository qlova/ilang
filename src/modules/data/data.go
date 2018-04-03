package data

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/types/list"
import "strconv"

func init() {
	ilang.RegisterToken([]string{"data"}, func(ic *ilang.Compiler) {
		var name = ic.Scan(ilang.Name)
		ic.Scan('=')
		ic.Scan('[')
		
		var t ilang.Type = ilang.Undefined
		
		var asm = "DATA "+name
		for {
			var value = ic.ScanExpression()
			if _, err := strconv.Atoi(string(value[0])); err != nil && value[0] != '-' {
				ic.RaiseError("Data values must be numeric types, ",value, " is not numeric!")
			}
			
			if t == ilang.Undefined {
				t = ic.ExpressionType
			}
			
			if !t.Equals(ic.ExpressionType) {
				ic.RaiseError("Inconsistent data value types! ", t.GetComplexName() ," and ", ic.ExpressionType.GetComplexName())
			}
			
			asm += " "+value
			
			var symbol = ic.Scan(0)
			if symbol == "]" {
				break
			} else if symbol != "," {
				ic.RaiseError("Expecting ','")
			}
		}
		ic.Assembly(asm)
		if t.Equals(ilang.Number) {
			ic.SetVariable(name, ilang.Array)
		} else {
			ic.SetVariable(name, list.Of(t))
		}
	})
}
