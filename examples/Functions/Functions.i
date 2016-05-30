function f() {
	output("hi\n")
}

function call( ()a ) {
	a()
}

software {
	b=f
	b()
	call(b)
}
