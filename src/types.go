package main

import (
	"text/scanner"
	"fmt"
	"io"
	"os"
)

type TYPE int

//These are the 4 types in I.
const (
	UNDEFINED TYPE = iota
	
	MULTIPLE
	
	ITYPE
	
	FUNCTION
	STRING
	ARRAY
	NUMBER
	LETTER
	FILE
	
	USER
	SOMETHING
)


func (t TYPE) String() string {
	if t >= USER {
		return DefinedTypes[t-USER].Name
	}
	return map[TYPE]string{
		FUNCTION:"function", 
		STRING: "text",
		ARRAY:  "array",
		LETTER: "letter",
		NUMBER:"number", 
		ITYPE: "type", 
		FILE: "pipe", 
		SOMETHING: "Something",
		UNDEFINED:"undefined",
	}[t]
}

func (t TYPE) Push() string {
	if t >= USER {
		return "SHARE"
	}
	return map[TYPE]string{
		FUNCTION:"RELAY", 
		ITYPE:"PUSH", 
		ARRAY:"SHARE", 
		STRING:"SHARE",
		LETTER:"PUSH", 
		NUMBER:"PUSH",
		FILE:"RELAY",
		UNDEFINED:"",
	}[t]
}

func (t TYPE) Pop() string {
	if t >= USER {
		return "GRAB"
	}
	return map[TYPE]string{
		FUNCTION:"TAKE", 
		STRING:"GRAB",
		NUMBER:"PULL",
		ITYPE: "PULL",
		ARRAY: "GRAB",
		LETTER: "PULL",
		FILE:"TAKE",
		UNDEFINED:"",
	}[t]
}

var MethodListHeaped = map[TYPE][]int{}

/*func GenMethodList(output io.Writer, t TYPE) {
	methods := t.methods()
	
	//Put the methods on the heap.
	if _, ok := MethodListHeaped[t]; !ok {
		for i := range methods {
			fmt.Fprintf(output, "SCOPE %s\nPUSH 0\nHEAPIT\n", methods[])
		}
	}
	
	unique++
	fmt.Fprintf(output, "ARRAY %s\n", unique)
	for i := range methods {
		fmt.Fprintf(output, "PUT %s\n", methods[])
	}
}

//Cache.
var TypeMethods = map[TYPE][]string{}
//This cannot be called before all methods have been defined!
func (t TYPE) Methods() (list []string) {
	if l, ok := TypeMethods[t]; ok {
		return l	
	}
	for key, f := range functions {
		split := strings.Split(key, "_m_")
		if len(split) > 0 {
			t2 := split[1]
			
			if t.String() == t2 {
				list = append(list, key)
			}
		} else if len(f.Args) == 1 && f.Args[0] == t {
			list = append(list, key)
		}
	}
	TypeMethods[t] = list
	return
}*/


var DefinedTypes = []UserType{
	UserType{},
	/*
		type Something {
			data, type, []methods
		}
	*/	
	UserType{
		Name: "Something",
		Elements: []TYPE{ UNDEFINED, ITYPE, UNDEFINED},
		Table: map[string]int{"type":1, "data":0},
	},
}
var StringToType map[string]TYPE = map[string]TYPE{
	"Something":SOMETHING, 
	"number":NUMBER, 
	"text":STRING, 
	"array":ARRAY, 
	"letter":LETTER,
	"pipe":FILE,
	"function":FUNCTION,
	"type": ITYPE,
}

//This is so methods know what type they are acting on.
var LastDefinedType TYPE

type UserType struct {
	Name string
	Elements []TYPE
	Table map[string]int
	SubElements map[int]TYPE
}

func NewType(name string) UserType {
	return UserType {
		Name: name,
		Table: make(map[string]int),
	}
}

func GetType(t TYPE) UserType {
	if t >= USER {
		return DefinedTypes[t-USER]
	} else {
		return UserType{}
		//panic(t.String()+" is a fundamental type, it cannot be indexed!")
	}
}

func AssignToUserType(s *scanner.Scanner, output io.Writer, name, element string) {
	structure := GetType(GetVariable(name))
	if index, ok := structure.Table[element]; !ok {
		RaiseError(s, name+" does not have an element named "+element)
	} else {

		//We are inferring the type now, we are being smarticles.
		if s.TokenText() == "is" {
			s.Scan()
			unique++
			value := fmt.Sprint("i+elem+",unique)
			fmt.Fprintf(output, "ARRAY %v\n", fmt.Sprint("i+elem+",unique))
			typ := ParseConstructor(s, output)
			
			if structure.Elements[index] == USER {
				structure.Elements[index] = typ
			}
			
			if typ != structure.Elements[index] {
				RaiseError(s, "Type mismatch! "+name+"."+element+" is a "+structure.Elements[index].String()+", not a "+typ.String())
			}
			
			unique++
			fmt.Fprintf(output, "SHARE %v\n PUSH 0\nHEAP\nPULL %v\n", value, fmt.Sprint("i+elem+",unique))
			fmt.Fprintf(output, "PLACE %v\nPUSH %v\nSET %v\n", name, index, fmt.Sprint("i+elem+",unique))
			return
		}
		
		if s.TokenText() == "." {
			s.Scan()
			fmt.Printf("%v Cannot assign to %s.%s.%s\n",s.Pos(),  name, element, s.TokenText())
			os.Exit(1)
		}
		
		if s.TokenText() != "=" {
			fmt.Println(s.Pos(), "Expecting = found ", s.TokenText())
			os.Exit(1)
		}
		s.Scan()
		
		value := expression(s, output)
		
		if ExpressionType != structure.Elements[index] {
			RaiseError(s, "Type mismatch! "+name+"."+element+" is a "+structure.Elements[index].String()+", not a "+ExpressionType.String())
		}
		
		switch structure.Elements[index] {
			case NUMBER, LETTER:
				fmt.Fprintf(output, "PLACE %v\nPUSH %v\nSET %v\n", name, index, value)
			case STRING, ARRAY:
				unique++
				fmt.Fprintf(output, "SHARE %v\n PUSH 0\nHEAP\nPULL %v\n", value, fmt.Sprint("i+elem+",unique))
				fmt.Fprintf(output, "PLACE %v\nPUSH %v\nSET %v\n", name, index, fmt.Sprint("i+elem+",unique))
			default:
				RaiseError(s, name+" cannot set "+element+", type is unsettable!!!")
		}
	}
}							
func IndexUserType(s *scanner.Scanner, output io.Writer, name, element string) string {
	var t UserType
	if GetVariable(name) > 0 {
		t = GetType(GetVariable(name))
	} else {
		t = GetType(ExpressionType)
	}
	
	//Deal with indexing Something types.
	if GetVariable(name) == SOMETHING {
		switch element {
			case "number":
				ExpressionType = NUMBER
				fmt.Fprintf(output, "PLACE %s\n", name)
				fmt.Fprintf(output, "PUSH 0\n")
				fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				return "i+user+"+fmt.Sprint(unique)
		}
	}
	
	if index, ok := t.Table[element]; ok {
	
		unique++
		ExpressionType = t.Elements[index]
	
		switch t.Elements[index] {
			case NUMBER, ITYPE, LETTER:
				fmt.Fprintf(output, "PLACE %s\n", name)
				fmt.Fprintf(output, "PUSH %v\n", index)
				fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				return "i+user+"+fmt.Sprint(unique)
			
			case STRING, USER, ARRAY:
				fmt.Fprintf(output, "PLACE %s\n", name)
				fmt.Fprintf(output, "PUSH %v\n", index)
				fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				
				fmt.Fprintf(output, "PUSH %s%v\n", "i+user+", unique)
				fmt.Fprintf(output, "HEAP\n")
				unique++
				fmt.Fprintf(output, "GRAB %s%v\n", "i+elem+", unique)
				
				return "i+elem+"+fmt.Sprint(unique)
				
				
			default:
				if t.Elements[index] >= USER {
					fmt.Fprintf(output, "PLACE %s\n", name)
					fmt.Fprintf(output, "PUSH %v\n", index)
					fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				
					fmt.Fprintf(output, "PUSH %s%v\n", "i+user+", unique)
					fmt.Fprintf(output, "HEAP\n")
					unique++
					fmt.Fprintf(output, "GRAB %s%v\n", "i+elem+", unique)
					
					Protected = true
				
					return "i+elem+"+fmt.Sprint(unique)
				}
				fmt.Println(s.Pos(), name+" cannot index "+element+", type is unindexable!!!")
				os.Exit(1)
		}
	} else {
		fmt.Println(s.Pos(), name+" does not have an element named "+element)
		os.Exit(1)
	}
	return ""
}

func ParseSymbolicType(s *scanner.Scanner) TYPE {
	var result TYPE
	var symbol = s.TokenText()
	switch symbol {
		case "{":
			result = USER
			s.Scan()
			Expecting(s, "}")
		case "[":
			result = ARRAY
			s.Scan()
			Expecting(s, "]")
		case "\"\"":
			result = STRING
		case "'":
			result = LETTER
			s.Scan()
			Expecting(s, "'")
		case "|":
			result = FILE
			s.Scan()
			Expecting(s, "|")
		case "(":
			result = FUNCTION
			s.Scan()
			Expecting(s, ")")
		case "<":
			result = ITYPE
			s.Scan()
			Expecting(s, ">")
		case ".":
			result = MULTIPLE
			s.Scan()
			Expecting(s, ".")
		default:
			result = NUMBER
			return result
	}
	s.Scan()
	return result
}

//Returns the type.
//Requires an "ARRAY name" before this is called.
func ParseConstructor(s *scanner.Scanner, output io.Writer) TYPE {
	stringtype := s.TokenText()
					
	if _, ok := StringToType[stringtype]; !ok {
		RaiseError(s, stringtype+" is an unrecognised type!")
	}
	
	s.Scan()
	//This is effectively a constructor.
	if s.TokenText() == "(" {
		var i int
		for {
			s.Scan()
			var add = expression(s, output)
			fmt.Fprintf(output, "PUT %s\n", add)
			if s.TokenText() == ")" {
				break
			} else if s.TokenText() != "," {
				RaiseError(s, "Expecting , found "+s.TokenText())
			}
			i++
		}
		for j := range DefinedTypes[StringToType[stringtype]-USER].Elements {
			if j > i {
				fmt.Fprintf(output, "PUT 0\n")
			}
		}
	} else {
		for range DefinedTypes[StringToType[stringtype]-USER].Elements {
			fmt.Fprintf(output, "PUT 0\n")
		}
	}
	return StringToType[stringtype]
}

//We have a new type.
func ParseTypeDef(s *scanner.Scanner, output io.Writer) {
	s.Scan()
		t := NewType(s.TokenText())
	s.Scan()
	if s.TokenText() == "{" {
		s.Scan()
		s.Scan()
	}
	//What are the element?
	for {
		if s.TokenText() == "}" {
			break
		}
		t.Elements = append(t.Elements, ParseSymbolicType(s))

		if s.TokenText() != "," {
			t.Table[s.TokenText()] = len(t.Elements)-1
			s.Scan()
		}
		if s.TokenText() == "}" {
			break
		}
		if s.TokenText() != "," && s.TokenText() != "\n" {
			fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
			os.Exit(1)
		}
		s.Scan()
		for s.TokenText() == "\n"  {
			s.Scan()
		}
	}
	DefinedTypes = append(DefinedTypes, t)
	StringToType[t.Name] = USER+TYPE(len(DefinedTypes)-1)

	LastDefinedType = USER+TYPE(len(DefinedTypes)-1)
	
	/*fmt.Println("NEW TYPE! ", t.Name)
	for _, v := range t.Elements {
		fmt.Println(v)
	}*/
}
