package types

type Meta struct {
	CharNames []string `json:"character_names"`
	ClanNames []string `json:"clan_names"`
}

type Link struct {
	Label string
	Url   string
}

type HomePage struct {
	Players []Link
	Clans   []Link
	Logs    []Link
}

type DetailPage struct {
	Logs []Link
}
