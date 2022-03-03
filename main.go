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
				Subcommands: []*cli.Command{
					{
						Name:    "production",
						Aliases: []string{"p"},
						Action: func(c *cli.Context) error {
							return trackInfo(c.Context, pkg, credsFile, c.Command.Name)
						},
					},
					{
						Name:    "alpha",
						Aliases: []string{"a"},
						Action: func(c *cli.Context) error {
							return trackInfo(c.Context, pkg, credsFile, c.Command.Name)
						},
					},
					{
						Name:    "beta",
						Aliases: []string{"b"},
						Action: func(c *cli.Context) error {
							return trackInfo(c.Context, pkg, credsFile, c.Command.Name)
						},
					},
					{
						Name:    "list",
						Usage:   "List all the tracks",
						Aliases: []string{"l"},
						Action: func(c *cli.Context) error {
							return listTracks(c.Context, pkg, credsFile)
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func buildService(ctx context.Context, pkg string, credsFile string) (*androidpublisher.Service, string, error) {
	androidpublisherService, err := androidpublisher.NewService(ctx, option.WithCredentialsFile(credsFile))
	if err != nil {
		return nil, "", err
	}

	var appEdit androidpublisher.AppEdit
	editCreate := androidpublisherService.Edits.Insert(pkg, &appEdit)
	editCreateRsp, err := editCreate.Do()
	if err != nil {
		return nil, "", err
	}
	return androidpublisherService, editCreateRsp.Id, nil
}

func marshallAndPrint(v interface{}) error {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}

func trackInfo(ctx context.Context, pkg string, credsFile string, track string) error {
	service, editId, err := buildService(ctx, pkg, credsFile)
	if err != nil {
		return err
	}

	c := service.Edits.Tracks.Get(pkg, editId, track)
	rsp, err := c.Do()
	if err != nil {
		return err
	}

	if len(rsp.Releases) != 1 {
		log.Fatal("Should have received one response")
	}
	return marshallAndPrint(rsp.Releases[0])
}

func listTracks(ctx context.Context, pkg string, credsFile string) error {
	service, editId, err := buildService(ctx, pkg, credsFile)
	if err != nil {
		return err
	}

	c := service.Edits.Tracks.List(pkg, editId)
	rsp, err := c.Do()
	if err != nil {
		return err
	}

	trackNames := []string{}
	for _, track := range rsp.Tracks {
		trackNames = append(trackNames, track.Track)
	}
	return marshallAndPrint(trackNames)
}
