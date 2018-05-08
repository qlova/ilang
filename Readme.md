## The 'i' programming language
![i](spotlight.png)  

'i' is a hardware-agnostic cross-platform creative programming language in heavy development.

Hello World:

		software {
			print("Hello World")
		}

## The Promise  
'i' offers two promises:

* Determinism
* Resilience

## Deterministic  
Given equivalent inputs, a compiled program in 'i' will provide identical outputs on any platform/target.  

## Resilient  
Given hardware with an infinite amount of memory, 'i' will never crash.
		
## Design
'i' is a language which is meant to be clean and concise.

```
	//Open a file and print its contents.
	software {
		file = open("input.txt")
		loop {
			print(file())
			issues {
				exit
			}
		}
	}
```

## Features

* Big integers.
* Clean.
* Simple.
* Statically typed.
* Feels dynamic.
* Cross-platform.
* Has cross-platform graphics support.

## Documentation

There is now a reddit page! http://reddit.com/r/ilang

Documentation is minimal but take a look at the [Wiki](https://github.com/Qlova/ilang/wiki).
The Rosetta Code [page](http://rosettacode.org/wiki/Category:I) may be helpful. 
Flick me an email (splizard @ splizard dot com) if you have any questions.

## DOWNLOAD AND INSTALL

### Windows
There is an alpha release for windows, it must be used from the command line.
You can find it at https://bitbucket.org/Splizard/ilang-release/downloads/it.exe

### Linux/Mac or Android (Termux)
Here are the hacking instructions:

		go get -u github.com/qlova/ilang/src/it
		echo "The binary is now located in:"
		echo "$GOPATH/bin/it.exe"

Please be aware that many features are missing or incomplete in i!

## Compilation

In order to compile run code written in i, place it in an isolated directory and run:

		cd /path/to/directory/
		it run

By default, the code is compiled into the Go programming language. (You will need Go installed to complete this process)
Otherwise you can target other langauges by providing their extension as an argument eg.

		it run py
		it run java
		it run js
		
You can export a distributable binary of the code by using the export command.

		it export py
		it export java
		it export js

A full set of supported languages can be found [here](http://github.com/qlova/uct)
