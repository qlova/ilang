
function create( []filename ) {
!	var file = open(filename) 
	output@file("")
	issues {
		print("Failed to create "+filename)
		return
	}
	print(filename+" created!")
	close(file)
}

software {
	create("output.txt")
	create("docs/")
	create("/output.txt")
	create("/docs/")
}
