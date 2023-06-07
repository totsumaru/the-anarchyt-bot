package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/message/gatcha"
	"github.com/techstart35/the-anarchy-bot/handler/message/invitation"
	"github.com/techstart35/the-anarchy-bot/handler/message/rule"
	"github.com/techstart35/the-anarchy-bot/handler/message/sneek_peek"
	"github.com/techstart35/the-anarchy-bot/handler/message/verify"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// メッセージが作成された時のハンドラーです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case internal.CMD_Send_Rule:
		if err := rule.SendRule(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("ルールを送信できません", err))
		}
	case internal.CMD_Send_gatcha_Add_Ticket_Role:
		if err := gatcha.AddRole(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("ロールを付与できません", err))
		}
	case internal.CMD_Send_gatcha_Panel:
		if err := gatcha.SendPanel(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("パネルを送信できません", err))
		}
	case internal.CMD_Send_verify_Panel:
		if err := verify.SendPanel(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("パネルを送信できません", err))
		}
	case internal.CMD_Create_Invitation:
		if err := invitation.CreateInvitation(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("招待を作成できません", err))
		}
	}

	if err := sneek_peek.Transfer(s, m); err != nil {
		errors.SendErrMsg(s, errors.NewError("スニークピークを転送できません", err))
	}
}
