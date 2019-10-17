package voice

import "sync/atomic"

// Speaking sends an Opcode 5 Speaking payload. This does nothing
// if the user is already in the given state.
func (vc *Connection) Speaking(s bool) error {
	// Return early if the user is already in the asked state.
	prev := atomic.LoadInt32(&vc.speaking)
	if (prev == 1) == s {
		return nil
	}

	if s {
		atomic.StoreInt32(&vc.speaking, 1)
	} else {
		atomic.StoreInt32(&vc.speaking, 0)
	}

	p := struct {
		Speaking bool   `json:"speaking"`
		Delay    int    `json:"delay"`
		SSRC     uint32 `json:"ssrc"`
	}{
		Speaking: s,
		Delay:    0,
		SSRC:     vc.ssrc,
	}

	if err := vc.sendPayload(voiceOpcodeSpeaking, p); err != nil {
		// If there is an error, reset our internal value to its previous
		// state because the update was not acknowledged by Discord.
		atomic.StoreInt32(&vc.speaking, prev)
		return err
	}

	return nil
}
