package gatcha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"math/rand"
	"time"
)

// 結果を送信します
func SendResult(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	time.Sleep(2 * time.Second)

	if isWinner() {
		return sendWinnerMessage(s, i)
	} else {
		return sendLoserMessage(s, i)
	}
}

// 当たりの場合のメッセージを送信します
func sendWinnerMessage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// ロールを付与します
	if err := addWinnerRole(s, i); err != nil {
		return errors.NewError("ロールを付与できません", err)
	}

	description := `
🎉🎉🎉🎉🎉🎉🎉
 「当たり」
🎉🎉🎉🎉🎉🎉🎉

おめでとうございます！
ロールを獲得しました！！
`
	embed := &discordgo.MessageEmbed{
		Description: description,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/1067807967950422096/1115604196473966633/2023-06-06_20.31.42.png",
		},
		Color: internal.ColorBlue,
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

// ハズレの場合のメッセージを送信します
func sendLoserMessage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	description := `

「ハズレ」

また明日チャレンジしてみてね！
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

// ロールを付与します
func addWinnerRole(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var hasRoleNum int
	for _, role := range i.Member.Roles {
		if role == internal.RoleID().PRIZE1 ||
			role == internal.RoleID().PRIZE2 ||
			role == internal.RoleID().PRIZE3 {
			hasRoleNum++
		}
	}

	switch hasRoleNum {
	case 0:
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE1); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	case 1:
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE2); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	case 2:
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE3); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	}

	return nil
}

// 当たり判定をします
func isWinner() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10) == 0
}
