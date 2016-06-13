
function fact(n, acc) r {
	if n = 0
		return acc
	end
	return fact(n-1, n*acc)
}

function factorial(n) r {
  return fact(n, 1)
}

software {
	output(text(factorial(0))&"\n")
	output(text(factorial(1))&"\n")
	output(text(factorial(2))&"\n")
	output(text(factorial(3))&"\n")
	output(text(factorial(22))&"\n")
}
