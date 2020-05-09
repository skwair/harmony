package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"sync"

	"github.com/skwair/harmony/voice"
	"github.com/skwair/harmony/voice/voiceutil"
)

// Possible states a player can be in.
const (
	stateIdle = iota
	statePlaying
	stateDestroyed
)

// Player wraps a voice connection and allows to play audio tracks on it.
// Only one player may be bound to a voice connection at a time.
type Player struct {
	mu sync.Mutex

	state int

	conn  *voice.Connection
	pcmIn chan []int16
	cmd   *exec.Cmd
	// This channel is created when the player enters the playing
	// state and gets closed when it leaves this state. Used when
	// destroying the player to know when it can close the pcmIn
	// chan without causing a panic due do a send on a closed chan.
	playDone chan struct{}
}

// NewPlayer returns a new player for the given voice connection.
// A Player should be destroyed when not needed anymore by calling its
// Destroy method.
func NewPlayer(conn *voice.Connection) (*Player, error) {
	p := &Player{
		conn:  conn,
		state: stateIdle,
	}

	var err error
	p.pcmIn, err = voiceutil.OpusEncoder(conn)
	if err != nil {
		return nil, fmt.Errorf("could not create the OpusEncoder: %w", err)
	}

	return p, nil
}

// Play starts playing the given track on the Player.
// It is an error to start playing a track if one is already playing.
func (p *Player) Play(track string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == stateDestroyed {
		return errors.New("player has been destroyed")
	}

	if p.state == statePlaying {
		return errors.New("already playing a track")
	}

	p.state = statePlaying
	p.playDone = make(chan struct{})

	p.cmd = exec.Command(
		"ffmpeg",
		"-i", track,
		"-f", "s16le",
		"-ar", strconv.Itoa(voiceutil.SampleRate),
		"-ac", strconv.Itoa(voiceutil.Channels),
		"-acodec", "pcm_s16le",
		"pipe:1",
	)

	ffmpegStdOut, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not get ffmpeg standard output: %w", err)
	}

	if err = p.cmd.Start(); err != nil {
		return fmt.Errorf("could not start ffmpeg process: %w", err)
	}

	if err = p.conn.SetSpeakingMode(voice.SpeakingModeMicrophone); err != nil {
		return fmt.Errorf("could not send speaking payload: %w", err)
	}

	go func() {
		defer close(p.playDone)
		defer func() { p.state = stateIdle }()

		for {
			pcm := make([]int16, voiceutil.FrameSize*voiceutil.Channels)
			err = binary.Read(ffmpegStdOut, binary.LittleEndian, &pcm)
			if err != nil {
				// The ffmpeg process is done.
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}

				p.conn.Logger().Errorf("could not read ffmpeg output: %v", err)
				return
			}
			p.pcmIn <- pcm
		}

		if err = p.conn.SetSpeakingMode(voice.SpeakingModeOff); err != nil {
			p.conn.Logger().Errorf("could not send speaking payload: %v", err)
			return
		}

	}()

	return nil
}

// Stop stops the currently played track. It's a no-op if nothing is being played.
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stop()
}

func (p *Player) stop() {
	if p.state != statePlaying {
		return
	}

	p.cmd.Process.Kill()

	p.state = stateIdle

	if err := p.conn.SetSpeakingMode(voice.SpeakingModeOff); err != nil {
		log.Printf("could not send speaking payload: %v", err)
		return
	}
}

// Destroy destroys the Player, freeing any allocated resources.
// No-op
func (p *Player) Destroy() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == stateDestroyed {
		return
	}

	p.stop()

	<-p.playDone
	close(p.pcmIn)

	p.state = stateDestroyed
}
