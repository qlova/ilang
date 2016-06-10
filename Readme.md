# i
i is a hardware-agnostic maths-based cross-platform programming language in heavy development.

Hello World:

		software {
			output("Hello World")
		}

There are no official releases yet as the language is in a aplha state but you can grab the source and start hacking!  
Hacking instructions: (Linux and Mac)
	
		#You need UCT.
		git clone https://github.com/Qlova/uct
		cd uct && make && sudo make install
		
		cd ../
		git clone https://github.com/Qlova/ilang
		cd ilang && make && sudo make install
		
		#Compile examples.
		cd examples && ./test.sh
		
		#You can now play around with the examples.
		#Standard building looks like this:
		ic File.i && uct -go File.u && go build File.go

Please be aware that many features are missing or incomplete in i!
		
