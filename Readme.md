# i
i is a hardware-agnostic maths-based cross-platform programming language in heavy development.

Hello World:

		software {
			print("Hello World")
		}
		
## Design
i is a language which is meant to be as clean and visually pleasing as possible, too many languages obfuscate code (Looking at you Java).
So I created i to solve this.

Here is a contrast between different languages and how visually pleasing they are at common tasks.

**File IO:**
Java:
```java
//This is ugly.
for (String filename : args) {
    try (FileReader fr = new FileReader(filename);BufferedReader br = new BufferedReader(fr)){
        String line;
        int lineNo = 0;
        while ((line = br.readLine()) != null) {
            processLine(++lineNo, line);
        }
    }
    catch (Exception x) {
        x.printStackTrace();
    }
}
```
Go:
```go
//Marginaly better
inputFile, err := os.Open("byline.go")
if err != nil {
	log.Fatal("Error opening input file:", err)
}
defer inputFile.Close()

scanner := bufio.NewScanner(inputFile)

for scanner.Scan() {
	fmt.Println(scanner.Text())
}

if err := scanner.Err(); err != nil {
	log.Fatal(scanner.Err())
}
```
Lua:
```lua
--Now we're talking! Still kinda obscure.
filename = "input.txt"
fp = io.open( filename, "r" )
 
for line in fp:lines() do
    print( line )
end
 
fp:close()
```
I:
```
	//Much better, well structured, clean and understandable.
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

So you can hopefully see some of the philosophy of i.

## Features

* Big integars.
* Clean.
* Simple.
* Statically typed.
* Feels dynamic.
* Cross-platform.
* Has a Gui library.
* Has a **WIP** game engine.

## Documentation

Documentation is minimal but take a look at the [Wiki](https://github.com/Qlova/ilang/wiki).
The Rosetta Code [page](http://rosettacode.org/wiki/Category:I) may be helpful. 
Flick me an email (splizard @ splizard dot com) if you have any questions =)

## Types
There are 13 types in the i language.
```
number()
rational()	\
decimal()   .
letter()	' '
array() 	[]
text() 		""
set()		<>

pipe 		||
function 	()

Thing() 	{}
List()  	..
Table()		:
Something() ?
```

**number**
Numbers in the i language can be integers of any size.

**decimal**
Decimals are fixed-precision numbers.

**letter**
Letters correspond to a character in a piece of text, for example "a", "!" etc.

**set**
A set is an unordered collection of arbitary labels.

**array**
An array is a list of numbers.

**text**
Text is a string of letters, these can form words, sentences and the like.

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

In order to compile code writtin in I, place it in an isolated directory and run:

		cd /path/to/directory/
		it build

By default, the code is compiled into the Go programming language. (You will need Go installed to complete this process)
Otherwise you can target other langauges by providing their extension as an argument eg.

		it build py
		it build java
		it build js

A full set of supported languages can be found [here](http://github.com/qlova/uct)
