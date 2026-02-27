package layout_model

type Layout struct {
	FormatVersion int     `json:"formatVersion"`
	Canvas        Canvas  `json:"canvas"`
	Walls         []Wall  `json:"walls"`
	Places        []Place `json:"places"`
}

type Canvas struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Wall struct {
	ID       string `json:"id"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Rotation int    `json:"rotation"`
}

type Place struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Rotation int    `json:"rotation"`
}
