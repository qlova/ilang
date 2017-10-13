package issues

import "github.com/qlova/ilang/src"

func init() {
	ilang.RegisterToken([]string{"issues"}, ScanIssues)
	ilang.RegisterToken([]string{"issue"}, ScanIssue)
	
	ilang.RegisterListener(Issues, IssuesEnd)
	ilang.RegisterListener(Issue, IssueEnd)
	
	ilang.RegisterToken([]string{"!"}, func(ic *ilang.Compiler) {
		ic.Assembly("ADD ERROR 0 0")
	})
}

var Issues = ilang.NewFlag()
var Issue = ilang.NewFlag()

func IssueEnd(ic *ilang.Compiler) {	
	ic.UnsetFlag(Issues)
	for i:=0; i < (ic.GetVariable("flag_nesting").Int); i++ {
		ic.Assembly("END")
	}
	ic.Assembly("END")
	ic.LoseScope()
	ic.Assembly("END")	
}

func IssuesEnd(ic *ilang.Compiler) {	
	ic.Assembly("END")	
}

func ScanIssues(ic *ilang.Compiler) {
	ic.Scan('{')
	ic.Assembly("IF ERROR")
	ic.GainScope()
	ic.Assembly("VAR issue")
	ic.Assembly("ADD issue ERROR 0")
	ic.Assembly("ADD ERROR 0 0")
	ic.SetFlag(Issues)
	
	var token string
	for {
		token = ic.Scan(0)
		if token != "\n" {
			if token != "issue" {
				ic.NextToken = token
			}
			break
		}
	}
	if token == "issue" {
	
		var expression = ic.ScanExpression()
		var condition = ic.Tmp("issue")
		ic.Assembly("VAR ", condition)
		ic.Assembly("SEQ %v %v %v", condition, expression, "issue")
		ic.Assembly("IF ",condition)
		ic.SetFlag(ilang.Type{Name: "flag_nesting", Int: 0})
		ic.GainScope()
		ic.SetFlag(Issue)
	}
}

func ScanIssue(ic *ilang.Compiler) {
	if !ic.GetFlag(Issues) {
		ic.RaiseError("'issue' must be within a 'issues' block!")
	}

	var expression = ic.ScanExpression()
	var condition = ic.Tmp("issue")
	
	ic.UnsetFlag(Issue)
	
	ic.UpdateVariable("flag_nesting", ilang.Type{Int: ic.GetVariable("flag_nesting").Int+1 })
	
	ic.LoseScope()
	
	ic.Assembly("ELSE")
	
	ic.Assembly("VAR ", condition)
	ic.Assembly("SEQ %v %v %v", condition, expression, "issue")
	ic.Assembly("IF ",condition)
	ic.GainScope()
	ic.SetFlag(Issue)
}
