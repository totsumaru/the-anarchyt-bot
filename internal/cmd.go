package internal

// メッセージコマンド
const (
	CMD_Send_Rule                   = "!an-rule"             // ルール
	CMD_Send_gatcha_Panel           = "!an-gatcha-panel"     // ガチャのパネルの送信（更新の場合はコマンドの後にURLを添付）
	CMD_Send_gatcha_Notice          = "!an-gatcha-notice"    // ガチャ通知の送信（毎朝7:00）
	CMD_Send_gatcha_Add_Ticket_Role = "!an-add-role"         // チケットルール付与
	CMD_Send_verify_Panel           = "!an-verify-panel"     // Verifyのパネル
	CMD_Create_Invitation           = "!an-invitation"       // 招待リンク発行(管理者)
	CMD_Send_Invitation_Panel       = "!an-invitation-panel" // 招待リンク発行のパネル送信
	CMD_News                        = "!an-news"             // ニュース
	CMD_Info                        = "!an-info"             // 公式情報を送信
	CMD_Info_Update                 = "!an-info-update"      // 公式情報を更新
	CMD_ADD_SLASH_COMMAND           = "!an-add-cmd"          // スラッシュコマンドの追加
	CMD_ADD_INVITE_ROLE             = "!an-add-invite-role"  // 招待券ロールを付与
)

// スラッシュコマンド
const (
	Slash_CMD_MyRoles = "my-roles"
)

// InteractionのカスタムID
const (
	// ガチャ
	Interaction_CustomID_gatcha_Go            = "gatcha-go"
	Interaction_CustomID_gatcha_Open          = "gatcha-open"
	Interaction_CustomID_gatcha_Notice        = "gatcha-notice"
	Interaction_CustomID_gatcha_Notice_Remove = "gatcha-notice-remove"
	// Verify
	Interaction_CustomID_Verify = "verify"
	// news
	Interaction_CustomID_News_Send   = "news-send"
	Interaction_CustomID_News_Cancel = "news-cancel"
	// Invitation
	Interaction_CustomID_Invitation = "invitation-create"
)
