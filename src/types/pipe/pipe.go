package pipe

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/types/letter"

var Type = ilang.NewType("pipe", "RELAY", "TAKE")

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	//Types.
	if token == "|" {
		ic.Scan('|')
		var name = ic.Tmp("pipe")
		ic.Assembly("PIPE ", name)
		
		ic.ExpressionType = Type
		
		return name
	}
	return ""
}

func init() {
	ilang.RegisterStatement(Type, ScanPipeStatement)
	ilang.RegisterSymbol("|", ScanPipeSymbol)
	ilang.RegisterExpression(ScanExpression)
	
	ilang.RegisterFunction("open_m_text", ilang.Method(Type, true, "OPEN"))
	ilang.RegisterFunction("open_m_number", ilang.Function{Exists:true, Returns:[]ilang.Type{Type}, Data: `
FUNCTION open_m_number
	PULL value
	ARRAY tmp
	PUT value
	SHARE tmp
	OPEN
RETURN
	`})
	ilang.RegisterFunction("read_m_pipe", ilang.InlineFunction(nil, "PUSH 0\nIN", []ilang.Type{ilang.Text}))
	
	ilang.RegisterFunction("close", ilang.Function{Exists:true, Args:[]ilang.Type{Type}, Data: `
FUNCTION close
	TAKE file
	CLOSE file
RETURN
	`})
	
	ilang.RegisterShunt("(", func(ic *ilang.Compiler, name string) string {
		//Calling pipes.
		if ic.ExpressionType == Type {
		
			token := ic.Scan(0)
			if token == ")" {
				//Read default from the pipe.
				var r = ic.Tmp("read")
				ic.Assembly("RELAY ", name)
				ic.Assembly("PUSH 0")
				ic.Assembly("IN")
				ic.Assembly("GRAB ", r)
				ic.ExpressionType = ilang.Text
				return ic.Shunt(r)	
			}
			
			ic.NextToken = token
							
			argument := ic.ScanExpression()
			
			switch ic.ExpressionType {
				case letter.Type:
					var r = ic.Tmp("reada")
					ic.Assembly("RELAY ", name)
					ic.Assembly("PUSH ", argument)
					ic.Assembly("RUN reada_m_pipe")
					ic.Assembly("GRAB ", r)
					ic.LoadFunction("reada_m_pipe")
					ic.ExpressionType = ilang.Text
					ic.Scan(')')
					return ic.Shunt(r)	
				case ilang.Number:
					var r = ic.Tmp("read")
					ic.Assembly("RELAY ", name)
					ic.Assembly("PUSH ", argument)
					ic.Assembly("IN")
					ic.Assembly("GRAB ", r)
					ic.ExpressionType = ilang.Text
					ic.Scan(')')
					return ic.Shunt(r)
				default:
					ic.RaiseError("Cannot call a pipe with a ", ic.ExpressionType.Name, " argument in an expression!")
			}

		}
		return ""
	})
}

func ScanPipeSymbol(ic *ilang.Compiler) ilang.Type {
	ic.Scan('|')
	return Type
}

/*
	Scan a pipe statement, eg.
		file("text to write")
		file = newfile
*/
func ScanPipeStatement(ic *ilang.Compiler) {
	var name = ic.Scan(ilang.Name)
	var token = ic.Scan(0)
	switch token {
		case "(":
			argument := ic.ScanExpression()
			ic.Scan(')')
			if ic.ExpressionType != ilang.Text && ic.ExpressionType != ilang.Array {
				if ic.ExpressionType == ilang.Number {
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
			if ic.ExpressionType != Type {
				ic.RaiseError("Only ",Type.Name," values can be assigned to ",name,".")
			}
			ic.Assembly("RELAY ", value)
			ic.Assembly("RELOAD ", name)
		default:
			ic.ExpressionType = Type
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != ilang.Undefined {
				ic.RaiseError("blank expression!")
			}
	}
}
