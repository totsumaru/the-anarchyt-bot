package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/message/gatcha"
	"github.com/techstart35/the-anarchy-bot/handler/message/invitation"
	"github.com/techstart35/the-anarchy-bot/handler/message/link"
	"github.com/techstart35/the-anarchy-bot/handler/message/news"
	"github.com/techstart35/the-anarchy-bot/handler/message/rule"
	"github.com/techstart35/the-anarchy-bot/handler/message/sneek_peek"
	"github.com/techstart35/the-anarchy-bot/handler/message/verify"
	"github.com/techstart35/the-anarchy-bot/internal"
	"strings"
)

// メッセージが作成された時のハンドラーです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case internal.CMD_Send_Rule:
		if err := rule.SendRule(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("ルールを送信できません", err))
		}
		return
	case internal.CMD_Send_gatcha_Add_Ticket_Role:
		if err := gatcha.AddRole(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("チケットロールを付与できません", err))
		}
		return
	case internal.CMD_Send_gatcha_Notice:
		if err := gatcha.SendNotice(s); err != nil {
			errors.SendErrMsg(s, errors.NewError("ガチャ通知を送信できません", err))
		}
		return
	case internal.CMD_Send_verify_Panel:
		if err := verify.SendPanel(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("パネルを送信できません", err))
		}
		return
	case internal.CMD_Create_Invitation:
		if err := invitation.CreateInvitationForAdmin(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("招待を作成できません", err))
		}
	case internal.CMD_Link:
		if err := link.SendPublicURL(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("公式リンクを送信できません", err))
		}
	case internal.CMD_Send_Invitation_Panel:
		if err := invitation.SendPanel(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("招待リンク発行のパネルを送信できません", err))
		}
		return
	}

	// ガチャパネル
	if strings.Contains(m.Content, internal.CMD_Send_gatcha_Panel) {
		// コマンドのみの場合は新規送信
		// URLがついている場合は更新
		if m.Content == internal.CMD_Send_gatcha_Panel {
			if err := gatcha.SendPanel(s, m, ""); err != nil {
				errors.SendErrMsg(s, errors.NewError("パネルを送信できません", err))
			}
			return
		} else {
			url := strings.Split(m.Content, " ")[1]
			if err := gatcha.SendPanel(s, m, url); err != nil {
				errors.SendErrMsg(s, errors.NewError("パネルを更新できません", err))
			}
			return
		}
	}

	// news
	if strings.Contains(m.Content, internal.CMD_News) {
		if err := news.Confirm(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("ニュースの確認を送信できません", err))
		}

		return
	}

	// sneak-peek
	if err := sneek_peek.Notice(s, m); err != nil {
		errors.SendErrMsg(s, errors.NewError("スニークピークを転送できません", err))
	}

	// news
	if err := news.Notice(s, m); err != nil {
		errors.SendErrMsg(s, errors.NewError("Newsを転送できません", err))
	}

	return
}
