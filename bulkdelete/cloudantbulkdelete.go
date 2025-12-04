package bulkdelete

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/IBM/cloudant-go-sdk/features"
)

// CloudantBulkDelete stores all the data we need to run a cloudantbulkdelete job.
type CloudantBulkDelete struct {
	appConfig *AppConfig             // our command-line options
	service   *cloudantv1.CloudantV1 // the Cloudant SDK client
}

// the format of the output object
type outputObject struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev"`
	Deleted bool   `json:"_deleted"`
}

// New creates a new CloudantBulkDelete struct, loading the CLI parameters
// and instantiating the Cloudant SDK client
func New() (*CloudantBulkDelete, error) {
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

// Run spools through the the chosen database's changes feed with the supplied
// selector filter. It outputs one line per matching document containing the
// _id/_rev pair + _deleted: true to stdout
func (cbd *CloudantBulkDelete) Run() error {

	// Required: the database name.
	postChangesOptions := cbd.service.NewPostChangesOptions(cbd.appConfig.DatabaseName)
	postChangesOptions.SetSince("0")
	postChangesOptions.SetFilter("_selector")

	// set the selector
	var selector map[string]interface{}
	err := json.Unmarshal([]byte(cbd.appConfig.SelectorString), &selector)
	if err != nil {
		return err
	}
	postChangesOptions.SetSelector(selector)

	// create a new changes follower
	follower, err := features.NewChangesFollower(cbd.service, postChangesOptions)
	if err != nil {
		return err
	}

	// start the follower, in one-off mode
	changesCh, err := follower.StartOneOff()
	if err != nil {
		return err
	}

	// range through the changes feed
	for changesItem := range changesCh {

		// changes item returns an error on failed requests
		item, err := changesItem.Item()
		if err != nil {
			continue
		}

		// build the output struct
		outputObj := outputObject{
			ID:      *item.ID,
			Rev:     *item.Changes[0].Rev,
			Deleted: true,
		}

		// output as JSON
		outputStr, err := json.Marshal(outputObj)
		if err != nil {
			return err
		}
		fmt.Println(string(outputStr))
	}
	return nil
}
