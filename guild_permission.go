package harmony

import "github.com/skwair/harmony/permission"

// computeBasePermissions returns the base permissions a member has in a given guild.
func computeBasePermissions(g *Guild, m *GuildMember) (permissions int) {
	if g.OwnerID == m.User.ID {
		return permission.All
	}

	// Role '@everyone' has the same ID as the guild ID
	// and is always present, so no need to check for nil.
	roleEveryone := roleByID(g.Roles, g.ID)
	permissions = roleEveryone.Permissions

	for _, id := range m.Roles {
		role := roleByID(g.Roles, id)
		if role != nil {
			permissions |= role.Permissions
		}
	}

	if permission.Contains(permissions, permission.Administrator) {
		return permission.All
	}
	return permissions
}

func computeOverwrites(ch *Channel, m *GuildMember, basePermissions int) (permissions int) {
	// Administrator can not be overridden.
	if permission.Contains(basePermissions, permission.Administrator) {
		return permission.All
	}

	permissions = basePermissions

	po := overwriteByID(ch.PermissionOverwrites, ch.GuildID)
	if po != nil {
		permissions &= ^po.Deny
		permissions |= po.Allow
	}

	pos := ch.PermissionOverwrites
	allow := permission.None
	deny := permission.None
	for _, id := range m.Roles {
		por := overwriteByID(pos, id)
		if por != nil {
			allow |= por.Allow
			deny |= por.Deny
		}
	}
	permissions &= ^deny
	permissions |= allow

	pom := overwriteByID(ch.PermissionOverwrites, m.User.ID)
	if pom != nil {
		permissions &= ^pom.Deny
		permissions |= pom.Allow
	}

	return permissions
}

func roleByID(roles []Role, id string) *Role {
	for i := 0; i < len(roles); i++ {
		if roles[i].ID == id {
			return &roles[i]
		}
	}
	return nil
}

func overwriteByID(po []permission.Overwrite, id string) *permission.Overwrite {
	for i := 0; i < len(po); i++ {
		if po[i].ID == id {
			return &po[i]
		}
	}
	return nil
}
