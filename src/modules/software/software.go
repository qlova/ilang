package software 

import "github.com/qlova/ilang/src"

func init() {
	ilang.RegisterToken([]string{"software", "ソフトウェア", "программного", "软件"}, ScanSoftware)
	ilang.RegisterListener(Software, SoftwareEnd)
	
	ilang.RegisterToken([]string{"exit"}, func(ic *ilang.Compiler) {
	
		ic.CollectGarbage()
		ic.Assembly("IF 1\nEXIT\nEND")
		
	})
}

var Software = ilang.NewFlag()

func SoftwareEnd(ic *ilang.Compiler) {
	ic.Assembly("EXIT")
}

func ScanSoftware(ic *ilang.Compiler) {

	//Russian check.
	if ic.LastToken == "программного" {
		if ic.Scan(0) != "обеспечения" {
			ic.RaiseError("ожидая обеспечения")
		}
	}
	
	ic.Header = false
	ic.Scan('{')
	ic.Assembly("SOFTWARE")
	ic.GainScope()
	ic.SetFlag(Software)
	ic.SoftwareBlockExists = true
	
	if ic.GUIExists && ic.GUIMainExists {
		ic.Assembly("SHARE gui_main")
		ic.Assembly("RUN gui")
		ic.LoadFunction("gui")
		ic.LoadFunction("output_m_pipe")
	}
}
