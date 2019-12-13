
set -e

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go build -o bin/posprinter-linux64 main_linux.go printer.go
env GOOS=darwin GOARCH=386 go build -o bin/posprinter-mac main_linux.go printer.go

#cp ~/amura/src/misc/posprinter.net/POSPrinter/bin/Debug/posprinter.exe bin/

mv bin/posprinter-mac ~/projects/posprinter/

#env GOOS=windows GOARCH=amd64 go build -o bin/posprinter.exe main_windows.go printer.go
#env GOOS=windows GOARCH=386 go build -o bin/posprinter32.exe main_windows.go printer.go 

