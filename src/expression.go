package ilang

import "strconv"
import "strings"
import "github.com/gedex/inflector"
import "fmt"
import "math"

func (ic *Compiler) expression() string {
	var token = ic.Scan(0)
	
	for token == "\n" {
		token = ic.Scan(0)
	}
	
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
	
	if token == "{" {
		var t = ic.Scan(0)
		if t == "}" {
			ic.ExpressionType = Something
		} else {
			ic.ExpressionType = ic.DefinedInterfaces[t].GetType()
			ic.Scan('}')
		}
		var tmp = ic.Tmp("something")
		ic.Assembly("ARRAY ", tmp)
		ic.Assembly("PUT 0")
		return tmp
	}
	
	if token == "something" {
		ic.Scan('(')
		ic.Scan(')')
		ic.ExpressionType = Something
		var tmp = ic.Tmp("something")
		ic.Assembly("ARRAY ", tmp)
		ic.Assembly("PUT 0")
		return tmp
	}
	
	if token == "?" {
		ic.ExpressionType = Something
		var tmp = ic.Tmp("something")
		ic.Assembly("ARRAY ", tmp)
		ic.Assembly("PUT 0")
		return tmp
	}
	
	//Text.
	if token[0] == '"' || token[0] == '`' {
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
	if token == "$" && ic.Peek() == "[" {
		return ic.ScanArray()
	}
	
	//Arrays.
	if token == "[" {
		return ic.ScanArray()
	}
	
	//Pipes.
	if token == "|" {
		var name = "open"
		if ic.Peek() != "|" {
			var arg = ic.ScanExpression()
			name += "_m_"+ic.ExpressionType.Name
			if f, ok := ic.DefinedFunctions[name]; ok {
				var tmp = ic.Tmp("open")
				ic.Assembly(ic.ExpressionType.Push, " ", arg)
				ic.Assembly(ic.RunFunction(name))
				ic.Assembly(f.Returns[0].Pop, " ", tmp)
				ic.ExpressionType = f.Returns[0]
				ic.Scan('|')
				return tmp
			} else {
				ic.RaiseError("Cannot create a pipe out of a ", ic.ExpressionType.Name)
			}
		} else {
			ic.Scan('|')
			var tmp = ic.Tmp("pipe")
			ic.Assembly("PIPE ", tmp)
			ic.ExpressionType = Pipe
			return tmp
			//ic.RaiseError("Blank pipe!")
		}
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
	
	//Decimal numbers.
	if strings.Contains(token, ".") {
		parts := strings.Split(token, ".")
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		ic.ExpressionType = Decimal
		return fmt.Sprint(a*1000000+b*int(math.Pow(10, 6-float64(len(parts[1])))))
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
		if tok := ic.Scan(0); tok == "[" || tok == "." {
			ic.NextToken = tok
			return ic.Shunt(token)
		} else {
			ic.NextToken = tok
		}
		return token
	}
	
	//Scope methods with multiple arguments inside the method.
	//eg. method Package.dosomething(22)
	// in a Package method, dosomething(22) should be local.
	if ic.GetFlag(InMethod) {
		if f, ok := ic.DefinedFunctions[token+"_m_"+ic.LastDefinedType.Name]; ok && len(f.Args) > 0 {
			if ic.LastDefinedType.Detail == nil || len(ic.LastDefinedType.Detail.Elements) > 0 {
				ic.Assembly("%v %v", ic.LastDefinedType.Push, ic.LastDefinedType.Name)
			}
			ic.ExpressionType = InFunction
			return token+"_m_"+ic.LastDefinedType.Name
		}
	}
	
	if ic.TypeExists(token) {
		ic.ExpressionType = ic.DefinedTypes[token]
		
		//This is a constructor. eg. var bug = Bug(); where Bug is a type.
		if ic.Peek() == "(" || ic.NextToken == "(" {
			ic.Scan('(')
			ic.Scan(')')
			
			return ic.CallType(token)
			
		//This is a type literal.
		} else if ic.Peek() == "{" {
			ic.NextToken = token
			variable := ic.ScanTypeLiteral()
				//TODO better gc protection.
			ic.SetVariable(variable, ic.DefinedTypes[token])
			ic.SetVariable(variable+"_use", Used)
			return variable
			
		} else if ic.GetFlag(InMethod) && ic.LastDefinedType.Super == token {
			ic.ExpressionType = ic.DefinedTypes[ic.LastDefinedType.Super]
			return ic.LastDefinedType.Name
		
		
		} else if len(ic.DefinedTypes[token].Detail.Elements) == 0 && ic.Peek() == "." {
			ic.Scan('.')
			ic.ExpressionType = InFunction
			var name = ic.Scan(Name)
			return name+"_m_"+token
		
		} else {
			ic.RaiseError()
		}
		
		
	}
	
	if token == "new" {
		var sort = ic.expression()
		if _, ok := ic.DefinedFunctions["new_m_"+ic.ExpressionType.Name]; !ok {
			ic.RaiseError("no new method found for ", ic.ExpressionType.Name)
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
	
	if ic.Translation && !ic.Translated {
		ic.Translated = true
		defer func() {
			ic.Translated = false
		}()
		var err error
		ic.NextToken, err = getTranslation(ic.Language, "en", token)
		if strings.Contains(ic.NextToken, "\n") {
			ic.NextToken = strings.Split(ic.NextToken, "\n")[0]
		}
		println(ic.NextToken, ic.Language)
		if err != nil {
			ic.RaiseError(err)
		}
		return ic.expression() 
	}
	
	token = inflector.Singularize(token)
	if t, ok := ic.DefinedTypes[token]; ok {
		ic.Scan('(')
		ic.Scan(')')
		return ic.NewListOf(t)
	}
	if t, ok := ic.DefinedInterfaces[token]; ok {
		ic.Scan('(')
		ic.Scan(')')
		return ic.NewListOf(t.GetType())
	}
	
	ic.NextToken = token
	ic.ExpressionType = Undefined
	//ic.RaiseError()
	return ""
}

func (ic *Compiler) ScanExpression() string {
	return ic.Shunt(ic.expression())
}
