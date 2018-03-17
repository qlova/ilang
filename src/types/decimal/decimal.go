package decimal

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/types/rational"
import (
	"strings"
	"fmt"
	"strconv"
	"math/big"
)

var Type = ilang.Type{Name: "decimal", Push: "PUSH", Pop: "PULL", Super: "1000000"}

func ScanStatement(ic *ilang.Compiler) bool {
	var name = ic.GetVariable(ic.LastToken).Name
	if name ==  "decimal" || len(name) > len("decimal") &&  name[:len("decimal")] == "decimal" {
		ic.NextToken = ic.LastToken
		ic.ExpressionType = ic.GetVariable(ic.LastToken)
		ic.ScanNumericStatement()
		return true
	}
	return false
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	if len(ic.LastToken) > 0 {
		var precision, err = strconv.Atoi(ic.LastToken[1:])
		if err != nil {
			ic.RaiseError("Decimals can only have numeric precisions!")
		}

		return GenerateTypeFor(ic, precision)
	} 
	return Type
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func ScanShunt(ic *ilang.Compiler, token string) string {
	if len(ic.ExpressionType.Name) >= len("decimal") && 
			ic.ExpressionType.Name[:len("decimal")] == "decimal" {
			
		var index_string = ic.Scan(0)
		var index, err = strconv.Atoi(index_string)
		if err != nil {
			ic.RaiseError("Decimals can only be precision-indexed.")
		}
		var cast = "1"+strings.Repeat("0", index)
		var difference = abs(len(cast)-len(ic.ExpressionType.Super))
		
		var tmp = ic.Tmp("decimalcast")
		ic.Assembly("VAR ", tmp)
		
		if len(cast) < len(ic.ExpressionType.Super) {
			ic.Assembly("DIV ", tmp, " ", token, " 1"+strings.Repeat("0", difference))
		} else {
			ic.Assembly("MUL ", tmp, " ", token, " 1"+strings.Repeat("0", difference))
		}
		
		ic.ExpressionType = GenerateTypeFor(ic, index)
		
		return tmp
	}
	return ""
}

func SpecialOperator(token string, a, b ilang.Type) (operator *ilang.Operator) {

	if len(a.Name) >= len("decimal") && a.Name[:len("decimal")] == "decimal" &&
		len(b.Name) >= len("decimal") && b.Name[:len("decimal")] == "decimal" {
		
		operator = new(ilang.Operator)
		var difference = abs(len(a.Super)-len(b.Super))
		
		switch token {
			case "+":
				if len(a.Super) < len(b.Super) {
					operator.ExpressionType = a
					operator.Assembly = "VAR %c\nVAR %t\nDIV %t %b 1"+strings.Repeat("0", difference)+"\nADD %c %a %t"
				} else {
					operator.ExpressionType = b
					operator.Assembly = "VAR %c\nVAR %t\nDIV %t %a 1"+strings.Repeat("0", difference)+"\nADD %c %b %t"
				}
			case "-":
				if len(a.Super) < len(b.Super) {
					operator.ExpressionType = a
					operator.Assembly = "VAR %c\nVAR %t\nDIV %t %b 1"+strings.Repeat("0", difference)+"\nSUB %c %a %t"
				} else {
					operator.ExpressionType = b
					operator.Assembly = "VAR %c\nVAR %t\nDIV %t %a 1"+strings.Repeat("0", difference)+"\nSUB %c %b %t"
				}
			default:
				println("Unsupported mixed-decimal operation, this feature will be coming soon!")
		}
	}
	
	return
}

var GeneratedTypes = make(map[int]ilang.Type)

func GenerateTypeFor(ic *ilang.Compiler, precision int) ilang.Type {
	
	//Check if we have already generated this type or not.
	if t, ok := GeneratedTypes[precision]; ok {
		return t
	}
	
	var Copy = Type
	var token = "decimal"
	
	if precision != 6 {
		token += strconv.Itoa(precision)
		Copy.Name = token
		Copy.Super = "1"+strings.Repeat("0", precision)
	}
	
	ilang.NewOperator(Copy, "=", Copy, "VAR %c\nSEQ %c %a %b", true, ilang.Number)
	ilang.NewOperator(Copy, "!=",Copy, "VAR %c\nSNE %c %a %b", true, ilang.Number)
	ilang.NewOperator(Copy, "<", Copy, "VAR %c\nSLT %c %a %b", true, ilang.Number)
	ilang.NewOperator(Copy, ">", Copy, "VAR %c\nSGT %c %a %b", true, ilang.Number)
	ilang.NewOperator(Copy, "<=",Copy, "VAR %c\nSLE %c %a %b", true, ilang.Number)
	ilang.NewOperator(Copy, ">=",Copy, "VAR %c\nSGE %c %a %b", true, ilang.Number)
	ilang.NewOperator(Copy, "mod", Copy, "VAR %c\nMOD %c %a %b", true)
	ilang.NewOperator(Copy, "^", ilang.Number, 
		"VAR %c\nPOW %c %a %b\nVAR %t\nPOW %t 100 %b\nMUL %t %t "+Copy.Super+"\nDIV %c %c %t\n", true, Copy)
	ilang.NewOperator(Copy, "+=", Copy, "ADD %a %a %b", false, ilang.Undefined)
	ilang.NewOperator(Copy, "-=", Copy, "SUB %a %a %b", false, ilang.Undefined)

	ilang.NewOperator(Copy, "*=", Copy, "MUL %a %a %b\nDIV %a %a "+Copy.Super, false, ilang.Undefined)
	ilang.NewOperator(Copy, "/=", Copy, "VAR %t\nMUL %t %a "+Copy.Super+"\nDIV %a %t %b", false, ilang.Undefined)

	ilang.NewOperator(Copy, "+", Copy, "VAR %c\nADD %c %a %b", false)
	ilang.NewOperator(Copy, "-", Copy, "VAR %c\nSUB %c %a %b", false)
	ilang.NewOperator(Copy, "/", Copy, "VAR %t\nVAR %c\nMUL %t %a "+Copy.Super+"\nDIV %c %t %b", true)
	ilang.NewOperator(Copy, "*", Copy, "VAR %c\nMUL %c %a %b\nDIV %c %c "+Copy.Super+"", true)
	
	//We want Normal numbers to work with decimal numbers.
	ilang.NewOperator(Copy, "+=", ilang.Number, "VAR %t\nMUL %t %b "+Copy.Super+"\nADD %a %a %t", false, ilang.Undefined)
	ilang.NewOperator(Copy, "-=", ilang.Number, "VAR %t\nMUL %t %b "+Copy.Super+"\nSUB %a %a %t", false, ilang.Undefined)
	ilang.NewOperator(Copy, "*=", ilang.Number, "MUL %a %a %b", false, ilang.Undefined)
	ilang.NewOperator(Copy, "/=", ilang.Number, "DIV %a %a %b", false, ilang.Undefined)
	
	ilang.NewOperator(Copy, "+", ilang.Number, "VAR %c\nVAR %t\nMUL %t %b "+Copy.Super+"\nADD %c %a %t", false)
	ilang.NewOperator(Copy, "-", ilang.Number, "VAR %c\nVAR %t\nMUL %t %b "+Copy.Super+"\nSUB %c %a %t", false)
	ilang.NewOperator(Copy, "/", ilang.Number, "VAR %c\nDIV %t %a %b", true)
	ilang.NewOperator(Copy, "*", ilang.Number, "VAR %c\nMUL %c %a %b", true)
	
	ilang.NewOperator(ilang.Number, "+", Copy, "VAR %c\nVAR %t\nMUL %t %a "+Copy.Super+"\nADD %c %t %b", false, Copy)
	ilang.NewOperator(ilang.Number, "-", Copy, "VAR %c\nVAR %t\nMUL %t %a "+Copy.Super+"\nSUB %c %t %b", false, Copy)
	
	ilang.NewOperator(ilang.Number, "/", Copy, "VAR %t\nVAR %c\nMUL %t %a "+Copy.Super+"\nMUL %t %a "+Copy.Super+"\nDIV %c %t %b", true, Copy)
	ilang.NewOperator(ilang.Number, "*", Copy, "VAR %c\nMUL %c %a %b", true, Copy)

	f := ilang.Function{Exists:true, Import:"i_base_number", Returns:[]ilang.Type{ilang.Text}, Data:`
FUNCTION text_m_`+token+`
	PULL value
	VAR i_operator2
	MOD i_operator2 value `+Copy.Super+`
	
	VAR test
	VAR test2
	SLT test value 0
	IF test
		SUB i_operator2 `+Copy.Super+` i_operator2 
	END
	
	SGT test2 value -`+Copy.Super+`
	IF test
	IF test2
		ADD value 0 0
	END
	END
	
	PUSH i_operator2
	PUSH 10
	RUN i_base_number
	GRAB i_result3
	SHARE i_result3
	GRAB decimal
	VAR i_operator4
	SNE i_operator4 #decimal `+strconv.Itoa(precision)+`
	IF i_operator4
		IF 1
		VAR each
		VAR i_backup6
		ADD i_backup6 0 #decimal
		ADD each 0 #decimal
		LOOP
			VAR i_over5
			SNE i_over5 each `+strconv.Itoa(precision)+`
			ADD each 0 i_backup6
			IF i_over5
				SLT i_over5 each `+strconv.Itoa(precision)+`
				IF i_over5
					ADD i_backup6 each 1
				ELSE
					SUB i_backup6 each 1
				END
				SEQ i_over5 #decimal `+strconv.Itoa(precision)+`
				IF i_over5
					BREAK
				END
			ELSE
				SEQ i_over5 #decimal `+strconv.Itoa(precision)+`
				IF i_over5
					ADD each each 1
		       ELSE
					BREAK
				END
			END
		
			ARRAY i_string7
			PUT 48
			ARRAY i_operator8
			JOIN i_operator8 i_string7 decimal
			SHARE i_operator8
			RENAME decimal
		REPEAT
		END
	END
	VAR i_operator9
	SUB i_operator9 #decimal 1
	IF 1
	VAR i
	VAR i_backup11
	ADD i_backup11 0 i_operator9
	ADD i 0 i_operator9
	LOOP
		VAR i_over10
		SNE i_over10 i 0
		ADD i 0 i_backup11
		IF i_over10
			SLT i_over10 i 0
			IF i_over10
				ADD i_backup11 i 1
			ELSE
				SUB i_backup11 i 1
			END
			SEQ i_over10 i_operator9 0
			IF i_over10
				BREAK
			END
		ELSE
			SEQ i_over10 i_operator9 0
			IF i_over10
				ADD i i 1
	       ELSE
				BREAK
			END
		END
	
		PLACE decimal
		PUSH i
		GET i_index12
		PUSH 48
		
		PULL i_result14
		VAR i_operator13
		SEQ i_operator13 i_index12 i_result14
		IF i_operator13
			PLACE decimal
			POP i_tmp16
			ADD i_tmp16 0 0
		ELSE
			BREAK
		END
	REPEAT
	END
	VAR i_operator17
	SGT i_operator17 #decimal 0
	IF i_operator17
		ARRAY i_string18
		PUT 46
		ARRAY i_operator19
		JOIN i_operator19 i_string18 decimal
		SHARE i_operator19
		RENAME decimal
	END
	VAR i_operator20
	DIV i_operator20 value `+Copy.Super+`
	PUSH i_operator20
	PUSH 10
	RUN i_base_number
	GRAB i_result21
	IF test
	IF test2
		POP n
		ADD n 0 0
		PUT 45
		PUT 48
	END
	END
	ARRAY i_operator22
	JOIN i_operator22 i_result21 decimal
	SHARE i_operator22
RETURN
`}

f2 := ilang.Function{Exists:true, Returns:[]ilang.Type{Copy}, Data:`
FUNCTION `+token+`_m_rational
	GRAB r
	PLACE r
	PUSH 0
	GET a
	PUSH 1
	GET b
	
	MUL a a `+Copy.Super+`
	DIV a a b
	
	PUSH a
END
`}

f3 := ilang.Function{Exists:true, Returns:[]ilang.Type{Copy}, Data:`
FUNCTION `+token+`_m_number
	PULL a
	MUL a a `+Copy.Super+`
	PUSH a
END
`}

f4 := ilang.Function{Exists:true, Returns:[]ilang.Type{ilang.Number}, Data:`
FUNCTION number_m_`+token+`
	PULL a
	DIV a a `+Copy.Super+`
	PUSH a
END
`}

f5 := ilang.Function{Exists:true, Import: "rational", Returns:[]ilang.Type{rational.Type}, Data:`
FUNCTION rational_m_`+token+`
	ARRAY basic
	PUT 1
	PUT 1
	
	PULL decimal
	
	ARRAY r
	PUT decimal
	PUT `+Copy.Super+`
	
	SHARE basic
	SHARE r
	RUN rational_times_rational
END
`}
	
	if ic == nil {
		ilang.RegisterFunction("text_m_"+token, f)
		ilang.RegisterFunction(token+`_m_rational`, f2)
		ilang.RegisterFunction(token+`_m_number`, f3)
		ilang.RegisterFunction("number_m_"+token, f4)
		ilang.RegisterFunction("rational_m_"+token, f5)
	} else {
		ic.DefinedFunctions["text_m_"+token] = f
		ic.DefinedFunctions[token+`_m_rational`] = f2
		ic.DefinedFunctions[token+`_m_number`] = f3
		ic.DefinedFunctions["number_m_"+token] = f4
		ilang.RegisterFunction("rational_m_"+token, f5)
	}	

	
	GeneratedTypes[precision] = Copy
	return Copy
}

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	
	if len(token) > len("decimal") && token[:len("decimal")] == "decimal" {
		precision, err := strconv.Atoi(token[len("decimal"):])
		if err != nil {
			ic.RaiseError(err)
		}
		ic.Scan('(')
		ic.Scan(')')
		
		var tmp = ic.Tmp("decimal")
		ic.Assembly("VAR ", tmp)
		
		ic.ExpressionType = GenerateTypeFor(ic, precision)
		
		return tmp
	}
	
	//Decimal numbers. TODO parsing errors.
	if strings.Contains(token, ".") {
	
		var precision = "1000000"

		if len(ic.ExpressionType.Name) > len("decimal") && 
			ic.ExpressionType.Name[:len("decimal")] == "decimal" {
			
			precision = ic.ExpressionType.Super
		} else {
			ic.ExpressionType = Type
		}
	
		parts := strings.Split(token, ".")
		
		var result = big.NewInt(10)
		result.Exp(result, big.NewInt(int64((len(precision)-1)-len(parts[1]))), nil)
		
		var p = big.NewInt(0)
		p.SetString(precision, 10)
		
		var big_a = big.NewInt(0)
		big_a.SetString(parts[0], 10)
		big_a.Mul(big_a, p)
		
		var big_b = big.NewInt(0)
		big_b.SetString(parts[1], 10)
		big_b.Mul(big_b, result)
		
		big_b.Add(big_a, big_b)
		
		return fmt.Sprint(big_b)
	}
	return ""
}

func init() {
	ilang.RegisterDefault(ScanStatement)	
	ilang.RegisterSymbol(".", ScanSymbol)
	ilang.RegisterShunt(".", ScanShunt)	
	
	ilang.RegisterFunction("decimal", ilang.Method(GenerateTypeFor(nil, 6), true, "PUSH 0"))
	
	ilang.RegisterExpression(ScanExpression)
	ilang.RegisterSpecialOperator(SpecialOperator)
}

