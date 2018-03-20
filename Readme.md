# 'i'
'i' is a hardware-agnostic maths-based cross-platform programming language in heavy development.

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
The result of 0/0 is considered as an input.

## Resilient  
Given hardware with an infinite amount of memory, 'i' will never crash.
		
## Design
'i' is a language which is meant to be clean and concise.

```
	//Open a file and print its contents.
	software {
		var file = open("input.txt")
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
* Has a Gui library.
* Has a **WIP** game engine.

## Documentation

There is now a reddit page! http://reddit.com/r/ilang

Documentation is minimal but take a look at the [Wiki](https://github.com/Qlova/ilang/wiki).
The Rosetta Code [page](http://rosettacode.org/wiki/Category:I) may be helpful. 
Flick me an email (splizard @ splizard dot com) if you have any questions.

## Types
These are all of the builtin types in the 'i' language. Each type has an associate set of symbols.
```
number()
rational()	\
decimal()	.
duplex()	±
letter()	' '
array() 	[]
text() 		""
collection()	<>

pipe 		||
function 	()

Thing() 	{}
List()  	..
Table()		:
Something() 	?
```

**number**
Numbers in the 'i' language can be integers of any size.

**rational**
Rational numbers are pairs of numbers representing a fractional value.

**decimal**
Decimals are fixed-precision decimal numbers.

**letter**
Letters correspond to a character in a piece of text, for example "a", "!" etc.

**array**
An array is a list of numbers.

**text**
Text is a string of letters, these can form words, sentences and the like.

**set**
A set is an unordered collection of arbitary labels.


**pipe**
A pipe is any object with an input and an output. These can be files, the internet etc.

**function**
A function is a special form of pipe, one which you can call.


**thing**
Any types which are added in a program are "things".

**Something** (WIP)
Something is the only builtin struture type, it can hold any value and needs to be type checked when used.  


## DOWNLOAD AND INSTALL

There is an alpha release for windows, it must be used from the command line.
You can find it at https://bitbucket.org/Splizard/ilang-release/downloads/it.exe

Otherwise here are the hacking instructions: (Linux and Mac)

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
