package domain

type User struct {
	Uid       string                    `json:"uid"`
	Username  string                    `json:"username"`
	Active    bool                      `json:"active"`
	AllHeroes map[string]map[string]int `json:"allHeroes"`
	Items     map[string]int            `json:"items"`
	Party     []string                  `json:"party"`
	Worlds    map[string]int            `json:"worlds"`
	CreatedAt int64                     `json:"createdAt"`
}
