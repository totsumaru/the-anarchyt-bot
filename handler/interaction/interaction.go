package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/gatcha"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/invitation"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/news"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/verify"
	"github.com/techstart35/the-anarchy-bot/internal"
)

// コマンドが実行された時のハンドラーです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case internal.Interaction_CustomID_gatcha_Go:
			if err := gatcha.SendCapsule(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("カプセルのメッセージを送信できません", err))
				return
			}
		case internal.Interaction_CustomID_gatcha_Open:
			if err := gatcha.SendResult(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("結果を送信できません", err))
				return
			}
		case internal.Interaction_CustomID_Verify:
			if err := verify.Verify(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("認証できません", err))
				return
			}
		case internal.Interaction_CustomID_News_Cancel:
			if err := news.Cancel(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("キャンセル処理を実行できません", err))
				return
			}
		case internal.Interaction_CustomID_News_Send:
			if err := news.Transfer(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("メッセージを転送できません", err))
				return
			}
		case internal.Interaction_CustomID_Invitation:
			if err := invitation.ReplyLink(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("招待リンクを発行できません", err))
				return
			}
		}
	}
}
