package main

import (
	"errors"
	"flag"
	"fmt"
)

// AppConfig contains the application configuration collected from command-line flags
type AppConfig struct {
	DatabaseName   string
	SelectorString string
}

func (ac AppConfig) Print() {
	fmt.Println("APP CONFIG")
	fmt.Println("----------")
	fmt.Printf("DatabaseName: %v\n", ac.DatabaseName)
	fmt.Printf("Selector: %v\n", ac.SelectorString)
}

func NewAppConfig() (*AppConfig, error) {
	appConfig := AppConfig{}

	// parse command-line options
	flag.StringVar(&appConfig.DatabaseName, "dbname", "", "The Cloudant database name to write to")
	flag.StringVar(&appConfig.DatabaseName, "db", "", "The Cloudant database name to write to")
	flag.StringVar(&appConfig.SelectorString, "selector", "", "The selector that defines the slice of data to delete")
	flag.StringVar(&appConfig.SelectorString, "s", "", "The selector that defines the slice of data to delete")
	flag.Parse()

	// if we don't have a database name after parsing
	if appConfig.DatabaseName == "" {
		return nil, errors.New("missing dbname/db")
	} else if appConfig.SelectorString == "" {
		return nil, errors.New("missing selector")
	} else {
		return &appConfig, nil
	}
}
