run:
	air --build.cmd "go build -o bin/main ." --build.bin "./bin/main"

b:
	CGO_ENABLED=0 go build -ldflags="-w -s" -gcflags=all=-l -o main .