package ilang

func (ic *Compiler) ScanStatement() {
	var token = ic.Scan(0)
	
	if t := ic.GetVariable(token); t != Undefined {
	
		ic.NextNextNextToken = ic.NextNextToken
		ic.NextNextToken = ic.NextToken
		ic.NextToken = token
		
		if t.User && t.Name != "something" {
			//A bit of a hack..
			t = string2type["thing"]
		}
		if t.Name == "list" {
			//A bit of a hack..
			t = string2type["list"]
		}
		
		if f, ok := Statements[t]; ok {
			f(ic)
			return
		}
		
		switch t {				
			case Number: 	ic.ScanNumericStatement()
				
			//This may become depreciated.
			case t.IsMatrix(): 					ic.ScanMatrixStatement()
				
			default:
				ic.RaiseError("Unsupported statement!")
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
}
