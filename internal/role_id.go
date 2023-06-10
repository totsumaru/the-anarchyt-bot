package internal

import "os"

type Role struct {
	GATCHA_TICKET string // ガチャチケット
	GATCHA_NOTICE string // ガチャ通知
	VERIFIED      string // 入場券
	PRIZE1        string // 当たり1
	PRIZE2        string // 当たり2
	PRIZE3        string // 当たり3
	INVITATION1   string // 招待券1
	INVITATION2   string // 招待券2
}

// ロールID
func RoleID() Role {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return Role{
			GATCHA_TICKET: "1115563590825553941",
			GATCHA_NOTICE: "1117087993375768646",
			VERIFIED:      "1112562612044038295",
			PRIZE1:        "1115569765101076510",
			PRIZE2:        "1115570175949951067",
			PRIZE3:        "1115570239183269930",
			INVITATION1:   "1116535909093998664",
			INVITATION2:   "1116536092573839400",
		}
	} else {
		// 本番環境
		return Role{
			GATCHA_TICKET: "1115563590825553941",
			GATCHA_NOTICE: "1117087993375768646",
			VERIFIED:      "1112562612044038295",
			PRIZE1:        "1115569765101076510",
			PRIZE2:        "1115570175949951067",
			PRIZE3:        "1115570239183269930",
			INVITATION1:   "1116535909093998664",
			INVITATION2:   "1116536092573839400",
		}
	}
}
