# parameters
GO-CMD=go
GO-BUILD=$(GO-CMD) build
GO-TEST=$(GO-CMD) test

# binary names
WINDOWS-BINARY=picli.exe
LINUX-BINARY=picli-linux
MAC-BINARY=picli-mac

# misc compiler flags
REMOVE-DEBUG-FLAGS=-ldflags "-s -w"

# cross compilation compiler flags
GOOS-WINDOWS=GOOS=windows
GOOS-MAC=GOOS=darwin
GOOS-LINUX=GOOS=linux

# compiling Pi-CLI for Windows
win:
	@echo "Compiling for Windows..."
	$(GOOS-WINDOWS) $(GO-BUILD) $(REMOVE-DEBUG-FLAGS) -o $(WINDOWS-BINARY)

# compiling Pi-CLI for Mac
mac:
	@echo "Compiling for Mac..."
	$(GOOS-MAC) $(GO-BUILD) $(REMOVE-DEBUG-FLAGS) -o $(MAC-BINARY)

# compiling Pi-CLI for Linux
linux:
	@echo "Compiling for Linux..."
	$(GOOS-LINUX) $(GO-BUILD) $(REMOVE-DEBUG-FLAGS) -o $(LINUX-BINARY)

# clean up
clean:
	@echo "Cleaning up..."
	@rm ./picli
	@rm ./picli.exe
