package domain

type Leaderboard struct {
	GlobalAverage  float64 `json:"globalAverage"`
	TotalUser      int     `json:"totalUser"`
	UserPercentile float64 `json:"userPercentile"`
	UserRank       int     `json:"userRank"`
	UserTotalPower float64 `json:"userTotalPower"`
}
