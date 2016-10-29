package main

//Plugins are like code injections or aspect-based programming.
//VERY POWERFUL.
type Plugin struct {
	File string
	Line int
	Tokens []string
}

func (ic *Compiler) InsertPlugins(name string) {
	if plugins, ok := ic.Plugins[name]; ok {
		for _, plugin := range plugins {
			ic.Insertion = append(ic.Insertion, plugin)
		}
	}
}

func (ic *Compiler) ScanPlugin() {
	var sort = ""
	var name = ic.Scan(Name)
	var token = ic.Scan(0)
	if token == "(" {
		sort = ic.Scan(Name)
		ic.Scan(')')
		ic.Scan('{')
	} else if token != "{" {
		ic.Expecting("{")
	}

	var plugin Plugin
	plugin.Line = ic.Scanner.Pos().Line

	var braces = 0
	for {
		var token = ic.Scan(0)
		if token == "}"  {
	 		if braces == 0 {
				break
			} else {
				braces--
			}
		}
		if token == "{" {
			braces++
		}
		plugin.Tokens = append(plugin.Tokens, token)
	}

	plugin.File = ic.Scanner.Pos().Filename
	
	if sort != "" {
		ic.Plugins[name+"_m_"+sort] = append(ic.Plugins[name+"_m_"+sort], plugin)
	} else {
		ic.Plugins[name] = append(ic.Plugins[name], plugin)
	}
}
