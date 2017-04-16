package main

import (
	"fmt"
	"os"

	"github.com/jeanlaurent/gomws/http"
	"github.com/jeanlaurent/gomws/mws"
)

//DEFAULTPORT for http server
const DEFAULTPORT int = 8080

func main() {
	var sellerID = getEnvOrFail("MWSSellerID")
	var accessKey = getEnvOrFail("MWSAccessKey")
	var secretKey = getEnvOrFail("MWSSecretKey")
	var seller = mws.Seller{
		ID:        sellerID,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	var server = http.NewServer(DEFAULTPORT, seller)

	server.Start()
}

func getEnvOrFail(key string) string {
	if os.Getenv(key) == "" {
		fmt.Println("Error : missing environment variable", key)
		os.Exit(-1)
	}
	return os.Getenv(key)
}
