package main

func (ic *Compiler) ScanForLoop() {
	
	var name = ic.Scan(Name)
	var name2 string
	
	if ic.Peek() == "," {
		ic.Scan(',')
		name2 = ic.Scan(Name)
	}
	
	switch ic.Scan(0) {
		case "=":
			ic.AssembleVar(name, ic.ScanExpression())
			ic.Scan(',')
			ic.Assembly("LOOP")
			ic.GainScope()
			condition := ic.ScanExpression()
			ic.Assembly("IF ", condition)
			ic.Assembly("	ADD ",name," ",name," 0")
			ic.Assembly("ELSE")
			ic.Assembly("	BREAK")
			ic.Assembly("END")
			ic.SetFlag(ForLoop)
		
		case "in":
			var array = ic.ScanExpression()
			
			var condition = ic.Tmp("in") 

			var i, v string
			if name2 != "" {
				i = name
				v = name2
			} else {
				i = ic.Tmp("i")
				v = name
			}
			backup := ic.Tmp("backup")
			
			ic.Assembly(`
VAR %v
VAR %v
LOOP
	VAR %v
	ADD %v 0 %v
	PLACE %v
	PUSH %v
	GET %v
	SGE %v %v #%v
	IF %v
		BREAK
	END
	ADD %v %v 1
`, i,backup, condition, i, backup, array, i, v, condition, i, array, condition, backup, i)

			ic.GainScope()
			ic.SetVariable(i, Number)
			ic.SetVariable(v, Number)
			ic.SetFlag(ForLoop)
			return
	
		case "over":
			ic.Scan('[')
			a := ic.ScanExpression()
			ic.Scan(',')
			b := ic.ScanExpression()
			ic.Scan(']')
			
			condition := ic.Tmp("over")
			backup := ic.Tmp("backup")
			
			ic.Assembly("",
				"VAR ",name,"\n",
				"VAR ",backup,"\n",
				"ADD ",backup," 0 ",a,"\n",
				"ADD ",name," 0 ",a,"\n",
				"LOOP\n",
				"	VAR ",condition,"\n",
				"	SNE ",condition," ",name," ",b,"\n",
				"	ADD ",name," 0 ",backup,"\n",
				"	IF ",condition,"\n",
				"		SLT ",condition," ",name," ",b,"\n",
				"		IF ",condition,"\n",
				"			ADD ",backup," ",name," 1\n",
				"		ELSE\n",
				"			SUB ",backup," ",name," 1\n",
				"		END\n",
				"		SEQ ",condition," ",a," ",b,"\n",
				"		IF ",condition,"\n",
				"			BREAK\n",
				"		END\n",
				"	ELSE\n",
				"		SEQ ",condition," ",a," ",b,"\n",
				"		IF ",condition,"\n",
				"			ADD ",name," ",name," 1\n",
				"       ELSE\n",
				"			BREAK\n",
				"		END\n",
				"	END\n",
			)
			ic.GainScope()
			ic.SetVariable(name, Number)
			ic.SetFlag(ForLoop)
			return
	}
	
}
