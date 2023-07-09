package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"strings"
)

// TwitterのURLを添付するともう1枚コインロールを付与します
//
// #はずれ町瓦版 に`#アナーキー`の入ったツイートURLを添付すると、@コイン,@2枚目付与済み ロールが付与されます。
//
// @2枚目付与済みを持っている人には付与されません。
func AddSecondCoinRoleForHazureUser(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.ChannelID != internal.ChannelID().HAZURE_TWEET {
		return nil
	}

	// 付与済みの人は除外する
	for _, role := range m.Member.Roles {
		if role == internal.RoleID().COIN_2_ADDED {
			return nil
		}
	}

	if !(strings.Contains(m.Content, "https://twitter.com") ||
		strings.Contains(m.Content, "https://mobile.twitter.com")) {

		return nil
	}

	// ガチャコインロールを付与
	if err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, internal.RoleID().GATCHA_COIN); err != nil {
		return errors.NewError("ガチャコインロールを付与できません", err)
	}

	// 付与済みロールを付与
	if err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, internal.RoleID().COIN_2_ADDED); err != nil {
		return errors.NewError("ガチャコインロールを付与できません", err)
	}

	// チェックの絵文字をつける
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅"); err != nil {
		return errors.NewError("リアクションを付与できません", err)
	}

	return nil
}
