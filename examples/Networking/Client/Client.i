
software {
	var server = open("tcp://localhost:8000")
	output@server("Hello World\n")
	close(server)
}
