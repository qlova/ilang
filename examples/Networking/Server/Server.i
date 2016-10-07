
function handler(||client) {
	var ip = info@client("ip")
	print("Connection from ", ip)
	output("new")
	loop
!		var message = reada@client('\n')
		issues {
			close(client)
			return
		}
		print(message)
	repeat
}


software {
	load("tcp://8000")
	issues {
		print("Could not open port 8000!")
		return
	}
	loop
		fork handler(open("tcp://8000"))
		issues {
			print("THERE WAS AN ERROR")
			return
		}
	repeat
}
