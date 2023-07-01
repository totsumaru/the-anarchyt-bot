package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"strings"
)

const (
	emojiRoleAdd = "✅"
)

// はずれた人に、コインロールを付与します
//
// #はずれ町一丁目に`#アナーキー`の入ったツイートURLを添付すると、コインロールが付与される
func AddSubCoinRoleForHazureUser(s *discordgo.Session, r *discordgo.MessageReactionAdd) error {
	if r.ChannelID != internal.ChannelID().HAZURE_MACHI_1 {
		return nil
	}

	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return errors.NewError("メッセージを取得できません", err)
	}

	if !(strings.Contains(message.Content, "https://twitter.com") ||
		strings.Contains(message.Content, "https://mobile.twitter.com")) {

		return nil
	}

	if r.Emoji.Name == emojiRoleAdd {
		if r.UserID == internal.UserID().MUG ||
			r.UserID == internal.UserID().TOTSUMARU {

			// ガチャコインロールを付与
			if err = s.GuildMemberRoleAdd(r.GuildID, message.Author.ID, internal.RoleID().GATCHA_TICKET); err != nil {
				return errors.NewError("ガチャコインロールを付与できません", err)
			}
		}
	}

	return nil
}
