
function thread() {
	print("This is a thread!")
}


software {
	fork thread()
	execute("sleep 0.5")
}
