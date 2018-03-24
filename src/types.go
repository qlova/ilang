package ilang

import "strings"
import "github.com/gedex/inflector"
import "fmt"

var TypeIota int

type Type struct {
	Name, Push, Pop string
	
	Int int
	User bool
	List bool
	
	Super string
	
	Decimal bool
	
	Plugin *Plugin
	Detail *UserType //This contains usertype information.
	SubType *Type //Subtype for recursive types such as lists.
	Class *Type
	
	Functions *[]Function
}

func (ic *Compiler) Cast(name string, t Type, b Type) string {
	
	for _, cast := range Casts {
		var r = cast(ic, name, t, b)
		if r != "" {
			return r
		}
	}
	
	var asm = ""
	asm += t.Push+" "+name+"\n"
	asm += ic.RunFunction(b.GetComplexName()+"_m_"+t.GetComplexName())
	
	return asm
}

func (ic *Compiler) CanCast(t Type, b Type) bool {
	for _, cast := range Casts {
		var r = cast(ic, "", t, b)
		if r != "" {
			return true
		}
	}
	
	_, ok := ic.DefinedFunctions[b.GetComplexName()+"_m_"+t.GetComplexName()]
	if ok {
		return true
	}
	
	//Maybe the function needs to be generated?
	for _, builder := range FunctionBuilders {
		var r = builder(b.GetComplexName()+"_m_"+t.GetComplexName())
		if r != nil {
			return true
		}
	}
	
	return false
}

func (t Type) Equals(b Type) bool {

	if t.Name == "list" {
		if b.Name == "list" {
			if t.SubType == nil || b.SubType == nil {
				return true
			}
		}
	}
	
	if t.Name == "list" && *t.SubType == Number && b.Name == "array" {
		return true
	}
	if b.Name == "list" && *b.SubType == Number && t.Name == "array" {
		return true
	}
	
	if t.Name == "list" && t.SubType == nil && b.Name == "array" {
		return true
	}
	if b.Name == "list" && b.SubType == nil && t.Name == "array" {
		return true
	}
	
	if t.Name == "function" {
		if b.Name == "function" {
			if t.Detail != nil && b.Detail != nil {
				if len(t.Detail.Elements) != len(b.Detail.Elements) {
					return false
				}
				
				for i := range t.Detail.Elements {
					if !t.Detail.Elements[i].Equals(b.Detail.Elements[i]) {
						return false
					}
				}
				
				return true
			}
		}
	}

	if t.Name != b.Name {
		if t.User {
			return false
		}
		
		if (t.Class == nil || t.Class.Name != b.Name) && (b.Class == nil || b.Class.Name != t.Name) {
			return false
		}
	}
	if t.SubType != nil && b.SubType != nil {
		return t.SubType.Equals(*b.SubType)
	}
	if t.SubType != nil || b.SubType != nil {
		return false
	}
	return true
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

//TODO this will break something types.
func (t Type) Empty() bool { 
	simple := t.Detail != nil && len(t.Detail.Elements) == 0
	if simple {
		return true
	}
	
	if t.Detail != nil && len(t.Detail.Elements) > 0 {
		for _, subtype := range t.Detail.Elements {
			if !subtype.Empty() {
				return false
			}
		}
		return true
	}
	
	return false
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

func GetType(name string) Type {
	return string2type[name]
}

var Undefined = NewType("undefined")
var Number = NewType("number", "PUSH", "PULL")

var Text = NewType("text", "SHARE", "GRAB")
var Array = NewType("array", "SHARE", "GRAB")
var Matrix = NewType("matrix", "SHARE", "GRAB")

var Variadic = NewFlag()

func (ic *Compiler) ScanSymbolicType() Type {
	var result Type = Undefined
	var symbol = ic.Scan(0)
	
	if symbol == "." && ic.Peek() == "." {
		ic.Scan('.')
		symbol = ".."
	}

	if ( symbol != ".." && len(symbol) > 0 && symbol[0] == '.') {
		symbol = "."
	}
	
	if f, ok := Symbols[symbol]; ok {
		return f(ic)
	}
	
	switch symbol {
		case `""`:
			result = Text
			
		//TODO move into a type module.
		case `[`:
			ic.Scan(']')
			result = Array
			
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


func (ic *Compiler) GetType(name string) Type {
	if t, ok := ic.DefinedTypes[name]; ok {
			return t
	}
	return string2type[name]
}

//This function generates the assembly for a type call.
//This does not set the ExpressionType.
//For example:
//	var x = text()
//	var y = UserType()
//	var z = number()
func (ic *Compiler) CallType(name string) string {
	if name == "text" || name == "list" {
		var array = ic.Tmp("user")
		ic.Assembly("ARRAY ", array)
		return array
	}
	
	var t = ic.GetType(name)
	if t == Undefined {
		ic.RaiseError(name, " is not a type!")
	}
	
	//Complex types.
	//eg. type TextList ..""
	if t.Class != nil {
			return ic.CallType(t.Class.Name)
	}
	
	//type modules.
	/*if f, ok := ic.DefinedFunctions[name]; ok {
		ic.Assembly(ic.RunFunction(name))
	
		if len(f.Returns) > 0 {
			id := ic.Tmp("result")
			
			var ReturnType = f.Returns[0]
			
			ic.Assembly("%v %v", ReturnType.Pop, id)
			
			return id
		}	
	}*/
	
	
	
	if t.User {
		if t.Empty() { //Empty types can be ignored.
			return ""
		} 
		
		//Create an array large enough to fit the type.
		var array = ic.Tmp("user")
		ic.Assembly("PUSH ", len(ic.DefinedTypes[name].Detail.Elements))
		ic.Assembly("MAKE")
		ic.Assembly("GRAB ", array)
		return array
	}

	ic.RaiseError("Complex type ", t.GetComplexName() ," needs to be implemented!")
	return ""
}

func (list Type) GetComplexName() string {
	if list.SubType == nil {
		if list.Name == "function" && list.Detail != nil && len(list.Detail.Elements) > 0 {
			var suffix string = "("
			for i := range list.Detail.Elements {
				suffix += list.Detail.Elements[i].GetComplexName()
				if i < len(list.Detail.Elements) -1 {
					suffix += ","
				}
			}
			suffix += ")"
			return "function"+suffix
		}
	
		return list.Name
	}
	
	if *list.SubType == list {
		panic("SELF REFERENCING LIST!")
	}

	var strings []string
	var t = list
	for {
		strings = append(strings, t.Name)
		if t.SubType == nil || (t == Type{}) {
			break
		}
		t = *t.SubType
	}
	var serialised string
	for i:=len(strings)-1; i >= 0; i-- {
		serialised += strings[i]
	}
	return serialised
}	

func (ic *Compiler) CreatePointer(t Type, pointer string, name string) string {
	var asm = ""
	if ic.ExpressionType.Push == "SHARE" {
		asm += fmt.Sprintln("PUSH 0")
		asm += fmt.Sprintln("SHARE", name)
		asm += fmt.Sprintln("HEAP")
		asm += fmt.Sprintln("PULL", pointer)
		asm += fmt.Sprintln("ADD",pointer,"0", pointer)
		return asm
	}
	if ic.ExpressionType.Push == "RELAY" {
		asm += fmt.Sprintln("PUSH 0")
		asm += fmt.Sprintln("RELAY", name)
		asm += fmt.Sprintln("HEAPIT")
		asm += fmt.Sprintln("PULL", pointer)
		return asm
	}
	asm += fmt.Sprintln("PUSH",name)
	asm += fmt.Sprintln("PULL",pointer)
	return asm
	
}	

//Get a numeric value which represents the type.
func (ic *Compiler) GetPointerTo(name string) string {
	var pointer = ic.Tmp("pointer")
	ic.Assembly(ic.CreatePointer(ic.ExpressionType, pointer, name))
	return pointer
}

func (ic *Compiler) Dereference(pointer string) string {
	if ic.ExpressionType.Push == "SHARE" {
		var value = ic.Tmp("deref")
		ic.Assembly("IF ", pointer)
		ic.Assembly("PUSH ", pointer)
		ic.Assembly("HEAP")
		ic.Assembly("ELSE ")
		
		var tmp = ic.CallType(ic.ExpressionType.Name)
		ic.Assembly("SHARE ", tmp)
		
		ic.Assembly("END")
		ic.Assembly("GRAB ", value)
		return value
	}
	
	if ic.ExpressionType.Push == "RELAY" {
		var value = ic.Tmp("deref")
		ic.Assembly("IF ", pointer)
		ic.Assembly("PUSH ", pointer)
		ic.Assembly("HEAPIT")
		ic.Assembly("ELSE ")

		var tmp = ic.Tmp("empty")
		ic.Assembly("ARRAY ", tmp)
		ic.Assembly("OPEN ", tmp)
		
		ic.Assembly("END")
		ic.Assembly("TAKE ", value)
		return value
	}
	return pointer
}

//This scans a new type definition and creates the type.
//eg. type Point { x, y }
func (ic *Compiler) ScanType() {
	var name = ic.Scan(Name)
	
	//This is for the grate engine.
	//Are we declaring a game?
	if name == "Graphics" {
		ic.Game = true
	}
	
	t := NewUserType(name)
	ic.DefinedTypes[name] = t
	
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
			
			//We assume that this is a complex type.
			//Introducing... ITDL or 'i' type description language.
			ic.NextToken = ic.LastToken
			var t = ic.ScanSymbolicType()
			var t2 = t
			t.Class = &t2
			t.Name = name
			ic.DefinedTypes[name] = t
			
			return
			
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

	//TODO depreciate this variable, rename it to CurrentMethodType or something.
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
