package main

import (
	"encoding/json"
	"fmt"
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
		}, Flags: []cli.Flag{ // global flags (config and logs purpose)
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
			&cli.Command{ // run command, she start the client
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
			}, &cli.Command{ // config-gen command, she generate an empty configuration file
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
			}, &cli.Command{ // version command output the current client version
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
	// we don't want remaining (non-parsed args)
	if c.NArg() != 0 {
		return fmt.Errorf("unknown remaining args: %q", strings.Join(c.Args().Slice(), " "))
	}
	return nil
}

func beforeCommandWhoNeedMergeConfiguration(c *cli.Context) (err error) {
	if err = beforeEveryCommand(c); err != nil {
		return err
	}
	// merge cli and file loaded configuration
	config.Merge()
	if err = config.Apply(); err != nil {
		return fmt.Errorf("configuration application failed: %v", err)
	}
	log.Logf(log.DEBUG, -1, "Configuration merged, validated and applied: %v", config.Config)
	return nil
}

func commandRun(_ *cli.Context) error {
	log.Infof("Starting Nebulo client build %s (%s): %s", BuildVersion, BuildTime, config.Config.Run.BaseURL)

	// try to reach the api server
	version, err := api.Initialize(BuildVersion, config.Config.Run.BaseURL, &config.Config.Run.TLS)
	if err != nil {
		return fmt.Errorf("unable to initialize API client: %v", err)
	}
	log.Infof("Using server API %q version: %s (%s)", config.Config.Run.BaseURL, version.Version, version.Time)

	// start the GUI
	return gui.GUI()
}

func commandConfigGen(c *cli.Context) error {
	if filepath := c.String("destination"); filepath != "" {
		config.Filepath = filepath
		if err := config.SaveFile(); err != nil {
			return fmt.Errorf("unable to write sql queries file: %v", err)
		}
	} else {
		conf, err := json.MarshalIndent(config.Config, "", "    ")
		if err != nil {
			return fmt.Errorf("unable to create json: %v", err)
		}
		fmt.Println(string(conf))
	}
	return nil
}

func commandVersion(_ *cli.Context) error {
	fmt.Printf("nebulo %s (%s)\n", BuildVersion, BuildTime)
	return nil
}
