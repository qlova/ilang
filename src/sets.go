package ilang

func (ic *Compiler) ScanSet() string {
	var id = ic.Tmp("set")
	
	ic.Assembly("VAR ", id)
	ic.Assembly("ADD ", id, " 1 1")
	
	for {
		var token = ic.Scan(0)
		if token == ">" {
			break
		}
		
		if prime, ok := ic.SetItems[token]; ok {
			ic.Assembly("MUL %s %s %v", id, id, prime)
		} else {
			ic.SetItems[token] = Primes[ic.SetItemCount]
			ic.SetItemCount++
			
			ic.Assembly("MUL %s %s %v", id, id, Primes[ic.SetItemCount])
		}
	}
	
	ic.ExpressionType = Set
	
	return id
}
