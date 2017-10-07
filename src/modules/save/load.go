package save

import "github.com/qlova/ilang/src"
import "fmt"

var GeneratedLoadMethods = make(map[string]string)

func init() {
	ilang.RegisterConstructor(func(ic *ilang.Compiler, t ilang.Type) {	
		ic.Library(GenerateLoadMethodFor(t))
		
		ic.DefinedFunctions[t.GetComplexName()+"_m_text"] = 
			ilang.Method(t, true, "RUN "+t.GetComplexName()+"_m_text")
	})

	ilang.RegisterFunction("i_serial_load", ilang.Function{Exists:false, Data: `
FUNCTION i_serial_load
	PULL sep
	GRAB string
	PULL i
	ARRAY i_string1
	SHARE i_string1
	GRAB result
	PUSH 32
	PULL value
	LOOP
		VAR i_operator2
		SGT i_operator2 i #string
		IF i_operator2
			SHARE result
			PUSH i
			RETURN
		END
		PLACE string
		PUSH i
		GET i_value3
		ADD value 0 i_value3
		VAR i_operator4
		SEQ i_operator4 value sep
		IF i_operator4
			VAR i_operator5
			SUB i_operator5 i 1
			PLACE string
			PUSH i_operator5
			GET i_value6
			VAR i_operator7
			SNE i_operator7 i_value6 92
			IF i_operator7
				SHARE result
				PUSH i
				RETURN
			ELSE
					PLACE result
					PUT value
			END
		ELSE
				PLACE result
				PUT value
		END
		ADD i i 1
	REPEAT
RETURN
	`})
}

//TODO peformance this function lol
func GenerateLoadMethodFor(t ilang.Type) string {
	if _, ok := GeneratedLoadMethods[t.GetComplexName()]; ok || t.Detail == nil {
		return ""
	}
	
	var assembly string
	
	assembly += "FUNCTION "+t.GetComplexName()+"_m_text\n"
	
	assembly += "GRAB value\n"
	
	assembly += "PUSH "+fmt.Sprint(len(t.Detail.Elements))+"\nMAKE\nGRAB a\n"
	
	//Get each type member.
	for element, _ := range t.Detail.Table {
		
		assembly += "ARRAY "+element+"\n"
		assembly += Puts(element)
		
	}
	
	assembly += "VAR i\nVAR condition\n"
	assembly += "LOOP\n"
	
		
		//Damn assembly gonna have to do it manually I think..
		
		assembly += "PUSH i\nPUSH 61\nSHARE value\nRUN i_serial_load\n"
		assembly += "GRAB test\n"
		assembly += "PULL i2\nADD i 0 i2\n"
		
		assembly += "SGT condition i #value\nIF condition\nBREAK\nEND\n"
		
	//Get each type member.
	var ends int
	for element, i := range t.Detail.Table {
		assembly += "SHARE test\nSHARE "+element+"\nRUN strings.equal\nPULL test_"+element+"\n"
		assembly += "IF test_"+element+"\n"
		
		assembly += "ADD i i 1\n"
		assembly += "PUSH i\nPUSH 44\nSHARE value\nRUN i_serial_load\n"
		assembly += "PULL i2\nADD i 0 i2\n"
		
		subelement := t.Detail.Elements[i]
		if subelement.Push == "PUSH" {
			assembly += "PUSH 10\nRUN i_base_string\n"
			assembly += "PULL n\n"
			
			assembly += "PUSH "+fmt.Sprint(i)+"\n"
			assembly += "PLACE a\nSET n\n"
		}
		
		if subelement == ilang.Text {
			assembly += "PUSH 0\nHEAP\n"
			assembly += "PULL n\n"
			
			assembly += "PUSH "+fmt.Sprint(i)+"\n"
			assembly += "PLACE a\nSET n\n"
		}
		
		assembly += "ELSE\n"
		ends++
	}
	
	assembly += "PUSH i\nPUSH 44\nSHARE value\nRUN i_serial_load\n"
	assembly += "GRAB test\n"
	assembly += "PULL i2\nADD i 0 i2\n"
	
	for i:=0; i < ends; i++ {
		assembly += "END\n"
	}

		assembly += "ADD i i 1\n"
	
	assembly += "REPEAT\n"
	
	GeneratedLoadMethods[t.GetComplexName()] = ""
	
	assembly += "SHARE a\n"
	
	assembly += "END\n"
	
	return assembly
}

