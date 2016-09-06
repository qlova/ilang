all:
	cd ./src && go build -o ../i

install:
	cp ./i /usr/bin/ic
	cp ./doc/i.lang /usr/share/gtksourceview-3.0/language-specs/i.lang
