package harmony

const (
	gatewayOpcodeDispatch = iota
	gatewayOpcodeHeartbeat
	gatewayOpcodeIdentify
	gatewayOpcodeStatusUpdate
	gatewayOpcodeVoiceStateUpdate
	gatewayOpcodeResume
	gatewayOpcodeReconnect
	gatewayOpcodeRequestGuildMembers
	gatewayOpcodeInvalidSession
	gatewayOpcodeHello
	gatewayOpcodeHeartbeatACK
)

//     harmony_test.go:31: could not connect to gateway: expected Opcode 0 Ready; got Opcode 9

const (
	voiceOpcodeIdentify = iota
	voiceOpcodeSelectProtocol
	voiceOpcodeReady
	voiceOpcodeHeartbeat
	voiceOpcodeSessionDescription
	voiceOpcodeSpeaking
	voiceOpcodeHeartbeatACK
	voiceOpcodeResume
	voiceOpcodeHello
	voiceOpcodeResumed
	voiceOpcodeClientDisconnect
)

// TODO: handle voice connection resumes.
var _ = voiceOpcodeResume // Make use of this constant so the CI doesn't complain.
