
function create( []filename ) {
	var file = open(filename) 
!	output@file("")
	issues {
		output("Failed to create "&filename&"\n")
		return
	}
	output(filename&" created!\n")
	close(file)
}

software {
	create("output.txt")
	create("docs/")
	create("/output.txt")
	create("/docs/")
}
