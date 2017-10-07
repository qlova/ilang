/*
	The save method serialises any type in the I language!
	(hopefully)
*/

package save

import "github.com/qlova/ilang/src"
import "fmt"

func init() {
	ilang.RegisterExpression(ScanSave)
}

func Puts(s string) string {
	var result string
	
	for _, v := range s {
		result += "PUT "+fmt.Sprint(int(v))+"\n"
	}
	
	return result
}

var GeneratedSaveMethods = make(map[string]string)

//TODO peformance this function lol
func GenerateSaveMethodFor(t ilang.Type) string {
	if _, ok := GeneratedSaveMethods[t.GetComplexName()]; ok {
		return ""
	}
	
	var assembly string
	
	assembly += "FUNCTION save_m_"+t.GetComplexName()+"\n"
	
	assembly += "GRAB value\n"
	
	assembly += "ARRAY a\n"
	
	//Get each type member.
	for element, i := range t.Detail.Table {
		
		assembly += Puts(element)
		assembly += "PUT 61\n" //=
		
		//Serialise subtypes.
		subelement := t.Detail.Elements[i]
		if subelement.Push == "PUSH" {
			assembly += "PUSH "+fmt.Sprint(i)+"\n"
			assembly += "PLACE value\n"
			assembly += "GET elem"+fmt.Sprint(i)+"\n"
			assembly += "PLACE a\n"
			assembly += "PUSH elem"+fmt.Sprint(i)+"\n"
			assembly += "PUSH 10\nRUN i_base_number\n"
			assembly += "GRAB txt"+fmt.Sprint(i)+"\n"
			assembly += "JOIN a a txt"+fmt.Sprint(i)+"\n"
			assembly += "PLACE a\n"
			
		} else if subelement == ilang.Text {
			assembly += "PUSH "+fmt.Sprint(i)+"\n"
			assembly += "PLACE value\n"
			assembly += "GET elem"+fmt.Sprint(i)+"\n"
			
			assembly += "PUSH elem"+fmt.Sprint(i)+"\n"
			assembly += "HEAP\n"
			assembly += "GRAB txt"+fmt.Sprint(i)+"\n"
	
			assembly += "JOIN a a txt"+fmt.Sprint(i)+"\n"
			assembly += "PLACE a\n"
		} else {
			assembly += "PUT 63"
		}
		assembly += "PUT 44\n" //,
		
	}
	
	
	assembly += "SHARE a\n"
	
	assembly += "END\n"
	
	GeneratedSaveMethods[t.GetComplexName()] = ""
	
	return assembly
}

func ScanSave(ic *ilang.Compiler) string {
	var token = ic.LastToken
	println(token)
	if token == "save" {	
		ic.Scan('(')
		ic.DisableOwnership = true
	
		value := ic.ScanExpression()
		if ic.ExpressionType.Detail == nil {
			ic.RaiseError("Cannot save values of type ", ic.ExpressionType.GetComplexName())
		}
	
		ic.Library(GenerateSaveMethodFor(ic.ExpressionType))
	
		ic.Assembly("SHARE ", value)
		ic.Assembly("RUN save_m_"+ic.ExpressionType.GetComplexName())
		ic.ExpressionType = ilang.Text
		
		ic.LoadFunction("i_serial_load")
		ic.LoadFunction("i_base_string")
		
		var tmp = ic.Tmp("save")
		ic.Assembly("GRAB ", tmp)
		
		ic.DisableOwnership = false
		ic.Scan(')')
		
		return tmp
	}
	return ""
}
