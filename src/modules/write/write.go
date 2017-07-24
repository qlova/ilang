package p

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/modules/method"

func init() {
	ilang.RegisterToken([]string{
		"write", 		//English 
	}, Scan)
}

func Scan(ic *ilang.Compiler) {
	ic.Scan('(')
	arg := ic.ScanExpression()
	if !ic.ExpressionType.Empty() {
		ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
	}
	method.Call(ic, "text", ic.ExpressionType)
	ic.Assembly("STDOUT")
	
	for {
		token := ic.Scan(0)
		if token != "," {
			if token != ")" {
				ic.RaiseError()
			}
			break
		}
		arg := ic.ScanExpression()
		if !ic.ExpressionType.Empty() {
			ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
		}
		method.Call(ic, "text", ic.ExpressionType)
		ic.Assembly("STDOUT")
	}
}
