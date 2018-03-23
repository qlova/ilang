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
			var pointer = ic.Tmp("cut")
			ic.Assembly("PLACE ", name)
			ic.Assembly("POP ", pointer)
			ic.Assembly(t.SubType.Free(pointer))
			ic.Assembly("MUL %s -1 %s ", pointer,  pointer)
			ic.Assembly("PUSH ", pointer)
			ic.Assembly("HEAP")
			
		case "+":
				var value string
			if ic.Peek() == "+" {
					//List++ 
					//We append to the list the zero value of it's SubType.
					ic.Scan('+')
					
					value = ic.CallType(t.SubType.Name)
					ic.ExpressionType = *t.SubType
			} else if ic.Peek() == "=" {
				ic.Scan('=')
				
				value = ic.ScanExpression()
				
				if t.SubType == nil {
					t.SubType = new(ilang.Type)
					*t.SubType = ic.ExpressionType
					ic.UpdateVariable(name, t)
				}
				
				if ic.ExpressionType != *t.SubType {
					
					if ic.CanCast(ic.ExpressionType, *t.SubType) {
						
						
						ic.Assembly(ic.Cast(value,  ic.ExpressionType, *t.SubType))
						
						var cast = ic.Tmp("cast")
						ic.Assembly(t.SubType.Pop, " ", cast)
						ic.ExpressionType = *t.SubType
						
						ic.Assembly("PLACE ", name)
						ic.Assembly("PUT ", ic.GetPointerTo(cast))
						ic.SetVariable(value+".", ilang.Protected)
						
						return
						
					} else {
						ic.RaiseError("Cannot add value of type '",ic.ExpressionType.Name,"' to a list of '",t.SubType.Name,"'")
					}
				}
			} else {
					ic.RaiseError("Expecting ++ or += got +", ic.Peek())
			}
			
			ic.Assembly("PLACE ", name)
			ic.Assembly("PUT ", ic.GetPointerTo(value))
			ic.SetVariable(value+".", ilang.Protected)
		
			
		case "[":
			
			index := ic.ScanExpression()
			ic.Scan(']')
			
			token = ic.Scan(0)
			if token == "=" {
				
				value := ic.ScanExpression()
				//TODO Garbage collection.
				ic.SetVariable(value+".", ilang.Protected)
				ic.Assembly("PUSH ", index)
				ic.Assembly("PLACE ", name)
				ic.Assembly("SET ", ic.GetPointerTo(value))
			
			} else if token == "[" {
				
				var pointer = ic.Tmp("pointer")
				ic.Assembly("PUSH ", index)
				ic.Assembly("PLACE ", name)
				ic.Assembly("GET ", pointer)
				
				ic.NextToken = ic.Dereference(pointer)
				ic.SetVariable(ic.NextToken, *t.SubType)
				ic.NextNextToken = token
				ic.ShuntStatement(*t.SubType)
				
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
	ListType.SubType  = new(ilang.Type)
	*ListType.SubType = ic.ScanSymbolicType()
	return ListType
}

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	
	if token == "list" {
		ic.Scan('(')
		var t = ic.Scan(0)
		ic.Scan(')')
		
		var ListType = Type
		ListType.SubType = new(ilang.Type)
		*ListType.SubType = ic.GetType(t)
		
		var tmp = ic.Tmp("newlist")
		ic.Assembly("ARRAY ", tmp)
		
		ic.ExpressionType = ListType
		
		return tmp
	}
	
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
	if ic.ExpressionType.Name == "list" || (ic.ExpressionType.Class != nil && ic.ExpressionType.Class.Name == "list") {
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

func Collection(ic *ilang.Compiler, t ilang.Type) {
	if t.Name != "list" {
		return
	}
	
	var scope = ic.Scope
	ic.GainScope()
	ic.Library("FUNCTION collect_m_", t.GetComplexName())
	ic.Library("GRAB list")
	
	
	ic.Library("VAR i")
	ic.Library("VAR condition")
	ic.Library("LOOP")
	ic.GainScope()
		ic.Library("SGE condition i #list")
		ic.Library("IF condition")
			ic.GainScope()
			ic.Library("BREAK")
			ic.LoseScope()
		ic.Library("END")
	
		ic.Library("PLACE list")
		ic.Library("PUSH i")
		ic.Library("GET pointer")
		ic.Library("ADD pointer pointer 0")
		
		ic.Library(t.SubType.Free("pointer"))
		
		ic.Library("ADD i i 1")
	ic.LoseScope()
	ic.Library("REPEAT")
	ic.Library("RETURN")
	ic.LoseScope()
	ic.Scope = scope
}

var Number = Type

func init() {
	Number.SubType = new(ilang.Type)
	*Number.SubType = ilang.Number

	ilang.RegisterStatement(Type, ScanStatement)	
	ilang.RegisterSymbol("..", ScanSymbol)
	ilang.RegisterShunt("[", Shunt)
	ilang.RegisterCollection(Collection)
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

