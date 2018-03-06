package ilang

/*
	Scans a numeric statement such as:
		number = 3
		variable += 4
		letter = 'a'
*/
func (ic *Compiler) ScanTextStatement() {
	var name = ic.Scan(0)
	
	if name == "error" {
		name = "ERROR"
	}
	
	var token = ic.Scan(0)
	
	switch token {
		case "=":
			value := ic.ScanExpression()
			if ic.ExpressionType != Text {
				ic.RaiseError("Only text values can assigned to ",name,".")
			}
				
			ic.Assembly("SHARE %v\nRENAME %v", value, name)
		default:
			ic.ExpressionType = Text
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != Undefined {
				ic.RaiseError("blank expression! hmm?")
			}
			ic.ExpressionType = Text
	}
}

func (ic *Compiler) ShuntStatement(t Type) {
		if t.User && t.Name != "something" {
			//A bit of a hack..
			t = string2type["thing"]
		}
		if t.Name == "list" {
			//A bit of a hack..
			t = string2type["list"]
		}
		if t.Name == "table" {
			//A bit of a hack..
			t = string2type["table"]
		}
		
		if f, ok := Statements[t]; ok {
			f(ic)
			return
		}
		
		switch t {				
			case Number: 	ic.ScanNumericStatement()
			
			case Text:
				ic.ScanTextStatement() 
				
			//This may become depreciated.
			case t.IsMatrix(): 					
				ic.ScanMatrixStatement()
				
			default:
				if t.Class != nil {
					ic.ShuntStatement(*t.Class)
				} else {
					ic.RaiseError("Unsupported statement!")
				}
		}
}

func (ic *Compiler) ScanStatement() {
	var token = ic.Scan(0)
	
	if t := ic.GetVariable(token); t != Undefined {
		
		ic.NextNextNextToken = ic.NextNextToken
		ic.NextNextToken = ic.NextToken
		ic.NextToken = token
		
		if t.Class != nil && ic.Peek() == "." {
			ic.ShuntStatement(string2type["thing"])
		} else {
			ic.ShuntStatement(t)
		}
	
		return
	//Function Calls and things.
	} else if _, ok := ic.DefinedFunctions[token]; ok {
		ic.Scan('(')
		ic.ScanFunctionCall(token)
		ic.Scan(')')
		return
	}
	
	ic.NextToken = token
	
	ic.ScanExpression()	
	if ic.ExpressionType == Undefined {
		ic.RaiseError(token, " undefined!")
	}
	
	ic.ShuntStatement(ic.ExpressionType)
}
