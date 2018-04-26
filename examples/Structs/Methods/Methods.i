
type Math {

	concept add(a, b) {
		return a + b
	}

	concept combine(a, b) {
		return add(a, b)
	}

}

software {
	print(Math.add(2, 2))
	print(Math.combine(2, 3))
}
