package skits

// SkitRepository holds all the methods needed to save, delete, load and list User objects.
type SkitRepository interface {
	Load(hash string) (*Skit, error)
	Delete(hash string) error
	Save(skit *Skit) error
	List(limit int) ([]*Skit, error)
	ListWithParent(parent string) ([]*Skit, error)
	ListWithUser(user string) ([]*Skit, error)
}
