function logic(a, b) {
	output(bool(a)&" and "&bool(b)&" is "&bool(a×b)&"\n")
	output(bool(a)&" or "&bool(b)&" is "&bool(a+b)&"\n")
	output(bool(a)&" xor "&bool(b)&" is "&bool(a-b)&"\n")
	output("not "&bool(a)&" is "&bool(a×0)&"\n")
}

software {
	logic(false, true)
}
