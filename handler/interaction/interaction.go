package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/address"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/gatcha"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/invitation"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/news"
	"github.com/techstart35/the-anarchy-bot/handler/interaction/roles"
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
				errors.SendErrMsg(s, errors.NewError("カプセルのメッセージを送信できません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_gatcha_Open:
			if err := gatcha.SendResult(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("結果を送信できません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_gatcha_Notice:
			if err := gatcha.SendNoticeRoleMessage(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("通知ボタンの処理ができません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_gatcha_Notice_Remove:
			if err := gatcha.RemoveNoticeRole(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("通知ロールを削除できません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_Verify:
			if err := verify.Verify(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("認証できません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_News_Cancel:
			if err := news.Cancel(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("キャンセル処理を実行できません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_News_Send:
			if err := news.Transfer(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("メッセージを転送できません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_Invitation:
			if err := invitation.ReplyLink(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("招待リンクを発行できません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_Address_Modal_Open:
			if err := address.OpenModal(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("Modalを開けません", err), i.Member.User)
				return
			}
		case internal.Interaction_CustomID_Address_Check:
			if err := address.Check(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("アドレスの確認ができません", err), i.Member.User)
				return
			}
		}
	case discordgo.InteractionApplicationCommand:
		name := i.Data.(discordgo.ApplicationCommandInteractionData).Name
		if name == internal.Slash_CMD_MyRoles {
			if err := roles.GetRoles(s, i); err != nil {
				errors.SendErrMsg(s, errors.NewError("ロールの一覧を取得できません", err), i.Member.User)
				return
			}
		}
	// Modal
	case discordgo.InteractionModalSubmit:
		if err := address.UpsertFromModal(s, i); err != nil {
			errors.SendErrMsg(s, errors.NewError("Modalの内容を処理できません", err), i.Member.User)
			return
		}
	}
}
