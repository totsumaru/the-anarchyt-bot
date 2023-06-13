package internal

import "os"

type Channel struct {
	LOGS            string
	TEST            string
	START           string // 最初に
	GATCHA          string // ロールガチャ
	CHAT            string // チャット
	SNEAK_PEEK      string // チラ見せ
	NEWS            string // ニュース
	INVITATION_LINK string // 招待リンク
}

func ChannelID() Channel {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return Channel{
			LOGS:            "1112548516955439194",
			TEST:            "1069459007321952316",
			START:           "1112585199373520916",
			GATCHA:          "1115532111693238272",
			CHAT:            "1112319028225130607",
			SNEAK_PEEK:      "1112524163643621379",
			NEWS:            "1116286793327841301",
			INVITATION_LINK: "1116549608663949393",
		}
	} else {
		// 本番環境
		return Channel{
			LOGS:            "1112548516955439194",
			TEST:            "1069459007321952316",
			START:           "1112585199373520916",
			GATCHA:          "1115532111693238272",
			CHAT:            "1112319028225130607",
			SNEAK_PEEK:      "1112524163643621379",
			NEWS:            "1116286793327841301",
			INVITATION_LINK: "1116549608663949393",
		}
	}
}
