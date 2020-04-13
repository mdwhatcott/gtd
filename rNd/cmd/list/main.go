package main

import (
	"os"

	"github.com/mdwhatcott/gtd/rNd"
)

func main() {
	var (
		handler    = rNd.NewTotallyFakeQueryHandler()
		controller = rNd.NewController(handler)
		model      = rNd.NewListProjectsInputModel()
		result     = rNd.NewCLIResult(os.Stdout, os.Stderr)
	)
	controller.ListProjects(model)
	model.Render(result)
}
