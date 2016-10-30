package main

import "strconv"

func (ic *Compiler) expression() string {
	var token = ic.Scan(0)
	
	switch token {
		case "true":
			ic.ExpressionType = Number
			return "1"
		case "false":
			ic.ExpressionType = Number
			return "0"
		case "error":
			ic.ExpressionType = Number
			return "ERROR"
	}
	
	//Text.
	if token[0] == '"' {
		ic.ExpressionType = Text
		return ic.ParseString(token)
	}
	
	//Letters.
	if token[0] == "'"[0] {
		if s, err := strconv.Unquote(token); err == nil {
			ic.ExpressionType = Letter
			return strconv.Itoa(int([]byte(s)[0]))
		} else {
			ic.RaiseError(err)
		}
	}
	
	//Hexadecimal.
	if len(token) > 2 && token[0] == '0' && token[1] == 'x' { 
		ic.ExpressionType = Number
		return token
	}
	
	//Arrays.
	if token == "[" {
		return ic.ScanArray()
	}
	
	//Subexpessions.
	if token == "(" {
		defer func() {
			ic.Scan(')')
		}()
		return ic.ScanExpression()
	}
	
	//Is it a literal number? Then just return it.
	if _, err := strconv.Atoi(token); err == nil{
		ic.ExpressionType = Number
		return token
	}
	
	//Minus.
	if token == "-" {
		ic.NextToken = token
		ic.ExpressionType = Number
		return ic.Shunt("0")
	}
	
	if t := ic.GetVariable(token); t != Undefined {
		ic.ExpressionType = t
		ic.SetVariable(token+"_use", Used)
		return token
	}
	
	if t, ok := ic.DefinedTypes[token]; ok {
		ic.ExpressionType = t
		
		if ic.Peek() == "(" {
			ic.Scan('(')
			if ic.Peek() == ")" {
				ic.Scan(')')
				
				var array = ic.Tmp("user")
				ic.Assembly("ARRAY ", array)
				for range ic.DefinedTypes[token].Detail.Elements {
					ic.Assembly("PUT 0")
				}
				return array
			}
			ic.RaiseError()	
		} else if ic.Peek() == "{" {
			ic.NextToken = token
			variable := ic.ScanConstructor()
				//TODO better gc protection.
			ic.SetVariable(variable, t)
			ic.SetVariable(variable+"_use", Used)
			return variable
		} else {
			ic.RaiseError()
		}
		
		
	}
	
	if token == "new" {
		var sort = ic.expression()
		if _, ok := ic.DefinedFunctions["new_m_"+ic.ExpressionType.Name]; !ok {
			ic.RaiseError()
		}
		var r = ic.Tmp("new")
		ic.Assembly(ic.ExpressionType.Push, " ", sort)
		ic.Assembly("RUN new_m_"+ic.ExpressionType.Name)
		ic.Assembly("GRAB ", r)
		return r
	}
	
	if _, ok := ic.DefinedFunctions[token]; ok {
		if ic.Peek() == "@" {
			ic.Scan('@')
			var variable = ic.expression()
			ic.Assembly("%v %v", ic.ExpressionType.Push, variable)
			var name = ic.ExpressionType.Name
			ic.ExpressionType = InFunction
			return token+"_m_"+name
		}
	
		if ic.Peek() != "(" {
			ic.ExpressionType = Func
			var id = ic.Tmp("func")
			ic.Assembly("SCOPE ", token)
			ic.Assembly("TAKE ", id)
		
			return id
		} else {
			ic.ExpressionType = InFunction
			return token
		}
	}
	
	ic.RaiseError()
	return ""
}

func (ic *Compiler) ScanExpression() string {
	return ic.Shunt(ic.expression())
}
