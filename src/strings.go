package ilang

import "strconv"

func (ic *Compiler) ParseString(s string) string {
	var err error
	s, err = strconv.Unquote(s)
	if err != nil {
		ic.RaiseError(err)
	}
	
	var r = ic.Tmp("string")
	ic.Assembly("ARRAY ", r)
	
	var b = []byte(s)
	
	for _, v := range b {
		ic.Assembly("PUT ", v)
	}
	
	return r
}
