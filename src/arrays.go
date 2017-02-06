package main

//Index a 1D array, returns the value at the index.
func (ic *Compiler) Index(array, index string) string {
	var result = ic.Tmp("index")
	ic.Assembly("PLACE ", array)
	ic.Assembly("PUSH ", index)
	ic.Assembly("GET ", result)
	return result
}

//Index a matrix at pos (x,y) returns the value.
func (ic *Compiler) IndexMatrix(matrix, x, y string) string {
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

	var result = ic.Tmp("result")
	ic.Assembly("PUSH ", ytmp)
	ic.Assembly("GET ", result)
	return result
}

//Set the value at pos (x,y) in the matrix.
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

//Set the value of an array at the specified index.
func (ic *Compiler) Set(array, index, value string) {
	ic.Assembly("PLACE ", array)
	ic.Assembly("PUSH ", index)
	if ic.GetVariable(array).List {
		
		//Garbage collect the previous item in the array.
		var collect = ic.Index(array, index)
		ic.Assembly("IF ", collect)
			ic.Assembly("PUSH ", collect)
			ic.Assembly("HEAP")
			ic.Assembly("RUN collect_m_", ic.GetVariable(array).Name)
			ic.Assembly("MUL ", collect, " ", collect, " -1")
			ic.Assembly("PUSH ", collect)
			ic.Assembly("HEAP")
		ic.Assembly("END")
		
		
		var tmp = ic.Tmp("index")
		ic.Assembly("SHARE ", value)
		ic.Assembly("PUSH 0")
		ic.Assembly("HEAP")
		ic.Assembly("PULL ", tmp)
		
		ic.Assembly("PLACE ", array)
		ic.Assembly("PUSH ", index)
		ic.Assembly("SET ", tmp)
		
	} else {
		ic.Assembly("SET ", value)
	}
}

//Scan an array for example: [0, 1, 2, 3, 4, 5]
func (ic *Compiler) ScanArray() string {
	var id = ic.Tmp("array")
	
	var result = Array
	
	//Decimal array $[]
	if ic.Peek() == "[" {
		ic.Scan('[')
		result.Decimal = true
	}
	
	//This is the size of the array, eg. 
	//		var a = [...50] 
	// (an array with 50 elements)
	if ic.Peek() == "." {
		ic.Scan('.')
		ic.Scan('.')
		size := ic.ScanExpression()
		ic.Scan(']')
		ic.Assembly("PUSH ", size)
		ic.Assembly("MAKE")
		ic.Assembly("GRAB ", id)
		ic.ExpressionType = result
		return id
	} else {
		ic.Assembly("ARRAY ", id)
	}
	
	ic.ExpressionType = result
		
	if ic.Peek() == "]" {
		ic.Scan(0)
		return id
	}

	//Push all the values.
	for {
		value := ic.ScanExpression()
		ic.Assembly("PLACE ", id)
		ic.Assembly("PUT ", value)
		
		token := ic.Scan(0)
		if token != "," {
			if token != "]" {
				ic.RaiseError()
			}
			break
		}
	}
	
	ic.ExpressionType = result
	return id
}
