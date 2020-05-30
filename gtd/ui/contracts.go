package ui

type Editor interface {
	EditTempFile(initialContent string) (resultContent string)
}

const TrackOutcomeTemplate = `# {TITLE}

> {EXPLANATION}


## Actions:

-  [ ] concurrent @home
1. [ ] sequential @home


## Support Materials:`
