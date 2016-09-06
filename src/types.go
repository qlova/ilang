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
	
	ITYPE
	
	FUNCTION
	STRING
	NUMBER
	FILE
	
	USER
	SOMETHING
)


func (t TYPE) String() string {
	if t >= USER {
		return DefinedTypes[t-USER].Name
	}
	return map[TYPE]string{FUNCTION:"function", STRING:"string",NUMBER:"number", ITYPE: "type", FILE:"file", UNDEFINED:"undefined"}[t]
}

func (t TYPE) Push() string {
	if t >= USER {
		return "SHARE"
	}
	return map[TYPE]string{FUNCTION:"RELAY", STRING:"SHARE",NUMBER:"PUSH",FILE:"RELAY",UNDEFINED:""}[t]
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
		Name: "something",
		Elements: []TYPE{ UNDEFINED, ITYPE, UNDEFINED},
		Table: map[string]int{"type":1, "data":0},
	},
}
var StringToType map[string]TYPE = map[string]TYPE{"something":SOMETHING}

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
		panic(t.String()+" is a fundamental type, it cannot be indexed!")
	}
}

func IndexUserType(s *scanner.Scanner, output io.Writer, name, element string) string {
	t := GetType(GetVariable(name))
	
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
			case NUMBER, ITYPE:
				fmt.Fprintf(output, "PLACE %s\n", name)
				fmt.Fprintf(output, "PUSH %v\n", index)
				fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				return "i+user+"+fmt.Sprint(unique)
			
			case STRING, USER:
				fmt.Fprintf(output, "PLACE %s\n", name)
				fmt.Fprintf(output, "PUSH %v\n", index)
				fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				
				fmt.Fprintf(output, "PUSH %s%v\n", "i+user+", unique)
				fmt.Fprintf(output, "HEAP\n")
				unique++
				fmt.Fprintf(output, "GRAB %s%v\n", "i+elem+", unique)
				
				return "i+elem+"+fmt.Sprint(unique)
				
			default:
				fmt.Println(s.Pos(), name+" cannot index "+element+", type is unindexable!!!")
				os.Exit(1)
		}
	} else {
		fmt.Println(s.Pos(), name+" does not have an element named "+element)
		os.Exit(1)
	}
	return ""
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
		//What type is the element?
		switch s.TokenText() {
			case "[":
				s.Scan()
				s.Scan()
				t.Elements = append(t.Elements, STRING)
			
			case "|":
				s.Scan()
				s.Scan()
				t.Elements = append(t.Elements, FILE)
			
			case "(":
				s.Scan()
				s.Scan()
				t.Elements = append(t.Elements, FUNCTION)
			
			default:
				t.Elements = append(t.Elements, NUMBER)
		}
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
