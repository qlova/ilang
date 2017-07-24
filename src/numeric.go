package ilang

/*
	Scans a numeric statement such as:
		number = 3
		variable += 4
		letter = 'a'
*/
func (ic *Compiler) ScanNumericStatement() {
	var name = ic.Scan(0)
	var numeric = ic.GetVariable(name)
	
	if name == "error" {
		name = "ERROR"
	}
	
	var token = ic.Scan(0)
	
	switch token {
		case "=":
			value := ic.ScanExpression()
			if ic.ExpressionType.Push != "PUSH" {
				ic.RaiseError("Only numeric values can assigned to ",name,".")
			}
			if ic.ExpressionType != numeric {
				ic.RaiseError("Cannot add %s to %s", ic.ExpressionType, numeric)
			}
				
			ic.Assembly("ADD %v %v %v", name, 0, value)
		default:
			ic.ExpressionType = numeric
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != Undefined {
				ic.RaiseError("blank expression! hmm?")
			}
			ic.ExpressionType = numeric
	}
}
