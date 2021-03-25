package harmony

import (
	"github.com/skwair/harmony/discord"
)

// SetBotStatus sets the bot's status. You need to be connected to the
// Gateway to call this method, else it will return ErrGatewayNotConnected.
func (c *Client) SetBotStatus(status *discord.BotStatus) error {
	if !c.isConnected() {
		return discord.ErrGatewayNotConnected
	}

	return c.sendPayload(c.ctx, gatewayOpcodeStatusUpdate, status)
}
