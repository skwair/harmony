package voice

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
	"net"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

// AudioPacket is a parsed and decrypted RTP frame containing Opus encoded audio.
type AudioPacket struct {
	Type      uint8
	Version   uint8
	Sequence  uint16
	Timestamp uint32
	SSRC      uint32
	Opus      []byte
}

// rtpFrame is a raw RTP frame, along with its size.
type rtpFrame struct {
	raw  []byte
	size int
}

// opusReceiver receives and decrypt audio packets, forwarding them to vc.Recv.
// The Opus encoded audio is not decoded.
func (vc *Connection) opusReceiver() {
	vc.opusReadinessWG.Done()
	defer vc.wg.Done()

	vc.logger.Debug("starting Opus receiver")
	defer vc.logger.Debug("stopped Opus receiver")

	rtpFrames := make(chan *rtpFrame)
	go vc.readUDP(rtpFrames)

	var nonce [24]byte
	for {
		select {
		case frame := <-rtpFrames:
			p := &AudioPacket{
				Type:      frame.raw[0],
				Version:   frame.raw[1],
				Sequence:  binary.BigEndian.Uint16(frame.raw[2:4]),
				Timestamp: binary.BigEndian.Uint32(frame.raw[4:8]),
				SSRC:      binary.BigEndian.Uint32(frame.raw[8:12]),
			}
			copy(nonce[:], frame.raw[0:12])
			p.Opus, _ = secretbox.Open(nil, frame.raw[12:frame.size], &nonce, &vc.secret)

			// If the RTP extension bit is set, we must remove the
			// first 8 bytes corresponding to the header extension,
			// else the opus signal will be invalid.
			if len(p.Opus) > 8 && frame.raw[0] == 0x90 {
				p.Opus = p.Opus[8:]
			}

			// Drop the packet if no one is receiving
			// on the other end of the channel.
			select {
			case vc.Recv <- p:
			default:
			}
		case <-vc.stop:
			return
		}
	}
}

// readUDP reads RTP frames from the voice connection's UDP socket
// and sends them through the given channel as they are read.
func (vc *Connection) readUDP(ch chan<- *rtpFrame) {
	defer vc.wg.Done()

	for {
		// NOTE: creating a new buffer for each packet to prevent
		// data races with the opusReceiver goroutine. Maybe there
		// is a better way to do this.
		buf := make([]byte, 1024)
		l, err := vc.udpConn.Read(buf)
		if err != nil {
			// Silently break out of this loop if
			// the connection was closed by the client.
			if isConnectionClosed(err) {
				break
			}
			// NOTE: this might be a bit too extreme ?
			vc.reportErr(err)
			return
		}

		// Handle UDP heartbeat ACK.
		if l == 8 {
			vc.lastUDPHeartbeatACK.Store(time.Now().UnixNano())

			// TODO: check the sequence number in the UDP heartbeat ?
			// udpSeq := binary.LittleEndian.Uint64(buf[:l])
			// // Since the UDP heartbeat sequence number is incremented right after it
			// // is sent, we should receive the number just before.
			// if udpSeq != atomic.LoadUint64(&vc.udpHeartbeatSequence)-1 {
			// }

			continue
		}
		// Skip non audio packets.
		if l < 12 || (buf[0] != 0x80 && buf[0] != 0x90) {
			continue
		}

		// Only send voice data through the channel if someone is listening
		// on the other side, else we'll just block forever.
		select {
		case ch <- &rtpFrame{raw: buf, size: l}:
		default:
		}
	}
}

// opusSender creates, encrypts and sends Opus encoded packets sent through the voice
// connection's Send channel.
func (vc *Connection) opusSender() {
	vc.opusReadinessWG.Done()
	defer vc.wg.Done()

	vc.logger.Debug("starting Opus sender")
	defer vc.logger.Debug("stopped Opus sender")

	const (
		sampleRate = 48000 // In Hz, the number of samples we take each second.
		frameSize  = 960   // This is the number of samples we send at each interval,
		// 960 samples at 48000Hz represents 20 milliseconds of audio.
		// 960/(48000/1000) = 20
	)

	// According to the RTP RFC, the initial value of the sequence number
	// SHOULD be random (unpredictable) to make known-plaintext attacks
	// on encryption more difficult.
	r, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint16/2))
	if err != nil {
		vc.reportErr(err)
		return
	}
	seq := uint16(r.Uint64())
	// Same goes for the timestamp field.
	r, err = rand.Int(rand.Reader, big.NewInt(math.MaxUint32/2))
	if err != nil {
		vc.reportErr(err)
		return
	}
	timestamp := uint32(r.Uint64())

	var nonce [24]byte
	rtpHeader := make([]byte, 12)
	// Set the static part of the RTP header.
	rtpHeader[0] = 0x80
	rtpHeader[1] = 0x78
	binary.BigEndian.PutUint32(rtpHeader[8:], vc.ssrc)

	ticker := time.NewTicker(time.Millisecond * time.Duration(frameSize/(sampleRate/1000)))
	defer ticker.Stop()

	for {
		select {
		case data := <-vc.Send:
			// Set the dynamic part of the RTP header.
			binary.BigEndian.PutUint16(rtpHeader[2:], seq)
			binary.BigEndian.PutUint32(rtpHeader[4:], timestamp)

			// Generate the nonce from the rtpHeader. Since the RTP header is only 12 bytes
			// long, it will leave the 12 trailing bytes of the nonce null, as specified by
			// https://discordapp.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice.
			copy(nonce[:], rtpHeader)

			buf := secretbox.Seal(rtpHeader, data, &nonce, &vc.secret)

			// Send voice packets at regular interval.
			// The ticker will drop ticks if we don't
			// consume them, like for example if we have
			// nothing to send.
			<-ticker.C

			_, err = vc.udpConn.Write(buf)
			if err != nil {
				// Silently break out of this loop because
				// the connection was closed by the client.
				if isConnectionClosed(err) {
					return
				}

				vc.reportErr(err)
				return
			}

			// Increase the sequence number. Since this is an unsigned
			// int16, it will reset to 0 when reaching its max value.
			seq++

			timestamp += frameSize

		case <-vc.stop:
			return
		}
	}
}

func isConnectionClosed(err error) bool {
	if e, ok := err.(*net.OpError); ok {
		// Ugly but : https://github.com/golang/go/issues/4373
		if e.Err.Error() == "use of closed network connection" {
			return true
		}
	}
	return false
}
