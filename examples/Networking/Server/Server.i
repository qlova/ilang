
function handler(|client) {
	var ip = info@client("ip")
	output("Connection from "&ip&"\n")
	output("new")
	loop
!		var message = reada@client('\n')
		issues {
			close(client)
			return
		}
		output(message&"\n")
	repeat
}


software {
	load("tcp://8000")
	issues {
		output("Could not open port 8000!\n")
		return
	}
	loop
		fork handler(open("tcp://8000"))
		issues {
			output("THERE WAS AN ERROR\n")
			return
		}
	repeat
}
