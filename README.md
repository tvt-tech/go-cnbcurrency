# CNB currency scrapper


**Build dll using GCC**
> [!NOTE]
> Use [codeblocks-20.03-32bit-mingw-32bit-nosetup](https://www.codeblocks.org/downloads/binaries/)
> to be able compile dll for winXP/Win7 32bit

**386**
```commandline
$env:GOOS="windows"
$env:GOARCH="386"
$env:CGO_ENABLED=1
go build -ldflags "-s -w" -buildmode=c-shared -o currency32.dll main.go
objdump -f currency32.dll
```

**amd64**
```commandline
$env:GOOS="windows"
$env:GOARCH="amd64"
$env:CGO_ENABLED=1
go build -ldflags "-s -w" -buildmode=c-shared -o currency64.dll main.go
objdump -f currency32.dll
```