package storage

type ProjectLoader interface {
	List() (names []string)
	Load(name string) (*Project, error)
	LoadFromTaskChecksum(checksum string) (*Project, error)
	Save(project *Project) error
	Delete(name string) error
}
