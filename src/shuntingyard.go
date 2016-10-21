package main

import (
	"strings"
)

func (ic *Compiler) Shunt(name string) string {
	var token = ic.Scan(0)
	
	switch token {
		case ")", ",", "\n", "]", ";", "{":
			ic.NextToken = token
			return name
		
		case ".":
			if ic.ExpressionType.IsUser() == Undefined {
				ic.RaiseError("Type '%v', cannot be indexed!", ic.ExpressionType.Name)
			}
			var index = ic.Scan(Name)
			return ic.Shunt(ic.IndexUserType(name, index))
		
		case ":":
			if ic.ExpressionType.Push == "PUSH" {
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
			
			return ic.Shunt(slice)
		
		
		case "(":
			if ic.ExpressionType != InFunction {
				ic.RaiseError("Cannot call "+name+", not a function! ("+ic.ExpressionType.Name+")")
			}
			var r = ic.ScanFunctionCall(name)
			ic.Scan(')')
			
			return ic.Shunt(r)
			
		case "[":
			var list bool
			var typename string
			if ic.ExpressionType.Push != "SHARE" {
				ic.RaiseError("Cannot index "+name+", not an array! ("+ic.ExpressionType.Name+")")
			}
			if ic.ExpressionType.List {
				list = true
				typename = ic.ExpressionType.Name
			}
			var index = ic.ScanExpression()
			ic.Scan(']')
			
			ic.ExpressionType = Number
			if ic.ExpressionType == Text {
				ic.ExpressionType = Letter
			}
			
			if !list {
				return ic.Shunt(ic.Index(name, index))
			} else {
				var listdex = ic.Tmp("listdex")
				ic.Assembly("PUSH ", ic.Index(name, index))
				ic.Assembly("HEAP")
				ic.Assembly("GRAB ", listdex)
				ic.ExpressionType = ic.DefinedTypes[typename]
				return ic.Shunt(listdex)
			}
		
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
				if operator, ok = GetOperator(token, ic.ExpressionType, Undefined); !ok {
			
					if OperatorPrecident(token) {
						next = ic.ScanExpression()
					} else {
						next = ic.expression()
					}
					B = ic.ExpressionType
					
					operator, ok = GetOperator(token, A, B)
					
					if token == "=" && A == Text && B == Text {
						ic.ExpressionType = Number
					}
				} else if token == "²" {
					next = name
					B = ic.ExpressionType
					operator, ok = GetOperator("*", A, B)
					token = "*"
				}

				if ok {
				
					asm := operator.Assembly
					asm = strings.Replace(asm, "%a", name, -1)
					asm = strings.Replace(asm, "%b", next, -1)
					asm = strings.Replace(asm, "%c", id, -1)
				
					ic.Assembly(asm)
				
					if operator.ExpressionType != Undefined {
						ic.ExpressionType = operator.ExpressionType
					}
				
					if !OperatorPrecident(token) {
						return ic.Shunt(id)
					}
					return id
				
				} else {
					ic.RaiseError("Invalid Operator Matchup! ", A.Name , token, B.Name, "(types do not support the opperator)")
				}
		}
	}
	
	ic.RaiseError()
	return ""
}
