
function numeric(n) {
!	num(n)
	issues {
		output(n&"  is not numeric!\n")
		return
	}
	output(n&"  is numeric :)\n")
}

software {
	numeric("1200")
	numeric("3.14")
	numeric("3/4")
	numeric("abcdefg")
}
