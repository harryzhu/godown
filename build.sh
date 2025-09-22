CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/macos_arm/godown -ldflags "-w -s" main.go
zip dist/macos_arm/godown_macos_arm.zip dist/macos_arm/godown

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macos_intel/godown -ldflags "-w -s" main.go
zip dist/macos_intel/godown_macos_intel.zip dist/macos_intel/godown


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux_amd64/godown -ldflags "-w -s" main.go
zip dist/linux_amd64/godown_linux_amd64.zip dist/linux_amd64/godown


CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows_amd64/godown.exe -ldflags "-w -s" main.go
zip dist/windows_amd64/godown_windows_amd64.zip dist/windows_amd64/godown.exe
