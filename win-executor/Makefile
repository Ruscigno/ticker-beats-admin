.PHONY: build
build:
	env GOOS=windows GOARCH=amd64 go build main.go
	cp ./main.exe ~/metatrader5/win-executor/win-executor.exe
	cp ./config.json ~/metatrader5/win-executor/config-template.json