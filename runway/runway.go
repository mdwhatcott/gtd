package runway

import "github.com/mdwhatcott/gtd/projects"

func GroupContextListings(input []projects.Project) map[string][]projects.Task {
	listings := map[string][]projects.Task{}
	for _, project := range input {
		for _, task := range project.Tasks {
			if !task.Complete {
				for _, context := range task.Contexts {
					listings[context] = append(listings[context], task)
				}
			}
		}
	}
	return listings
}

func IdentifyCompletedTasks(input []projects.Project) (tasks []projects.Task) {
	for _, project := range input {
		for _, task := range project.Tasks {
			if task.Complete {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

func IdentifyStalledProjects(input []projects.Project) (listing []projects.Project) {
	for _, project := range input {
		hasIncompleteTasks := false
		stalled := true
		for _, task := range project.Tasks {
			if task.Complete {
				continue
			}
			hasIncompleteTasks = true
			if len(task.Contexts) > 0 {
				stalled = false
			}
		}

		if stalled && hasIncompleteTasks {
			listing = append(listing, project)
		}
	}
	return listing
}

func IdentifyFinishedProjects(input []projects.Project) (listing []projects.Project) {
	for _, project := range input {
		hasIncompleteTasks := false
		for _, task := range project.Tasks {
			if !task.Complete {
				hasIncompleteTasks = true
			}
		}
		if !hasIncompleteTasks {
			listing = append(listing, project)
		}
	}
	return listing
}
