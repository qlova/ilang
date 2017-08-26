
type Math {}

method Math.add(a, b) r {
	return a + b
}

method Math.combine(a, b) r {
	return add(a, b)
}

software {
	print(Math.add(2, 2))
	print(Math.combine(2, 3))
}
