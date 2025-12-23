package bulkdelete

import (
	"errors"
	"flag"
)

// AppConfig contains the application configuration collected from command-line flags
type AppConfig struct {
	DatabaseName   string
	SelectorString string
}

// NewAppConfig creates a new AppConfig struct by parsing the incoming command-line parameters
// using the flag module.
func NewAppConfig() (*AppConfig, error) {
	appConfig := AppConfig{}

	// parse command-line options
	flag.StringVar(&appConfig.DatabaseName, "db", "", "The Cloudant database name to write to")
	flag.StringVar(&appConfig.SelectorString, "selector", "", "The selector that defines the slice of data to delete")
	flag.Parse()

	// if we don't have a database name after parsing
	if appConfig.DatabaseName == "" {
		return nil, errors.New("missing db - please supply the Cloudant database name to delete from")
	} else if appConfig.SelectorString == "" {
		return nil, errors.New("missing selector parameter - please supply Cloudant Query selector")
	}
	return &appConfig, nil
}
