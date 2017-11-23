all:
	cd ./src/it && go build -o ../../it

install:
	cp ./ic /usr/bin/ic
	cp ./it /usr/bin/it
	cp ./doc/i.lang /usr/share/gtksourceview-3.0/language-specs/i.lang

windows:
	cd ./src/ic && GOOS=windows go build -o ../../ic.exe
	cd ./src/it && GOOS=windows go build -o ../../it.exe
