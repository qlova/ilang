#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

if [ "$1" = "" ]; then
	LANGUAGE=py
else
	LANGUAGE=$1
fi

function runit {
	case $LANGUAGE in
		py)
			python3 $1.py <<< $(echo -e "$2")
		;;
		go)
			go build && ./$1 <<< $(echo -e "$2")
		;;
		bash) 
			./$1.bash <<< $(echo -e "$2")
		;;
		java) 
			javac $1.java && java $1 <<< $(echo -e "$2")
		;;
		cs) 
			mcs -nowarn:414 /r:mscorlib.dll /r:System.Numerics.dll $1.cs && mono $1.exe <<< $(echo -e "$2")
		;;
		rb)
			ruby $1.rb <<< $(echo -e "$2")
		;;
		lua)
			lua $1.lua <<< $(echo -e "$2")
		;;
	esac
}

function BasicTest {
	cd ../$1/ && ic $1.i && uct -$LANGUAGE $1.u
	if [ "$?" -eq "1" ]; then
		echo -e "$1 \e[31mFAILED!\e[0m to compile D:"
		exit 1
	fi
	local OUTPUT=$(runit $1 "$3")
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

cd Plus
if [ "$2" != "" ]; then
	for l in rb py go java lua rb bash cs; do
		LANGUAGE=$l
		echo $l
		BasicTest "$1" "$2" "$3" 
	done
	exit
fi

BasicTest HelloWorld "Hello World"
BasicTest Chars "97\na"
BasicTest Arrays "2\n4"
BasicTest Concat "This string is joined!"
BasicTest Functions "hi\nhi\ncba"
BasicTest Length "2"
BasicTest OrderOfOperation "405"
BasicTest Plus "4" "2 2\n"
BasicTest Arithmetic "Sum: 8\nDifference: -2\nProduct: 15\nQuotient: 0\nModulus: 3\nExponent: 243" "3 5\n"
BasicTest FileExists "input.txt exists\n/input.txt does not exist\ndocs exists\n/docs does not exist"
BasicTest Maths "d d b"
BasicTest Conditionals "3=3\n3!=2\nverified"
BasicTest Copy "2\n1"
BasicTest Variables "2"
BasicTest CreateFile "output.txt created!\ndocs/ created!\nFailed to create /output.txt\nFailed to create /docs/"
BasicTest DivideByZero "5/0 is a divivision by zero.\n5/2 is not divivision by zero.\n0/0 is a divivision by zero."
BasicTest Strconv "1200  is numeric :)\n3.14  is numeric :)\n3/4  is numeric :)\nabcdefg  is not numeric!"
BasicTest EmptyString "Empty string!"
BasicTest Enviroment "$HOME\n$USER\n$PATH"
BasicTest Empty "\n"
BasicTest Factorial "1\n1\n2\n6\n1124000727777607680000"
BasicTest FizzBuzz "1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n16\n17\nFizz\n19\nBuzz\nFizz\n22\n23\nFizz\nBuzz\n26\nFizz\n28\n29\nFizzBuzz\n31\n32\nFizz\n34\nBuzz\nFizz\n37\n38\nFizz\nBuzz\n41\nFizz\n43\n44\nFizzBuzz\n46\n47\nFizz\n49\nBuzz\nFizz\n52\n53\nFizz\nBuzz\n56\nFizz\n58\n59\nFizzBuzz\n61\n62\nFizz\n64\nBuzz\nFizz\n67\n68\nFizz\nBuzz\n71\nFizz\n73\n74\nFizzBuzz\n76\n77\nFizz\n79\nBuzz\nFizz\n82\n83\nFizz\nBuzz\n86\nFizz\n88\n89\nFizzBuzz\n91\n92\nFizz\n94\nBuzz\nFizz\n97\n98\nFizz\nBuzz\n"
BasicTest Define "6"
BasicTest Greatest "122"
BasicTest IncrementString "101"
BasicTest Input "\n"
BasicTest IntCompare "2 is smaller than 3" "2 3\n"
BasicTest Bases "68\n68\n68"
BasicTest Wierd "true\nfalse\nfalse\nfalse"
BasicTest Logic "false and true is false\nfalse or true is true\nfalse xor true is true\nnot false is true"
BasicTest DownwardFor "10\n9\n8\n7\n6\n5\n4\n3\n2\n1\n0"
BasicTest For "*\n**\n***\n****\n*****"
BasicTest ForStep "2,4,6,8,10,"
BasicTest Sort "2\n3\n4\n5\n6"
BasicTest Binary "111111011001000\n0\n-1"
BasicTest Split "192\n168\n1\n70"
BasicTest Links "1\n2\n3\n"
BasicTest Slices "Hello\nbob\nmy\n"
BasicTest Universal "Hello World\n"
BasicTest Issues "Issue 2\n"
BasicTest Import "Did something!\n"
BasicTest Constant "42\n"
BasicTest ReadFile "This is the contents of the file!"
