package config

import (
	"fmt"

	"github.com/krostar/nebulo-golib/log"
	validator "gopkg.in/validator.v2"
)

// Apply validate configuration and initialize needed package with
// values from configuration
func Apply() (err error) {
	// check the configuration
	if err = validator.Validate(Config); err != nil {
		return err
	}

	if err = ApplyLoggingOptions(&Config.Global.Logging); err != nil {
		return fmt.Errorf("apply logging configuration failed: %v", err)
	}

	return nil
}

// ApplyLoggingOptions apply configuration on log package
func ApplyLoggingOptions(lc *logOptions) (err error) {
	if lc.Verbose != "" {
		log.Verbosity = log.VerboseMapping[lc.Verbose]
	}
	if lc.File != "" {
		if err = log.SetOutputFile(lc.File); err != nil {
			return fmt.Errorf("unable to set log outputfile: %v", err)
		}
	}
	return nil
}
