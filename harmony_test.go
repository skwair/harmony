package harmony_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/skwair/harmony"
	"github.com/skwair/harmony/channel"
	"github.com/skwair/harmony/permission"
	"github.com/skwair/harmony/role"
)

func TestHarmony(t *testing.T) {
	token := os.Getenv("HARMONY_TEST_BOT_TOKEN")
	if token == "" {
		t.Fatal("environment variable HARMONY_TEST_BOT_TOKEN not set")
	}

	guildID := os.Getenv("HARMONY_TEST_GUILD_ID")
	if guildID == "" {
		t.Fatal("environment variable HARMONY_TEST_GUILD_ID not set")
	}

	client, err := harmony.NewClient(harmony.WithBotToken(token))
	if err != nil {
		t.Fatalf("could not create harmony client: %v", err)
	}

	if err = client.Connect(); err != nil {
		t.Fatalf("could not connect to gateway: %v", err)
	}
	defer client.Disconnect()

	// Purge existing channels.
	chs, err := client.Guild(guildID).Channels()
	if err != nil {
		t.Fatalf("could not get guild channels: %v", err)
	}
	for _, ch := range chs {
		if _, err = client.Channel(ch.ID).Delete(); err != nil {
			t.Fatalf("could not delete channel %q: %v", ch.Name, err)
		}
	}

	var (
		txtCh     *harmony.Channel
		lastMsgID string
	)

	currentUserID := client.State.CurrentUser().ID

	t.Run("create channels", func(t *testing.T) {
		// Create a new channel category.
		settings := channel.NewSettings(
			channel.WithName("test-category"),
			channel.WithType(channel.TypeGuildCategory),
		)
		cat, err := client.Guild(guildID).NewChannel(settings)
		if err != nil {
			t.Fatalf("could not create channel category: %v", err)
		}

		// Create a new text channel in this category.
		settings = channel.NewSettings(
			channel.WithName("test-text-channel"),
			channel.WithType(channel.TypeGuildText),
			channel.WithParent(cat.ID), // Set this channel as a child of the new category.
		)
		txtCh, err = client.Guild(guildID).NewChannel(settings)
		if err != nil {
			t.Fatalf("could not create text channel: %v", err)
		}
	})

	t.Run("send messages", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			content := fmt.Sprintf("foobar %d", i)
			msg, err := client.Channel(txtCh.ID).SendMessage(content)
			if err != nil {
				t.Fatalf("could not send message (%d): %v", i, err)
			}
			lastMsgID = msg.ID
		}
	})

	t.Run("get messages", func(t *testing.T) {
		msgs, err := client.Channel(txtCh.ID).Messages("<"+lastMsgID, 0)
		if err != nil {
			t.Fatalf("could not retrieve text channel messages: %v", err)
		}

		if len(msgs) != 4 {
			t.Fatalf("expected to retrieve %d messages; got %d", 4, len(msgs))
		}
	})

	t.Run("add reactions", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).AddReaction(lastMsgID, "👍"); err != nil {
			t.Fatalf("could not add reaction to last message: %v", err)
		}

		if err = client.Channel(txtCh.ID).AddReaction(lastMsgID, "👎"); err != nil {
			t.Fatalf("could not add reaction to last message: %v", err)
		}
	})

	t.Run("remove reactions", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).RemoveReaction(lastMsgID, "👎"); err != nil {
			t.Fatalf("could not remove reaction to last message: %v", err)
		}
	})

	t.Run("get reactions", func(t *testing.T) {
		users, err := client.Channel(txtCh.ID).GetReactions(lastMsgID, "👍", 0, "", "")
		if err != nil {
			t.Fatalf("could not get reactions to last message: %v", err)
		}

		if len(users) != 1 {
			t.Fatalf("expected to have %d user with this reaction; got %d", 0, len(users))
		}

		if users[0].ID != currentUserID {
			t.Fatalf("expected the ID of the user to be %s; got %s", currentUserID, users[0].ID)
		}

		users, err = client.Channel(txtCh.ID).GetReactions(lastMsgID, "👎", 0, "", "")
		if err != nil {
			t.Fatalf("could not get reactions to last message: %v", err)
		}

		if len(users) != 0 {
			t.Fatalf("expected to have %d user with this reaction; got %d", 0, len(users))
		}
	})

	t.Run("pin message", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).PinMessage(lastMsgID); err != nil {
			t.Fatalf("could not pin last message: %v", err)
		}
	})

	t.Run("get pins", func(t *testing.T) {
		pins, err := client.Channel(txtCh.ID).Pins()
		if err != nil {
			t.Fatalf("could not pin last message: %v", err)
		}

		if len(pins) != 1 {
			t.Fatalf("expected to have %d pins; got %d", 1, len(pins))
		}

		if pins[0].ID != lastMsgID {
			t.Fatalf("expected pinned message ID to be %s; got %s", lastMsgID, pins[0].ID)
		}
	})

	t.Run("remove pin", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).UnpinMessage(lastMsgID); err != nil {
			t.Fatalf("could not unpin last message: %v", err)
		}
	})

	var testRole *harmony.Role

	t.Run("new role", func(t *testing.T) {
		perms := permission.ReadMessageHistory | permission.SendMessages

		settings := role.NewSettings(
			role.WithName("test-role"),
			role.WithColor(0x336677),
			role.WithHoist(true),
			role.WithMentionable(true),
			role.WithPermissions(perms),
		)
		testRole, err = client.Guild(guildID).NewRole(settings)
		if err != nil {
			t.Fatalf("could not create new role: %v", err)
		}
	})

	t.Run("add role", func(t *testing.T) {
		if err = client.Guild(guildID).AddMemberRole(currentUserID, testRole.ID); err != nil {
			t.Fatalf("could not add new role to user: %v", err)
		}
	})

	t.Run("get guild member#01", func(t *testing.T) {
		member, err := client.Guild(guildID).Member(currentUserID)
		if err != nil {
			t.Fatalf("could not get guild member: %v", err)
		}

		if !member.HasRole(testRole.ID) {
			t.Fatal("guild member should have test role")
		}
	})

	t.Run("remove role", func(t *testing.T) {
		if err = client.Guild(guildID).RemoveMemberRole(currentUserID, testRole.ID); err != nil {
			t.Fatalf("could not remove role from user: %v", err)
		}
	})

	t.Run("get guild member#02", func(t *testing.T) {
		member, err := client.Guild(guildID).Member(currentUserID)
		if err != nil {
			t.Fatalf("could not get guild member: %v", err)
		}

		if member.HasRole(testRole.ID) {
			t.Fatal("guild member should not have test role anymore")
		}
	})

	t.Run("delete role", func(t *testing.T) {
		if err = client.Guild(guildID).DeleteRole(testRole.ID); err != nil {
			t.Fatalf("could not delete test role: %v", err)
		}
	})
}
