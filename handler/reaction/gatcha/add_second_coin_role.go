package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"strings"
)

const (
	emojiCheckGreen = "✅"
)

// はずれた人に、コインロールを付与します
//
// #はずれ町瓦版 に`#アナーキー`の入ったツイートURLを添付すると、@コイン,@2枚目付与済み ロールが付与されます。
//
// @2枚目付与済みを持っている人には付与されません。
func AddSubCoinRoleForHazureUser(s *discordgo.Session, r *discordgo.MessageReactionAdd) error {
	if r.ChannelID != internal.ChannelID().HAZURE_TWEET {
		return nil
	}

	// 付与済みの人は除外する
	for _, role := range r.Member.Roles {
		if role == internal.RoleID().COIN_2_ADDED {
			return nil
		}
	}

	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return errors.NewError("メッセージを取得できません", err)
	}

	if !(strings.Contains(message.Content, "https://twitter.com") ||
		strings.Contains(message.Content, "https://mobile.twitter.com")) {

		return nil
	}

	if r.Emoji.Name == emojiCheckGreen {
		if r.UserID == internal.UserID().MUG ||
			r.UserID == internal.UserID().TOTSUMARU {

			// ガチャコインロールを付与
			if err = s.GuildMemberRoleAdd(r.GuildID, message.Author.ID, internal.RoleID().GATCHA_TICKET); err != nil {
				return errors.NewError("ガチャコインロールを付与できません", err)
			}

			// 付与済みロールを付与
			if err = s.GuildMemberRoleAdd(r.GuildID, message.Author.ID, internal.RoleID().COIN_2_ADDED); err != nil {
				return errors.NewError("ガチャコインロールを付与できません", err)
			}
		}
	}

	return nil
}
