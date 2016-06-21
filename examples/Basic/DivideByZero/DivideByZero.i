function isdivbyzero(a, b) {
	var c = a/b
	if (c = 0) or ( a = 0  and  c >= 0)
		output(text(a)&"/"&text(b)&" is a divivision by zero.\n")
		return
	end
	output(text(a)&"/"&text(b)&" is not divivision by zero.\n")
}

software {
	isdivbyzero(5, 0)
	isdivbyzero(5, 2)
	isdivbyzero(0, 0)
}
