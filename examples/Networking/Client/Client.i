
software {
	var server = open("tcp://localhost:8000")
	server("Hello World\n")
	close(server)
}
