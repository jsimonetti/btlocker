package main

import (
	"github.com/jsimonetti/btlocker/bt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var rssi_threshold = -15
var threshold_duration = 5 * time.Second
var poll_interval = 2 * time.Second
var lockCmd = "/usr/bin/loginctl"
var lockArgs = "lock-session"
var unlockCmd = "/usr/bin/loginctl"
var unlockArgs = "unlock-session"
var unlockEnabled = false
var neighbor = "00:00:00:f3:6d:6b"

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	timer := time.Tick(poll_interval)
	lastInRange := time.Now()

	prevState := disconnected
	locked := false

	for {
		select {
		case <-sigs:
			os.Exit(0)

		case <-timer:
			if info, err := bt.GetConnInfo(bt.NeighborFromString(neighbor)); err == nil {
				switch curState(info) {
				case inRange:
					lastInRange = time.Now()
					if prevState != inRange {
						prevState = inRange
						log.Print("device " + neighbor + " in range")
						unlock()
						locked = false
					}
				case disconnected:
					if prevState != disconnected {
						prevState = disconnected
						log.Print("device " + neighbor + "disconnected")
						lock()
						locked = true
					}
				case outOfRange:
					if prevState == outOfRange {
						if lastInRange.Before(time.Now().Add(threshold_duration)) {
							if !locked {
								locked = true
								log.Print("device " + neighbor + " out of range for " + threshold_duration.String() + "")
								lock()
							}
						} else {
							log.Print("device " + neighbor + " out of range")
						}
					} else {
						prevState = outOfRange
					}
				}
			}
		}
	}
}

type state int

const (
	inRange state = iota
	outOfRange
	disconnected
)

func curState(info bt.ConnInfo) state {
	if info.RSSI >= rssi_threshold && info.MAXTXPower != 0 {
		return inRange
	}

	if info.RSSI < rssi_threshold && info.MAXTXPower != 0 {
		return outOfRange
	}

	return disconnected
}

func lock() {
	if err := exec.Command(lockCmd, lockArgs).Run(); err != nil {
		log.Print("error locking session: " + err.Error())
	} else {
		log.Print("locked session")
	}
}

func unlock() {
	if unlockEnabled {
		if err := exec.Command(unlockCmd, unlockArgs).Run(); err != nil {
			log.Print("error unlocking session: " + err.Error())
		} else {
			log.Print("unlocked session")
		}
	}
}
