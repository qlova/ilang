function thread(||p) {
	var message = p()
	print(message)
}

software {
	var c = ||
	fork thread(c)
	c("Hello World")
}
