
concept fact(n, acc) {
	if n = 0
		return acc
	end
	return fact(n-1, n*acc)
}

concept factorial(n) {
  return fact(n, 1)
}

software {
	print(factorial(0))
	print(factorial(1))
	print(factorial(2))
	print(factorial(3))
	print(factorial(22))
}
