package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	_ "github.com/krostar/nebulo-client-desktop/validator" // used to init custom validators before using them
	"github.com/krostar/nebulo-golib/tools"
)

type globalOptions struct {
	Logging logOptions `json:"log"`
}

type logOptions struct {
	Verbose string `json:"verbose" validate:"regexp=^(quiet|critical|error|warning|info|request|debug)?$"`
	File    string `json:"file" validate:"file=omitempty+writable"`
}

type runOptions struct {
	TLS          TLSOptions `json:"tls"`
	BaseURL      string     `json:"baseurl" validate:"string=nonempty"`
	ContactsFile string     `json:"contacts_file" validate:"string=nonempty"`
}

// TLSOptions store required TLS options
type TLSOptions struct {
	Key           string `json:"key" validate:"file=omitempty+readable"`
	KeyPassword   string `json:"key_password"`
	ClientsCACert string `json:"clients_ca_cert" validate:"file=readable"`
	Cert          string `json:"cert" validate:"string=nonempty"`
}

// Options list all the available configurations
type Options struct {
	Global globalOptions `json:"global"`
	Run    runOptions    `json:"run"`
}

var (
	// Config store the active merge configuration
	Config = &Options{}
	// CLI store the configuration fetched from the console line
	CLI = &Options{}
	// File store the configuration fetched from an optional file
	File = &Options{}

	// Filepath is the path of the loaded configuration file
	Filepath string
)

// LoadFile fill config.File with the configuration parsed from path
func LoadFile(filepath string) (err error) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("unable to read file %q: %v", filepath, err)
	}

	if err = json.Unmarshal(raw, File); err != nil {
		return fmt.Errorf("unable to parse json file: %v", err)
	}

	Filepath = filepath
	return nil
}

// SaveFile save the current configuration to a file
func SaveFile() (err error) {
	if Filepath == "" {
		return nil
	}
	conf, err := json.MarshalIndent(Config, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to create json: %v", err)
	}
	if err := ioutil.WriteFile(Filepath, conf, 0600); err != nil {
		return fmt.Errorf("unable to write configuration file %q: %v", Filepath, err)
	}
	return nil
}

// Merge fill config.Config based on config.CLI and config.File
// File < CLI
func Merge() {
	mergeRecursive(reflect.ValueOf(CLI).Elem(), reflect.ValueOf(File).Elem(), reflect.ValueOf(Config).Elem())
}

func mergeRecursive(cli, file, config reflect.Value) {
	switch config.Kind() {
	case reflect.Struct: // nested struct, we want to go deeper
		for i := 0; i < config.NumField(); i++ {
			mergeRecursive(cli.Field(i), file.Field(i), config.Field(i))
		}
	default: // everything else, we want to copy/merge
		if !tools.IsZeroOrNil(cli) && cli.String() != "" {
			config.Set(cli)
		} else if !tools.IsZeroOrNil(file) && file.String() != "" {
			config.Set(file)
		}
	}
}
