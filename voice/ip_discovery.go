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
	b := make([]byte, 74)
	binary.BigEndian.PutUint16(b, 1)        // Packet type (0x1 is request, 0x2 is response)
	binary.BigEndian.PutUint16(b[2:], 70)   // Packet length (excluding type and length fields)
	binary.BigEndian.PutUint32(b[4:], ssrc) // The SSRC code from the Op 2 VoiceConnection event
	if _, err = conn.Write(b); err != nil {
		return "", 0, err
	}

	b = make([]byte, 74)
	l, _, err := conn.ReadFromUDP(b)
	if err != nil {
		return "", 0, err
	}
	if l < 74 {
		return "", 0, errors.New("ipDiscovery: did not receive enough bytes")
	}

	s := strings.Builder{}
	for i := 8; b[i] != 0; i++ {
		if err = s.WriteByte(b[i]); err != nil {
			return "", 0, err
		}
	}

	return s.String(), binary.LittleEndian.Uint16(b[72:74]), nil
}
