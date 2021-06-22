# Parameters
GO-CMD=go
GO-BUILD=$(GO-CMD) build
GO-TEST=$(GO-CMD) test

# Binary names
WINDOWS-BINARY=picli.exe
LINUX-BINARY=picli-linux
MAC-BINARY=picli-mac

# Misc. compiler flags
REMOVE-DEBUG-FLAGS=-ldflags "-s -w"

# Cross-compilation compiler operating system flags
GOOS-WINDOWS=GOOS=windows
GOOS-MAC=GOOS=darwin
GOOS-LINUX=GOOS=linux

# Compiling Pi-CLI for Windows
win:
	@echo "Compiling for Windows..."
	$(GOOS-WINDOWS) $(GO-BUILD) $(REMOVE-DEBUG-FLAGS) -o $(WINDOWS-BINARY)

# Compiling for Windows *on* a Windows machine. This is required as
# placing environment variables inline with the go compiler
# (or any other executable for that matter), is not supported on Windows.
# This allows compilation via PowerShell (or cmd if you're a masochist)
# ------------
# As an aside, the Linux and Mac compilation options can be ran from a windows host by passing the make command
# through to WSL first.
# I.e. `wsl make mac` will allow the prefixed args to work as expected.
win-self:
	@echo "Compiling for Windows (on a Windows dev machine)..."
	$(GO-BUILD) $(REMOVE-DEBUG-FLAGS) -o $(WINDOWS-BINARY)

# Compiling Pi-CLI for Mac
mac:
	@echo "Compiling for Mac..."
	$(GOOS-MAC) $(GO-BUILD) $(REMOVE-DEBUG-FLAGS) -o $(MAC-BINARY)

# Compiling Pi-CLI for Linux
linux:
	@echo "Compiling for Linux..."
	$(GOOS-LINUX) $(GO-BUILD) $(REMOVE-DEBUG-FLAGS) -o $(LINUX-BINARY)
