# MWS BRIDGE

Use the [mws](https://developer.amazonservices.com/) api in a simpler way.

* [API](https://developer.amazonservices.fr/gp/mws/docs.html)

## prerequisite

* golang 1.7
* docker 1.12

## build

`make`

## run

```
export MWSSellerID=zzzz
export MWSAccessKey=yyyy
export MWSSecretKey='xxxx'
make run
```

## usage via http server

### Get stock for given SKU
```
GET http://localhost:8080/stock?skus=KH-VGFD-TGLN,P0-PDD6-VHFT
```
returns
```
[
	{
		SellerSKU: "KH-VGFD-TGLN",
		ASIN: "B01HH3HTFQ",
		Quantity: 16,
		InStockSupplyQuantity: 15
	},
	{
		SellerSKU: "P0-PDD6-VHFT",
		ASIN: "B01H7ALDZ6",
		Quantity: 40,
		InStockSupplyQuantity: 40
	}
]
```
### Sent a package
```
POST /sentViaAmazon
{
	"id":"4567",
	"items" : [{"sku":"GDTE-DJSB-SNSB", quantity: 3}, {"sku":"MSNS-SNSN-KSJW", quantity: 2}],
	"shippingAddress" : {
		"name": "Robert Polka",
		"line1": "88, rue edimbourg",
		"line2": "",
		"city": "Paris",
		"countryCode": "FR",
		"postalCode": "78800"
	},
	"comment":"Thank you !"
}
```
returns
```
{"RequestID":"dalfsdfsdghe"}
```

## Usage via API
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

	requestID, err := api.CreateFulfillmentOrder("test-14", items, address, "Thank you for your order")
	if err != nil {
		fmt.Println("error !")
		fmt.Println(err)
	}
	fmt.Println(requestID)
}

```
