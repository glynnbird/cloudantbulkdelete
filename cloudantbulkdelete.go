package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/IBM/cloudant-go-sdk/features"
)

const changesBatchSize int64 = 10000

type CloudantBulkDelete struct {
	appConfig *AppConfig             // our command-line options
	service   *cloudantv1.CloudantV1 // the Cloudant SDK client
}

type outputObject struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev"`
	Deleted bool   `json:"_deleted"`
}

// NewCloudantBulkDelete creates a new CloudantBulkDelete struct, loading the CLI parameters
// and instantiating the Cloudant SDK client
func NewCloudantBulkDelete() (*CloudantBulkDelete, error) {
	// load the CLI parameters
	appConfig, err := NewAppConfig()
	if err != nil {
		return nil, err
	}

	// set up the Cloudant service
	service, err := cloudantv1.NewCloudantV1UsingExternalConfig(&cloudantv1.CloudantV1Options{})
	if err != nil {
		return nil, err
	}
	service.EnableRetries(3, 5*time.Second)

	cbd := CloudantBulkDelete{
		appConfig: appConfig,
		service:   service,
	}

	return &cbd, nil
}

func (cbd *CloudantBulkDelete) Run() error {

	// Required: the database name.
	postChangesOptions := cbd.service.NewPostChangesOptions(cbd.appConfig.DatabaseName)
	postChangesOptions.SetLimit(changesBatchSize)
	postChangesOptions.SetSince("0")
	postChangesOptions.SetFilter("_selector")

	// set the selector
	var selector map[string]interface{}
	err := json.Unmarshal([]byte(cbd.appConfig.SelectorString), &selector)
	if err != nil {
		return err
	}
	postChangesOptions.SetSelector(selector)

	// Required: the Cloudant service client instance and an instance of PostChangesOptions
	follower, err := features.NewChangesFollower(cbd.service, postChangesOptions)
	if err != nil {
		return err
	}

	// start the follower
	changesCh, err := follower.StartOneOff()
	if err != nil {
		return err
	}

	for changesItem := range changesCh {
		// changes item returns an error on failed requests
		item, err := changesItem.Item()
		if err != nil {
			return err
		}

		// do something with changes
		outputObj := outputObject{
			ID:      *item.ID,
			Rev:     *item.Changes[0].Rev,
			Deleted: true,
		}
		outputStr, err := json.Marshal(outputObj)
		if err != nil {
			return err
		}
		fmt.Println(string(outputStr))
	}
	return nil
}
