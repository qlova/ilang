package ilang

import (
	"strings"
)

func (ic *Compiler) Shunt(name string) string {
	var token = ic.Scan(0)
	
	if list, ok := Shunts[token]; ok {
		for _, f := range list {
			result := f(ic, name)
			if result != "" {
				return ic.Shunt(result)
			}
		}
	}
	
	if token != "." && token[0] == '.' {
		ic.NextToken = token[1:]
		
		if list, ok := Shunts["."]; ok {
			for _, f := range list {
				result := f(ic, name)
				if result != "" {
					return ic.Shunt(result)
				}
			}
		}
		println("oi", name, ic.ExpressionType.Name)
		ic.NextToken = ""
	}
	
	switch token {
		case ")", ",", "\n", "]", ";", "{", "}", "|":
			ic.NextToken = token
			return name
		
		case ".":
			var index = ic.Scan(Name)
			
			if ic.ExpressionType.IsUser() == Undefined {
				ic.RaiseError("Type ", ic.ExpressionType.GetComplexName(), ", cannot be indexed with ", index ,"!")
			}

			if ic.Peek() != "." && !ic.DisableOwnership {
				ic.TakingExpression = true
			} else {
				ic.TakingExpression = false
			}
			
			defer func() {
				ic.TakingExpression = false
			}()
			
			return ic.Shunt(ic.IndexUserType(name, index))
		
		case ":":
			/*if ic.ExpressionType.Push == "PUSH" {
				ic.NextToken = token
				return name
			}
			if ic.ExpressionType.Push != "SHARE" {
				ic.RaiseError("Cannot index "+name+", not an array! ("+ic.ExpressionType.Name+")")
			}
			
			var original = ic.ExpressionType
			
			var slice = ic.Tmp("slice")
			ic.Assembly("SHARE ", name)
			
			var low,high string
			if tok := ic.Scan(0); tok != ":" {
				ic.NextToken = tok
				low = ic.ScanExpression()
				ic.Scan(':')
			} else {
				low = "0"
			}
			
			if tok := ic.Scan(0); tok != ":" {
				ic.NextToken = tok
				high = ic.ScanExpression()
				ic.Scan(':')
			} else {
				high = "#"+name
			}
			
			ic.Assembly("PUSH ", high)
			ic.Assembly("PUSH ", low)
			
			ic.Assembly("SLICE")
			
			ic.Assembly("GRAB ", slice)
			
			ic.ExpressionType = original
			
			return ic.Shunt(slice)*/
		

			
		case "[":
			if ic.ExpressionType != Text {
				ic.RaiseError("Cannot index type ", ic.ExpressionType)
			}
			
			index := ic.ScanExpression()
			ic.Scan(']')
			
			var value = ic.Tmp("value")
			ic.Assembly("PLACE ", name)
			ic.Assembly("PUSH ", index)
			ic.Assembly("GET ", value)
				
			ic.ExpressionType = GetType("letter")		
			return ic.Shunt(value)

		
		default:
			
			if IsOperator(token+ic.Peek()) {
				token += ic.Peek()
				ic.Scan(0)
			}
		
			if IsOperator(token) {
				id := ic.Tmp("operator")
			
				var operator Operator
				var ok bool
			
				var A = ic.ExpressionType
				var B Type
				var next string
				if operator, ok = GetOperator(token, ic.ExpressionType, Undefined); token == "²" {
					next = name
					B = ic.ExpressionType
					operator, ok = GetOperator("*", A, B)
					token = "*"
				} else if !ok {
			
					if OperatorPrecident(token) {
						next = ic.ScanExpression()
					} else {
						next = ic.expression()
					}
					B = ic.ExpressionType
					
					operator, ok = GetOperator(token, A, B)
					ic.ExpressionType = operator.ExpressionType
				}
				
				if !ok {
					for _, f := range SpecialOperators {
						o := f(token, A, B)
						if o != nil {
							operator = *o
							ok = true
							ic.ExpressionType = operator.ExpressionType
							break
						}
					}
				}

				if ok {
				
					asm := operator.Assembly
					asm = strings.Replace(asm, "%a", name, -1)
					asm = strings.Replace(asm, "%b", next, -1)
					asm = strings.Replace(asm, "%c", id, -1)
					asm = strings.Replace(asm, "$Super", A.Super, -1)
					
					if strings.Contains(asm, "%t") {
						asm = strings.Replace(asm, "%t", ic.Tmp("tmp"), -1)
					}
				
					ic.Assembly(asm)
				
					ic.ExpressionType = operator.ExpressionType
				
					if !OperatorPrecident(token) {
						return ic.Shunt(id)
					}
					return id
				
				} else {
				
					ic.RaiseError("Invalid Operator Matchup! ", A.Name , token, B.Name, "(types do not support the opperator)")
				}
		}
	}
	
	//ic.RaiseError()
	ic.NextToken = token
	return name
}
