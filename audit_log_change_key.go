package harmony

type changeKey string

const (
	changeKeyName                       changeKey = "name"
	changeKeyIconHash                   changeKey = "icon_hash"
	changeKeySplashHash                 changeKey = "splash_hash"
	changeKeyOwnerID                    changeKey = "owner_id"
	changeKeyRegion                     changeKey = "region"
	changeKeyAFKChannelID               changeKey = "afk_channel_id"
	changeKeyAFKTimeout                 changeKey = "afk_timeout"
	changeKeyMFALevel                   changeKey = "mfa_level"
	changeKeyVerificationLevel          changeKey = "verification_level"
	changeKeyExplicitContentFilter      changeKey = "explicit_content_filter"
	changeKeyDefaultMessageNotification changeKey = "default_message_notifications"
	changeKeyVanityURLCode              changeKey = "vanity_url_code"
	changeKeyAddRole                    changeKey = "$add"
	changeKeyRemoveRole                 changeKey = "$remove"
	changeKeyPruneDeleteDays            changeKey = "prune_delete_days"
	changeKeyWidgetEnabled              changeKey = "widget_enabled"
	changeKeyWidgetChannelID            changeKey = "widget_channel_id"

	changeKeyPosition             changeKey = "position"
	changeKeyTopic                changeKey = "topic"
	changeKeyBitrate              changeKey = "bitrate"
	changeKeyRateLimitPerUser     changeKey = "rate_limit_per_user"
	changeKeyPermissionOverwrites changeKey = "permission_overwrites"
	changeKeyNFSW                 changeKey = "nsfw"
	changeKeyApplicationID        changeKey = "application_id"

	changeKeyPermissions changeKey = "permissions"
	changeKeyColor       changeKey = "color"
	changeKeyHoist       changeKey = "hoist"
	changeKeyMentionable changeKey = "mentionable"
	changeKeyAllow       changeKey = "allow"
	changeKeyDeny        changeKey = "deny"

	changeKeyCode      changeKey = "code"
	changeKeyChannelID changeKey = "channel_id"
	changeKeyInviterID changeKey = "inviter_id"
	changeKeyMaxUses   changeKey = "max_uses"
	changeKeyUses      changeKey = "uses"
	changeKeyMaxAge    changeKey = "max_age"
	changeKeyTemporary changeKey = "temporary"

	changeKeyDeaf       changeKey = "deaf"
	changeKeyMute       changeKey = "mute"
	changeKeyNick       changeKey = "nick"
	changeKeyAvatarHash changeKey = "avatar_hash"

	changeKeyID   changeKey = "id"
	changeKeyType changeKey = "type"
)
