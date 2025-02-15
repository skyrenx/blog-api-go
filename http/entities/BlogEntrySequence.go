package entities

type BlogEntrySequence struct {
	NextId int `json:"NextId" db:"next_id"`
}
