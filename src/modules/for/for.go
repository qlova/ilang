package f

import "fmt"
import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/types/list"

var ForLoop = ilang.NewFlag()
var Delete = ilang.NewFlag()

func init() {
	ilang.RegisterToken([]string{"for"}, ScanFor)
	ilang.RegisterListener(ForLoop, EndForLoop)
	ilang.RegisterFunction("delete", ilang.Function{
		Inline: true,
		Method: true,
		Assemble: func(ic *ilang.Compiler) string {
			if !ic.GetFlag(ForLoop) {
				ic.RaiseError("delete with zero arguments must be called from within a for loop!")
			}
			
			//Delete things in a for loop.
			return fmt.Sprint(	"PLACE ", ic.GetVariable("i_for_delete").Name, "\n",
								"PUT ", ic.GetVariable("i_for_id").Name)
		},	
	})
}

func EndForLoop(ic *ilang.Compiler) {
	if ic.LastToken != "end" {
		ic.RaiseError("For Loops must end with a 'end' token, found '", ic.LastToken, "'")
	}
	ic.Assembly("REPEAT")
	
	if ic.GetScopedFlag(TypeLoop) {
		EndTypeLoop(ic)
	}
	
	if ic.GetScopedFlag(Delete) {
		var array = ic.GetVariable("i_for_array").Name
		var del = ic.GetVariable("i_for_delete").Name
		
		var collect = ""
		if ic.GetVariable(array).SubType != nil {
			collect = ic.GetVariable(array).SubType.Free("i_pointer")
		}
		
		ic.Assembly(`
	VAR i_i
	VAR i_test
	LOOP
		SGE i_test i_i #`+del+`
		IF i_test
			BREAK
		END
	
		IF #`+array+`
		ELSE
			BREAK
		END
	
		PLACE `+del+`
		PUSH i_i
		GET i_index
		
		PLACE `+array+`
		PUSH i_index
		GET i_pointer
		`+collect+`
		
		PLACE `+array+`		
		PUSH -1
		GET i_swapdex
		
		PUSH i_index
		SET i_swapdex

		POP end
		ADD end 0 0
		
		ADD i_i i_i 1
	REPEAT
				`)
		ic.LoseScope()
	}
	ic.Assembly("END")
}

func ScanFor(ic *ilang.Compiler) {
	
	var name = ic.Scan(ilang.Name)
	
	var name2 string
	
	if ic.Peek() == "," {
		ic.Scan(',')
		name2 = ic.Scan(ilang.Name)
	}

	var OverList = false
	
	switch ic.Scan(0) {
		case "=":
			ic.AssembleVar(name, ic.ScanExpression())
			ic.Scan(',')
			ic.Assembly("IF 1")
			ic.Assembly("LOOP")
			ic.GainScope()
			condition := ic.ScanExpression()
			ic.Assembly("IF ", condition)
			ic.Assembly("	ADD ",name," ",name," 0")
			ic.Assembly("ELSE")
			ic.Assembly("	BREAK")
			ic.Assembly("END")
			ic.SetFlag(ForLoop)
	
		case "over":
			var token = ic.Scan(0)
			if token == "[" {
				a := ic.ScanExpression()
				ic.Scan(',')
				b := ic.ScanExpression()
				ic.Scan(']')
			
				condition := ic.Tmp("over")
				backup := ic.Tmp("backup")
			
				ic.Assembly("IF 1\n",	
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
				if name != "each" { 
					ic.SetVariable(name, ilang.Number)
				}
				ic.SetFlag(ForLoop)
				return
			}
			ic.NextToken = token
			OverList = true
			fallthrough
		case "in":
			var peek = ic.Scan(0)
			if t, ok := ic.DefinedTypes[peek]; ok {
				ScanTypeLoop(ic, t, name, name2)
				return
			} else {
				ic.NextToken = peek
			}
			var array = ic.ScanExpression()
			
			if ic.ExpressionType == list.Type {
				//We can ignore this loop.
				for {
					var token = ic.Scan(0)
					if token == "end" {
						return
					}
				}
			}
			
			if ic.ExpressionType.Push != "SHARE" {
				ic.RaiseError("Cannot iterate over "+array+" (", ic.ExpressionType.Name, ")")
			}
			
			var condition = ic.Tmp("in") 

			var i, v, vo string
			if name2 != "" {
				i = name
				v = name2
			} else {
				i = ic.Tmp("i")
				v = name
			}
			
			if OverList {
				i = name
			}
			
			vo = v
			if ic.ExpressionType.Name == list.Type.Name && ic.ExpressionType.SubType.Push == "SHARE" {
				v += "_address"
			}
			
			backup := ic.Tmp("backup")
			del := ic.Tmp("delete")
			
			if OverList {
			ic.Assembly(`
IF 1
ARRAY %v
VAR %v
VAR %v
LOOP
	VAR %v
	ADD %v 0 %v
	SGE %v %v #%v
	IF %v
		BREAK
	END
	ADD %v %v 1
`, del, i,backup, condition, i, backup,  condition, i, array, condition, backup, i)			
			} else {
			ic.Assembly(`
IF 1
ARRAY %v
VAR %v
VAR %v
LOOP
	VAR %v
	ADD %v 0 %v
	SGE %v %v #%v
	IF %v
		BREAK
	END
	PLACE %v
	PUSH %v
	GET %v
	ADD %v %v 1
`, del, i,backup, condition, i, backup,  condition, i, array, condition, array, i, v, backup, i)
		if ic.ExpressionType == ilang.Array || ic.ExpressionType == ilang.Text {
			//ic.Assembly("ADD %v %v %v", name, 0, i)
		}
	}

			if ic.ExpressionType.Name == list.Type.Name && ic.ExpressionType.SubType.Push == "SHARE" {
				ic.Assembly("PUSH ", v)
				ic.Assembly("HEAP")
				ic.Assembly("GRAB ", vo)
			}
			
			ic.GainScope()
			ic.SetFlag(Delete)
			ic.SetVariable("i_for_delete", ilang.Type{Name:del})
			ic.SetVariable("i_for_id", ilang.Type{Name:i})
			ic.SetVariable("i_for_array", ilang.Type{Name:array})
			
			ic.GainScope()
			ic.SetVariable(i, ilang.Number)
			
			if !OverList {
			
			
			if ic.ExpressionType.Name == list.Type.Name  {
				ic.SetVariable(vo, *ic.ExpressionType.SubType)
				ic.SetVariable(vo+".", ilang.Protected)
				
				
			} else if ic.ExpressionType == ilang.Text {
				ic.SetVariable(vo, ilang.GetType("letter"))
			} else {
				ic.RaiseError("Cannot find values inside ", ic.ExpressionType)
			}
			}
			
			ic.SetFlag(ForLoop)
			return
	}
}
