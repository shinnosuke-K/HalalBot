package ocr

import "strings"

type ResultsWords struct {
	BoundingBox [][]float64
	Text        string
	Confidence  float64
	LineBreak   bool
}

type Results struct {
	Words []ResultsWords `json:"words"`
}

var ngFood = []string{"ワイン", "みりん", "日本酒", "ビール", "ラム酒", "料理酒", "豚", "ポーク", "ゼラチン", "ラード"}

func foodContain(text string) bool {
	for _, ng := range ngFood {
		if strings.Contains(text, ng) {
			return true
		}
	}
	return false
}
