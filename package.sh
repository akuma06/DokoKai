#!/usr/bin/env bash
# Helper script to ease building binary packages for multiple targets.
# Requires the linux64 and mingw64 gcc compilers and zip.
# On Debian-based distros install mingw-w64.

declare -a OSes
OSes[0]='linux;x86_64-linux-gnu-gcc'
OSes[1]='windows;x86_64-w64-mingw32-gcc'
mkdir -p dist

for i in "${OSes[@]}"; do
	arr=(${i//;/ })
	os=${arr[0]}
	cc=${arr[1]}
	rm -f DokoKai DokoKai.exe
	echo -e "\nBuilding $os..."
	echo GOOS=$os GOARCH=amd64 CC=$cc CGO_ENABLED=1 go build -v -ldflags="-X  github.com/akuma06/DokoKai/app.Version=`git describe --tags --always --dirty --match=v*` -X  github.com/akuma06/DokoKai/app.Commit=`git rev-parse HEAD`"
	GOOS=$os GOARCH=amd64 CC=$cc CGO_ENABLED=1 go build -v -ldflags="-X  github.com/akuma06/DokoKai/app.Version=$(git describe --tags --always --dirty --match=v*) -X  github.com/akuma06/DokoKai/app.Commit=$(git rev-parse HEAD)"
	zip -9 -r dist/DokoKai-${version}_${os}_amd64.zip os public templates service/user/locale *.md DokoKai DokoKai.exe
done