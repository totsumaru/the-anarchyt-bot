package gatcha

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/utils"
	"github.com/techstart35/the-anarchy-bot/internal"
	"math/rand"
	"net/url"
	"time"
)

const CurrentRankRoleNone = "none"

// 結果を送信します
func SendResult(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	editFunc, err := utils.SendInteractionWaitingMessage(s, i, true, true)
	if err != nil {
		return errors.NewError("Waitingメッセージが送信できません")
	}

	isWin, err := isWinner(i.Member)
	if err != nil {
		return errors.NewError("当たり判定に失敗しました", err)
	}

	var embed *discordgo.MessageEmbed

	if isWin {
		embed = createWinnerMessage()

		// 当たりロールを付与します
		if err = addWinnerRole(s, i); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}

		// ハズレ町民ロールを削除します
		for _, role := range i.Member.Roles {
			if role == internal.RoleID().HAZURE {
				if err = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, internal.RoleID().HAZURE); err != nil {
					return errors.NewError("ハズレ町民ロールを削除できません", err)
				}
			}
		}

		// 当たりLogを送信します
		// ここでエラーが出ても処理は継続します。
		if err = sendAtariLog(s, i.Member.User); err != nil {
			errors.SendErrMsg(s, errors.NewError(
				"当たりログを送信できませんが、処理は継続します",
				err,
			), i.Member.User)
		}
	} else {
		embed = createLoserMessage(i.Member.Roles)

		// ハズレロールを付与します
		if err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().HAZURE); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	}

	btn1 := discordgo.Button{
		Label: "ツイートする",
		Style: discordgo.LinkButton,
		URL:   "https://twitter.com/intent/tweet?text=" + url.QueryEscape("#アナーキー"),
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
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

// 当たりの場合のメッセージを送信します
func createWinnerMessage() *discordgo.MessageEmbed {
	description := `
🎉🎉🎉🎉🎉🎉🎉
 「当たり」
🎉🎉🎉🎉🎉🎉🎉

おめでとうございます！
ロールを獲得しました！！

---
<#%s>で**2人まで**お友達を招待できるよ。
招待券は、上書きされる前に早めに使ってね。
`
	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			description,
			internal.ChannelID().INVITATION_LINK,
		),
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/1103240223376293938/1116312750277263390/atari.png",
		},
		Color: internal.ColorBlue,
	}

	return embed
}

// ハズレの場合のメッセージを送信します
func createLoserMessage(roles []string) *discordgo.MessageEmbed {
	description := `
「ハズレ」

また明日チャレンジしてみてね！
もしよければ、<#%s>にもコメントしてね👋

-------------------------
**もう1枚コインをもらえるチャンス!**

1. Twitterで「#アナーキー」のタグをつけて投稿
2. そのURLを <#%s> で共有
`

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			description,
			internal.ChannelID().CHAT,
			internal.ChannelID().HAZURE_TWEET,
		),
		Image: &discordgo.MessageEmbedImage{
			URL: randFailureImageURL(roles),
		},
		Color: internal.ColorBlue,
	}

	return embed
}

// 当たった場合のランクロールの変更を取得します
//
// 現在のランクロール: 次のランクロール
var nextRankRoleID = map[string]string{
	CurrentRankRoleNone:        internal.RoleID().AL,
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
		prizeRoleNum      int                   // 当たりロールの数
		currentRankRoleID = CurrentRankRoleNone // 現在のAL,Gold,Silver...ロールのID
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
		nextRank, ok := nextRankRoleID[currentRankRoleID]

		// 現在のランクがMAXの場合は以下の処理をスキップします
		if ok {
			// 当たりロールを削除します
			// 現在のランクロールがMAXの場合(nextRankRoleIDが無い場合)は、当たりを削除しません
			if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE1); err != nil {
				return errors.NewError("ロールを削除できません", err)
			}
			if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE2); err != nil {
				return errors.NewError("ロールを削除できません", err)
			}

			// 次のランクのロールを付与します
			if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, nextRank); err != nil {
				return errors.NewError("新しいランクロールを付与できません", err)
			}

			// 現在のランクロールを削除します
			// 現在のランクロールがnone,ALロールの場合は、削除するロールがないためこの処理をスキップします
			switch currentRankRoleID {
			case CurrentRankRoleNone, internal.RoleID().AL: // 何もしない
			default:
				if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, currentRankRoleID); err != nil {
					return errors.NewError("現在のランクロールを削除できません", err)
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
func randFailureImageURL(roles []string) string {
	urls := []string{
		"https://cdn.discordapp.com/attachments/1103240223376293938/1116312806598389771/hazure.png",
		"https://cdn.discordapp.com/attachments/1103240223376293938/1118010136762519642/hazure_02.png",
		"https://cdn.discordapp.com/attachments/1103240223376293938/1119037463650914344/hazure_03.png",
		"https://cdn.discordapp.com/attachments/1103240223376293938/1120971084737560627/hazure_04.png",
		"https://cdn.discordapp.com/attachments/1103240223376293938/1125368074208555078/hazure_05.png",
	}

	for _, role := range roles {
		// ALを持っている人はBonsaiコラボ画像を追加
		switch role {
		case internal.RoleID().AL:
			const hazureImageURLBonsai = "https://cdn.discordapp.com/attachments/1103240223376293938/1130011522043760670/hazure_06.png"
			urls = append(urls, hazureImageURLBonsai)
		// ブロンズを持っている人はOtoさんコラボ画像を追加
		case internal.RoleID().BRONZE:
			hazureBronze := []string{
				"https://media.discordapp.net/attachments/1103240223376293938/1135858239175659610/hazure_07.png",
				"https://cdn.discordapp.com/attachments/1103240223376293938/1143139074253795400/hazure_08.png",
			}

			urls = append(urls, hazureBronze...)
		}
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

// 当たりログを送信します
func sendAtariLog(s *discordgo.Session, user *discordgo.User) error {
	now := time.Now()
	formattedTime := now.Format("2006-01-02T15:04:05Z")

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Username,
			IconURL: user.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "ユーザーID",
				Value: user.ID,
			},
			{
				Name:  "日時",
				Value: formattedTime,
			},
		},
	}

	data := &discordgo.MessageSend{
		Content: fmt.Sprintf("0x%s", user.ID),
		Embed:   embed,
	}

	if _, err := s.ChannelMessageSendComplex(
		internal.ChannelID().ATARI_LOG, data,
	); err != nil {
		return errors.NewError("埋め込みメッセージを送信できません", err)
	}

	return nil
}
