package entity

type Userscore struct {
	Mid		int64	`json:"mid"`
	PositiveScore	int32	`json:"positive_score"`
	NegativeScore	int32	`json:"negative_score"`
}
