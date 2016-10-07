package main

func (ic *Compiler) Index(name, index string) string {
	var result = ic.Tmp("index")
	ic.Assembly("PLACE ", name)
	ic.Assembly("PUSH ", index)
	ic.Assembly("GET ", result)
	return result
}

func (ic *Compiler) Set(name, index, value string) {
	ic.Assembly("PLACE ", name)
	ic.Assembly("PUSH ", index)
	ic.Assembly("SET ", value)
}

func (ic *Compiler) ScanArray() string {
	var id = ic.Tmp("array")

	ic.Assembly("ARRAY ", id)
	
	if ic.Peek() == "]" {
		ic.ExpressionType = Array
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
	
	ic.ExpressionType = Array
	return id
}
