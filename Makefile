BIN="sfserver"
BIN_WINDOWS="sfserver.exe"
BUILD_FLAGS="-s -w"
INSTALL_DIR="/usr/local/bin"

build:
	go build -o $(BIN)

compile:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN) -ldflags $(BUILD_FLAGS)
	env GOOS=windows GOARCH=amd64 go build -o $(BIN_WINDOWS) -ldflags $(BUILD_FLAGS)

clean:
	go clean

run: 	build
	./$(BIN)

install:
	cp $(BIN) $(INSTALL_DIR)

all:	build
