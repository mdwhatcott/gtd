package main

import (
	"io"
	"log"
	"os"

	"github.com/mdwhatcott/gtd/gtd/storage/csv"
	"github.com/mdwhatcott/gtd/gtd/storage/json"
)

func main() {
	log.SetFlags(log.Lshortfile)

	log.Println(os.Remove("/Users/mike/src/github.com/mdwhatcott/gtd/events.csv"))

	reader := open("/Users/mike/src/github.com/mdwhatcott/gtd/events.json")
	defer reader.Close()

	writer := open("/Users/mike/src/github.com/mdwhatcott/gtd/events.csv")
	defer writer.Close()

	decoder := json.NewDecoder(reader, json.Registry())
	encoder := csv.NewEncoder(writer, csv.EncoderRegistry())

	x := 0
	for ; ; x++ {
		event, err := decoder.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		err = encoder.Encode(event)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Converted:", x)
}

func open(_path string) *os.File {
	FILE, ERR := os.OpenFile(_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if ERR != nil {
		panic(ERR)
	}
	return FILE
}
