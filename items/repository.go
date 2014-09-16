package items

// ItemRepository holds all the methods needed to save, delete, load and list User objects.
type ItemRepository interface {
	Load(hash string) (*Item, error)
	Delete(hash string) error
	Save(item *Item) error
	List(limit int) ([]*Item, error)
	ListWithParent(parent string) ([]*Item, error)
	ListWithUser(user string) ([]*Item, error)
	ListParents(hash string) ([]*Item, error)
	Archive(item *Item) error
	UnArchive(item *Item) error
}
