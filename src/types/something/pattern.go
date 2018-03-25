package something

import "github.com/qlova/ilang/src/modules/method"
import "github.com/qlova/ilang/src"
import "fmt"

func init() {
	ilang.RegisterCast(Cast)
}

func Cast(ic *ilang.Compiler, name string, a ilang.Type, b ilang.Type) string {
	if b.Class != nil && b.Class.Name == "something" && len(*b.Functions) > 0 {
		
		var asm = ""
		
		var array = ic.Tmp("array")
		var pointer =  ic.Tmp("pointer")
		
		asm += fmt.Sprintln("ARRAY", array)
		
		if a.Empty() {
			asm += fmt.Sprintln("PLACE", array)
			asm += fmt.Sprintln("PUT 0")
		} else {
			asm += ic.CreatePointer(a, pointer, name)
			asm += fmt.Sprintln("PLACE", array)
			asm += fmt.Sprintln("PUT", pointer)
		}
		
		if a.Empty() {
			asm += fmt.Sprintln("PUT 0")
		} else if a.Push == "PUSH"  {
			asm += fmt.Sprintln("PUT 1")
		} else if a.Push == "SHARE" {
			asm += fmt.Sprintln("PUT 2")
		} else if a.Push == "RELAY" {
			asm += fmt.Sprintln("PUT 3")
		}
		
		//Add each method as a pointer value in the array.
		for _, InterfaceMethod := range *b.Functions {
			
			//Check if the type has the method.
			if f := method.Get(ic, a, InterfaceMethod.Name); f != nil {
				
				if f.Inline && f.Data != "" {
					ic.ExportInlineFunction(f.Name, *f)
				}
				
				if len(f.Returns) != len(InterfaceMethod.Returns) {
					ic.RaiseError("Cannot convert to ", b.Name, ". The ",InterfaceMethod.Name, " method has the wrong number of return values!")
				}
				
				if len(f.Returns) > 0 && len(InterfaceMethod.Returns) > 0 && f.Returns[0] != InterfaceMethod.Returns[0] {
					ic.RaiseError("Cannot convert to ", b.Name, ". The ",InterfaceMethod.Name, " method of the "+a.GetComplexName()+" type returns a ", f.Returns[0].GetComplexName(), "!")
				}
				
				ic.LoadFunction(f.Name)
				
				var address = ic.Tmp("address")
				
				if f.Inline && f.Data == "" {
					asm += fmt.Sprintln("PUT 0")
				} else {
					asm += fmt.Sprintln("SCOPE", f.Name)
					asm += fmt.Sprintln("PUSH 0")
					asm += fmt.Sprintln("HEAPIT")
					asm += fmt.Sprintln("PULL", address)
					asm += fmt.Sprintln("PUT", address)
				}
				
				
				
			} else { 
				ic.RaiseError("Cannot convert to ", b.Name, ". There is no ", InterfaceMethod.Name, " method for the ", a.Name, " type.")
			}
			
		}
		
		asm += fmt.Sprintln("SHARE", array)
		
		return asm
		
	}
	return ""
}

func ScanInterface(ic *ilang.Compiler) ilang.Type {
	
	var af = make([]ilang.Function, 0)
	var intf = ilang.Type{Name:"something", Push: "SHARE", Pop: "GRAB", Functions: &af, Super: ic.Tmp("interface")}
	
	ic.Scan('{')
	
	for {
		var method = ic.Scan(0)
		
		if method == "}" {
			break
		}
		if method == "\n" {
			continue
		}
		

		var f = ilang.Function{Name: method}
		
		ic.Scan('(')
		//Scan argument types.
		if ic.Peek() != ")" {
			for {
				var a = ic.ScanSymbolicType()
				f.Args = append(f.Args, a)
				if ic.Peek() == ")" {
					break
				} else {
 					ic.Scan(',')
				}
			}
		}
		ic.Scan(')')
		var r = ic.ScanSymbolicType()
		f.Returns = append(f.Returns, r)
		
		*intf.Functions = append(*intf.Functions, f)
		
		ic.DefinedFunctions[f.Name+"_m_"+intf.Super] = ilang.Function{
			Name: f.Name+"_m_"+intf.Super,
			Data: `
FUNCTION `+f.Name+"_m_"+intf.Super+`,
	GRAB interface
	PUSH 1
	PLACE interface
	GET t
	
	VAR condition
	SEQ condition t 0
	IF condition
		
	ELSE
		SEQ condition t 1
		
		IF condition
		
			PUSH 0
			GET v
			PUSH v
			
		ELSE
		
			SEQ condition t 2
			
			IF condition
			
				PUSH 0
				GET address
				PUSH address
				HEAP
			
			
			ELSE
			
				PUSH 0
				GET address
				PUSH address
				HEAP
			END
		
		END
	END
	
	PUSH `+fmt.Sprint(len(*intf.Functions)+1)+`
	GET method_address
	
	IF method_address
		
		PUSH method_address
		HEAPIT
		TAKE method
		
		EXE method
	END
RETURN
`}
		
		var token = ic.Scan(0)
		if token == "}" {
			break
		} else if token != "," && token != "\n" {
			ic.RaiseError()
		}
	}
	
	return intf
}
	
