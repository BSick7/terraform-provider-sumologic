package main

import (
	"encoding/json"
	"fmt"
	"github.com/BSick7/sumologic-sdk-go/api"
)

func main() {
	session := api.DefaultSession()
	client := api.NewClient(session)
	client.Discover()

	collectors, err := client.Collectors().List(0, 5)
	raw, _ := json.Marshal(collectors)
	fmt.Println(string(raw), err)

	sources, err := client.Collectors().Sources(collectors[0].ID).List()
	raw, _ = json.Marshal(sources)
	fmt.Println(string(raw), err)
}
