package layout_schema

import _ "embed"

//go:embed layout_schema.json
var LayoutSchemaData string

type LayoutSchema struct {
	FormatVersion int     `json:"formatVersion"`
	Canvas        canvas  `json:"canvas"`
	Walls         []wall  `json:"walls"`
	Places        []place `json:"places"`
}

type canvas struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type wall struct {
	ID       string `json:"id"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Rotation int    `json:"rotation"`
}

type place struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Rotation int    `json:"rotation"`
}
