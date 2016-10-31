plugin Moddable {
	""message
}

plugin new(Moddable) {
	message = "Hello World"
}

plugin run(Moddable) {
	print(message)
}
