package reaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/the-anarchy-bot/errors"
	"github.com/techstart35/the-anarchy-bot/handler/reaction/gatcha"
)

// リアクションが付与された時のハンドラーです
func ReactionCreateHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if err := gatcha.AddSubCoinRoleForHazureUser(s, r); err != nil {
		errors.SendErrMsg(s, errors.NewError("リアクションを付与できません", err), r.Member.User)
	}
}
