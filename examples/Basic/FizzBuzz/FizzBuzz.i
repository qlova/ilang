software {
	var i = 0
	loop
		i = i + 1
		
		if i mod 15 = 0
			output("FizzBuzz\n")
		elseif i mod 3 = 0
			output("Fizz\n")
		elseif i mod 5 = 0
			output("Buzz\n")
		else
			output(text(i)&"\n")
		end
		
		if i >= 100
			break
		end
	repeat
}
