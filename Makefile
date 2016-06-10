all:
	cd ./src && go build -o ../i

install:
	cp ./i /usr/bin/ic
