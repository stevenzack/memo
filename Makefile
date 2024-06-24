run:
	air --build.cmd "go build -o bin/main ." --build.bin "./bin/main"

b:
	mkdir -p bin && CGO_ENABLED=0 go build -ldflags="-w -s" -gcflags=all=-l -o bin/main .