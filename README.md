# MWS API

Usage
```
package main

import "fmt"

func main() {
	api := newAmazonMWSAPI(Seller{
		ID:        "", // SellerID
		AccessKey: "",
		SecretKey: "",
	})

	result, err := api.ListInventorySupply([]string{"KH-VGFD-TGLN", "P0-PDD6-VHFT", "D7-DOTV-ID1S", "AJ-UP55-SU54"})

	if err != nil {
		fmt.Println("error !")
		fmt.Println(err)
	}
	fmt.Println(result)

	address := Address{Name: "JL", Line1: "my street address", City: "City", CountryCode: "FR", PostalCode: "999999"}
	item := FulfillmentItem{SellerSKU: "KH-VGFD-TGLN", Quantity: 1}
	items := []FulfillmentItem{item}

	requestID, err := api.CreateFulfillmentOrder("test-14", items, address)
	if err != nil {
		fmt.Println("error !")
		fmt.Println(err)
	}
	fmt.Println(requestID)
}

```
