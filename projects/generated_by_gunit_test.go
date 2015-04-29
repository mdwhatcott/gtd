//////////////////////////////////////////////////////////////////////////////
// Generated Code ////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

package projects

import (
	"testing"

	"github.com/smartystreets/gunit"
)

///////////////////////////////////////////////////////////////////////////////

func TestChecksummerFixture(t *testing.T) {
	fixture := gunit.NewFixture(t)
	defer fixture.Finalize()

	test0 := &ChecksummerFixture{Fixture: fixture}
	test0.RunTestCase__(test0.TestFirstListingIsNeverDirty, "Test first listing is never dirty")

	test1 := &ChecksummerFixture{Fixture: fixture}
	test1.RunTestCase__(test1.TestSameListingIsNotDirty, "Test same listing is not dirty")

	test2 := &ChecksummerFixture{Fixture: fixture}
	test2.RunTestCase__(test2.TestDifferentListingSecondIsDirty, "Test different listing second is dirty")

	test3 := &ChecksummerFixture{Fixture: fixture}
	test3.RunTestCase__(test3.TestRemovedFileIsDirty, "Test removed file is dirty")
}

func (self *ChecksummerFixture) RunTestCase__(test func(), description string) {
	self.Describe(description)
	self.Setup()
	test()
}

///////////////////////////////////////////////////////////////////////////////

func init() {
	gunit.Validate("b74169e708cb3f495b872d98fde7271a")
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////// Generated Code //
///////////////////////////////////////////////////////////////////////////////
