
concept largest(l) {
	i = 0
	r = l[0]
	loop {
		if i = l.size()
			return r 
		end
		
		if number(l[i]) > number(r)
			r = l[i]
		end
		
		i = i + 1
	}
}

software {
	letters = "The Quick Brown Fox, Jumped Over The Lazy Log!"
	
	print(largest(letters))
}
