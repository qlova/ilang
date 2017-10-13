package f

import "github.com/qlova/ilang/src"

func init() {
	ilang.RegisterListener(TypeLoop, EndTypeLoop)
}

var TypeLoop = ilang.NewFlag()

func EndTypeLoop(ic *ilang.Compiler) {
	if ic.LastToken != "end" {
		ic.RaiseError("For Loops must end with a 'end' token, found '", ic.LastToken, "'")
	}
	
	var index = ic.GetVariable("flag_index")
	var t = ic.GetVariable("flag_type")

	ic.LoseScope()
	ic.Assembly("END")

	if index.Int < len(t.Detail.Elements) {
		NextTypeLoop(ic, t, index.Int, index.Push, index.Pop, *index.Plugin)
	}
	
}

func NextTypeLoop(ic *ilang.Compiler, t ilang.Type, id int, index, value string, plugin ilang.Plugin) {
	ic.Assembly("IF 1")
	ic.GainScope()
	
	ic.SetVariable("flag_type", t)
	ic.SetVariable("flag_type.", ilang.Protected)
	ic.SetFlag(ilang.Type{Name: "flag_index", Int: id+1, Plugin: &plugin, Push: index, Pop: value})
	
	ic.GainScope()
	ic.SetFlag(TypeLoop)
	
	for name, i := range t.Detail.Table {
		if i == id {
			ic.Aliases[index] = name
			
			if value != "" {
				ic.Assembly("SHARE ", ic.ParseString("\""+t.Detail.Elements[i].Name+"\""))
				ic.Assembly("GRAB ", value)
				ic.SetVariable(value, ilang.Text)
			}
			
			ic.SetVariable("element_type", t.Detail.Elements[i])
			ic.SetVariable("element_type.", ilang.Protected)	
			
			ic.Assembly("SHARE ", ic.ParseString("\""+name+"\""))
			ic.Assembly("GRAB ", name)
			ic.SetVariable(name, ilang.Text)
			break
		}
	}

	ic.Insertion = append(ic.Insertion, plugin)
}

func ScanTypeLoop(ic *ilang.Compiler, t ilang.Type, index string, value string) {

	//Scan the code in the for loop, this code will be generated multiple times.
	var plugin ilang.Plugin
	plugin.Line = ic.Scanner.Pos().Line

	var braces = 0
	for {
		var token = ic.Scan(0)
		if token == "end"  {
	 		if braces == 0 {
	 			plugin.Tokens = append(plugin.Tokens, token)
				break
			} else {
				braces--
			}
		}
		if token == "if" || token == "for" {
			braces++
		}
		plugin.Tokens = append(plugin.Tokens, token)
	}

	plugin.File = ic.Scanner.Pos().Filename
	
	NextTypeLoop(ic, t, 0, index, "type", plugin)
}
