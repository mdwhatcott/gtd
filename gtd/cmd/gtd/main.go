package main

import "fmt"

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
	fmt.Println("Hello, world!")
}
