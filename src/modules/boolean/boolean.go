package boolean

import "github.com/qlova/ilang/src"

func init() {
	ilang.RegisterToken([]string{"boolean"}, func(ic *ilang.Compiler) {
		var positive = ic.Scan(0)
		ic.Scan('/')
		var negative = ic.Scan(0)
		
		ic.SetVariable(positive, ilang.Number)
		ic.SetVariable(positive+"_bool", ilang.Type{Name: "i_boolean", Push: positive, Pop: negative})
		ic.SetVariable(negative+"_bool", ilang.Type{Name: "i_boolean", Push: positive, Pop: negative})
		
		var next = ic.Scan(0)
		if next == "=" {
			
			var expr = ic.ScanExpression()
			if ic.ExpressionType != ilang.Number {
				ic.RaiseError("Only numbers can be assigned to boolean variables")
			}
			ic.Assembly("PUSH ", expr)
			ic.Assembly("PULL ", positive)
			
		} else if (next != "\n" || next != ";") {
			ic.NextToken = next
			ic.Assembly("VAR ", positive)
			return
		}
	})
	
	ilang.RegisterExpression(func(ic *ilang.Compiler) string {
		if ic.GetVariable(ic.LastToken+"_bool").Name == "i_boolean" {
			var negate = ic.Tmp("negate")
			ic.Assembly("VAR ", negate)
			ic.Assembly("SEQ ", negate, " ", ic.GetVariable(ic.LastToken+"_bool").Push, " 0")
			
			ic.ExpressionType = ilang.Number
			return negate
		}
		return ""
	})
	
	ilang.RegisterDefault(func(ic *ilang.Compiler) bool{
		
		if ic.GetVariable(ic.LastToken+"_bool").Name == "i_boolean" {
			var name = ic.LastToken
			ic.Scan('(')
			ic.Scan(')')
			
			if ic.GetVariable(name+"_bool").Push == name {
				ic.Assembly("ADD ", ic.GetVariable(name+"_bool").Push, " 0 1")
			} else if ic.GetVariable(name+"_bool").Pop == name {
				ic.Assembly("ADD ", ic.GetVariable(name+"_bool").Push, " 0 0")
			} else {
				ic.RaiseError("This is a compiler bug, please report this. (boolean module) ", name)
			}
			return true
		}
		
		return false
	})
}
