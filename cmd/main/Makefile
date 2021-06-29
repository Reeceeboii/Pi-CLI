# REQUIREMENTS
# 	Make
#	go toolchain
#	git
#	gcc
#	gcc-mingw-w64-x86-64 (cross compiler required if compiling for windows on Mac or Linux
#	MinGW (for Windows)

# Parameters
GO-CMD=go
GO-BUILD=$(GO-CMD) build

# Binary names
WINDOWS-BINARY=picli.exe
LINUX-BINARY=picli-linux
MAC-BINARY=picli-mac

# Current Git commit hash
CURR-GIT-HASH=$(shell git rev-parse --short HEAD)

# Linker flags
LINKER_FLAGS=-ldflags "-s -w -X 'main.PICLIVersion=$(v)' -X 'main.GitHash=$(CURR-GIT-HASH)'"

# Compiling for Windows from Linux requires the MinGW-w64 x86 cross compiler
CGO_ENABLED=CGO_ENABLED=1
MINGW_CROSS_COMPILER=CC=x86_64-w64-mingw32-gcc

# Cross-compilation compiler operating system flags
GOOS-WINDOWS=GOOS=windows
GOOS-MAC=GOOS=darwin
GOOS-LINUX=GOOS=linux

# If a version number is not provided to the makefile, exit. This is a requirement
ifndef v
$(error v is not set)
endif

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
ifeq ($(OS), Windows_NT)
	$(error This target cannot be ran via Windows yet. Please use WSL instead...)
	exit
endif
	@echo "Compiling for Mac..."
	$(GOOS-MAC) $(GO-BUILD) $(LINKER_FLAGS) -o $(MAC-BINARY)

# Compiling Pi-CLI for Linux
linux:
ifeq ($(OS), Windows_NT)
	$(error This target cannot be ran via Windows yet. Please use WSL instead...)
	exit
endif
	@echo "Compiling for Linux..."
	$(GOOS-LINUX) $(GO-BUILD) $(LINKER_FLAGS) -o $(LINUX-BINARY)
