package main

import (
	"encoding/json"
	"fmt"

	"github.com/BSick7/sumologic-sdk-go/api"
)

func main() {
	session := api.DefaultSession()
	session.Discover()
	client := api.NewClient(session)

	collectors, err := client.Collectors().List(0, 5)
	raw, _ := json.Marshal(collectors)
	fmt.Println(string(raw), err)

	sources, err := client.Collectors().Sources(collectors[0].ID).List()
	raw, _ = json.Marshal(sources)
	fmt.Println(string(raw), err)

	/*
	collector, err := client.Collectors().Create(&api.CollectorCreate{
		CollectorType: "Hosted",
		Name:          "abc-collector",
		Description:   "",
		Category:      "",
	})
	fmt.Printf("%+v\n%s\n", collector, err)
	*/
}
