package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/krostar/nebulo-client-desktop/api"
	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/gui"
	"github.com/krostar/nebulo-golib/log"

	cli "gopkg.in/urfave/cli.v2"
)

var (
	// BuildTime represent the time when the binary has been created
	BuildTime = "undefined"
	// BuildVersion is the version of the binary (git tag or revision)
	BuildVersion = "undefined"
)

func main() {
	app := &cli.App{
		Name:        "Nebulo desktop client",
		Usage:       "encrypted chat",
		HideVersion: true,
		Before: func(c *cli.Context) (err error) {
			if err = config.ApplyLoggingOptions(&config.Config.Global.Logging); err != nil {
				return fmt.Errorf("unable to apply logging configuration: %v", err)
			}
			if configFile := c.String("config"); configFile != "" {
				if err = config.LoadFile(configFile); err != nil {
					return fmt.Errorf("unable to load configuration file %q:%v", configFile, err)
				}
			}
			return nil
		}, Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to the configuration file",
			}, &cli.StringFlag{
				Name:        "log",
				Aliases:     []string{"l"},
				Usage:       "path to a file where the logs will be writted",
				DefaultText: "standart output",
				Destination: &config.CLI.Global.Logging.File,
			}, &cli.StringFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "level of informations to write (quiet, critical, error, warning, info, request, debug)",
				DefaultText: "debug",
				Destination: &config.CLI.Global.Logging.Verbose,
			},
		}, Commands: []*cli.Command{
			&cli.Command{
				Name:  "run",
				Usage: "start the client",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "baseurl",
						Aliases:     []string{"b"},
						Usage:       "base url to use to contact API server",
						Destination: &config.CLI.Run.BaseURL,
					}, &cli.StringFlag{
						Name:        "tls-crt",
						Usage:       "* tls certificate file used to encrypt communication (https)",
						Destination: &config.CLI.Run.TLS.Cert,
					}, &cli.StringFlag{
						Name:        "tls-key",
						Usage:       "* tls certificate key used with --tls-crt",
						Destination: &config.CLI.Run.TLS.Key,
					}, &cli.StringFlag{
						Name:        "tls-clients-ca",
						Usage:       "* tls certification authority used to validate clients certificate for the tls mutual authentication",
						Destination: &config.CLI.Run.TLS.ClientsCACert,
					},
				}, Before: beforeCommandWhoNeedMergeConfiguration,
				Action: commandRun,
			}, &cli.Command{
				Name:  "config-gen",
				Usage: "generate a configuration file and quit",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "destination",
						Aliases:     []string{"d"},
						Usage:       "path to a file where the configuration will be writted",
						DefaultText: "standart output",
					},
				}, Before: beforeEveryCommand,
				Action: commandConfigGen,
			}, &cli.Command{
				Name:   "version",
				Usage:  "display the version",
				Before: beforeEveryCommand,
				Action: commandVersion,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Criticalf("unable to run app: %v", err)
		os.Exit(1)
	}
}

func beforeEveryCommand(c *cli.Context) (err error) {
	if c.NArg() != 0 {
		return fmt.Errorf("unknown remaining args: %q", strings.Join(c.Args().Slice(), " "))
	}
	return nil
}

func beforeCommandWhoNeedMergeConfiguration(c *cli.Context) (err error) {
	if err = beforeEveryCommand(c); err != nil {
		return err
	}
	config.Merge()
	if err = config.Apply(); err != nil {
		return fmt.Errorf("configuration application failed: %v", err)
	}
	log.Logf(log.DEBUG, -1, "Configuration merged, validated and applied: %v", config.Config)
	return nil
}

func commandRun(_ *cli.Context) error {
	log.Infof("Starting Nebulo client build %s (%s): %s", BuildVersion, BuildTime, config.Config.Run.BaseURL)

	_, version, err := api.New(BuildVersion, config.Config.Run.BaseURL, &config.Config.Run.TLS)
	if err != nil {
		return fmt.Errorf("unable to initialize API client: %v", err)
	}
	log.Infof("Using server API %q version: %s (%s)", config.Config.Run.BaseURL, version.Version, version.Time)

	return gui.GUI()
}

func commandConfigGen(c *cli.Context) error {
	conf, err := json.MarshalIndent(config.Config, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to create json: %v", err)
	}
	if filepath := c.String("destination"); filepath != "" {
		if err := ioutil.WriteFile(filepath, conf, 0644); err != nil {
			return fmt.Errorf("unable to write sql queries file: %v", err)
		}
	} else {
		fmt.Println(string(conf))
	}
	return nil
}

func commandVersion(_ *cli.Context) error {
	fmt.Printf("nebulo %s (%s)\n", BuildVersion, BuildTime)
	return nil
}
