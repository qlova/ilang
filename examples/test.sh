#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

#Test A+B
cd A+B/ && ic a+b.i && uct -py a+b.u
if [ "$?" -eq "1" ]; then
	echo -e "A+B \e[31mFAILED!\e[0m to compile D:"
	exit 1
fi
OUTPUT=$(python3 a+b.py <<< $(echo -e "2 2\n"))
if [ "$OUTPUT" = "4" ]; then
	echo -e "A+B \e[32mPASSED!\e[0m"
else
	echo -e "A+B \e[31mFAILED!\e[0m"
	exit 1
fi

#Test Arithmetic
cd ../Arithmetic/ && ic Arithmetic.i && uct -py Arithmetic.u
if [ "$?" -eq "1" ]; then
	echo -e "Arithmetic \e[31mFAILED!\e[0m to compile D:"
	exit 1
fi
OUTPUT=$(python3 Arithmetic.py <<< $(echo -e "3 5\n"))
DEFINED=$(echo -e "Sum: 8\nDifference: -2\nProduct: 15\nQuotient: 0\nModulus: 3\nExponent: 243")
if [ "$OUTPUT" = "$DEFINED" ]; then
	echo -e "Arithmetic \e[32mPASSED!\e[0m!"
else
	echo -e "Arithmetic \e[31mFAILED!\e[0m!"
	echo "$OUTPUT"
	exit 1
fi

function BasicTest {
	cd ../$1/ && ic $1.i && uct -py $1.u
	if [ "$?" -eq "1" ]; then
		echo -e "$1 \e[31mFAILED!\e[0m to compile D:"
		exit 1
	fi
	local OUTPUT=$(python3 $1.py)
	local DEFINED=$(echo -e "$2")
	if [ "$OUTPUT" = "$DEFINED" ]; then
		echo -e "$1 \e[32mPASSED!\e[0m"
	else
		echo -e "$1 \e[31mFAILED!\e[0m Got:"
		echo	 "$OUTPUT"
		echo "(Expecting)"
		echo	 "$DEFINED"
		exit 1
	fi
}

BasicTest Arrays "2\n4"
BasicTest Chars "97\na"
BasicTest Concat "This string is joined!"
BasicTest Functions "hi\nhi\ncba"
BasicTest HelloWorld "Hello World"
BasicTest Length "2"
BasicTest OrderOfOperation "405"
BasicTest FileExists "input.txt exists\n/input.txt does not exist\ndocs exists\n/docs does not exist"
BasicTest Maths "d d b"
