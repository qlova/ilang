function exists([]filename) {
	file = open(filename)
	issues {
		output(filename&" does not exist"&"\n")
		return
	}
	output(filename&" exists"&"\n")
	close(file)
}

software {
	exists("input.txt")
	exists("/input.txt")
	exists("docs")
	exists("/docs")
}
