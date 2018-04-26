
concept numeric(n) {
	number(n)
	if error
		print(n, "  is not numeric!")
		return
	end
	print(n, "  is numeric :)")
}

software {
	numeric("1200")
	numeric("3.14")
	numeric("3/4")
	numeric("abcdefg")
}
