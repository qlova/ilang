package rational

import "github.com/qlova/ilang/src"

var Type = ilang.NewType("duplex", "SHARE", "GRAB")

func ScanStatement(ic *ilang.Compiler) {
	var name = ic.Scan(ilang.Name)
	var token = ic.Scan(0)
	
	switch token {
		case "=":
			value := ic.ScanExpression()
			if ic.ExpressionType != Type {
				ic.RaiseError("Type mismatch! Type '%s' cannot be assigned to %s which is a duplex number.", ic.ExpressionType, name)
			}
			
			ic.Assembly("PLACE ", value)
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
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	return Type
}

func init() {
	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol("±", ScanSymbol)	

	ilang.RegisterFunction("text_m_duplex", ilang.Function{Exists:true, Returns:[]ilang.Type{ilang.Text}, Data: `
FUNCTION text_m_duplex
		ARRAY i_string1
		SHARE i_string1
		GRAB result
		ARRAY i_newlist2
		SHARE i_newlist2
		GRAB tasks
		GRAB du
		PLACE du
		PUSH 0
		GET base
		PLACE tasks
		PUT 1
		PLACE tasks
		PUT base
		LOOP
			VAR i_operator3
			SEQ i_operator3 #tasks 0
			IF i_operator3
				BREAK
			END
			PLACE tasks
			POP sum
			POP index
			PLACE du
			PUSH index
			GET v
			VAR i_operator5
			SUB i_operator5 #du 1
			VAR i_operator4
			SGE i_operator4 index i_operator5
			IF i_operator4
				VAR i_operator7
				ADD i_operator7 sum v
				PUSH i_operator7
				PUSH 10
				RUN i_base_number
				GRAB i_result8
				JOIN result result i_result8
				ARRAY i_string10
				PUT 44
				JOIN result result i_string10
				VAR i_operator12
				SUB i_operator12 sum v
				PUSH i_operator12
				PUSH 10
				RUN i_base_number
				GRAB i_result13
				JOIN result result i_result13
				ARRAY i_string15
				PUT 44
				JOIN result result i_string15
			ELSE
					VAR i_operator16
					ADD i_operator16 index 1
					PLACE tasks
					PUT i_operator16
					VAR i_operator17
					ADD i_operator17 sum v
					PLACE tasks
					PUT i_operator17
					VAR i_operator18
					ADD i_operator18 index 1
					PLACE tasks
					PUT i_operator18
					VAR i_operator19
					SUB i_operator19 sum v
					PLACE tasks
					PUT i_operator19
			END
		REPEAT
		SHARE result
RETURN 

	`})
	
	
	ilang.NewOperator(ilang.Number, "±", ilang.Number, "ARRAY %c\nPUT %a\nPUT %b", false, Type)
	ilang.RegisterShunt("±", func(ic *ilang.Compiler, name string) string {
		if ic.ExpressionType != Type {
			return ""
		}
		var value = ic.ScanExpression()
		if ic.ExpressionType != ilang.Number {
			ic.RaiseError("Only numbers can be duplexed!")
		}

		ic.Assembly("PLACE ", name)
		ic.Assembly("PUT ", value)

		ic.ExpressionType = Type
		return name
	})
}