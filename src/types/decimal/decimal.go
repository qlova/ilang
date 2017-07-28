package decimal

import "github.com/qlova/ilang/src"
import (
	"strings"
	"fmt"
	"strconv"
	"math"
)

var Type = ilang.NewType("decimal", "PUSH", "PULL")

func ScanStatement(ic *ilang.Compiler) {
	ic.ScanNumericStatement()
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	return Type
}

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	
	//Decimal numbers.
	if strings.Contains(token, ".") {
		parts := strings.Split(token, ".")
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		ic.ExpressionType = Type
		return fmt.Sprint(a*1000000+b*int(math.Pow(10, 6-float64(len(parts[1])))))
	}
	return ""
}

func init() {
	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol(".", ScanSymbol)	
	
	
	ilang.NewOperator(Type, "=", Type, "VAR %c\nSEQ %c %a %b", true, ilang.Number)
	ilang.NewOperator(Type, "!=",Type, "VAR %c\nSNE %c %a %b", true, ilang.Number)
	ilang.NewOperator(Type, "<", Type, "VAR %c\nSLT %c %a %b", true, ilang.Number)
	ilang.NewOperator(Type, ">", Type, "VAR %c\nSGT %c %a %b", true, ilang.Number)
	ilang.NewOperator(Type, "<=",Type, "VAR %c\nSLE %c %a %b", true, ilang.Number)
	ilang.NewOperator(Type, ">=",Type, "VAR %c\nSGE %c %a %b", true, ilang.Number)
	ilang.NewOperator(Type, "mod", Type, "VAR %c\nMOD %c %a %b", true)
	ilang.NewOperator(Type, "^", ilang.Number, 
		"VAR %c\nPOW %c %a %b\nVAR %t\nPOW %t 100 %b\nMUL %t %t 1000000\nDIV %c %c %t", true, Type)
	ilang.NewOperator(Type, "+=", Type, "ADD %a %a %b", false, ilang.Undefined)
	ilang.NewOperator(Type, "-=", Type, "SUB %a %a %b", false, ilang.Undefined)
	ilang.NewOperator(Type, "*=", Type, "MUL %c %a %b\nDIV %c %c 1000000", false, ilang.Undefined)
	ilang.NewOperator(Type, "/=", Type, "VAR %t\nMUL %t %a 1000000\nDIV %c %t %b", false, ilang.Undefined)
	ilang.NewOperator(Type, "+", Type, "VAR %c\nADD %c %a %b", false)
	ilang.NewOperator(Type, "-", Type, "VAR %c\nSUB %c %a %b", false)
	ilang.NewOperator(Type, "/", Type, "VAR %t\nVAR %c\nMUL %t %a 1000000\nDIV %c %t %b", true)
	ilang.NewOperator(Type, "*", Type, "VAR %c\nMUL %c %a %b\nDIV %c %c 1000000", true)
	
	ilang.RegisterFunction("decimal", ilang.Method(Type, true, "PUSH 0"))
	ilang.RegisterExpression(ScanExpression)
}

