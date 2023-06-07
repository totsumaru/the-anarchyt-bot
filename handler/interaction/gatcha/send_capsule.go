package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"time"
)

// カプセルのメッセージを送信します
func SendCapsule(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if !hasTicketRole(i.Member.Roles) {
		if err := sendHasNotTicketErr(s, i); err != nil {
			return errors.NewError("チケット未保持エラーを送信できません", err)
		}
		return nil
	}

	time.Sleep(1 * time.Second)

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
			URL: "https://cdn.discordapp.com/attachments/1067807967950422096/1115604195991629874/2023-06-06_20.28.23.png",
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
	if err := s.GuildMemberRoleRemove(
		i.GuildID,
		i.Member.User.ID,
		internal.RoleID().TICKET,
	); err != nil {
		return errors.NewError("チケットロールを削除できません", err)
	}

	return nil
}

// チケットロールを保持しているか確認します
func hasTicketRole(roles []string) bool {
	for _, role := range roles {
		if role == internal.RoleID().TICKET {
			return true
		}
	}

	return false
}

// チケット未保持エラーを送信します
func sendHasNotTicketErr(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	description := `
チケットロールがありません。

毎日1枚もらえるので、また明日チャレンジしてみてね！
`
	embed := &discordgo.MessageEmbed{
		Title:       "ERROR",
		Description: description,
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
