package internal

// メッセージコマンド
const (
	CMD_Send_Rule                   = "!an-rule"             // ルール
	CMD_Send_gatcha_Panel           = "!an-gatcha-panel"     // ガチャのパネルの送信（更新の場合はコマンドの後にURLを添付）
	CMD_Send_gatcha_Add_Ticket_Role = "!an-add-role"         // チケットルール付与
	CMD_Send_verify_Panel           = "!an-verify-panel"     // Verifyのパネル
	CMD_Create_Invitation           = "!an-invitation"       // 招待リンク発行(管理者)
	CMD_Send_Invitation_Panel       = "!an-invitation-panel" // 招待リンク発行のパネル送信
	CMD_News                        = "!an-news"             // ニュース
	CMD_Link                        = "!an-link"             // 公式リンクを送信
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
