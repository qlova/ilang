package ilang

/*
	Shunts a table eg.
		table["key"]
*/
func (ic *Compiler) ShuntTable(name string) string {
	var index = ic.ScanExpression()
	ic.Scan(']')
	if ic.ExpressionType != Text {
		ic.RaiseError("A Table must have a text index!")
	}
	var tableval = ic.Tmp("tableval")
	
	ic.Assembly("SHARE %s", index)
	ic.Assembly("PUSH %s", name)
	ic.Assembly(ic.RunFunction("table_get"))
	ic.Assembly("PULL %s", tableval)
	
	ic.ExpressionType = Number
	
	return ic.Shunt(tableval)
}

/*
	Scans a table statement such as:
		table["key"] = value
*/
func (ic *Compiler) ScanTableStatement() {
	var table = ic.Scan(0)
	
	ic.Scan('[')
	var index = ic.ScanExpression()
	if ic.ExpressionType != Text {
		ic.RaiseError("Table must have text index.")
	}
	ic.Scan(']')
	ic.Scan('=')
	var value = ic.ScanExpression()
	if ic.ExpressionType != Number {
		ic.RaiseError("Table can only take numbers.")
	}
	
	var tmp = ic.Tmp("newtableref")
	
	ic.Assembly("PUSH %s", table)
	ic.Assembly("SHARE %s", index)
	ic.Assembly("PUSH %s", value)
	ic.Assembly(ic.RunFunction("table_set"))
	ic.Assembly("PULL %s", tmp)
	ic.Assembly("ADD %s %s %v", table, tmp, 0)
}
