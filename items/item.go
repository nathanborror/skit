package items

import "time"

// Item defines a blob
type Item struct {
	Hash     string    `json:"hash"`
	Parent   string    `json:"parent"`
	Root     string    `json:"root"`
	User     string    `json:"user"`
	Text     string    `json:"text"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// ChildCount returns children count
func (i Item) ChildCount() int {
	results, err := repo.ListWithParent(i.Hash)
	if err != nil {
		return 0
	}
	return len(results)
}
