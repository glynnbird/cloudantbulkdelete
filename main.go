package main

import (
	"github.com/glynnbird/cloudantbulkdelete/bulkdelete"
)

func main() {

	// create cloudant bulk delete
	cloudantBulkDelete, err := bulkdelete.New()
	if err != nil {
		panic(err)
	}

	// run it
	err = cloudantBulkDelete.Run()
	if err != nil {
		panic(err)
	}
}
