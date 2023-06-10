package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 通知ロールを削除します
func RemoveNoticeRole(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// 通知ロールを削除
	if err := s.GuildMemberRoleRemove(
		i.GuildID,
		i.Member.User.ID,
		internal.RoleID().GATCHA_NOTICE,
	); err != nil {
		return errors.NewError("ロールを削除できません", err)
	}

	// 確認メッセージを送信
	description := `
通知ロールを削除しました！

もう一度ボタンを押せば、再度ロールを取得できます。
`

	embed := &discordgo.MessageEmbed{
		Description: description,
		Color:       internal.ColorBlue,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}

	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}
