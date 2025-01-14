package domain

type Hero struct {
	Name       string `json:"name"`
	Level      int    `json:"level"`
	CurrentExp int    `json:"currentExp"`
	ExpToLevel int    `json:"expToLevel"`
}
