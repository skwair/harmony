package main

import (
	"context"
	"fmt"
	"os"

	"github.com/skwair/harmony"
	"github.com/skwair/harmony/resource/guild"
	"github.com/skwair/harmony/resource/guild/audit"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Environment variable BOT_TOKEN must be set.")
		return
	}

	// The guild ID you want to get the audit log from.
	// Requires the bot to have the 'VIEW_AUDIT_LOG' permission.
	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		fmt.Fprintln(os.Stderr, "Environment variable GUILD_ID must be set.")
		return
	}

	client, err := harmony.NewClient(token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	log, err := client.Guild(guildID).AuditLog(context.Background(), guild.WithAuditLimit(25))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	printRoleEntries(log)
}

func printRoleEntries(log *audit.Log) {
	// The audit log is composed of entries.
	for _, entry := range log.Entries {
		// Each entry has a type, matching the type of event this entry describes.
		switch e := entry.(type) {
		case *audit.RoleCreate:
			// Here, the entry will contain all the settings this role was created with.
			fmt.Printf("role %q was created", e.Name)

		case *audit.RoleUpdate:
			fmt.Printf("role with ID %q was updated", e.ID)

			// Fields that are of type *StringValues, *IntValues, *BoolValues
			// are settings that have potentially been modified. If they are non-nil,
			// it means they were and they will hold the old as well as the new value.
			if e.Name != nil {
				fmt.Printf("name changed from %q to %q", e.Name.Old, e.Name.New)
			}

		case *audit.RoleDelete:
			// Here, the entry will contain all the settings this role had before being deleted.
			fmt.Printf("role %q was deleted", e.Name)
		}
	}
}
