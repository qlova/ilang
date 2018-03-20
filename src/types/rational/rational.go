package rational

import "github.com/qlova/ilang/src"

//TODO optimise this so it's not a pointer type.
var Type = ilang.NewType("rational", "SHARE", "GRAB")

func ScanStatement(ic *ilang.Compiler) {
	var name = ic.Scan(ilang.Name)
	var token = ic.Scan(0)
	
	switch token {
		case "=":
			value := ic.ScanExpression()
			if ic.ExpressionType != Type {
				ic.RaiseError("Type mismatch! Type '%s' cannot be assigned to %s which is a rational number.", ic.ExpressionType, name)
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
	
	ic.LoadFunction("rational")
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	return Type
}

func init() {
	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol("\\", ScanSymbol)	
	ilang.NewOperator(Type, "+", Type, "SHARE %a\nSHARE %b\nRUN rational_plus_rational\nGRAB %c", false)
	ilang.NewOperator(Type, "-", Type, "SHARE %a\nSHARE %b\nRUN rational_minus_rational\nGRAB %c", false)
	ilang.NewOperator(Type, "*", Type, "SHARE %a\nSHARE %b\nRUN rational_times_rational\nGRAB %c", false)
	ilang.NewOperator(Type, "/", Type, "SHARE %a\nSHARE %b\nRUN rational_div_rational\nGRAB %c", false)
	
	ilang.NewOperator(Type, "+=", Type, "SHARE %a\nSHARE %b\nRUN rational_plus_rational\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "-=", Type, "SHARE %a\nSHARE %b\nRUN rational_minus_rational\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "*=", Type, "SHARE %a\nSHARE %b\nRUN rational_times_rational\nRENAME %a", false, ilang.Undefined)
	ilang.NewOperator(Type, "/=", Type, "SHARE %a\nSHARE %b\nRUN rational_div_rational\nRENAME %a", false, ilang.Undefined)
	
	ilang.NewOperator(ilang.Number, "\\", ilang.Number, "ARRAY %c\nPUT %a\nPUT %b", false, Type)
	
	ilang.RegisterFunction("number_m_rational", ilang.Function{Exists:true, 
		Returns:[]ilang.Type{ilang.Number}, 
		Args:[]ilang.Type{Type}, 
		Data:`
FUNCTION number_m_rational
	GRAB r
	PLACE r
	PUSH 0
	GET a
	PUSH 1
	GET b
	DIV a a b
	PUSH a
RETURN
	`})
	
	ilang.RegisterFunction("rational_m_number", ilang.Function{Exists:true, 
		Returns:[]ilang.Type{Type}, 
		Args:[]ilang.Type{ilang.Number}, 
		Data:`
FUNCTION rational_m_number
	PULL numerator
	ARRAY r
	PUT numerator
	PUT 1
	SHARE r
RETURN
	`})
	
	ilang.RegisterFunction("text_m_rational", ilang.Function{Exists:true, Returns:[]ilang.Type{ilang.Text}})
	ilang.RegisterFunction("rational", ilang.Function{Exists:true, Method: true, Returns:[]ilang.Type{Type}, Args:[]ilang.Type{}, Data:`
FUNCTION rational
	PUSH 2
	MAKE
RETURN

FUNCTION i_gcd
	PULL b
	PULL a
	VAR i_operator1
	SEQ i_operator1 b 0
	IF i_operator1
		PUSH a
		RETURN
	ELSE
		PUSH b
		VAR i_operator2
		MOD i_operator2 a b
		PUSH i_operator2
		RUN i_gcd
		PULL i_result3
		PUSH i_result3
		RETURN
	END
RETURN
FUNCTION rational_times_rational

	GRAB b
	GRAB a
	ARRAY c
	
	PUT 0
	
	PUT 0
	
	PLACE a
	PUSH 0
	GET i_index4
	PLACE b
	PUSH 0
	GET i_index6
	VAR i_operator5
	MUL i_operator5 i_index4 i_index6
	PLACE c
	PUSH 0
	SET i_operator5
	PLACE a
	PUSH 1
	GET i_index7
	PLACE b
	PUSH 1
	GET i_index9
	VAR i_operator8
	MUL i_operator8 i_index7 i_index9
	PLACE c
	PUSH 1
	SET i_operator8
	PLACE c
	PUSH 0
	GET i_index10
	PUSH i_index10
	PLACE c
	PUSH 1
	GET i_index11
	PUSH i_index11
	RUN i_gcd
	PULL i_result12
	PUSH i_result12
	PULL g
	PLACE c
	PUSH 0
	GET i_index13
	VAR i_operator14
	DIV i_operator14 i_index13 g
	PLACE c
	PUSH 0
	SET i_operator14
	PLACE c
	PUSH 1
	GET i_index15
	VAR i_operator16
	DIV i_operator16 i_index15 g
	PLACE c
	PUSH 1
	SET i_operator16
SHARE c
RETURN
FUNCTION rational_plus_rational

	GRAB b
	GRAB a
	ARRAY c
	
	PUT 0
	
	PUT 0
	
	PLACE a
	PUSH 0
	GET i_index17
	PLACE b
	PUSH 1
	GET i_index19
	VAR i_operator18
	MUL i_operator18 i_index17 i_index19
	PLACE b
	PUSH 0
	GET i_index21
	PLACE a
	PUSH 1
	GET i_index23
	VAR i_operator22
	MUL i_operator22 i_index21 i_index23
	VAR i_operator20
	ADD i_operator20 i_operator18 i_operator22
	PLACE c
	PUSH 0
	SET i_operator20
	PLACE a
	PUSH 1
	GET i_index24
	PLACE b
	PUSH 1
	GET i_index26
	VAR i_operator25
	MUL i_operator25 i_index24 i_index26
	PLACE c
	PUSH 1
	SET i_operator25
	PLACE c
	PUSH 0
	GET i_index27
	PUSH i_index27
	PLACE c
	PUSH 1
	GET i_index28
	PUSH i_index28
	RUN i_gcd
	PULL g
	PLACE c
	PUSH 0
	GET i_index30
	VAR i_operator31
	DIV i_operator31 i_index30 g
	PLACE c
	PUSH 0
	SET i_operator31
	PLACE c
	PUSH 1
	GET i_index32
	VAR i_operator33
	DIV i_operator33 i_index32 g
	PLACE c
	PUSH 1
	SET i_operator33
SHARE c
RETURN
FUNCTION rational_minus_rational

	GRAB b
	GRAB a
	ARRAY c
	
	PUT 0
	
	PUT 0
	
	PLACE a
	PUSH 0
	GET i_index34
	PLACE b
	PUSH 1
	GET i_index36
	VAR i_operator35
	MUL i_operator35 i_index34 i_index36
	PLACE b
	PUSH 0
	GET i_index38
	PLACE a
	PUSH 1
	GET i_index40
	VAR i_operator39
	MUL i_operator39 i_index38 i_index40
	VAR i_operator37
	SUB i_operator37 i_operator35 i_operator39
	PLACE c
	PUSH 0
	SET i_operator37
	PLACE a
	PUSH 1
	GET i_index41
	PLACE b
	PUSH 1
	GET i_index43
	VAR i_operator42
	MUL i_operator42 i_index41 i_index43
	PLACE c
	PUSH 1
	SET i_operator42
	PLACE c
	PUSH 0
	GET i_index44
	PUSH i_index44
	PLACE c
	PUSH 1
	GET i_index45
	PUSH i_index45
	RUN i_gcd
	PULL g
	PLACE c
	PUSH 0
	GET i_index47
	VAR i_operator48
	DIV i_operator48 i_index47 g
	PLACE c
	PUSH 0
	SET i_operator48
	PLACE c
	PUSH 1
	GET i_index49
	VAR i_operator50
	DIV i_operator50 i_index49 g
	PLACE c
	PUSH 1
	SET i_operator50
SHARE c
RETURN
FUNCTION rational_div_rational

	GRAB b
	GRAB a
	ARRAY c
	
	PUT 0
	
	PUT 0
	
	PLACE a
	PUSH 0
	GET i_index51
	PLACE b
	PUSH 1
	GET i_index53
	VAR i_operator52
	MUL i_operator52 i_index51 i_index53
	PLACE c
	PUSH 0
	SET i_operator52
	PLACE a
	PUSH 1
	GET i_index54
	PLACE b
	PUSH 0
	GET i_index56
	VAR i_operator55
	MUL i_operator55 i_index54 i_index56
	PLACE c
	PUSH 1
	SET i_operator55
	PLACE c
	PUSH 0
	GET i_index57
	PUSH i_index57
	PLACE c
	PUSH 1
	GET i_index58
	PUSH i_index58
	RUN i_gcd
	PULL i_result59
	PUSH i_result59
	PULL g
	PLACE c
	PUSH 0
	GET i_index60
	VAR i_operator61
	DIV i_operator61 i_index60 g
	PLACE c
	PUSH 0
	SET i_operator61
	PLACE c
	PUSH 1
	GET i_index62
	VAR i_operator63
	DIV i_operator63 i_index62 g
	PLACE c
	PUSH 1
	SET i_operator63
SHARE c
RETURN
FUNCTION text_m_rational
	GRAB Type
	PLACE Type
	PUSH 0
	GET i_index64
	PUSH i_index64
	PULL numer
	PUSH numer
	PUSH 10
	RUN i_base_number
	GRAB i_result65
	ARRAY i_string67
	PUT 47
	PLACE Type
	PUSH 1
	GET i_index69
	PUSH i_index69
	PULL denom
	PUSH denom
	PUSH 10
	RUN i_base_number
	GRAB i_result70
	ARRAY i_operator68
	JOIN i_operator68 i_string67 i_result70
	ARRAY i_operator66
	JOIN i_operator66 i_result65 i_operator68
	SHARE i_operator66
RETURN
`})
}

