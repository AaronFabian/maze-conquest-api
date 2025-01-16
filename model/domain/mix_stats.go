package domain

type MixStats struct {
	Uid   string `json:"uid" firestore:"uid"`
	Power int    `json:"power" firestore:"power"`
}
