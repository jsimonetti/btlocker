# Readme
BTlocker is a bluetooth proximity locker. It works by running in the background and measuring the RSSI of a paired bluetooth device.
As soon as the device gets out of range (or disconnects), it locks the computer.
This can work nice with (for example) a cell phone or a smartwatch paired to your laptop.

No longer accidentally leaving your laptop unlocked.

All settings are currently hardcoded into the software. I might add configuration support at a later stage.

# Instructions
1. Install package:
`> go get https://github.com/jsimonetti/btlocker`

2. Change settings in cmd/locker.go

3. Compile binary
`> go build -o bt_locker cmd/locker.go`

4. Set capabilities to avoid having to run this as root:
`> sudo setcap cap_net_admin+eip ./bt_locker`