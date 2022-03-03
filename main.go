package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"

	"github.com/urfave/cli/v2"
)

func main() {
	var credsFile string
	var pkg string

	app := &cli.App{
		Name:  "Android Publisher",
		Usage: "Access the Android Publisher API",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cred",
				Aliases:     []string{"c"},
				Value:       "google-play-api-key.json",
				Usage:       "Load Google Play API Credentials from `FILE`",
				Destination: &credsFile,
				EnvVars:     []string{"GOOGLE_PLAY_API_CREDENTIALS_FILE"},
			},
			&cli.StringFlag{
				Name:        "package",
				Aliases:     []string{"p"},
				Usage:       "App Package ID",
				Destination: &pkg,
				EnvVars:     []string{"APB_PACKAGE"},
				Required:    true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "trackInfo",
				Usage: "Display Info about a track",
				Action: func(c *cli.Context) error {
					return trackInfo(c.Context, pkg, credsFile)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func trackInfo(ctx context.Context, pkg string, credsFile string) error {
	androidpublisherService, err := androidpublisher.NewService(ctx, option.WithCredentialsFile(credsFile))
	if err != nil {
		return err
	}

	var appEdit androidpublisher.AppEdit
	editCreate := androidpublisherService.Edits.Insert(pkg, &appEdit)
	editCreateRsp, err := editCreate.Do()
	if err != nil {
		return err
	}
	editId := editCreateRsp.Id

	c := androidpublisherService.Edits.Tracks.Get(pkg, editId, "production")
	rsp, err := c.Do()
	if err != nil {
		return err
	}

	if len(rsp.Releases) != 1 {
		log.Fatal("Should have received one response")
	}
	release := rsp.Releases[0]

	json, err := json.MarshalIndent(release, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(json))

	return nil
}
