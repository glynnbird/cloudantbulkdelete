package main

import (
	"fmt"
	"os"

	"github.com/glynnbird/cloudantbulkdelete/bulkdelete"
)

func main() {

	// create cloudant bulk delete
	cloudantBulkDelete, err := bulkdelete.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// run it
	err = cloudantBulkDelete.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
