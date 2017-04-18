#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

if [ "$2" = "" ]; then
	LANGUAGE=py
else
	LANGUAGE=$2
fi

function runit {
	case $LANGUAGE in
		py)
			cd ./.it && python3 $1.py <<< $(echo -e "$2")
		;;
		go)
			cd ./.it && go build -o ../$1 && cd .. && ./$1 <<< $(echo -e "$2")
		;;
		
		rs)
			./$1 <<< $(echo -e "$2")
		;;
		
		bash) 
			./$1.bash <<< $(echo -e "$2")
		;;
		java) 
			cd ./.it && javac $1.java && java $1 <<< $(echo -e "$2")
		;;
		cs) 
			mcs -nowarn:414 /r:mscorlib.dll /r:System.Numerics.dll $1.cs && mono $1.exe <<< $(echo -e "$2")
		;;
		rb)
			cd ./.it && ruby $1.rb <<< $(echo -e "$2")
		;;
		lua)
			cd ./.it && lua $1.lua <<< $(echo -e "$2") && cd .. 
		;;
		js) cd ./.it && nodejs $1.js <<< $(echo -e "$2") && cd .. 
		;;
	esac
}

function BasicTest {
	cd ./$1/ && it build -$LANGUAGE
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
	cd ..
}

function TESTING {
	if [ "$2" != "" ]; then
		for l in rb py go java lua rb js bash cs; do
			LANGUAGE=$l
			echo $l
			BasicTest "$1" "$2" "$3" 
		done
		exit
	fi
}

export -f TESTING
export -f BasicTest
export -f runit

if [ "$1" != "" ]; then
	cd ./$1 && ./test.sh $2 $3 $4
else
	cd Structs && ./test.sh && cd ../Basic && ./test.sh
fi
