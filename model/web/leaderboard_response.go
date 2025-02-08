package web

import "maze-conquest-api/model/domain"

type LeaderboardResponse struct {
	Code   int                 `json:"code"`
	Status string              `json:"status"`
	Data   *domain.Leaderboard `json:"data"`
}
