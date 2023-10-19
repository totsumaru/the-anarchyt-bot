package message

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/message/address"
	"github.com/techstart35/the-anarchy-bot/handler/message/alert"
	"github.com/techstart35/the-anarchy-bot/handler/message/command"
	"github.com/techstart35/the-anarchy-bot/handler/message/gatcha"
	"github.com/techstart35/the-anarchy-bot/handler/message/info"
	"github.com/techstart35/the-anarchy-bot/handler/message/invitation"
	"github.com/techstart35/the-anarchy-bot/handler/message/news"
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
			errors.SendErrMsg(s, errors.NewError("ルールを送信できません", err), m.Author)
		}
		return
	case internal.CMD_Send_gatcha_Add_Coin_Role:
		if err := gatcha.AddCoinRole(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("ガチャコインロールを付与できません", err), m.Author)
		}
		return
	case internal.CMD_Send_gatcha_Add_Bonus_Coin_Role:
		if err := gatcha.AddBonusCoinRole(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("ボーナスコインロールを付与できません", err), m.Author)
		}
		return
	case internal.CMD_Send_gatcha_Notice:
		if err := gatcha.SendNotice(s); err != nil {
			errors.SendErrMsg(s, errors.NewError("ガチャ通知を送信できません", err), m.Author)
		}
		return
	case internal.CMD_Send_verify_Panel:
		if err := verify.SendPanel(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("パネルを送信できません", err), m.Author)
		}
		return
	case internal.CMD_Create_Invitation:
		if err := invitation.CreateInvitationForAdmin(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("招待を作成できません", err), m.Author)
		}
		return
	case internal.CMD_Info_Update:
		if err := info.UpdatePublicInfos(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("公式情報を更新できません", err), m.Author)
		}
		return
	case internal.CMD_Send_Invitation_Panel:
		if err := invitation.SendPanel(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("招待リンク発行のパネルを送信できません", err), m.Author)
		}
		return
	case internal.CMD_ADD_SLASH_COMMAND:
		if err := command.RegisterSlashCommand(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("スラッシュコマンドを登録できません", err), m.Author)
		}
		return
	case internal.CMD_ADD_INVITE_ROLE:
		if err := invitation.AddInvitationRole(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("ロールを付与できません", err), m.Author)
		}
		return
	case internal.CMD_Send_Address_Panel:
		if err := address.SendAddressPanel(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError(
				"アドレス集計のパネルを送信できません", err,
			), m.Author)
		}
	}

	// ガチャパネルコマンド
	if strings.Contains(m.Content, internal.CMD_Send_gatcha_Panel) {
		// コマンドのみの場合は新規送信
		// URLがついている場合は更新
		if m.Content == internal.CMD_Send_gatcha_Panel {
			if err := gatcha.SendPanel(s, m, ""); err != nil {
				errors.SendErrMsg(s, errors.NewError("パネルを送信できません", err), m.Author)
			}
			return
		} else {
			url := strings.Split(m.Content, " ")[1]
			if err := gatcha.SendPanel(s, m, url); err != nil {
				errors.SendErrMsg(s, errors.NewError("パネルを更新できません", err), m.Author)
			}
			return
		}
	}

	// newsコマンド
	if strings.Contains(m.Content, internal.CMD_News) {
		if err := news.Confirm(s, m); err != nil {
			errors.SendErrMsg(s, errors.NewError("ニュースの確認を送信できません", err), m.Author)
		}

		return
	}

	// sneak-peekをchatに転送
	if err := sneek_peek.Notice(s, m); err != nil {
		errors.SendErrMsg(s, errors.NewError("スニークピークを転送できません", err), m.Author)
	}

	// newsをchatに転送
	if err := news.Notice(s, m); err != nil {
		errors.SendErrMsg(s, errors.NewError("Newsを転送できません", err), m.Author)
	}

	// ハズレの人に2枚目のコインロールを付与
	if err := gatcha.AddSecondCoinRoleForHazureUser(s, m); err != nil {
		errors.SendErrMsg(s, errors.NewError("2枚目のコインロールを付与できません", err), m.Author)
	}

	// アラートを送信
	if err := alert.SendAlert(s, m); err != nil {
		errors.SendErrMsg(s, errors.NewError("アラートを送信できません", err), m.Author)
	}

	// TODO: テスト用のため、以下は削除
	{
		if strings.Contains(m.Content, "!test") {
			targetID := strings.Split(m.Content, " ")[1]
			target, err := s.GuildMember(m.GuildID, targetID)
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("テスト用の確認メッセージを送信できません", err), m.Author)
			}

			if _, err = s.ChannelMessageSend(internal.ChannelID().TEST, target.AvatarURL("")); err != nil {
				errors.SendErrMsg(s, errors.NewError("テスト用の確認メッセージを送信できません", err), m.Author)
			}
		}
	}

	return
}
