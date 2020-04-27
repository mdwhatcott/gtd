package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/wireup"
	"github.com/mdwhatcott/gtd/gtd/storage/eventstore"
)

func main() {
	// STARTUP:
	//   for each file in event store folder
	//     apply events to aggregate to build in-memory repo
	//     apply events to projection to build in-memory repo
	//   for each file in projection folder
	//     unmarshal projection into review queue
	//   for each projection in review queue
	//     compare with canonical projection loaded from event store
	//     if different, publish (apply and store) events representing the diff
	//   All stored aggregates and projections should be up to date w/ actual on disk

	storage := eventstore.NewReadWriter(eventstore.Dependencies{
		// TODO: storage wireup...
	})

	command := commands.TrackOutcome{Title: "App Finished"}

	wireup.BuildHandler(storage, storage).Handle(command)

	fmt.Println("ID: ", command.Result.ID)
	fmt.Println("Err:", command.Result.Error)
}
