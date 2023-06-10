package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 通知ボタンが押された時の処理です
func SendNoticeMessage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if hasGatchaRole(i.Member) {
		// 通知ロールを持っている場合
		// 確認メッセージを送信
		btn1 := discordgo.Button{
			Label:    "通知ロールを削除する",
			Style:    discordgo.PrimaryButton,
			CustomID: internal.Interaction_CustomID_gatcha_Notice_Remove,
		}

		actions := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{btn1},
		}

		description := `
現在通知ロールを持っています。

通知ロールを削除しますか？
`

		embed := &discordgo.MessageEmbed{
			Description: description,
			Color:       internal.ColorBlue,
		}

		resp := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Components: []discordgo.MessageComponent{actions},
				Embeds:     []*discordgo.MessageEmbed{embed},
				Flags:      discordgo.MessageFlagsEphemeral,
			},
		}

		if err := s.InteractionRespond(i.Interaction, resp); err != nil {
			return errors.NewError("レスポンスを送信できません", err)
		}
	} else {
		// 通知ロールを持っていない場合
		// 通知ロールを付与
		if err := s.GuildMemberRoleAdd(
			i.GuildID,
			i.Member.User.ID,
			internal.RoleID().GATCHA_NOTICE,
		); err != nil {
			return errors.NewError("ガチャの通知ロールを付与できません", err)
		}

		// 確認メッセージを送信
		description := `
<@&%s> ロールを付与しました！

不要な場合は、もう一度ボタンを押せばロールを削除できます。
`

		embed := &discordgo.MessageEmbed{
			Description: description,
			Color:       internal.ColorBlue,
		}

		resp := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		}

		if err := s.InteractionRespond(i.Interaction, resp); err != nil {
			return errors.NewError("レスポンスを送信できません", err)
		}
	}

	return nil
}

// ガチャ通知ロールを持っているかを確認します
func hasGatchaRole(member *discordgo.Member) bool {
	for _, roleID := range member.Roles {
		if roleID == internal.RoleID().GATCHA_NOTICE {
			return true
		}
	}

	return false
}
