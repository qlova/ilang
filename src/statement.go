package ilang

func (ic *Compiler) ScanStatement() {
	var token = ic.Scan(0)
	
	if t := ic.GetVariable(token); t != Undefined {
	
		ic.NextNextNextToken = ic.NextNextToken
		ic.NextNextToken = ic.NextToken
		ic.NextToken = token
		switch t {
			case Table: 						ic.ScanTableStatement()
				
			case Number, Decimal, Letter, Set: 	ic.ScanNumericStatement()
				
			//This may become depreciated.
			case t.IsMatrix(): 					ic.ScanMatrixStatement()
			
			case Array, Text, t.IsArray():		ic.ScanArrayStatement()
				
			case List, t.IsList():				ic.ScanListStatement()
				
			case Pipe:							ic.ScanPipeStatement()
				
			case Func:							ic.ScanFuncStatement()
			
			case t.IsSomething():				ic.ScanSomethingStatement()
				
			case User, t.IsUser():				ic.ScanUserStatement()
				
			default:
				ic.RaiseError()
		}
	
	//Function Calls and things.
	} else if _, ok := ic.DefinedFunctions[token]; ok {
		var check = ic.Scan(0)
		if check == "(" {
			ic.ScanFunctionCall(token)
			ic.Scan(')')
		} else if check == "@" {
			var variable = ic.expression()
			ic.Assembly("%v %v", ic.ExpressionType.Push, variable)
			ic.Scan('(')
			ic.ScanFunctionCall(token+"_m_"+ic.ExpressionType.Name)
			ic.Scan(')')
		} else {
			ic.RaiseError()
		}
	
	} else if ic.GetFlag(InMethod) {
		ic.NextToken = token
		ic.ScanExpression()	
		if ic.ExpressionType == Undefined {
			ic.RaiseError(token, " undefined!")
		}
	
	} else {
		ic.RaiseError()
	}
}
