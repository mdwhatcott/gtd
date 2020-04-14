package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash/fnv"
)

func main() {
	seen := make(map[string]struct{})
	for x := 0; ; x++ {
		id := identifier()
		if len(seen)%100_000 == 0 {
			fmt.Println(id, len(seen))
		}
		_, found := seen[id]
		if found {
			fmt.Println("Collision:", id, len(seen))
			break
		}
		seen[id] = struct{}{}
	}
}

func identifier() string {
	h := fnv.New64()
	b := make([]byte, 1024)
	_, _ = rand.Read(b)
	_, _ = h.Write(b)
	sum := h.Sum(nil)
	id := hex.EncodeToString(sum)
	return fmt.Sprintf("%s-%s-%s", id[:6], id[6:10], id[10:])
}
