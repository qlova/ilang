package c

import "github.com/qlova/ilang/src"


func init() {
	ilang.RegisterToken([]string{"const"}, ScanConst)
}

func ScanConst(ic *ilang.Compiler) {
	if len(ic.Scope) > 1 {
		ic.RaiseError("Constants must be declared outside of functions.")
	}
	
	var name = ic.Scan(ilang.Name)
	ic.Scan('=')
	var value = ic.ScanExpression()
	
	if ic.ExpressionType.Push != "PUSH" {
		ic.RaiseError("Constant must be a numerical value! (",ic.ExpressionType.Name,")")
	} 
	
	ic.Assembly(".const %v %v", name, value)
	ic.SetVariable(name, ic.ExpressionType)
}
