concept f() {
	print("hi")
}

concept call(a) {
	a()
}

concept add(a,b) {
	return a+b
}

concept PrintLetters(x...) {
	print(text(x))
}

software {
	f()
	b = f
	b()
	call(b)
	
	PrintLetters(add(40, 59), 98, 97)
}
