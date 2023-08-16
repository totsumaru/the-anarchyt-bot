package gatcha

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/utils"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// カプセルのメッセージを送信します
func SendCapsule(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	editFunc, err := utils.SendInteractionWaitingMessage(s, i, false, true)
	if err != nil {
		return errors.NewError("Waitingメッセージが送信できません")
	}

	coinRoles := hasGatchaCoin(i.Member.Roles)

	if len(coinRoles) == 0 {
		if err = sendHasNotTicketErr(i, editFunc); err != nil {
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
			// 1. 最初の画像
			//URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1116312679884259338/open.png",
			// 2. オレンジ&緑
			//URL: "https://cdn.discordapp.com/attachments/1067807967950422096/1141186462184902806/open_02.jpg",
			// 3. 黄色&緑
			//URL: "https://cdn.discordapp.com/attachments/1067807967950422096/1141186462537232504/open_03.jpg",
			// 4. 青
			//URL: "https://cdn.discordapp.com/attachments/1067807967950422096/1141186462965055519/open_04.jpg",
			// 5. GIFモノクロ
			URL: "https://cdn.discordapp.com/attachments/1067807967950422096/1141186523010711573/Open_05.gif",
		},
		Color: internal.ColorBlue,
	}

	// チケットロールを削除
	//
	// ガチャコインの削除が優先
	{
		var removeRoleID string

		for _, coinRole := range coinRoles {
			switch coinRole {
			case internal.RoleID().GATCHA_COIN:
				removeRoleID = coinRole
				break
			case internal.RoleID().BONUS_COIN:
				removeRoleID = coinRole
			}
		}

		if removeRoleID != "" {
			if err = s.GuildMemberRoleRemove(
				i.GuildID,
				i.Member.User.ID,
				removeRoleID,
			); err != nil {
				return errors.NewError("コインロールを削除できません", err)
			}
		}
	}

	webhook := &discordgo.WebhookEdit{
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &[]discordgo.MessageComponent{actions},
	}
	if _, err = editFunc(i.Interaction, webhook); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
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
func sendHasNotTicketErr(
	i *discordgo.InteractionCreate,
	editFunc utils.EditFunc,
) error {
	description := `
<@&%s>をもっていません。

毎日1枚もらえるので、また明日チャレンジしてみてね！
`
	embed := &discordgo.MessageEmbed{
		Title:       "ERROR",
		Description: fmt.Sprintf(description, internal.RoleID().GATCHA_COIN),
		Color:       internal.ColorRed,
	}

	webhook := &discordgo.WebhookEdit{
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &[]discordgo.MessageComponent{},
	}
	if _, err := editFunc(i.Interaction, webhook); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}
