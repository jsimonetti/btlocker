package bt

import (
	"errors"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"os"
)

const (
	HCI_DEV_NONE = 0xffff
)

type ConnInfo struct {
	RSSI       int
	TXPower    int
	MAXTXPower int
}

func GetConnInfo(neigh Neighbor) (ConnInfo, error) {
	info, err := getConnInfo(neigh.Bytes())
	if err != nil {
		return ConnInfo{}, noSuchNeighbor
	}
	return *info, nil
}

func getConnInfo(neigh []byte) (*ConnInfo, error) {
	// connect a socket to AF_BLUETOOTH
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW|unix.SOCK_CLOEXEC, unix.BTPROTO_HCI)
	if err != nil {
		return nil, err
	}
	defer unix.Close(fd)

	// use the HCI control channel and special device NONE
	sa := unix.SockaddrHCI{
		Dev:     HCI_DEV_NONE,
		Channel: unix.HCI_CHANNEL_CONTROL,
	}

	// bind the socket
	err = unix.Bind(fd, &sa)
	if err != nil {
		return nil, err
	}

	// extract a connection from the fd
	conn := os.NewFile(uintptr(fd), "AF_BLUETOOTH")

	// build command structure to get connection info
	wr := []byte{0x31, 0x00, 0x00, 0x00, 0x07, 0x00}
	for i := len(neigh) - 1; i >= 0; i-- {
		wr = append(wr, neigh[i])
	}
	wr = append(wr, 0x00)
	conn.Write(wr)

	b := make([]byte, 32)
	n, err := conn.Read(b)
	if err != nil {
		if err != io.EOF {
			log.Print(err)
			return nil, err
		}
	}

	// the reply should be 19 bytes
	if n != 19 {
		return nil, errors.New("wrong reply size")
	}

	rssi := int(b[16])
	if rssi > 0 {
		rssi = rssi - 255
	}
	tx := int(b[17])
	max := int(b[18])

	return &ConnInfo{rssi, tx, max}, nil
}
