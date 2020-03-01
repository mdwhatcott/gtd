package main

import (
	"os"

	"github.com/mdwhatcott/gtd/gtd"
)

func main() {
	var (
		handler    = gtd.NewTotallyFakeQueryHandler()
		controller = gtd.NewController(handler)
		model      = gtd.NewListProjectsInputModel()
		result     = gtd.NewCLIResult(os.Stdout, os.Stderr)
	)
	controller.ListProjects(model)
	model.Render(result)
}
