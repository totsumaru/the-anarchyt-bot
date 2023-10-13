package invitation

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// 招待リンクを発行し、送信します
func ReplyLink(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// 招待券を確認します
	roleID, err := hasInvitationTicketRole(s, i.Member)
	if err != nil {
		return errors.NewError("招待券の保持を確認できません", err)
	}
	if roleID == "" {
		if err = sendHasNotTicketErr(s, i); err != nil {
			return errors.NewError("招待券未保持エラーを送信できません", err)
		}
		return nil
	}

	// 招待リンクの作成
	url, err := createURL(s)
	if err != nil {
		return errors.NewError("招待リンクを作成できません", err)
	}

	// 返信を送信
	description := `
招待リンクを発行しました。
※再発行はできませんので、すぐにコピーしておきましょう。

%s

✅️招待リンクについて
- **1人のみ** 招待可
- **7日間** 有効
`

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(description, url),
		Color:       internal.ColorBlue,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}

	if err = s.InteractionRespond(i.Interaction, resp); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	// 招待券を削除
	if err = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID); err != nil {
		return errors.NewError("招待券ロールを削除できません", err)
	}

	// ログを送信
	if err = sendLog(s, i, url); err != nil {
		return errors.NewError("ログを送信できません", err)
	}

	return nil
}

// 招待リンクを作成します
func createURL(s *discordgo.Session) (string, error) {
	iv := discordgo.Invite{
		MaxAge:    604800, // 7日間
		MaxUses:   1,      // 1回限り
		Unique:    true,
		Temporary: false,
	}

	invite, err := s.ChannelInviteCreate(internal.ChannelID().START, iv)
	if err != nil {
		return "", errors.NewError("招待リンクを作成できません", err)
	}

	url := fmt.Sprintf("https://discord.gg/%s", invite.Code)

	return url, nil
}

// 招待券を持っているかを確認します
//
// 持っている招待券のロールIDを1つ返します
func hasInvitationTicketRole(s *discordgo.Session, member *discordgo.Member) (string, error) {
	for _, roleID := range member.Roles {
		if roleID == internal.RoleID().INVITATION1 || roleID == internal.RoleID().INVITATION2 {
			return roleID, nil
		}
	}

	return "", nil
}

// 招待券未保持エラーを送信します
func sendHasNotTicketErr(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	description := `
<@&%s>をもっていません。

<#%s>で当選するか、その他イベントでもらうことができます！
ぜひ色々チャレンジしてみてね！
`
	embed := &discordgo.MessageEmbed{
		Title: "ERROR",
		Description: fmt.Sprintf(
			description,
			internal.RoleID().INVITATION1,
			internal.ChannelID().GATCHA,
		),
		Color: internal.ColorRed,
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

// ログを送信します
func sendLog(s *discordgo.Session, i *discordgo.InteractionCreate, url string) error {
	description := `
<@%s>
ユーザー名: %s
招待リンク: %s
`
	embed := &discordgo.MessageEmbed{
		Title: "招待リンクの発行（ログ）",
		Description: fmt.Sprintf(
			description,
			i.Member.User.ID,
			i.Member.User.Username,
			url,
		),
		Color: internal.ColorGreen,
	}

	_, err := s.ChannelMessageSendEmbed(internal.ChannelID().LOGS, embed)
	if err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}
