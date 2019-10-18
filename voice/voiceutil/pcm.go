package voiceutil

import (
	opus "layeh.com/gopus"

	"github.com/skwair/harmony/voice"
)

const (
	sampleRate = 48000
	channels   = 2
	frameSize  = 960
)

// OpusEncoder is an adapter that allows to send a PCM signal through the
// returned channel, it will encode it with Opus and send it through the given
// Discord voice connection.
//
// The returned channel must be closed when the adapter is not needed
// anymore in order to free allocated resources.
//
// Only one OpusEncoder is meant to be used at once on the same voice connection.
func OpusEncoder(conn *voice.Connection) (pcmIn chan []int16, err error) {
	enc, err := opus.NewEncoder(sampleRate, channels, opus.Audio)
	if err != nil {
		return nil, err
	}
	pcmIn = make(chan []int16, 512)

	go func() {
		for {
			pcm, ok := <-pcmIn
			if !ok { // Input chan has been closed.
				return
			}

			opusEncoded, err := enc.Encode(pcm, frameSize, frameSize*2*2)
			if err != nil {
				conn.Logger().Errorf("could not encode PCM data: %v", err)
				return
			}
			conn.Send <- opusEncoded
		}
	}()

	return pcmIn, nil
}

// OpusDecoder is an adapter that allows to read the incoming voice data
// on conn as PCM, sent through the returned channel.
//
// Disconnecting the VoiceConnection will close the decoder and free any
// allocated resources. If the adapter is not needed anymore but the
// VoiceConnection is, the free returned channel can be closed to free
// those resources.
//
// Only one OpusDecoder is meant to be used at once on the same voice connection.
func OpusDecoder(conn *voice.Connection) (pcmOut chan []int16, free chan struct{}) {
	speakers := make(map[uint32]*opus.Decoder)
	pcmOut = make(chan []int16)

	go func() {
		var ok bool
		var packet *voice.AudioPacket
		for {
			select {
			case packet, ok = <-conn.Recv:
				if !ok { // The underlying voice connection has been closed.
					return
				}

			case _, ok := <-free:
				if !ok { // The free chan was closed.
					return
				}
			}

			if _, ok = speakers[packet.SSRC]; !ok {
				var err error
				speakers[packet.SSRC], err = opus.NewDecoder(sampleRate, channels)
				if err != nil {
					conn.Logger().Errorf("could not create Opus decoder: %v", err)
					return
				}
			}

			pcm, err := speakers[packet.SSRC].Decode(packet.Opus, frameSize, false)
			if err != nil {
				conn.Logger().Errorf("could not decode Opus data: %v", err)
				return
			}
			pcmOut <- pcm
		}
	}()

	return pcmOut, free
}
