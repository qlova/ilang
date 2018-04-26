concept bool(n) {
	if n return "true" else return "false" end
}

concept logic(a, b) {
	print(bool(a), " and ", bool(b) ," is ", bool(a×b))
	print(bool(a), " or ", bool(b) ," is ", bool(a+b))
	print(bool(a), " xor ", bool(b) ," is ", bool(a-b))
	print("not ", bool(a) ," is ", bool(not(a)))
}

software {
	logic(false, true)
}
