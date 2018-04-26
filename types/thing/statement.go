package thing

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

func ScanStatement(c *compiler.Compiler, t compiler.Type) compiler.Type {
	if NotThing(t) {
		panic("Illegal call to thing.ScanStatement")
	}
	
	var offset = 0
	for {
		var element = c.Token()
		
		if this_offset, ok := t.Data.(Data).Offsets[element]; ok {
			var subtype =  t.Data.(Data).Elements[t.Data.(Data).Map[element]]
			
			offset += this_offset
			
			if !NotThing(subtype) {
				if c.Peek() == symbols.Index {
					t = subtype
					c.Scan()
					c.Scan()
					continue
				}
			}
			
			c.Int(int64(offset))
			
			return subtype
			
		} else {
			
			c.RaiseError(errors.NoSuchElement(element, t))
		}
	}
}
