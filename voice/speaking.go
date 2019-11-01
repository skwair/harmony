package voice

import (
	"context"
)

// SpeakingMode is the type for modes that can be used as a bitwise mask for SetSpeakingMode.
type SpeakingMode uint32

const (
	// Normal transmission of audio.
	SpeakingModeVoice SpeakingMode = 0x1
	// Transmission of context audio for video, no speaking indicator.
	SpeakingModeSoundshare SpeakingMode = 0x2
	// Priority speaker, lowering audio of other speakers.
	SpeakingModePriority SpeakingMode = 0x4
	// No audio transmission.
	SpeakingModeOff SpeakingMode = 0x0
)

// SetSpeakingMode sends an Opcode 5 Speaking payload. This does nothing
// if the user is already in the given state.
func (vc *Connection) SetSpeakingMode(ctx context.Context, mode SpeakingMode) error {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	// Return early if the user is already in the asked state.
	if mode == vc.speakingMode {
		return nil
	}

	p := struct {
		Speaking uint32 `json:"speaking"`
		Delay    int    `json:"delay"`
		SSRC     uint32 `json:"ssrc"`
	}{
		Speaking: uint32(mode),
		Delay:    0,
		SSRC:     vc.ssrc,
	}

	if err := vc.sendPayload(ctx, voiceOpcodeSpeaking, p); err != nil {
		return err
	}

	vc.speakingMode = mode

	return nil
}
