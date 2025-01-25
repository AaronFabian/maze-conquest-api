package domain

type MixStats struct {
	Uid           string `json:"uid" firestore:"uid"`
	Power         int    `json:"power" firestore:"power"`
	OwnerUsername string `json:"ownerUsername" firestore:"ownerUsername"`
	PhotoUrl      string `json:"photoUrl" firestore:"photoUrl"`
}
