package projects

import "os"

type Checksummer struct {
	previous int64
}

func NewChecksummer() *Checksummer {
	return &Checksummer{previous: -1}
}

func (self *Checksummer) IsDirty(listing []os.FileInfo) bool {
	if self.previous == -1 {
		defer self.remember(self.sum(listing))
		return false
	}

	current := self.sum(listing)
	defer self.remember(current)
	return current != self.previous
}

func (self *Checksummer) remember(current int64) {
	self.previous = current
}

func (self *Checksummer) sum(listing []os.FileInfo) int64 {
	var total int64
	for _, item := range listing {
		if item.IsDir() {
			continue
		}
		total += item.Size()
		total += item.ModTime().Unix()
	}
	return total
}
