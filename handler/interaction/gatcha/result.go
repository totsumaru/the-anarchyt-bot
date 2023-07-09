package gatcha

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/internal"
	"math/rand"
	"time"
)

// 結果を送信します
func SendResult(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	time.Sleep(1 * time.Second)

	isWin, err := isWinner(i.Member)
	if err != nil {
		return errors.NewError("当たり判定に失敗しました", err)
	}

	if isWin {
		return sendWinnerMessage(s, i)
	} else {
		return sendLoserMessage(s, i)
	}
}

// 当たりの場合のメッセージを送信します
func sendWinnerMessage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	description := `
🎉🎉🎉🎉🎉🎉🎉
 「当たり」
🎉🎉🎉🎉🎉🎉🎉

おめでとうございます！
ロールを獲得しました！！

---
招待券は、上書きされる前に早めに使ってね。
<#%s>で**2人まで**お友達を招待できるよ。
`
	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(description, internal.ChannelID().INVITATION_LINK),
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1116312750277263390/atari.png",
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

	// 当たりロールを付与します
	if err := addWinnerRole(s, i); err != nil {
		return errors.NewError("ロールを付与できません", err)
	}

	// ハズレ町民ロールを削除します
	for _, role := range i.Member.Roles {
		if role == internal.RoleID().HAZURE {
			if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, internal.RoleID().HAZURE); err != nil {
				return errors.NewError("ハズレ町民ロールを削除できません", err)
			}
		}
	}

	return nil
}

// ハズレの場合のメッセージを送信します
func sendLoserMessage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	description := `

「ハズレ」

また明日チャレンジしてみてね！
もしよければ、<#%s>にもコメントしてね👋

-------------------------
**もう1枚コインをもらえるチャンス!**

1. Twitterで「#アナーキー」のタグをつけて投稿
2. そのURLを <#%s> で共有
3. 運営が確認し ✅のリアクションが付けば、もう1枚コインGET!!!
`

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			description,
			internal.ChannelID().CHAT,
			internal.ChannelID().HAZURE_TWEET,
		),
		Image: &discordgo.MessageEmbedImage{
			URL: randFailureImageURL(),
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

	// ハズレロールを付与します
	if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().HAZURE); err != nil {
		return errors.NewError("ロールを付与できません", err)
	}

	return nil
}

// 当たった場合のランクロールの変更を取得します
//
// 現在のランクロール: 次のランクロール
var nextRankRoleID = map[string]string{
	"none":                     internal.RoleID().AL,
	internal.RoleID().AL:       internal.RoleID().BRONZE,
	internal.RoleID().BRONZE:   internal.RoleID().SILVER,
	internal.RoleID().SILVER:   internal.RoleID().GOLD,
	internal.RoleID().GOLD:     internal.RoleID().PLATINUM,
	internal.RoleID().PLATINUM: internal.RoleID().DIAMOND,
	internal.RoleID().DIAMOND:  internal.RoleID().CRAZY,
}

// 当たりロールを付与します
//
// 当たり,招待券を付与します。
func addWinnerRole(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var (
		prizeRoleNum      int      // 当たりロールの数
		currentRankRoleID = "none" // 現在のAL,Gold,Silver...ロールのID
	)

	for _, role := range i.Member.Roles {
		switch role {
		case internal.RoleID().PRIZE1, internal.RoleID().PRIZE2:
			prizeRoleNum++
		case internal.RoleID().AL,
			internal.RoleID().BRONZE,
			internal.RoleID().SILVER,
			internal.RoleID().GOLD,
			internal.RoleID().PLATINUM,
			internal.RoleID().DIAMOND,
			internal.RoleID().CRAZY:
			currentRankRoleID = role
		}
	}

	// 当たりロールを正しい状態に変更します
	switch prizeRoleNum {
	case 0:
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE1); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	case 1:
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE2); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	case 2:
		if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE1); err != nil {
			return errors.NewError("ロールを削除できません", err)
		}
		if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE2); err != nil {
			return errors.NewError("ロールを削除できません", err)
		}

		// ランクロールを正しい状態に変更します
		{
			if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, currentRankRoleID); err != nil {
				return errors.NewError("現在のランクロールを削除できません", err)
			}

			nextRank, ok := nextRankRoleID[currentRankRoleID]
			if ok {
				if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, nextRank); err != nil {
					return errors.NewError("新しいランクロールを付与できません", err)
				}
			}
		}
	}

	// 招待券を付与します
	{
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().INVITATION1); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().INVITATION2); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	}

	return nil
}

// 当たり判定をします
func isWinner(member *discordgo.Member) (bool, error) {
	rand.Seed(time.Now().UnixNano())

	// 当たりの回数
	var prizedNum int

	for _, roleID := range member.Roles {
		if roleID == internal.RoleID().PRIZE1 ||
			roleID == internal.RoleID().PRIZE2 ||
			roleID == internal.RoleID().AL {
			prizedNum++
		}

		// 検証用ロールの場合は、必ず当たり
		if roleID == internal.RoleID().FOR_TEST_ATARI {
			return true, nil
		}
	}

	switch prizedNum {
	case 0:
		// 当たりなしで参加から2週間以上経過している人は2/3で当たり
		if isTwoWeeksOrMoreBefore(member.JoinedAt) {
			rand.Seed(time.Now().UnixNano())
			return rand.Float64() < 2.0/3.0, nil // 0.0 <= x < 1.0 の範囲でランダムな値を生成するため、この条件は2/3の確率でtrueになります。
		}

		// 当たりなし -> 1/5
		return rand.Intn(5) == 0, nil
	case 1:
		// 当たり1回 -> 1/11
		return rand.Intn(11) == 0, nil
	default:
		// 当たり2回 -> 1/12
		return rand.Intn(12) == 0, nil
	}
}

// ハズレの画像URLをランダムに取得します
func randFailureImageURL() string {
	urls := []string{
		"https://cdn.discordapp.com/attachments/1103240223376293938/1116312806598389771/hazure.png",
		"https://cdn.discordapp.com/attachments/1103240223376293938/1118010136762519642/hazure_02.png",
		"https://cdn.discordapp.com/attachments/1103240223376293938/1119037463650914344/hazure_03.png",
		"https://cdn.discordapp.com/attachments/1103240223376293938/1120971084737560627/hazure_04.png",
		"https://cdn.discordapp.com/attachments/1103240223376293938/1125368074208555078/hazure_05.png",
	}

	rand.Seed(time.Now().UnixNano())

	return urls[rand.Intn(len(urls))]
}

// 指定した日時が今日より2週間以上前であればtrueを返します
func isTwoWeeksOrMoreBefore(t time.Time) bool {
	now := time.Now()
	twoWeeksAgo := now.AddDate(0, 0, -14)
	return t.Before(twoWeeksAgo)
}
