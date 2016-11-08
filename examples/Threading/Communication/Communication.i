function spinoff() {
	execute("sleep 2")
	outbox("done")
}

software {
	fork spinoff()
	var m = inbox()
	
	print(m) //This will print done.
}
