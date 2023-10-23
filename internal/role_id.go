package internal

import "os"

type Role struct {
	GATCHA_COIN    string // ガチャコイン
	BONUS_COIN     string // ボーナスコイン
	GATCHA_NOTICE  string // ガチャ通知
	VERIFIED       string // 入場券
	PRIZE1         string // 当たり1
	PRIZE2         string // 当たり2
	INVITATION1    string // 招待券1
	INVITATION2    string // 招待券2
	TEST           string // 検証用ロール
	AL             string // AL
	HAZURE         string // はずれ町民
	COIN_2_ADDED   string // コイン2枚目付与済み
	BRONZE         string // ブロンズガチャーキー
	SILVER         string // シルバーガチャーキー
	GOLD           string // ゴールドガチャーキー
	PLATINUM       string // プラチナガチャーキー
	DIAMOND        string // ダイヤモンドガチャーキー
	CRAZY          string // クレイジーガチャーキー
	FOR_TEST_ATARI string // [検証用]ガチャ当たり100%
	TOKYO_ANARCHY  string // 東京アナーキー
	CHAINSAW_CLUB  string // チェンソー倶楽部
	SUBMITTED      string // アドレス提出済み
}

// ロールID
func RoleID() Role {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return Role{
			GATCHA_COIN:    "1115563590825553941",
			BONUS_COIN:     "1127490096707424337",
			GATCHA_NOTICE:  "1117087993375768646",
			VERIFIED:       "1112562612044038295",
			PRIZE1:         "1115569765101076510",
			PRIZE2:         "1115570175949951067",
			INVITATION1:    "1116535909093998664",
			INVITATION2:    "1116536092573839400",
			TEST:           "1117102782185492560",
			AL:             "1118720616258863226",
			HAZURE:         "1123476359675650078",
			COIN_2_ADDED:   "1124556325867765871",
			BRONZE:         "1125684880710303805",
			SILVER:         "1125685256960356413",
			GOLD:           "1125685475835904001",
			PLATINUM:       "1125685655951904839",
			DIAMOND:        "1125685776299081769",
			CRAZY:          "1125686122454982718",
			FOR_TEST_ATARI: "1127478160460624042",
			TOKYO_ANARCHY:  "1112319701742260284",
			CHAINSAW_CLUB:  "1112319985960882296",
			SUBMITTED:      "1166034627820015770",
		}
	} else {
		// 本番環境
		return Role{
			GATCHA_COIN:    "1115563590825553941",
			BONUS_COIN:     "1127490096707424337",
			GATCHA_NOTICE:  "1117087993375768646",
			VERIFIED:       "1112562612044038295",
			PRIZE1:         "1115569765101076510",
			PRIZE2:         "1115570175949951067",
			INVITATION1:    "1116535909093998664",
			INVITATION2:    "1116536092573839400",
			TEST:           "1117102782185492560",
			AL:             "1118720616258863226",
			HAZURE:         "1123476359675650078",
			COIN_2_ADDED:   "1124556325867765871",
			BRONZE:         "1125684880710303805",
			SILVER:         "1125685256960356413",
			GOLD:           "1125685475835904001",
			PLATINUM:       "1125685655951904839",
			DIAMOND:        "1125685776299081769",
			CRAZY:          "1125686122454982718",
			FOR_TEST_ATARI: "1127478160460624042",
			TOKYO_ANARCHY:  "1112319701742260284",
			CHAINSAW_CLUB:  "1112319985960882296",
			SUBMITTED:      "1166034627820015770",
		}
	}
}
