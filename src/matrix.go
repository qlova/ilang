package ilang

/* 
	Scans a matrix, like:
		matrix[x][y] = value
		matrix = newmatrix
*/
func (ic *Compiler) ScanMatrixStatement() {
	var name = ic.Scan(0)
	var t = ic.GetVariable(name)
	
	var token = ic.Scan(0)
	switch token {
		case "[":
			var x = ic.ScanExpression()
			ic.Scan(']')
			ic.Scan('[')
			var y = ic.ScanExpression()
			ic.Scan(']')
			ic.Scan('=')
			var value = ic.ScanExpression()
			
			ic.SetMatrix(name, x, y, value)
		
		case "=":
		
			value := ic.ScanExpression()
			if ic.ExpressionType != t {
				ic.RaiseError("Only ",t.Name," values can be assigned to ",name,".")
			}
			
			if _, ok := ic.LastDefinedType.Detail.Table[name]; ic.GetFlag(InMethod) && ok {
				ic.SetUserType(ic.LastDefinedType.Name, name, value)	
			} else {									
				ic.Assembly("PLACE ", value)
				ic.Assembly("RENAME ", name)
			}
		default:
			ic.ExpressionType = t
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != Undefined {
				ic.RaiseError("blank expression!")
			}
	}	
}

//Set the value at pos (x,y) in the matrix.
//Value must be numeric.
func (ic *Compiler) SetMatrix(matrix, x, y, value string) {
	var width = ic.Tmp("width")
	ic.Assembly("PLACE ", matrix)
	ic.Assembly("PUSH 0")
	ic.Assembly("GET ", width)
	
	var height = ic.Tmp("height")
	ic.Assembly("PLACE ", matrix)
	ic.Assembly("PUSH 1")
	ic.Assembly("GET ", height)

	var ytmp = ic.Tmp("y")
	var xtmp = ic.Tmp("x")
	ic.Assembly("VAR %v\nVAR %v", xtmp, ytmp)
	ic.Assembly("MOD %v %v %v", xtmp, x, width)
	ic.Assembly("ADD %v %v %v", xtmp, xtmp, 2)

	ic.Assembly("MOD %v %v %v", ytmp, y, height)
	
	ic.Assembly("MUL %v %v %v", ytmp, ytmp, width)

	ic.Assembly("ADD %v %v %v", ytmp, ytmp, xtmp)

	ic.Assembly("PUSH ", ytmp)
	ic.Assembly("SET ", value)
}
