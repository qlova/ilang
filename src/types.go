package ilang

import "strings"
import "github.com/gedex/inflector"

var TypeIota int

type Type struct {
	Name, Push, Pop string
	Int int
	User bool
	List bool
	
	Super string
	
	Decimal bool
	
	Detail *UserType
	Interface *Interface
}

func (t Type) DefaultValue() string {
	switch t.Push {
		case "PUSH":
			return "0"
		case "SHARE":
			return "backup"
		case "RELAY":
			return "open"
	}
	return ""
}

func (t Type) IsUser() Type {
	if t.User {
		return t
	} else {
		return Undefined
	}
}


func (t Type) IsMatrix() Type {
	if t.Name == "matrix" {
		return t
	} else {
		return Undefined
	}
}

func (t Type) IsArray() Type {
	if t.Name == "array" {
		return t
	} else {
		return Undefined
	}
}

func (t Type) IsList() Type {
	if t.List {
		return t
	} else {
		return Undefined
	}
}

func (t Type) Empty() bool { 
	return t.Detail != nil && len(t.Detail.Elements) == 0 && t != Something
}

type UserType struct {	
	Elements []Type
	Table map[string]int
	SubElements map[int]Type
}

func NewUserType(name string) Type {
	t := NewType(name, "SHARE", "GRAB")
	t.User = true
	t.Detail = new(UserType)
	t.Detail.Table = make(map[string]int)
	return t
}

var string2type = map[string]Type{}

func NewType(name string, options ...string) Type {
	var t Type
	t.Name = name
	
	if len(options) == 2 {
		t.Pop = options[1]
		t.Push = options[0]
	}
	
	t.Int = TypeIota
	TypeIota++
	
	string2type[name] = t
	
	return t
}

var Undefined = NewType("undefined")
var Number = NewType("number", "PUSH", "PULL")
var Decimal = NewType("decimal", "PUSH", "PULL")
var Letter = NewType("letter", "PUSH", "PULL")
var Text = NewType("text", "SHARE", "GRAB")
var Array = NewType("array", "SHARE", "GRAB")
var Matrix = NewType("matrix", "SHARE", "GRAB")

var TextArray = Text

var Itype = NewType("type", "PUSH", "PULL")
var User = NewType("usertype", "SHARE", "GRAB")
var List = NewType("list", "SHARE", "GRAB")
var Pipe = NewType("pipe", "RELAY", "TAKE")
var Func = NewType("function", "RELAY", "TAKE")
var Something = NewUserType("Something")

var Variadic = NewFlag()

func init() {
	TextArray.List = true
}

func (ic *Compiler) ScanSymbolicType() Type {
	var result Type = Undefined
	var symbol = ic.Scan(0)
	switch symbol {
		case "{":
			result = User
			t := ic.Scan(0)
			if t == "." {	
				result = List
				ic.Scan('.')
				ic.Scan('}')
			} else if t != "}" {
				ic.RaiseError()
			}
		case "[":
			result = List
			ic.Scan(']')
			if tok := ic.Scan(0); tok == "[" {
				result = Matrix
				ic.Scan(']')
			} else {
				ic.NextToken = tok
			}
		case "$":
			result = ic.ScanSymbolicType()
			result.Decimal = true
		case `""`:
			result = Text
		case "' '":
			result = Letter
		case "|":
			result = Pipe
			ic.Scan('|')
		case "(":
			result = Func
			ic.Scan(')')
		case "<":
			result = Itype
			ic.Scan('>')
		case ".":
			if tok := ic.Scan(0); tok == "." {
				result = Variadic
			} else {
				ic.NextToken = tok
				result = Decimal
			}
		default:
			result = Number
			ic.NextToken = symbol
			return result
	}
	return result
}

//Check if the given type exists or not.
func (ic *Compiler) TypeExists(name string) bool {
	_, ok := ic.DefinedTypes[name]
	return ok
}

func (ic *Compiler) CallType(name string) string {
	if ic.DefinedTypes[name].Empty() {
		return ""
	} else {
		var array = ic.Tmp("user")
		ic.Assembly("ARRAY ", array)
		for range ic.DefinedTypes[name].Detail.Elements {
			ic.Assembly("PUT 0")
		}
		return array
	}
}

//This scans a new type definition and creates the type.
//eg. type Point { x, y }
func (ic *Compiler) ScanType() {
	var name = ic.Scan(Name)
	
	//This is for the grate engine.
	//Are we declaring a game?
	if name == "Game" {
		ic.Game = true
	}
	
	t := NewUserType(name)
	
	switch ic.Scan(0) {
		case "{":
		case "is": //Inheritance eg. type WeightedPoint is Point { weight }
			super := ic.Scan(Name)
			t = ic.DefinedTypes[super]
			t.Super = t.Name
			t.Name = name
			switch ic.Scan(0) {
				case "\n":
					ic.DefinedTypes[name] = t
					ic.LastDefinedType = t
					return
				case "{":
				default:
					ic.RaiseError()
			}
		default:
			ic.RaiseError()
	}
		
	ic.InsertPlugins(name)
	//What are the elements?
	for {
		var token = ic.Scan(0)
		if token == "}" {
			break
		}
		if token != "," && token != "\n" {
			ic.NextToken = token
		}
		
		MemberType := ic.ScanSymbolicType()
		
		ident := ic.Scan(Name)
		if ident == "}" {
			break
		}
		
		//Embedded structs which are inferred.
		//eg.
		/*
			type Member {}
			type Base {
				Member() //This will be accessed as 'member'.
			}
		*/
		if ic.Peek() == "(" {
			if MemberType != Number {
				ic.RaiseError("Unexpected (")
			}
			ic.Scan('(')
			ic.Scan(')')
			var ok bool
			MemberType, ok = ic.DefinedTypes[ident]
			if !ok {
				ident := inflector.Singularize(ident)
				MemberType, ok = ic.DefinedTypes[ident]
				MemberType.List = true
				MemberType.User = false
				if !ok {
					ic.RaiseError("No such type! ", ident)
				}
			}
			ident = strings.ToLower(ident)
		}
		
		
		if ident != "\n" { 
			t.Detail.Elements = append(t.Detail.Elements, MemberType)
			t.Detail.Table[ident] = len(t.Detail.Elements)-1
		}
		
	}
	
	ic.DefinedTypes[name] = t

	ic.LastDefinedType = t
}

//This scans a type literal.
// eg. 
/*
	type Object {value}
	
	software { var o = Object{22} }
*/
func (ic *Compiler) ScanTypeLiteral() string {
	return ic.ScanConstructor()
}

func (ic *Compiler) ScanConstructor() string {
	var name = ic.Scan(Name)
					
	if _, ok := ic.DefinedTypes[name]; !ok {
		ic.RaiseError(name+" is an unrecognised type!")
	}
	
	var token = ic.Scan(0)
	
	/*if ic.Peek() == ")" && token == "(" {
		ic.ExpressionType = InFunction
		ic.NextToken = "("
		return name
	}*/
	
	var array = ic.Tmp("constructor")
	
	ic.Assembly("ARRAY ", array)
	//This is effectively a constructor.
	if token == "{" {
		var i int
		for {
			
			var expr = ic.ScanExpression()
			ic.Assembly("PLACE ", array)
			if ic.ExpressionType.Push == "PUSH" {
				ic.Assembly("PUT %v", expr)
			} else {
				var tmp = ic.Tmp("heap")
				ic.Assembly(ic.ExpressionType.Push," ", expr)
				ic.Assembly("PUSH 0")
				if ic.ExpressionType.Push == "RELAY" {
					ic.Assembly("HEAPIT")
				} else {
					ic.Assembly("HEAP")
				}
				ic.Assembly("PULL ", tmp)
				ic.Assembly("PUT ", tmp)
			}
			if i >= len(ic.DefinedTypes[name].Detail.Elements) {
				ic.RaiseError("Too many arguments passed to constructor!")
			}
			if ic.ExpressionType != ic.DefinedTypes[name].Detail.Elements[i] {
				ic.RaiseError("Mismatched types! Argument (%v) of constructor should be '%v'", i+1, 
					ic.DefinedTypes[name].Detail.Elements[i])
			}
			token = ic.Scan(0)
			for token == "\n" {
				token = ic.Scan(0)
			}
			if token == "}" {
				break
			} else if token != "," {
				ic.Expecting(",")
			}
			i++
		}
		for j := range ic.DefinedTypes[name].Detail.Elements {
			if j > i {
				ic.Assembly("PUT 0")
			}
		}
	} else if token == "\n" || token == ")" {
		for range ic.DefinedTypes[name].Detail.Elements {
			ic.Assembly("PUT 0")
		}
		if token == ")" {
			ic.NextToken = ")"
		}	
	} else {
		ic.RaiseError()
	}
	ic.ExpressionType = ic.DefinedTypes[name]
	return array
}

func (ic *Compiler) IndexUserType(name, element string) string {
	var t UserType
	if ic.GetVariable(name) != Undefined {
		t = *ic.GetVariable(name).Detail
		ic.SetVariable(name+"_use", Used)
	} else {
		t = *ic.ExpressionType.Detail
	}
	
	//Deal with indexing Something types.
	/*if GetVariable(name) == SOMETHING {
		switch element {
			case "number":
				ExpressionType = NUMBER
				fmt.Fprintf(output, "PLACE %s\n", name)
				fmt.Fprintf(output, "PUSH 0\n")
				fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				return "i+user+"+fmt.Sprint(unique)
		}
	}*/
	
	if index, ok := t.Table[element]; !ok {
		ic.RaiseError(name+" does not have an element named "+element)
	} else {
	
		var tmp = ic.Tmp("index")
		ic.ExpressionType = t.Elements[index]
	
		switch t.Elements[index].Push {
			case "PUSH":
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUSH ", index)
				ic.Assembly("GET ", tmp)
				return tmp
			
			case "SHARE", "RELAY":
				ic.Assembly("PLACE ", name) //The array we are indexing, the place.
				ic.Assembly("PUSH ", index) //Push the index onto the stack.
				ic.Assembly("GET ", tmp)	//Get the value of the array at the index on the stack.
				ic.Assembly("IF ",tmp)		//If there is a valid address, (greater than zero)
				ic.GainScope()
				
				//Retrieve the array.
				ic.Assembly("PUSH ", tmp)
				if t.Elements[index].Push == "RELAY" {
					ic.Assembly("HEAPIT")
				} else {
					ic.Assembly("HEAP")
				}
				tmp = ic.Tmp("index")
				ic.Assembly(t.Elements[index].Pop, " ", tmp)
				ic.Assembly(t.Elements[index].Push, " ", tmp)
				ic.LoseScope()
				
				ic.Assembly("ELSE") //We will return a new array.
				ic.GainScope()
				ic.Assembly("ARRAY ", tmp)
				if t.Elements[index].User {
				for range t.Elements[index].Detail.Elements {
					ic.Assembly("PUT 0")
				}
				}
				ic.Assembly("SHARE ", tmp)
				if t.Elements[index].Push == "RELAY" {
					ic.Assembly("OPEN")
				}
				ic.LoseScope()
				ic.Assembly("END")
				ic.Assembly(t.Elements[index].Pop, " ", tmp)
				
				return tmp
				
			default:
				ic.RaiseError(name+" cannot index "+element+", type is unindexable!!!")
		}
	}
	return ""
}

func (ic *Compiler) SetUserType(name, element, value string) {
	var t UserType
	if ic.GetVariable(name) != Undefined {
		t = *ic.GetVariable(name).Detail
		ic.SetVariable(name+"_use", Used)
	} else {
		ic.RaiseError("Cannot set type without type identity!")
	}
	
	if index, ok := t.Table[element]; !ok {
		ic.RaiseError(name+" does not have an element named "+element)
	} else {
	
		if t.Elements[index] == User || (t.Elements[index].List && ic.ExpressionType.Push == "SHARE") || ic.ExpressionType.Name == "matrix" {
			t.Elements[index] = ic.ExpressionType
			
			if  ic.GetFlag(InMethod) {
				ic.Assembly("PLACE ", value)
				ic.Assembly("RENAME ", element)
				//ic.SetVariable(element, ic.ExpressionType)
			}
		}
	
		if ic.ExpressionType != t.Elements[index] {
			ic.RaiseError("Type mismatch, cannot assign '",ic.ExpressionType.Name,"', to a element of type '",t.Elements[index].Name,"'")		
		}

		switch t.Elements[index].Push {
			case "PUSH":
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUSH ", index)
				ic.Assembly("SET ", value)
			
			case "SHARE", "RELAY":
				
				//TODO garbage collect
				var tmp = ic.Tmp("index")
				ic.Assembly(t.Elements[index].Push, " ", value)
				ic.Assembly("PUSH 0")
				if t.Elements[index].Push == "RELAY" {
					ic.Assembly("HEAPIT")
				} else {
					ic.Assembly("HEAP")
				}
				ic.Assembly("PULL ", tmp)
				
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUSH ", index)
				ic.Assembly("SET ", tmp)
				
			default:
				ic.RaiseError(name+" cannot index "+element+", type is unindexable!!!")
		}
	}
}
