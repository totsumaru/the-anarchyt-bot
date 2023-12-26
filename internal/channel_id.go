package internal

import "os"

type Channel struct {
	TEAM           string
	LOGS           string
	TEST           string
	START          string // 最初に
	GATCHA         string // ロールガチャ
	CHAT           string // チャット
	SNEAK_PEEK     string // チラ見せ
	NEWS           string // ニュース
	HAZURE_MACHI_1 string // はずれ町一丁目
	HAZURE_TWEET   string // はずれ町瓦版
	PUBLIC_INFO    string // 公式情報
	BOT_COMMAND    string // botコマンド
	ATARI_LOG      string // 当たりログ
}

func ChannelID() Channel {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return Channel{
			TEAM:           "",
			LOGS:           "1112548516955439194",
			TEST:           "1069459007321952316",
			START:          "1112585199373520916",
			GATCHA:         "1115532111693238272",
			CHAT:           "1112319028225130607",
			SNEAK_PEEK:     "1112524163643621379",
			NEWS:           "1116286793327841301",
			HAZURE_MACHI_1: "1123476926305157243",
			HAZURE_TWEET:   "1125221784795480135",
			PUBLIC_INFO:    "1116472032738152588",
			BOT_COMMAND:    "1127463906676330506",
			ATARI_LOG:      "1143727344616489021",
		}
	} else {
		// 本番環境
		return Channel{
			TEAM:           "1067807967950422096",
			LOGS:           "1112548516955439194",
			TEST:           "1069459007321952316",
			START:          "1112585199373520916",
			GATCHA:         "1115532111693238272",
			CHAT:           "1112319028225130607",
			SNEAK_PEEK:     "1112524163643621379",
			NEWS:           "1116286793327841301",
			HAZURE_MACHI_1: "1123476926305157243",
			HAZURE_TWEET:   "1125221784795480135",
			PUBLIC_INFO:    "1116472032738152588",
			BOT_COMMAND:    "1127463906676330506",
			ATARI_LOG:      "1143727344616489021",
		}
	}
}
