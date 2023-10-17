package rule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
)

// ルールを送信します
func SendRule(s *discordgo.Session, m *discordgo.MessageCreate) error {
	ruleTmpl := `
居心地の良いプロジェクトを作り上げていくため、

1️⃣お互いに敬意を払い
2️⃣歓迎する気持ちを大切に

楽しいコミュニティにしていきましょう！

以下の行為は禁止しています。
1.荒らし行為、ヘイトスピーチまたはそれに準ずる行為
2.他人への誹謗中傷
3.ハラスメント、性差別、人種差別
4.悪質な広告やスパムメール

これらの行為が見受けられた場合、その他運営が悪質と判断した場合、強制退会などの対応をとらせて頂くことがあります。
ご理解ご了承よろしくお願い致します。
`

	embed := &discordgo.MessageEmbed{
		Title:       "Rules",
		Description: ruleTmpl,
	}

	edit := &discordgo.MessageEdit{
		ID:      "1115241135359668234",
		Channel: "1112544665942634595",
		Embed:   embed,
	}

	if _, err := s.ChannelMessageEditComplex(edit); err != nil {
		return errors.NewError("メッセージを送信できません", err)
	}

	return nil
}
