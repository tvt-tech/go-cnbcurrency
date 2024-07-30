# CNB currency scrapper


**Build dll using GCC**
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