concept bool(n) {
	if n
		return "true"
	else
		return "false"
	end
}

software {
	print(bool(0/0))
	print(bool(0*0))
	print(bool(0^3))
	print(bool(0^8))
}
