package gatcha

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/utils"
	"github.com/techstart35/the-anarchy-bot/internal"
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
		addedRoleID, err := addWinnerRole(s, i)
		if err != nil {
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
		if err = sendAtariLog(s, i.Member.User, addedRoleID); err != nil {
			errors.SendErrMsg(s, errors.NewError(
				"当たりログを送信できませんが、処理は継続します", err,
			), i.Member.User)
		}
	} else {
		embed = createLoserMessage(i.Member.Roles)

		// ハズレロールを付与します
		if err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().HAZURE); err != nil {
			return errors.NewError("ロールを付与できません", err)
		}
	}

	btn2 := discordgo.Button{
		Emoji: discordgo.ComponentEmoji{
			Name: "✅",
		},
		Label: "招待タグをつけてポスト",
		Style: discordgo.LinkButton,
		URL:   "https://twitter.com/intent/tweet?text=" + url.QueryEscape("#アナーキー #アナーキーにおいでよ"),
	}

	btn1 := discordgo.Button{
		Label: "Xにポスト",
		Style: discordgo.LinkButton,
		URL:   "https://twitter.com/intent/tweet?text=" + url.QueryEscape("#アナーキー"),
	}

	components := []discordgo.MessageComponent{btn1}
	for _, role := range i.Member.Roles {
		if role == internal.RoleID().INVITATION1 ||
			role == internal.RoleID().INVITATION2 {
			components = []discordgo.MessageComponent{btn2, btn1}
		}
	}

	actions := discordgo.ActionsRow{
		Components: components,
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
`
	embed := &discordgo.MessageEmbed{
		Description: description,
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
**ワンモアチャンス!**

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
	internal.RoleID().CRAZY:    internal.RoleID().FUCKIN,
}

// 当たりロールを付与します
//
// 当たり,招待券を付与します。
// 新規で付与したロールIDを返します。
func addWinnerRole(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	var (
		currentPrizeRoleNum int                   // 現在の当たりロールの数
		currentRankRoleID   = CurrentRankRoleNone // 現在のAL,Gold,Silver...ロールのID
	)

	for _, role := range i.Member.Roles {
		switch role {
		case internal.RoleID().PRIZE1,
			internal.RoleID().PRIZE2,
			internal.RoleID().PRIZE3,
			internal.RoleID().PRIZE4,
			internal.RoleID().PRIZE5,
			internal.RoleID().PRIZE6,
			internal.RoleID().PRIZE7,
			internal.RoleID().PRIZE8,
			internal.RoleID().PRIZE9,
			internal.RoleID().PRIZE10,
			internal.RoleID().PRIZE11,
			internal.RoleID().PRIZE12,
			internal.RoleID().PRIZE13,
			internal.RoleID().PRIZE14:
			currentPrizeRoleNum++
		case internal.RoleID().AL:
			// ALは他のランクロールと共存する可能性があるため、
			// currentRankRoleIDが`none`の場合のみ設定できる
			if currentRankRoleID == CurrentRankRoleNone {
				currentRankRoleID = role
			}
		case internal.RoleID().BRONZE,
			internal.RoleID().SILVER,
			internal.RoleID().GOLD,
			internal.RoleID().PLATINUM,
			internal.RoleID().DIAMOND,
			internal.RoleID().CRAZY,
			internal.RoleID().FUCKIN:
			currentRankRoleID = role
		}
	}

	var addedRoleID string

	// 当たりロールを正しい状態に変更します
	switch currentPrizeRoleNum {
	case 0:
		addRoleID := internal.RoleID().PRIZE1
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 1:
		addRoleID := internal.RoleID().PRIZE2
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 2:
		// 現在CRAZYの場合は、当たり3を付与します
		if currentRankRoleID == internal.RoleID().CRAZY {
			// 次のランクのロールを付与します
			if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE3); err != nil {
				return "", errors.NewError("新しいランクロールを付与できません", err)
			}
			break
		}

		nextRank, ok := nextRankRoleID[currentRankRoleID]

		// 現在のランクがMAXの場合は以下の処理をスキップします
		if ok {
			// 当たりロールを削除します
			// 現在のランクロールがMAXの場合(nextRankRoleIDが無い場合)は、当たりを削除しません
			if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE1); err != nil {
				return "", errors.NewError("ロールを削除できません", err)
			}
			if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, internal.RoleID().PRIZE2); err != nil {
				return "", errors.NewError("ロールを削除できません", err)
			}

			// 次のランクのロールを付与します
			if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, nextRank); err != nil {
				return "", errors.NewError("新しいランクロールを付与できません", err)
			}
			addedRoleID = nextRank

			// 現在のランクロールを削除します
			// 現在のランクロールがnone,ALロールの場合は、削除するロールがないためこの処理をスキップします
			switch currentRankRoleID {
			case CurrentRankRoleNone, internal.RoleID().AL: // 何もしない
			default:
				if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, currentRankRoleID); err != nil {
					return "", errors.NewError("現在のランクロールを削除できません", err)
				}
			}
		}
	case 3:
		addRoleID := internal.RoleID().PRIZE4
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 4:
		addRoleID := internal.RoleID().PRIZE5
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 5:
		addRoleID := internal.RoleID().PRIZE6
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 6:
		addRoleID := internal.RoleID().PRIZE7
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 7:
		addRoleID := internal.RoleID().PRIZE8
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 8:
		addRoleID := internal.RoleID().PRIZE9
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 9:
		addRoleID := internal.RoleID().PRIZE10
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 10:
		addRoleID := internal.RoleID().PRIZE11
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 11:
		addRoleID := internal.RoleID().PRIZE12
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 12:
		addRoleID := internal.RoleID().PRIZE13
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 13:
		addRoleID := internal.RoleID().PRIZE14
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID
	case 14:
		addRoleID := internal.RoleID().FUCKIN
		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, addRoleID); err != nil {
			return "", errors.NewError("ロールを付与できません", err)
		}
		addedRoleID = addRoleID

		// 当たりロールを全て削除します
		prizeRoles := []string{
			currentRankRoleID, // CRAZYロール
			internal.RoleID().PRIZE1,
			internal.RoleID().PRIZE2,
			internal.RoleID().PRIZE3,
			internal.RoleID().PRIZE4,
			internal.RoleID().PRIZE5,
			internal.RoleID().PRIZE6,
			internal.RoleID().PRIZE7,
			internal.RoleID().PRIZE8,
			internal.RoleID().PRIZE9,
			internal.RoleID().PRIZE10,
			internal.RoleID().PRIZE11,
			internal.RoleID().PRIZE12,
			internal.RoleID().PRIZE13,
			internal.RoleID().PRIZE14,
		}
		for _, prizeRoleID := range prizeRoles {
			if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, prizeRoleID); err != nil {
				return "", errors.NewError("ロールを削除できません", err)
			}
		}
	}

	return addedRoleID, nil
}

// 当たり判定をします
func isWinner(member *discordgo.Member) (bool, error) {
	rand.Seed(time.Now().UnixNano())

	// 当たりの回数
	prizedNum := 0
	hasAL := false

	for _, roleID := range member.Roles {
		switch roleID {
		case internal.RoleID().PRIZE1, internal.RoleID().PRIZE2:
			prizedNum++
		case internal.RoleID().AL:
			hasAL = true
		case internal.RoleID().FOR_TEST_ATARI:
			// 検証用ロールの場合は、必ず当たり。ここで終了
			return true, nil
		}
	}

	// ALを持っている人は一律 1/12
	if hasAL {
		return rand.Intn(12) == 0, nil
	}

	// ALを持っていない人
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
		// 当たり1回 -> 1/8 -> 1/6(AL登録前なので、1/6に変更)
		return rand.Intn(6) == 0, nil
	default:
		// 当たり2回 -> 1/10
		return rand.Intn(8) == 0, nil
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
		case internal.RoleID().BRONZE:
			hazureImages := []string{
				"https://media.discordapp.net/attachments/1103240223376293938/1135858239175659610/hazure_07.png",
				"https://cdn.discordapp.com/attachments/1103240223376293938/1143139074253795400/hazure_08.png",
			}
			urls = append(urls, hazureImages...)
		case internal.RoleID().GOLD,
			internal.RoleID().PLATINUM,
			internal.RoleID().DIAMOND,
			internal.RoleID().CRAZY:
			hazureImages := []string{
				"https://cdn.discordapp.com/attachments/1103240223376293938/1155802039071285248/hazure_09.png",
			}
			urls = append(urls, hazureImages...)
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
func sendAtariLog(s *discordgo.Session, user *discordgo.User, grantedRoleID string) error {
	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Username,
			IconURL: user.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "user id",
				Value: user.ID,
			},
			{
				Name:  "date",
				Value: formattedTime,
			},
			{
				Name:  "granted role",
				Value: fmt.Sprintf("<@&%s>", grantedRoleID),
			},
		},
	}

	data := &discordgo.MessageSend{
		Content: fmt.Sprintf("atari%s", user.ID),
		Embed:   embed,
	}

	if _, err := s.ChannelMessageSendComplex(
		internal.ChannelID().ATARI_LOG, data,
	); err != nil {
		return errors.NewError("埋め込みメッセージを送信できません", err)
	}

	return nil
}
