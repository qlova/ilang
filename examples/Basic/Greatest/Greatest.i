
function largest([]l) r {
	var i = 0
	var r = l.0
	loop
		if i >= len(l)
			return r 
		end
		
		if l.i > r
			r = l.i
		end
		
		i = i + 1
	repeat
}

software {
	var []list = "The Quick Brown Fox, Jumped Over The Lazy Log!"
	
	output(text(largest(list))&"\n")
}