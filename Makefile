run:
	air --build.cmd "go build -o bin/memo ." --build.bin "./bin/memo"

b:
	mkdir bin && CGO_ENABLED=0 go build -ldflags="-w -s" -gcflags=all=-l -o bin/memo .