function f() {
	print("hi")
}

function call( ()a ) {
	a()
}

function add(a,b) r {
	return a+b
}

function printchars( ..x ) {
	output(x)
}

software {
	var b = f
	b()
	call(b)
	printchars(add(40, 59), 98, 97, 10)
}
