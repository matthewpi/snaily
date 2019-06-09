package hypixel

type playerResponse struct {
	Success bool   `json:"success"`
	Player  Player `json:"player"`
}

type Player struct {
	ID                 string        `json:"_id"`
	UUID               string        `json:"uuid"`
	FirstLogin         int64         `json:"firstLogin"`
	Playername         string        `json:"playername"`
	LastLogin          int64         `json:"lastLogin"`
	Displayname        string        `json:"displayname"`
	KnownAliases       []string      `json:"knownAliases"`
	KnownAliasesLower  []string      `json:"knownAliasesLower"`
	Stats              Stats         `json:"stats"`
	LastLogout         int64         `json:"lastLogout"`
	AchievementPoints  int64         `json:"achievementPoints"`
	NetworkExp         float64       `json:"networkExp"`
	Karma              int64         `json:"karma"`
	NewPackageRank     string        `json:"newPackageRank"`
	FriendRequestsUUID []interface{} `json:"friendRequestsUuid"`
	MostRecentGameType string        `json:"mostRecentGameType"`
}

type Stats struct {
	Pit Pit `json:"Pit"`
}

type Pit struct {
	Profile     PitProfile       `json:"profile"`
	PitStatsPtl map[string]int64 `json:"pit_stats_ptl"`
}

type PitProfile struct {
	Renown                     int64             `json:"renown"`
	OutgoingOffers             []interface{}     `json:"outgoing_offers"`
	LastSave                   int64             `json:"last_save"`
	Prestiges                  []Prestige        `json:"prestiges"`
	TradeTimestamps            []interface{}     `json:"trade_timestamps"`
	ZeroPointThreeGoldTransfer bool              `json:"zero_point_three_gold_transfer"`
	RenownUnlocks              []RenownUnlock    `json:"renown_unlocks"`
	Unlocks1                   []RenownUnlock    `json:"unlocks_1"`
	InvEnderchest              DeathRecaps       `json:"inv_enderchest"`
	DeathRecaps                DeathRecaps       `json:"death_recaps"`
	Cash                       float64           `json:"cash"`
	LastMidfightDisconnect     int64             `json:"last_midfight_disconnect"`
	LeaderboardStats           map[string]int64  `json:"leaderboard_stats"`
	SelectedPerk3              interface{}       `json:"selected_perk_3"`
	SelectedPerk2              string            `json:"selected_perk_2"`
	InvArmor                   DeathRecaps       `json:"inv_armor"`
	SelectedPerk1              string            `json:"selected_perk_1"`
	SelectedPerk0              string            `json:"selected_perk_0"`
	ItemStash                  DeathRecaps       `json:"item_stash"`
	GoldTransactions           []GoldTransaction `json:"gold_transactions"`
	LoginMessages              []interface{}     `json:"login_messages"`
	HotbarFavorites            []int64           `json:"hotbar_favorites"`
	RecentKills                []RecentKill      `json:"recent_kills"`
	InvContents                DeathRecaps       `json:"inv_contents"`
	XP                         int64             `json:"xp"`
	Bounties                   []interface{}     `json:"bounties"`
	Unlocks                    []RenownUnlock    `json:"unlocks"`
	CashDuringPrestige1        float64           `json:"cash_during_prestige_1"`
	CashDuringPrestige0        float64           `json:"cash_during_prestige_0"`
}

type DeathRecaps struct {
	Type int64   `json:"type"`
	Data []int64 `json:"data"`
}

type GoldTransaction struct {
	Amount    int64 `json:"amount"`
	Timestamp int64 `json:"timestamp"`
}

type Prestige struct {
	Index        int64 `json:"index"`
	XPOnPrestige int64 `json:"xp_on_prestige"`
	Timestamp    int64 `json:"timestamp"`
}

type RecentKill struct {
	Victim    string `json:"victim"`
	Timestamp int64  `json:"timestamp"`
}

type RenownUnlock struct {
	Tier        int64  `json:"tier"`
	AcquireDate int64  `json:"acquireDate"`
	Key         string `json:"key"`
}
