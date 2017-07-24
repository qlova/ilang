package test

import ".."
import "testing"
import "strings"
import "io/ioutil"

func TestScanUserStatement(t *testing.T) {
	var ic = ilang.NewCompiler(strings.NewReader(`
		type User {
			element
		}
		
		software {
			var u = {}
			u = User()
			print(u.element)
		}
	`))
	ic.Output 	= ioutil.Discard
	ic.Lib 		= ioutil.Discard
	ic.Compile()
}
