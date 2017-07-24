package block

import "github.com/qlova/ilang/src"

func init() {
	ilang.RegisterToken([]string{"{"}, ScanBlock)
	ilang.RegisterListener(Block, BlockEnd)
}

var Block = ilang.NewFlag()

func ScanBlock(ic *ilang.Compiler) {
	ic.Assembly("IF 1")
	ic.GainScope()
	ic.SetFlag(Block)
}

func BlockEnd(ic *ilang.Compiler) {
	ic.Assembly("END")
}
