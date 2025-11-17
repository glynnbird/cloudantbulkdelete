package main

func main() {

	cloudantBulkDelete, err := NewCloudantBulkDelete()
	if err != nil {
		panic(err)
	}

	err = cloudantBulkDelete.Run()
	if err != nil {
		panic(err)
	}
}
