package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: packageId credentialsFilePath")
		os.Exit(1)
	}
	pkg := args[0]
	credsFile := args[1]

	ctx := context.Background()
	androidpublisherService, err := androidpublisher.NewService(ctx, option.WithCredentialsFile(credsFile))
	if err != nil {
		log.Fatal(err)
	}

	var appEdit androidpublisher.AppEdit
	c1 := androidpublisherService.Edits.Insert(pkg, &appEdit)
	rsp, err := c1.Do()
	if err != nil {
		log.Fatal(err)
	}

	editId := rsp.Id

	c3 := androidpublisherService.Edits.Tracks.Get(pkg, editId, "production")
	rsp3, err := c3.Do()
	if err != nil {
		log.Fatal(err)
	}

	if len(rsp3.Releases) != 1 {
		log.Fatal("Should have received one response")
	}
	release := rsp3.Releases[0]

	json, err := json.MarshalIndent(release, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
}
