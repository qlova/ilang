package complex

import "github.com/qlova/ilang/src"
import "strconv"

//TODO optimise this so it's not a pointer type.
var Type = ilang.NewType("complex", "SHARE", "GRAB")

func ScanStatement(ic *ilang.Compiler) {
	var name = ic.Scan(ilang.Name)
	var token = ic.Scan(0)
	
	switch token {
		case "=":
			value := ic.ScanExpression()
			if ic.ExpressionType != Type {
				ic.RaiseError("Type mismatch! Type '%s' cannot be assigned to %s which is a complex number.", ic.ExpressionType, name)
			}
			
			ic.Assembly("SHARE ", value)
			ic.Assembly("RENAME ", name)
			
		default:
			ic.ExpressionType = Type
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != ilang.Undefined {
				ic.RaiseError("blank expression!")
			}
			ic.ExpressionType = Type
	}
	
	ic.LoadFunction("complex")
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	return Type
}

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	
	if _, err := strconv.Atoi(token); err == nil{
		if ic.Peek() == "i" {
			ic.Scan(0)
			
			ic.LoadFunction("complex")
			
			ic.ExpressionType = Type
			var tmp = ic.Tmp("imaginary")
			ic.Assembly("ARRAY ", tmp)
			ic.Assembly("PUT 0")
			ic.Assembly("PUT ", token)
			return tmp
			
		}
	}
	return ""
}

func init() {
	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol("%", ScanSymbol)	
	ilang.RegisterExpression(ScanExpression)
	
	ilang.NewOperator(Type, "+", Type, "SHARE %a\nSHARE %b\nRUN complex_plus_complex\nGRAB %c", false)
	ilang.NewOperator(Type, "-", Type, "SHARE %a\nSHARE %b\nRUN complex_minus_complex\nGRAB %c", false)
	ilang.NewOperator(Type, "*", Type, "SHARE %a\nSHARE %b\nRUN complex_times_complex\nGRAB %c", false)
	ilang.NewOperator(Type, "/", Type, "SHARE %a\nSHARE %b\nRUN complex_div_complex\nGRAB %c", false)
	
	ilang.NewOperator(Type, "+", ilang.Number, "SHARE %a\nPUSH %b\nRUN complex_m_number\nRUN complex_plus_complex\nGRAB %c", false)
	ilang.NewOperator(Type, "-", ilang.Number, "SHARE %a\nPUSH %b\nRUN complex_m_number\nRUN complex_minus_complex\nGRAB %c", false)
	ilang.NewOperator(Type, "*", ilang.Number, "SHARE %a\nPUSH %b\nRUN complex_m_number\nRUN complex_times_complex\nGRAB %c", false)
	ilang.NewOperator(Type, "/", ilang.Number, "SHARE %a\nPUSH %b\nRUN complex_m_number\nRUN complex_div_complex\nGRAB %c", false)
	
	ilang.NewOperator(ilang.Number, "+", Type, "PUSH %a\nRUN complex_m_number\nSHARE %b\nRUN complex_plus_complex\nGRAB %c", false, Type)
	ilang.NewOperator(ilang.Number, "-", Type, "PUSH %a\nRUN complex_m_number\nSHARE %b\nRUN complex_minus_complex\nGRAB %c", false, Type)
	ilang.NewOperator(ilang.Number, "*", Type, "PUSH %a\nRUN complex_m_number\nSHARE %b\nRUN complex_times_complex\nGRAB %c", false, Type)
	ilang.NewOperator(ilang.Number, "/", Type, "PUSH %a\nRUN complex_m_number\nSHARE %b\nRUN complex_div_complex\nGRAB %c", false, Type)
	
	ilang.NewOperator(Type, "+=", Type, "SHARE %a\nSHARE %b\nRUN complex_plus_complex\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "-=", Type, "SHARE %a\nSHARE %b\nRUN complex_minus_complex\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "*=", Type, "SHARE %a\nSHARE %b\nRUN complex_times_complex\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "/=", Type, "SHARE %a\nSHARE %b\nRUN complex_div_complex\nRENAME %a", false, ilang.Undefined)
	
	ilang.NewOperator(Type, "+=", ilang.Number, "SHARE %a\nPUSH %b\nRUN complex_m_number\nRUN complex_plus_complex\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "-=", ilang.Number, "SHARE %a\nPUSH %b\nRUN complex_m_number\nRUN complex_minus_complex\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "*=", ilang.Number, "SHARE %a\nPUSH %b\nRUN complex_m_number\nRUN complex_times_complex\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "/=", ilang.Number, "SHARE %a\nPUSH %b\nRUN complex_m_number\nRUN complex_div_complex\nRENAME %a", false, ilang.Undefined)
	
	//ilang.NewOperator(ilang.Number, "\\", ilang.Number, "ARRAY %c\nPUT %a\nPUT %b", false, Type)
	
	ilang.RegisterFunction("number_m_complex", ilang.Function{Exists:true, 
		Returns:[]ilang.Type{ilang.Number}, 
		Args:[]ilang.Type{Type}, 
		Data:`
FUNCTION number_m_complex
	GRAB c
	PLACE c
	PUSH 0
	GET a
	PUSH a
RETURN
	`})
	
	ilang.RegisterFunction("complex_m_number", ilang.Function{Exists:true, 
		Returns:[]ilang.Type{Type}, 
		Args:[]ilang.Type{ilang.Number}, 
		Data:``})
	
	ilang.RegisterFunction("text_m_complex", ilang.Function{Exists:true, Returns:[]ilang.Type{ilang.Text}})
	ilang.RegisterFunction("complex", ilang.Function{Exists:true, Method: true, Returns:[]ilang.Type{Type}, Args:[]ilang.Type{}, Data:`
FUNCTION complex
	PUSH 2
	MAKE
RETURN

FUNCTION complex_m_number
	PULL n
	ARRAY c
	PUT n
	PUT 0
	SHARE c
RETURN

FUNCTION complex_times_complex

	GRAB b
	GRAB a
	ARRAY c
	
	PUT 0
	
	PUT 0
	
	PLACE a
	PUSH 0
	GET i_index1
	PLACE b
	PUSH 0
	GET i_index3
	VAR i_operator2
	MUL i_operator2 i_index1 i_index3
	PLACE a
	PUSH 1
	GET i_index5
	PLACE b
	PUSH 1
	GET i_index7
	VAR i_operator6
	MUL i_operator6 i_index5 i_index7
	VAR i_operator4
	SUB i_operator4 i_operator2 i_operator6
	PLACE c
	PUSH 0
	SET i_operator4
	PLACE a
	PUSH 0
	GET i_index8
	PLACE b
	PUSH 1
	GET i_index10
	VAR i_operator9
	MUL i_operator9 i_index8 i_index10
	PLACE b
	PUSH 0
	GET i_index12
	PLACE a
	PUSH 1
	GET i_index14
	VAR i_operator13
	MUL i_operator13 i_index12 i_index14
	VAR i_operator11
	ADD i_operator11 i_operator9 i_operator13
	PLACE c
	PUSH 1
	SET i_operator11
SHARE c
RETURN
FUNCTION complex_plus_complex

	GRAB b
	GRAB a
	ARRAY c
	
	PUT 0
	
	PUT 0
	
	PLACE a
	PUSH 0
	GET i_index15
	PLACE b
	PUSH 0
	GET i_index17
	VAR i_operator16
	ADD i_operator16 i_index15 i_index17
	PLACE c
	PUSH 0
	SET i_operator16
	PLACE a
	PUSH 1
	GET i_index18
	PLACE b
	PUSH 1
	GET i_index20
	VAR i_operator19
	ADD i_operator19 i_index18 i_index20
	PLACE c
	PUSH 1
	SET i_operator19
SHARE c
RETURN
FUNCTION complex_minus_complex

	GRAB b
	GRAB a
	ARRAY c
	
	PUT 0
	
	PUT 0
	
	PLACE a
	PUSH 0
	GET i_index21
	PLACE b
	PUSH 0
	GET i_index23
	VAR i_operator22
	SUB i_operator22 i_index21 i_index23
	PLACE c
	PUSH 0
	SET i_operator22
	PLACE a
	PUSH 1
	GET i_index24
	PLACE b
	PUSH 1
	GET i_index26
	VAR i_operator25
	SUB i_operator25 i_index24 i_index26
	PLACE c
	PUSH 1
	SET i_operator25
SHARE c
RETURN
FUNCTION complex_div_complex

	GRAB b
	GRAB a
	ARRAY c
	
	PUT 0
	
	PUT 0
	
	PLACE b
	PUSH 0
	GET i_index27
	VAR i_operator28
	MUL i_operator28 i_index27 i_index27
	PLACE b
	PUSH 1
	GET i_index30
	VAR i_operator31
	MUL i_operator31 i_index30 i_index30
	VAR i_operator29
	ADD i_operator29 i_operator28 i_operator31
	PUSH i_operator29
	PULL d
	PLACE a
	PUSH 0
	GET i_index32
	PLACE b
	PUSH 0
	GET i_index34
	VAR i_operator33
	MUL i_operator33 i_index32 i_index34
	PLACE a
	PUSH 1
	GET i_index36
	PLACE b
	PUSH 1
	GET i_index38
	VAR i_operator37
	MUL i_operator37 i_index36 i_index38
	VAR i_operator35
	ADD i_operator35 i_operator33 i_operator37
	VAR i_operator39
	DIV i_operator39 i_operator35 d
	PLACE c
	PUSH 0
	SET i_operator39
	PLACE a
	PUSH 1
	GET i_index40
	PLACE b
	PUSH 0
	GET i_index42
	VAR i_operator41
	MUL i_operator41 i_index40 i_index42
	PLACE a
	PUSH 0
	GET i_index44
	PLACE b
	PUSH 1
	GET i_index46
	VAR i_operator45
	MUL i_operator45 i_index44 i_index46
	VAR i_operator43
	SUB i_operator43 i_operator41 i_operator45
	VAR i_operator47
	DIV i_operator47 i_operator43 d
	PLACE c
	PUSH 1
	SET i_operator47
SHARE c
RETURN
FUNCTION text_m_complex
	GRAB Complex
	PLACE Complex
	PUSH 0
	GET i_index48
	PUSH i_index48
	PULL real
	PUSH real
	PUSH 10
	RUN i_base_number
	GRAB i_result49
	ARRAY i_string51
	PUT 32
	PUT 43
	PUT 32
	PLACE Complex
	PUSH 1
	GET i_index53
	PUSH i_index53
	PULL imag
	PUSH imag
	PUSH 10
	RUN i_base_number
	GRAB i_result54
	ARRAY i_string56
	PUT 105
	ARRAY i_operator55
	JOIN i_operator55 i_result54 i_string56
	ARRAY i_operator52
	JOIN i_operator52 i_string51 i_operator55
	ARRAY i_operator50
	JOIN i_operator50 i_result49 i_operator52
	SHARE i_operator50
RETURN
`})
}

