package voice

import (
	"context"
)

// SpeakingFlag is the type for flags that can be used as a bitwise mask for SetSpeaking.
type SpeakingFlag uint32

const (
	// Normal transmission of audio.
	SpeakingFlagVoice SpeakingFlag = 0x1
	// Transmission of context audio for video, no speaking indicator.
	SpeakingFlagSoundshare SpeakingFlag = 0x2
	// Priority speaker, lowering audio of other speakers.
	SpeakingFlagPriority SpeakingFlag = 0x4
	// No audio transmission.
	SpeakingFlagOff SpeakingFlag = 0x0
)

// Speaking sends an Opcode 5 Speaking payload. This does nothing
// if the user is already in the given state.
func (vc *Connection) SetSpeaking(ctx context.Context, speaking SpeakingFlag) error {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	// Return early if the user is already in the asked state.
	if speaking == vc.speaking {
		return nil
	}

	p := struct {
		Speaking uint32 `json:"speaking"`
		Delay    int    `json:"delay"`
		SSRC     uint32 `json:"ssrc"`
	}{
		Speaking: uint32(speaking),
		Delay:    0,
		SSRC:     vc.ssrc,
	}

	if err := vc.sendPayload(ctx, voiceOpcodeSpeaking, p); err != nil {
		return err
	}

	vc.speaking = speaking

	return nil
}
