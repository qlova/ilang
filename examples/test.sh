#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

if [[ -z "$2" ]]; then
	export LANGUAGE=py
else
	export LANGUAGE=$2
fi

function runit {
	case $LANGUAGE in
		py)
			python3 ./.it/main.py <<< $(echo -e "$2")
		;;
		go)
			go run ./.it/main.go  <<< $(echo -e "$2")
		;;
		
		rs)
			./$1 <<< $(echo -e "$2")
		;;
		
		bash) 
			./$1.bash <<< $(echo -e "$2")
		;;
		java) 
			rm Runtime.java 2> /dev/null
			mv main.java Runtime.java 2> /dev/null
			cd ./.it && javac Runtime.java && java Runtime <<< $(echo -e "$2")
		;;
		cs) 
			cd ./.it && mcs -nowarn:414 /r:mscorlib.dll /r:System.Numerics.dll $1.cs stack.cs > /dev/null && mono $1.exe <<< $(echo -e "$2")
		;;
		rb)
			cd ./.it && ruby $1.rb <<< $(echo -e "$2")
		;;
		lua)
			cd ./.it && lua $1.lua <<< $(echo -e "$2") && cd .. 
		;;
		js) cd ./.it && nodejs $1.js <<< $(echo -e "$2") && cd .. 
		;;
		sh) cd ./.it && bash $1.sh <<< $(echo -e "$2") && cd .. 
		;;
	esac
}

function BasicTest {
	cd "${1}" && it build $LANGUAGE
	if [ "$?" -eq "1" ]; then
		echo -e "$1 \e[31mFAILED!\e[0m to compile D:"
		exit 1
	fi
	local OUTPUT=$(runit "$1" "$3")
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
	rm -r ./.it
	cd ..
}

function GraphicsTest {
	echo -n "$1"
	echo -e " NOTSURE!"
}

function Passed {
	echo -e " \e[32mPASSED!\e[0m"
}

function FakeTest {
	echo -n "$1"
	echo -e " \e[32mFAILED!\e[0m"
}

function Failed {
	echo -e " \e[31mFAILED!\e[0m"
}

function IgnoreTest {
	cd "${1}" && it build $LANGUAGE
	runit $1 "$3"
	echo -n "$1"
	cd ..
}

function TESTING {
	if [ ! -z "$2" ]; then
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
export -f IgnoreTest
export -f FakeTest
export -f GraphicsTest
export -f runit
export -f Passed
export -f Failed

if [[ ! -z "$1" ]]; then
	cd ./$1 && ./test.sh $2 $3 $4
else
	cd Structs && ./test.sh && cd ../Basic && ./test.sh
fi
