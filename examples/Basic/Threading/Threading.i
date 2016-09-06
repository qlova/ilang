
function loopy(n) n {
	var outbox = open("outbox") 
	for i = 0; i < n; i = i + 1
		output(text(i)&"\n")
	repeat
	output@outbox(text(6)&"\n")
}


software {
	var inbox = open("inbox") 
	fork loopy(100)
	fork loopy(100)
	loop
		var message = reada@inbox('\n')
		output(message&"\n")
	repeat
}
