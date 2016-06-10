#! /bin/bash

#Test A+B
cd A+B/ && ic a+b.i && uct -py a+b.u
OUTPUT=$(python3 a+b.py <<< $(echo -e "2 2\n"))
if [ "$OUTPUT" = "4" ]; then
	echo "A+B PASSED!"
else
	echo "A+B FAILED!"
	exit 1
fi

#Test Arithmetic
cd ../Arithmetic/ && ic Arithmetic.i && uct -py Arithmetic.u
OUTPUT=$(python3 Arithmetic.py <<< $(echo -e "3 5\n"))
DEFINED=$(echo -e "Sum: 8\nDifference: -2\nProduct: 15\nQuotient: 0\nModulus: 3\nExponent: 243")
if [ "$OUTPUT" = "$DEFINED" ]; then
	echo "Arithmetic PASSED!"
else
	echo "Arithmetic FAILED!"
	echo "$OUTPUT"
	exit 1
fi

function BasicTest {
	cd ../$1/ && ic $1.i && uct -py $1.u
	local OUTPUT=$(python3 $1.py)
	local DEFINED=$(echo -e "$2")
	if [ "$OUTPUT" = "$DEFINED" ]; then
		echo "$1 PASSED!"
	else
		echo "$1 FAILED! Got:"
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
