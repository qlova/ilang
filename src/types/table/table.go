package table

import "github.com/qlova/ilang/src"
import "strconv"

var Type = ilang.NewType("table", "PUSH", "PULL")

func ScanStatement(ic *ilang.Compiler) {
	var table = ic.Scan(0)
    var t  = ic.GetVariable(table)
	
	ic.Scan('[')
	var index = ic.ScanExpression()
	if ic.ExpressionType != ilang.Text {
		ic.RaiseError("Table must have text index.")
	}
	ic.Scan(']')
	ic.Scan('=')
	var value = ic.ScanExpression()
	
    if t.SubType == nil {
        t.SubType = new(ilang.Type)
        *t.SubType = ic.ExpressionType
        ic.UpdateVariable(table, t)
    }
    
    if !ic.ExpressionType.Equals(*t.SubType) {
        ic.RaiseError("Cannot add value of type '",ic.ExpressionType.GetComplexName(),"' to a table of '",t.SubType.GetComplexName(),"'")
    }
	
	var tmp = ic.Tmp("newtableref")
	
	ic.Assembly("IF ", table)
		ic.Assembly("PUSH %s", table)
		ic.Assembly("SHARE %s", index)
		ic.Assembly("PUSH %s", ic.GetPointerTo(value))
		ic.Assembly(ic.RunFunction("table_set"))
		ic.Assembly("PULL %s", tmp)
		ic.Assembly("ADD %s %s %v", table, tmp, 0)
	ic.Assembly("ELSE")
		ic.Assembly("ADD ERROR 1 0")
	ic.Assembly("END")
	
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	var TableType = Type
	TableType.SubType = new(ilang.Type)
	*TableType.SubType = ic.ScanSymbolicType()
	return TableType
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

func Shunt(ic *ilang.Compiler, name string) string {
	if ic.ExpressionType.Name == "table" {
        var table = ic.ExpressionType
        
		var index = ic.ScanExpression()
		ic.Scan(']')
		if ic.ExpressionType != ilang.Text {
			ic.RaiseError("A Table must have a text index! Found ", ic.ExpressionType.Name)
		}
		
		var tableval = ic.Tmp("tableval")
		
		ic.Assembly("IF ", name)
			ic.Assembly("SHARE %s", index)
			ic.Assembly("PUSH %s", name)
			ic.Assembly(ic.RunFunction("table_get"))
		ic.Assembly("ELSE")
			ic.Assembly("ADD ERROR 1 0")
			ic.Assembly("PUSH 0")
		ic.Assembly("END")
		
		ic.Assembly("PULL %s", tableval)
	
		ic.ExpressionType = *table.SubType
	
		return ic.Shunt(ic.Dereference(tableval))
	}
	return ""
}

func init() {
	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol(":", ScanSymbol)
	ilang.RegisterShunt("[", Shunt)
	
	ilang.RegisterFunction("table", ilang.Method(Type, true, "PUSH 64\nMAKE\nPUSH 0\nHEAP\n"))
	
	ilang.RegisterFunction("table_set", ilang.Function{Exists:true, Args:[]ilang.Type{ilang.Number, ilang.Text, ilang.Number}, Returns:[]ilang.Type{ilang.Number}, Data:`
FUNCTION table_set
	PULL value
	GRAB key
	
	PULL ref
	PUSH ref
	HEAP
	GRAB t
	
	SHARE key
	PUSH 255
	RUN i_hash
	PULL hash

	PLACE t
	PUSH hash
	GET entry
	
	PUSH 0
	PULL tooheavy

	PUSH -1
	PULL negativeone
	LOOP
		PLACE t
		PUSH 0
		GET firstindex
		
		#First index stores information on how heavy the table is.
		IF firstindex
			PUSH firstindex
			HEAP
			GRAB bucket
			PLACE bucket
			PUSH 1
			GET i_index5
			VAR i_operator6
			ADD i_operator6 i_index5 1
			PLACE bucket
			PUSH 1
			SET i_operator6
			PLACE bucket
			PUSH 1
			GET i_index7
			VAR i_operator10
			DIV i_operator10 #t 4
			VAR i_operator9
			SUB i_operator9 #t i_operator10
			VAR i_operator8
			SGT i_operator8 i_index7 i_operator9
			ADD tooheavy 0 i_operator8
			BREAK
		ELSE
			ARRAY a
			PUT negativeone
			PUT 0
			SHARE a
			PUSH 0
			HEAP
			PULL pointer
			PLACE t
			PUSH 0
			SET pointer
		END
	REPEAT
	IF entry
		PUSH entry
		HEAP
		GRAB bucket
		
		#Need to check if we are already in the bucket!
		PUSH 0
		PULL i
		LOOP
			PLACE bucket
			PUSH i
			GET i_index29
			VAR i_operator30
			SEQ i_operator30 i_index29 hash
			IF i_operator30
				VAR i_operator31
				ADD i_operator31 i 1
				PLACE bucket
				PUSH i_operator31
				SET value
				BREAK
			END
			ADD i i 2
			VAR i_operator34
			SGT i_operator34 i #bucket
			IF i_operator34
			
				#Add it to the end!
				PLACE bucket
				PUT hash
				PLACE bucket
				PUT value
				BREAK
			END
		REPEAT
		
	ELSE
		ARRAY a
		PUT hash
		PUT value
		SHARE a
		PUSH 0
		HEAP
		PULL pointer
		PLACE t
		PUSH hash
		SET pointer
	END
	IF tooheavy
		VAR i_operator13
		MUL i_operator13 #t 2
		VAR i_operator14
		SUB i_operator14 i_operator13 1
		PUSH i_operator14
		PULL newlength
		PUSH newlength
		MAKE
		PUSH 0
		HEAP
		PULL newref
		
		IF 1
		ARRAY i_delete18
		VAR i_i16
		VAR i_backup17
		LOOP
			VAR i_in15
			ADD i_i16 0 i_backup17
			SGE i_in15 i_i16 #t
			IF i_in15
				BREAK
			END
			PLACE t
			PUSH i_i16
			GET anentry
			ADD i_backup17 i_i16 1
		
			IF anentry
				PUSH anentry
				HEAP
				GRAB bucket
				PUSH 0
				PULL i
				LOOP
					PUSH newref
					ARRAY i_array19
					PLACE bucket
					PUSH i
					GET i_index20
					PLACE i_array19
					PUT i_index20
					SHARE i_array19
					
					GRAB i_result21
					SHARE i_result21
					VAR i_operator22
					ADD i_operator22 i 1
					PLACE bucket
					PUSH i_operator22
					GET i_index23
					PUSH i_index23
					RUN table_set
					PULL i_result24
					ADD i i 2
					VAR i_operator26
					SGT i_operator26 i #bucket
					IF i_operator26
						BREAK
					END
				REPEAT
			END
		REPEAT
		
			VAR ii_i8
			VAR ii_backup9
			LOOP
				VAR ii_in7
				ADD ii_i8 0 ii_backup9
				SGE ii_in7 ii_i8 #i_delete18
				IF ii_in7
					BREAK
				END
				PLACE i_delete18
				PUSH ii_i8
				GET i_v
				ADD ii_backup9 ii_i8 1
		
				VAR ii_operator11
				SUB ii_operator11 #t 1
				PLACE t
				PUSH ii_operator11
				GET ii_index12
				PLACE t
				PUSH i_v
				SET ii_index12
				PLACE t
				POP n
				ADD n 0 0
			REPEAT
								
		END
		PUSH newref
		RETURN
	END
	PUSH ref
RETURN
`})

ilang.RegisterFunction("table_get", ilang.Function{Exists:true, Args:[]ilang.Type{ilang.Number, ilang.Text}, Returns:[]ilang.Type{ilang.Number}, Data:`
FUNCTION table_get
	GRAB key
	PULL ref
	PUSH ref
	HEAP
	GRAB t
	SHARE key
	PUSH 255
	RUN i_hash
	PULL i_operator27
	
	PUSH i_operator27
	PULL hash
	PLACE t
	PUSH hash
	GET i_index28
	PUSH i_index28
	PULL entry
	IF entry
		PUSH entry
		HEAP
		GRAB bucket
		PUSH 0
		PULL i
		LOOP
			PLACE bucket
			PUSH i
			GET i_index29
			VAR i_operator30
			SEQ i_operator30 i_index29 hash
			IF i_operator30
				VAR i_operator31
				ADD i_operator31 i 1
				PLACE bucket
				PUSH i_operator31
				GET i_index32
				PUSH i_index32
				RETURN
			END
			ADD i i 2
			VAR i_operator34
			SGT i_operator34 i #bucket
			IF i_operator34
				ADD ERROR 0 404
				PUSH 0
				RETURN
			END
		REPEAT
	ELSE
		ADD ERROR 0 404
		PUSH 0
		RETURN
	END
RETURN
`})
}

