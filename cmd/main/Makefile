# REQUIREMENTS
# 	Make
#	go toolchain
# LINUX SPECIFIC
#	gcc
#	gcc-mingw-w64-x86-64
# WINDOWS SPECIFIC
#	MinGW

# Parameters
GO-CMD=go
GO-BUILD=$(GO-CMD) build

# Binary names
WINDOWS-BINARY=picli.exe
LINUX-BINARY=picli-linux
MAC-BINARY=picli-mac

# Linker flags
LINKER_FLAGS=-ldflags "-s -w"

# Compiling for Windows from Linux requires the MinGW-w64 x86 cross compiler
CGO_ENABLED=CGO_ENABLED=1
MINGW_CROSS_COMPILER=CC=x86_64-w64-mingw32-gcc

# Cross-compilation compiler operating system flags
GOOS-WINDOWS=GOOS=windows
GOOS-MAC=GOOS=darwin
GOOS-LINUX=GOOS=linux

# Compiling Pi-CLI for Windows
win: $(objects)
	@echo "Compiling for Windows..."
ifeq ($(OS), Windows_NT)
	$(GO-BUILD) $(LINKER_FLAGS) -o $(WINDOWS-BINARY)
else
	$(GOOS-WINDOWS) $(MINGW_CROSS_COMPILER) $(CGO_ENABLED) $(GO-BUILD) $(LINKER_FLAGS) -o $(WINDOWS-BINARY)
endif

# Compiling Pi-CLI for Mac
mac:
	@echo "Compiling for Mac..."
	$(GOOS-MAC) $(GO-BUILD) $(LINKER_FLAGS) -o $(MAC-BINARY)

# Compiling Pi-CLI for Linux
linux:
	@echo "Compiling for Linux..."
	$(GOOS-LINUX) $(GO-BUILD) $(LINKER_FLAGS) -o $(LINUX-BINARY)
