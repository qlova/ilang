package ilang

/*
	Scan a pipe statement, eg.
		file("text to write")
		file = newfile
*/
func (ic *Compiler) ScanPipeStatement() {
	var name = ic.Scan(0)
	var token = ic.Scan(0)
	switch token {
		case "(":
			argument := ic.ScanExpression()
			ic.Scan(')')
			if ic.ExpressionType != Text && ic.ExpressionType != Array {
				if ic.ExpressionType == Number {
					ic.Assembly("RELAY ", name)
					if argument != "" {
						ic.Assembly("PUSH ", argument)
					} else {
						ic.Assembly("PUSH 0")
					}
					ic.Assembly("IN")
					ic.Assembly("GRAB ", ic.Tmp("discard"))
					return
				}
				ic.RaiseError("Only text and number values can be passed to a pipe call (outside of an expression).")
			}
			ic.Assembly("RELAY ", name)
			ic.Assembly("SHARE ", argument)
			ic.Assembly("OUT")
		case "=":
			value := ic.ScanExpression()
			if ic.ExpressionType != Pipe {
				ic.RaiseError("Only ",Func.Name," values can be assigned to ",name,".")
			}
			ic.Assembly("RELAY ", value)
			ic.Assembly("RELOAD ", name)
		default:
			ic.ExpressionType = Pipe
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != Undefined {
				ic.RaiseError("blank expression!")
			}
	}
}
