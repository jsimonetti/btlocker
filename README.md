# Readme
This simple program runs in the background and polls for a specified Bluetooth device.

All settings are currently hardcoded into the software. I might add configuration support at a later stage.

# Instructions
1. Install package:
`> go get https://github.com/jsimonetti/btlocker`

2. Change settings in cmd/locker.go

3. Compile binary
`> go build -o bt_locker cmd/locker.go`

4. Set capabilities to avoid having to run this as root:
`> sudo setcap cap_net_admin+eip ./bt_locker`