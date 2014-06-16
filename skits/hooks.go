package skits

import (
	"github.com/jmoiron/modl"
	"time"
)

// PreInsert sets the Created and Modified time before Skit is saved.
func (s *Skit) PreInsert(modl.SqlExecutor) error {
	if s.Created.IsZero() {
		s.Created = time.Now()
	}
	if s.Modified.IsZero() {
		s.Modified = time.Now()
	}
	return nil
}

// PreUpdate updates the Modified time before Skit is updated.
func (s *Skit) PreUpdate(modl.SqlExecutor) error {
	s.Modified = time.Now()
	return nil
}
