package storage

type Task struct {
	text             string
	completed        bool
	project          string
	contexts         []string
	previousChecksum string
	currentChecksum  string
}
