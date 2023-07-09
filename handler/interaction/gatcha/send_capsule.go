package gatcha

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// カプセルのメッセージを送信します
func SendCapsule(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	coinRoles := hasGatchaCoin(i.Member.Roles)

	if len(coinRoles) == 0 {
		if err := sendHasNotTicketErr(s, i); err != nil {
			return errors.NewError("チケット未保持エラーを送信できません", err)
		}
		return nil
	}

	btn1 := discordgo.Button{
		Label:    "カプセルを開ける",
		Style:    discordgo.PrimaryButton,
		CustomID: internal.Interaction_CustomID_gatcha_Open,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
	}

	description := `
ガラガラ...

ゴトン。
`

	embed := &discordgo.MessageEmbed{
		Description: description,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1116312679884259338/open.png",
		},
		Color: internal.ColorBlue,
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

	// チケットロールを削除
	for _, coinRole := range coinRoles {
		if coinRole != internal.RoleID().FOR_TEST_ATARI {
			if err := s.GuildMemberRoleRemove(
				i.GuildID,
				i.Member.User.ID,
				coinRole,
			); err != nil {
				return errors.NewError("チケットロールを削除できません", err)
			}

			// 1枚のみ削除するため、1枚を削除したら処理を終了します。
			return nil
		}
	}

	return nil
}

// ガチャコインを保持しているか確認します
//
// 持っているコインロールのIDを返します。
// ボーナスコインがあるので、スライスで返します。
func hasGatchaCoin(roles []string) []string {
	res := make([]string, 0)

	for _, role := range roles {
		switch role {
		case internal.RoleID().GATCHA_COIN,
			internal.RoleID().FOR_TEST_ATARI,
			internal.RoleID().BONUS_COIN:

			res = append(res, role)
		}
	}

	return res
}

// チケット未保持エラーを送信します
func sendHasNotTicketErr(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	description := `
<@&%s>をもっていません。

毎日1枚もらえるので、また明日チャレンジしてみてね！
`
	embed := &discordgo.MessageEmbed{
		Title:       "ERROR",
		Description: fmt.Sprintf(description, internal.RoleID().GATCHA_COIN),
		Color:       internal.ColorRed,
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

	return nil
}
