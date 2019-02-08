package harmony_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/skwair/harmony"
	"github.com/skwair/harmony/channel"
	"github.com/skwair/harmony/invite"
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

	client, err := harmony.NewClient(token)
	if err != nil {
		t.Fatalf("could not create harmony client: %v", err)
	}

	if err = client.Connect(context.TODO()); err != nil {
		t.Fatalf("could not connect to gateway: %v", err)
	}
	defer client.Disconnect()

	// Purge existing channels.
	chs, err := client.Guild(guildID).Channels(context.TODO())
	if err != nil {
		t.Fatalf("could not get guild channels: %v", err)
	}
	for _, ch := range chs {
		if _, err = client.Channel(ch.ID).Delete(context.TODO()); err != nil {
			t.Fatalf("could not delete channel %q: %v", ch.Name, err)
		}
	}

	var txtCh *harmony.Channel

	t.Run("create channels", func(t *testing.T) {
		// Create a new channel category.
		settings := channel.NewSettings(
			channel.WithName("test-category"),
			channel.WithType(channel.TypeGuildCategory),
		)
		cat, err := client.Guild(guildID).NewChannel(context.TODO(), settings)
		if err != nil {
			t.Fatalf("could not create channel category: %v", err)
		}

		// Create a new text channel in this category.
		settings = channel.NewSettings(
			channel.WithName("test-text-channel"),
			channel.WithType(channel.TypeGuildText),
			channel.WithParent(cat.ID), // Set this channel as a child of the new category.
		)
		txtCh, err = client.Guild(guildID).NewChannel(context.TODO(), settings)
		if err != nil {
			t.Fatalf("could not create text channel: %v", err)
		}
	})

	t.Run("create a channel invite", func(t *testing.T) {
		settings := invite.NewSettings(
			invite.WithMaxUses(1),
		)

		i, err := client.Channel(txtCh.ID).NewInvite(context.TODO(), settings)
		if err != nil {
			t.Fatalf("could not create new invitation: %v", err)
		}

		if i.MaxUses != 1 {
			t.Fatalf("expected to new invite to have %d max uses; got %d", 1, i.MaxUses)
		}
	})

	var (
		firstMsgIDs []string
		lastMsgID   string
	)

	t.Run("send messages", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			content := fmt.Sprintf("foobar %d", i)
			msg, err := client.Channel(txtCh.ID).SendMessage(context.TODO(), content)
			if err != nil {
				t.Fatalf("could not send message (%d): %v", i, err)
			}

			if i == 4 {
				lastMsgID = msg.ID
			} else {
				firstMsgIDs = append(firstMsgIDs, msg.ID)
			}
		}
	})

	t.Run("get messages", func(t *testing.T) {
		msgs, err := client.Channel(txtCh.ID).Messages(context.TODO(), "<"+lastMsgID, 0)
		if err != nil {
			t.Fatalf("could not retrieve text channel messages: %v", err)
		}

		if len(msgs) != 4 {
			t.Fatalf("expected to retrieve %d messages; got %d", 4, len(msgs))
		}
	})

	t.Run("edit message", func(t *testing.T) {
		if _, err = client.Channel(txtCh.ID).EditMessage(context.TODO(), lastMsgID, "foobar edited"); err != nil {
			t.Fatalf("could not edit message: %v", err)
		}
	})

	t.Run("get single message", func(t *testing.T) {
		msg, err := client.Channel(txtCh.ID).Message(context.TODO(), lastMsgID)
		if err != nil {
			t.Fatalf("coult not get single message: %v", err)
		}

		if msg.Content != "foobar edited" {
			t.Fatalf("expected message content to be %q; got %q", "foobar edited", msg.Content)
		}
	})

	t.Run("add reactions", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).AddReaction(context.TODO(), lastMsgID, "👍"); err != nil {
			t.Fatalf("could not add reaction to last message: %v", err)
		}

		if err = client.Channel(txtCh.ID).AddReaction(context.TODO(), lastMsgID, "👎"); err != nil {
			t.Fatalf("could not add reaction to last message: %v", err)
		}
	})

	t.Run("remove reaction", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).RemoveReaction(context.TODO(), lastMsgID, "👎"); err != nil {
			t.Fatalf("could not remove reaction to last message: %v", err)
		}
	})

	currentUserID := client.State.CurrentUser().ID

	t.Run("get reactions", func(t *testing.T) {
		users, err := client.Channel(txtCh.ID).GetReactions(context.TODO(), lastMsgID, "👍", 0, "", "")
		if err != nil {
			t.Fatalf("could not get reactions to last message: %v", err)
		}

		if len(users) != 1 {
			t.Fatalf("expected to have %d user with this reaction; got %d", 1, len(users))
		}

		if users[0].ID != currentUserID {
			t.Fatalf("expected the ID of the user to be %s; got %s", currentUserID, users[0].ID)
		}

		users, err = client.Channel(txtCh.ID).GetReactions(context.TODO(), lastMsgID, "👎", 0, "", "")
		if err != nil {
			t.Fatalf("could not get reactions to last message: %v", err)
		}

		if len(users) != 0 {
			t.Fatalf("expected to have %d user with this reaction; got %d", 0, len(users))
		}
	})

	t.Run("remove all reactions", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).RemoveAllReactions(context.TODO(), lastMsgID); err != nil {
			t.Fatalf("could not remove all reactions to last message: %v", err)
		}
	})

	t.Run("pin message", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).PinMessage(context.TODO(), lastMsgID); err != nil {
			t.Fatalf("could not pin last message: %v", err)
		}
	})

	t.Run("get pins", func(t *testing.T) {
		pins, err := client.Channel(txtCh.ID).Pins(context.TODO())
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
		if err = client.Channel(txtCh.ID).UnpinMessage(context.TODO(), lastMsgID); err != nil {
			t.Fatalf("could not unpin last message: %v", err)
		}
	})

	t.Run("delete single message", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).DeleteMessage(context.TODO(), lastMsgID); err != nil {
			t.Fatalf("could not delete single message: %v", err)
		}
	})

	t.Run("delete messages", func(t *testing.T) {
		if err = client.Channel(txtCh.ID).DeleteMessageBulk(context.TODO(), firstMsgIDs); err != nil {
			t.Fatalf("could not delete messages: %v", err)
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
		testRole, err = client.Guild(guildID).NewRole(context.TODO(), settings)
		if err != nil {
			t.Fatalf("could not create new role: %v", err)
		}
	})

	t.Run("add role", func(t *testing.T) {
		if err = client.Guild(guildID).AddMemberRole(context.TODO(), currentUserID, testRole.ID); err != nil {
			t.Fatalf("could not add new role to user: %v", err)
		}
	})

	t.Run("get guild member#01", func(t *testing.T) {
		member, err := client.Guild(guildID).Member(context.TODO(), currentUserID)
		if err != nil {
			t.Fatalf("could not get guild member: %v", err)
		}

		if !member.HasRole(testRole.ID) {
			t.Fatal("guild member should have test role")
		}
	})

	t.Run("remove role", func(t *testing.T) {
		if err = client.Guild(guildID).RemoveMemberRole(context.TODO(), currentUserID, testRole.ID); err != nil {
			t.Fatalf("could not remove role from user: %v", err)
		}
	})

	t.Run("get guild member#02", func(t *testing.T) {
		member, err := client.Guild(guildID).Member(context.TODO(), currentUserID)
		if err != nil {
			t.Fatalf("could not get guild member: %v", err)
		}

		if member.HasRole(testRole.ID) {
			t.Fatal("guild member should not have test role anymore")
		}
	})

	t.Run("delete role", func(t *testing.T) {
		if err = client.Guild(guildID).DeleteRole(context.TODO(), testRole.ID); err != nil {
			t.Fatalf("could not delete test role: %v", err)
		}
	})
}
