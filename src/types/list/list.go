package list

import "github.com/qlova/ilang/src"

var Type = ilang.NewType("list", "SHARE", "GRAB")

func ScanStatement(ic *ilang.Compiler) {
	var name = ic.Scan(ilang.Name)
	var t  = ic.GetVariable(name)
	var token = ic.Scan(0)
	
	//TODO CLEAN THIS UP!
	switch token {
		case "-":
			ic.Scan('-')
			ic.Assembly("PLACE ", name)
			ic.Assembly("POP ", ic.Tmp("cut"))
			
		case "+":
			ic.Scan('=')
			
			value := ic.ScanExpression()
			
			if t.SubType == nil {
				t.SubType = new(ilang.Type)
				*t.SubType = ic.ExpressionType
				ic.UpdateVariable(name, t)
			}
			
			if ic.ExpressionType != *t.SubType {
				ic.RaiseError("Cannot add value of type '",ic.ExpressionType.Name,"' to a list of '",t.SubType.Name,"'")
			}
			
			ic.Assembly("PLACE ", name)
			ic.Assembly("PUT ", ic.GetPointerTo(value))
		
			
		case "[":
			
			index := ic.ScanExpression()
			ic.Scan(']')
			
			token = ic.Scan(0)
			if token == "=" {
				
				value := ic.ScanExpression()
				ic.Assembly("PUSH ", index)
				ic.Assembly("PLACE ", name)
				ic.Assembly("SET ", ic.GetPointerTo(value))
				
			} else {
				var pointer = ic.Tmp("pointer")
				ic.Assembly("PUSH ", index)
				ic.Assembly("PLACE ", name)
				ic.Assembly("GET ", pointer)
				
				ic.NextToken = pointer
				ic.NextNextToken = token
				ic.ScanExpression()
			}
			
		case "=":
			
			
		default:
			ic.ExpressionType = t
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != ilang.Undefined {
				ic.RaiseError("blank expression!")
			}
	}
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	var ListType = Type
	ListType.SubType = new(ilang.Type)
	*ListType.SubType = ic.ScanSymbolicType()
	return ListType
}

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	//Types.
	if token == "[" {
					
		var tmp = ic.Tmp("newlist")
		ic.Assembly("ARRAY ", tmp)
		
		var ListType = ilang.Undefined
		
		if tok := ic.Scan(0); tok == "]" {
			ic.ExpressionType = Type
			return tmp
		} else {
			ic.NextToken = tok
		}
		
		for {			
			value := ic.ScanExpression()
			if (ListType == ilang.Undefined) {
				ListType = ic.ExpressionType
			}
			
			//Oh great, this needs to become a something list... #TODO
			if ic.ExpressionType != ListType {
				ic.RaiseError("Consistency error! There is a ", ic.ExpressionType.Name, 
				" value in the list.\nThis is inconsistant with the previous elements which are ",ListType.Name," values.")
			}

			ic.Assembly("PLACE ", tmp)
			ic.Assembly("PUT ", ic.GetPointerTo(value))
			
			token = ic.Scan(0)
			if token != "," {
				if token == "]" {
					ic.ExpressionType = Type
					ic.ExpressionType.SubType = &ListType
					break
				}
				ic.RaiseError("Expecting ,")
			}
		}
		
		return tmp
	}
	
	return ""
}

func Shunt(ic *ilang.Compiler, name string) string {
	if ic.ExpressionType.Name == "list" {
		var list = ic.ExpressionType
	
		index := ic.ScanExpression()
		ic.Scan(']')
		if ic.ExpressionType != ilang.Number {
			ic.RaiseError("Index must be a number.")
		}
		
		if list.SubType == nil {
			ic.RaiseError("List elements are undefined!")
		}
		
		var pointer = ic.Tmp("index")
		ic.Assembly("PUSH ", index)
		ic.Assembly("PLACE ", name)
		ic.Assembly("GET ", pointer)
		
		ic.ExpressionType = *list.SubType
		
		return ic.Dereference(pointer)
	}
	return ""
}

var Number = Type

func init() {
	Number.SubType = new(ilang.Type)
	*Number.SubType = ilang.Number

	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol("..", ScanSymbol)
	ilang.RegisterShunt("[", Shunt)
	ilang.RegisterExpression(ScanExpression)
	
	ilang.RegisterFunction("table", ilang.Method(Type, true, "PUSH 64\nMAKE\nPUSH 0\nHEAP\n"))
	ilang.RegisterFunction("text_m_numberlist", ilang.BlankMethod(ilang.Text))
	
	ilang.RegisterFunction("copy_m_numberlist", ilang.Function{Exists:true, Returns:[]ilang.Type{Number}, Data: `
	#Compiled with IC.
	FUNCTION copy_m_numberlist
		GRAB array
		ARRAY c
	
		VAR i
		LOOP
			VAR i+shunt+1
			SGE i+shunt+1 i #array
			IF i+shunt+1
				SHARE c
				RETURN
			END
			PLACE array 
				PUSH i 
				GET i+shunt+3
			VAR v
			ADD v 0 i+shunt+3
			PLACE c
				PUT v
			ADD i i 1
		REPEAT
	RETURN
	`})
}

