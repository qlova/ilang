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
	
	ic.LoadFunction("duplex")
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	return Type
}

func init() {
	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol("±", ScanSymbol)

	ilang.RegisterFunction("duplex", ilang.Function{Exists:true, Returns:[]ilang.Type{Type}, Data: `
	FUNCTION duplex
		ARRAY d
		PUT 0
		SHARE d
	RETURN
	
	FUNCTION duplex_plus_duplex
        GRAB b
        GRAB a
        ARRAY r
        
        VAR i
        VAR a[i]+b[i]
        VAR i>=len(a)
        VAR i>=len(b)
    
        LOOP
            SGE i>=len(a) i #a
            SGE i>=len(b) i #b
            
            PLACE a
            PUSH i
            GET a[i]
            
            IF i>=len(a)
                ADD a[i] 0 0
            END
            
            PLACE b
            PUSH i
            GET b[i]
            
            IF i>=len(b)
                ADD b[i] 0 0
                IF i>=len(a)
                    BREAK
                END
            END
            
            ADD a[i]+b[i] a[i] b[i]
            
            PLACE r
            PUT a[i]+b[i]
            
            ADD i i 1
        REPEAT 
        SHARE r
	RETURN
	
	FUNCTION duplex_minus_duplex
        GRAB b
        GRAB a
        ARRAY r
        
        VAR i
        VAR a[i]+b[i]
        VAR i>=len(a)
        VAR i>=len(b)
    
        LOOP
            SGE i>=len(a) i #a
            SGE i>=len(b) i #b
            
            PLACE a
            PUSH i
            GET a[i]
            
            IF i>=len(a)
                ADD a[i] 0 0
            END
            
            PLACE b
            PUSH i
            GET b[i]
            
            IF i>=len(b)
                ADD b[i] 0 0
                IF i>=len(a)
                    BREAK
                END
            END
            
            SUB a[i]+b[i] a[i] b[i]
            
            PLACE r
            PUT a[i]+b[i]
            
            ADD i i 1
        REPEAT 
        SHARE r
	RETURN
	
	FUNCTION duplex_times_number
        GRAB a
        PULL n
        ARRAY r
        
        VAR i
        VAR a[i]+n
        VAR i>=len(a)
    
        LOOP
            SGE i>=len(a) i #a
            
            PLACE a
            PUSH i
            GET a[i]
            
            IF i>=len(a)
               BREAK
            END
            
            MUL a[i]+n a[i] n
            
            PLACE r
            PUT a[i]+n
            
            ADD i i 1
        REPEAT 
        SHARE r
	RETURN
	
	FUNCTION duplex_div_number
        GRAB a
        PULL n
        ARRAY r
        
        VAR i
        VAR a[i]+n
        VAR i>=len(a)
    
        LOOP
            SGE i>=len(a) i #a
            
            PLACE a
            PUSH i
            GET a[i]
            
            IF i>=len(a)
               BREAK
            END
            
            DIV a[i]+n a[i] n
            
            PLACE r
            PUT a[i]+n
            
            ADD i i 1
        REPEAT 
        SHARE r
	RETURN
	
	FUNCTION duplex_plus_number
        GRAB a
        PULL n
        ARRAY r
        
        VAR i
        VAR a[0]+n
        VAR i>=len(a)
    
        LOOP
            SGE i>=len(a) i #a
            
            PLACE a
            PUSH i
            GET a[i]
            
            IF i>=len(a)
               BREAK
            END
            
            PLACE r
            IF i
                PUT a[i]
            ELSE
                ADD a[0]+n a[i] n
                PUT a[0]+n
            END
            
            ADD i i 1
        REPEAT 
        SHARE r
	RETURN
	
	FUNCTION duplex_minus_number
        GRAB a
        PULL n
        ARRAY r
        
        VAR i
        VAR a[0]+n
        VAR i>=len(a)
    
        LOOP
            SGE i>=len(a) i #a
            
            PLACE a
            PUSH i
            GET a[i]
            
            IF i>=len(a)
               BREAK
            END
            
            PLACE r
            IF i
                PUT a[i]
            ELSE
                SUB a[0]+n a[i] n
                PUT a[0]+n
            END
            
            ADD i i 1
        REPEAT 
        SHARE r
	RETURN
	`})

	ilang.RegisterFunction("array_m_duplex", ilang.Function{Exists:true, Returns:[]ilang.Type{ilang.Array}, Data: `
FUNCTION array_m_duplex
	PUSH 0
	MAKE
	GRAB i_result11
	SHARE i_result11
	GRAB result
	ARRAY i_newlist12
	SHARE i_newlist12
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
		VAR i_operator13
		SEQ i_operator13 #tasks 0
		IF i_operator13
			BREAK
		END
		PLACE tasks
		POP sum
		POP index
		PLACE du
		PUSH index
		GET v
		VAR i_operator15
		SUB i_operator15 #du 1
		VAR i_operator14
		SGE i_operator14 index i_operator15
		IF i_operator14
			VAR i_operator16
			ADD i_operator16 sum v
			PUSH i_operator16
			PULL a
			VAR i_operator17
			SUB i_operator17 sum v
			PUSH i_operator17
			PULL b
			PLACE result
			PUT a
			PUT b
		ELSE
				VAR i_operator18
				ADD i_operator18 index 1
				PLACE tasks
				PUT i_operator18
				VAR i_operator19
				ADD i_operator19 sum v
				PLACE tasks
				PUT i_operator19
				VAR i_operator20
				ADD i_operator20 index 1
				PLACE tasks
				PUT i_operator20
				VAR i_operator21
				SUB i_operator21 sum v
				PLACE tasks
				PUT i_operator21
		END
	REPEAT
	SHARE result
RETURN
	`})

	ilang.RegisterFunction("text_m_duplex", ilang.Function{Exists:true, Returns:[]ilang.Type{ilang.Text}, Data: `
FUNCTION text_m_duplex
	ARRAY i_string1
	SHARE i_string1
	GRAB result
	GRAB du
	PUSH 0
	PULL i
	PLACE du
	PUSH i
	GET v
	PUSH v
	PUSH 10
	RUN i_base_number
	GRAB i_result3
	JOIN result result i_result3
	ADD i i 1
	LOOP
		VAR i_operator5
		SGE i_operator5 i #du
		IF i_operator5
			BREAK
		END
		ARRAY i_string7
		PUT 194
		PUT 177
		JOIN result result i_string7
		PLACE du
		PUSH i
		GET v2
		PUSH v2
		PUSH 10
		RUN i_base_number
		GRAB i_result9
		JOIN result result i_result9
		ADD i i 1
	REPEAT
	SHARE result
RETURN
	`})
	
	
    ilang.NewOperator(Type, "+", Type, "SHARE %a\nSHARE %b\nRUN duplex_plus_duplex\nGRAB %c", false)
    ilang.NewOperator(Type, "-", Type, "SHARE %a\nSHARE %b\nRUN duplex_minus_duplex\nGRAB %c", false)
    
    ilang.NewOperator(Type, "+", ilang.Number, "SHARE %a\nPUSH %b\nRUN duplex_plus_number\nGRAB %c", false)
    ilang.NewOperator(Type, "-", ilang.Number, "SHARE %a\nPUSH %b\nRUN duplex_minus_number\nGRAB %c", false)
    ilang.NewOperator(ilang.Number, "+", Type, "PUSH %a\nSHARE %b\nRUN duplex_plus_number\nGRAB %c", false)
    ilang.NewOperator(ilang.Number, "-", Type, "PUSH %a\nSHARE %b\nRUN duplex_minus_number\nGRAB %c", false)
    
    ilang.NewOperator(Type, "*", ilang.Number, "SHARE %a\nPUSH %b\nRUN duplex_times_number\nGRAB %c", false)
    ilang.NewOperator(Type, "/", ilang.Number, "SHARE %a\nPUSH %b\nRUN duplex_div_number\nGRAB %c", false)
    ilang.NewOperator(ilang.Number, "*", Type, "SHARE %b\nPUSH %a\nRUN duplex_times_number\nGRAB %c", false)
    ilang.NewOperator(ilang.Number, "/", Type, "SHARE %b\nPUSH %a\nRUN duplex_div_number\nGRAB %c", false)
    
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
