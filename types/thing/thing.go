package thing

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

import "github.com/qlova/ilang/types/text"
import "github.com/qlova/ilang/types/number"

type Base struct{}

func (Base) Push(c *compiler.Compiler, data string) {
	c.Push(data+"_pointer")
	c.PushList(data)
}

func (Base) Pull(c *compiler.Compiler, data string) {
	c.Pull(data+"_pointer")
	c.PullList(data)
}


func (Base) Drop(c *compiler.Compiler) {
	c.Drop()
	c.DropList()
}


func (Base) Free(c *compiler.Compiler) {
	c.Free()
	c.FreeList()
}

//Note this only works on standalone things!
func (Base) Attach(c *compiler.Compiler) {
	compiler.LIST.Attach(c)
}

func (Base) Detach(c *compiler.Compiler) {
	compiler.LIST.Detach(c)
}

var Type = compiler.Type {
	Name: compiler.Translatable{
		compiler.English: "thing",
	},
	
	Base: Base{},
}

type Data struct {
	Size int
	
	Elements []compiler.Type
	Map map[string]int
	Offsets map[string]int
	
	Concepts []*compiler.Function
}

func (d Data) Name(compiler.Language) string {
	return ""
}

func (d Data) Equals(b compiler.Data) bool {
	return false
}

var Shunts = compiler.Shunt {
	symbols.Index: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			c.RaiseError(errors.Single(Type, symbols.Plus,t))
		}
		
		
		
		return Type
	},
}

func init() {
	Type.Shunts = Shunts
}

func Collect(c *compiler.Compiler,  t compiler.Type) {
	if NotThing(t) {
		panic("Illegal call to ScanEmbeddedLiteral")
	}
	
	for element, offset := range t.Data.(Data).Offsets {
		var subtype = t.Data.(Data).Elements[t.Data.(Data).Map[element]]
		
		if subtype.Base == compiler.LIST && NotThing(subtype) {
			c.Copy()
			c.Int(int64(offset))
			c.Add()
			c.Get()
			c.Flip()
			c.HeapList()
		} else if !NotThing(subtype) {
			
			c.Int(int64(offset))
			Collect(c, subtype)
			
		}
	}
	
	c.Drop()
}
		

func ScanEmbeddedLiteral(c *compiler.Compiler, t compiler.Type) {
	
	if NotThing(t) {
		panic("Illegal call to ScanEmbeddedLiteral")
	}
	
	//We need to copy the Map of Data so that we can keep track of which 
	//elements have been substantiated and which one's havn't.
	
	var TrackingMap = make(map[string]int)
	for pair, value := range t.Data.(Data).Map {
		TrackingMap[pair] = value
	}	
	
	for {
		var element = c.Scan()
		
		if element == "\n" {
			continue
		}
		
		if element == symbols.CodeBlockEnd {
			break
		}
		
		c.Expecting("=")
		
		if index, ok := TrackingMap[element]; ok {
			
			var ExpectedType = t.Data.(Data).Elements[index]
			var Offset = t.Data.(Data).Offsets[element]
			
			var value = c.ScanExpression()
			
			if !value.Equals(ExpectedType) {
				c.RaiseError(errors.AssignmentMismatch(value, ExpectedType))
			}

			c.Int(int64(Offset))
			c.Add()
			
			if value.Equals(text.Type) {

				c.Int(0)
				c.HeapList()
				c.Set()
			
			} else if value.Equals(number.Type) {
				
				c.Set()
				
			} else {
				c.Unimplemented()
			}
			
		} else {
			c.RaiseError(errors.NoSuchElement(element, t))
		}
	}
}

func GetElementOffsetFromIndex(t compiler.Type, i int) int {
	if NotThing(t) {
		panic("Illegal call to GetElementOffsetFromIndex")
	}
	
	return t.Data.(Data).Offsets[GetElementNameFromIndex(t, i)]
}

func GetElementNameFromIndex(t compiler.Type, i int) string {
	if NotThing(t) {
		panic("Illegal call to GetElementNameFromIndex")
	}
	
	for name, index := range t.Data.(Data).Map {
		if index == i {
			return name
		}
	}
	
	panic("Element index out of range!")
}

func NotThing(t compiler.Type) bool {
	data := t.Data
	
	if data == nil {
		return true
	}
	
	if _, ok := data.(Data); ok {
		return false
	}
	
	return true
}

func Is(t compiler.Type) bool {
	data := t.Data
	
	if data == nil {
		return false
	}
	
	if _, ok := data.(Data); ok {
		return true
	}
	
	return false
}
