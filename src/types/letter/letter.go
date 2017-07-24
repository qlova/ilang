package letter

import "github.com/qlova/ilang/src"
import "strconv"

var Type = ilang.NewType("letter", "PUSH", "PULL")

func ScanStatement(ic *ilang.Compiler) {
	ic.ScanNumericStatement()
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	return Type
}

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	//Types.
	if token[0] == "'"[0] {
		if s, err := strconv.Unquote(token); err == nil {
			ic.ExpressionType = Type
			return strconv.Itoa(int([]byte(s)[0]))
		} else {
			ic.RaiseError(err)
		}
	}
	return ""
}


func init() {
	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol("' '", ScanSymbol)
	ilang.RegisterExpression(ScanExpression)
	
	ilang.RegisterFunction("letter", ilang.Method(Type, true, "PUSH 0"))
	ilang.RegisterFunction("letter_m_number", ilang.BlankMethod(Type))
	
	
	ilang.NewOperator(Type, "=", Type, "VAR %c\nSEQ %c %a %b", true, ilang.Number)
	ilang.NewOperator(Type, "!=",Type, "VAR %c\nSNE %c %a %b", true, ilang.Number)
	ilang.NewOperator(ilang.Text, "+=", Type, "PLACE %a\nPUT %b", false, ilang.Undefined)
}

