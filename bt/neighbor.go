package bt

import (
	"errors"
	"net"
)

type Neighbor net.HardwareAddr

func (n Neighbor) String() string {
	return net.HardwareAddr(n).String()
}

func (n Neighbor) Bytes() []byte {
	return n
}

func NeighborFromString(addr string) Neighbor {
	hw, err := net.ParseMAC(addr)
	if err != nil {
		return Neighbor{}
	}
	return Neighbor(hw)
}

var noSuchNeighbor = errors.New("No such neighbor")
