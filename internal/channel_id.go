package internal

import "os"

type Channel struct {
	LOGS            string
	TEST            string
	INVITATION_LINK string
	GATCHA          string
	CHAT            string
	SNEAK_PEEK      string
	NEWS            string
}

func ChannelID() Channel {
	if os.Getenv("ENV") == "dev" {
		// テスト環境
		return Channel{
			LOGS:            "1112548516955439194",
			TEST:            "1069459007321952316",
			INVITATION_LINK: "1112585199373520916",
			GATCHA:          "1115532111693238272",
			CHAT:            "1112319028225130607",
			SNEAK_PEEK:      "1112524163643621379",
			NEWS:            "1116286793327841301",
		}
	} else {
		// 本番環境
		return Channel{
			LOGS:            "1112548516955439194",
			TEST:            "1069459007321952316",
			INVITATION_LINK: "1112585199373520916",
			GATCHA:          "1115532111693238272",
			CHAT:            "1112319028225130607",
			SNEAK_PEEK:      "1112524163643621379",
			NEWS:            "1116286793327841301",
		}
	}
}
