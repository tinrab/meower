package db

type MemoryRepository struct {
	meows []*Meow
}

func (r *MemoryRepository) InsertMeow(meow Meow) error {
	r.meows = append(r.meows, &meow)
	return nil
}
