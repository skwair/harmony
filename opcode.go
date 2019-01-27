package harmony

const (
	gatewayOpcodeDispatch            = 0
	gatewayOpcodeHeartbeat           = 1
	gatewayOpcodeIdentify            = 2
	gatewayOpcodeStatusUpdate        = 3
	gatewayOpcodeVoiceStateUpdate    = 4
	gatewayOpcodeResume              = 6
	gatewayOpcodeReconnect           = 7
	gatewayOpcodeRequestGuildMembers = 8
	gatewayOpcodeInvalidSession      = 9
	gatewayOpcodeHello               = 10
	gatewayOpcodeHeartbeatACK        = 11
)

const (
	voiceOpcodeIdentify           = 0
	voiceOpcodeSelectProtocol     = 1
	voiceOpcodeReady              = 2
	voiceOpcodeHeartbeat          = 3
	voiceOpcodeSessionDescription = 4
	voiceOpcodeSpeaking           = 5
	voiceOpcodeHeartbeatACK       = 6
	voiceOpcodeResume             = 7
	voiceOpcodeHello              = 8
	voiceOpcodeResumed            = 9
	voiceOpcodeClientDisconnect   = 13
)

// TODO: handle voice connection resumes.
var _ = voiceOpcodeResume // Make use of this constant so the CI doesn't complain.
