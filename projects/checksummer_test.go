package projects

//go:generate gunit

import (
	"os"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ChecksummerFixture struct {
	*gunit.Fixture

	checksummer *Checksummer
}

func (self *ChecksummerFixture) Setup() {
	self.checksummer = NewChecksummer()
}

func (self *ChecksummerFixture) TestFirstListingIsNeverDirty() {
	self.So(self.checksummer.IsDirty(twoFiles), should.BeFalse)
}
func (self *ChecksummerFixture) TestSameListingIsNotDirty() {
	self.checksummer.IsDirty(twoFiles)
	self.So(self.checksummer.IsDirty(twoFiles), should.BeFalse)
}
func (self *ChecksummerFixture) TestDifferentListingSecondIsDirty() {
	self.checksummer.IsDirty(twoFiles)
	self.So(self.checksummer.IsDirty(addedFile), should.BeTrue)
}
func (self *ChecksummerFixture) TestRemovedFileIsDirty() {
	self.checksummer.IsDirty(twoFiles)
	self.So(self.checksummer.IsDirty(removedFile), should.BeTrue)
}

var (
	fileA = FileInfo{
		name:     "a",
		isDir:    false,
		size:     1,
		modified: time.Unix(2, 0),
	}
	fileB = FileInfo{
		name:     "b",
		isDir:    false,
		size:     3,
		modified: time.Unix(4, 0),
	}
	fileC = FileInfo{
		name:     "c",
		isDir:    false,
		size:     5,
		modified: time.Unix(6, 0),
	}
	dir1 = FileInfo{
		name:     "1",
		isDir:    true,
		size:     42,
		modified: time.Now(),
	}

	twoFiles    = []os.FileInfo{dir1, fileA, fileB}
	addedFile   = append(twoFiles, fileC)
	removedFile = twoFiles[:2]
)

////////////////////////////////////////////////////////////

type FileInfo struct {
	name     string
	isDir    bool
	size     int64
	modified time.Time
}

func (self FileInfo) Name() string       { return self.name }
func (self FileInfo) Size() int64        { return self.size }
func (self FileInfo) ModTime() time.Time { return self.modified }
func (self FileInfo) IsDir() bool        { return self.isDir }
func (self FileInfo) Mode() os.FileMode  { return os.FileMode(0) }
func (self FileInfo) Sys() interface{}   { return nil }
