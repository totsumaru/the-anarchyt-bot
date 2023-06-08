package internal

// メッセージコマンド
const (
	CMD_Send_Rule                   = "!an-rule"         // ルール
	CMD_Send_gatcha_Panel           = "!an-gatcha-panel" // ガチャのパネルの送信（更新の場合はコマンドの後にURLを添付）
	CMD_Send_gatcha_Add_Ticket_Role = "!an-add-role"     // チケットルール付与
	CMD_Send_verify_Panel           = "!an-verify-panel" // Verifyのパネル
	CMD_Create_Invitation           = "!an-invitation"   // 招待リンク発行
	CMD_News                        = "!an-news"         // ニュース
)

// InteractionのカスタムID
const (
	// ガチャ
	Interaction_CustomID_gatcha_Go   = "gatcha-go"
	Interaction_CustomID_gatcha_Open = "gatcha-open"
	// Verify
	Interaction_CustomID_Verify = "verify"
	// news
	Interaction_CustomID_News_Send   = "news-send"
	Interaction_CustomID_News_Cancel = "news-cancel"
)
