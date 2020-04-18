package storage

type ListStoredProjectsQuery struct {
	Result []StoredProject
}

type StoredProject struct {
	ID      string
	Name    string
	Outcome string
}
