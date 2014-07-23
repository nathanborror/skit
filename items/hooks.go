package items

import (
	"github.com/jmoiron/modl"
	"time"
)

// PreInsert sets the Created and Modified time before Item is saved.
func (i *Item) PreInsert(modl.SqlExecutor) error {
	if i.Created.IsZero() {
		i.Created = time.Now()
	}
	if i.Modified.IsZero() {
		i.Modified = time.Now()
	}
	return nil
}

// PreUpdate updates the Modified time before Item is updated.
func (i *Item) PreUpdate(modl.SqlExecutor) error {
	i.Modified = time.Now()
	return nil
}
