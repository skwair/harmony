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
