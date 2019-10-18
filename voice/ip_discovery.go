package voice

import (
	"encoding/binary"
	"errors"
	"net"
	"strings"
)

// ipDiscovery uses Discord's IP discovery service to get the external ip and port the
// given UDP connection is using.
func ipDiscovery(conn *net.UDPConn, ssrc uint32) (ip string, port uint16, err error) {
	b := make([]byte, 70)
	binary.BigEndian.PutUint32(b, ssrc)
	if _, err = conn.Write(b); err != nil {
		return "", 0, err
	}

	b = make([]byte, 70)
	l, err := conn.Read(b)
	if err != nil {
		return "", 0, err
	}
	if l < 70 {
		return "", 0, errors.New("ipDiscovery: did not receive enough bytes")
	}

	s := strings.Builder{}
	for i := 4; b[i] != 0; i++ {
		if err = s.WriteByte(b[i]); err != nil {
			return "", 0, err
		}
	}
	return s.String(), binary.LittleEndian.Uint16(b[68:70]), nil
}
