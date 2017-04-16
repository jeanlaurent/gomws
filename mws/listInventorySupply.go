package mws

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

// ISO8601 Date Pattern
const ISO8601 = "2006-01-02T15:04:05-0700"

// ListInventorySupplyServiceStatus returns the status of the ListInventorySupplyService MWS Api
func (api AmazonMWSAPI) ListInventorySupplyServiceStatus() (string, error) {
	return api.mwsCall("/FulfillmentInventory", "GetServiceStatus", make(map[string]string))
}

// ProductStock represents a Product with it's associated stock
type ProductStock struct {
	SellerSKU             string `xml:"sellerSKU"`
	ASIN                  string `xml:"asin"`
	Quantity              int    `xml:"totalSupplyQuantity"`
	InStockSupplyQuantity int    `xml:"inStockSupplyQuantity"`
}

// ListInventorySupply returns the corresponding product stock for a given list of product skus
func (api AmazonMWSAPI) ListInventorySupply(skus []string) ([]ProductStock, error) {
	params := make(map[string]string)

	for index, sku := range skus {
		key := fmt.Sprintf("SellerSkus.member.%d", (index + 1))
		params[key] = sku
	}

	params["ResponseGroup"] = "Basic"

	var xml, err = api.mwsCall("/FulfillmentInventory", "ListInventorySupply", params)
	if err != nil {
		return []ProductStock{}, err
	}
	if strings.HasPrefix(xml, "<ErrorResponse") {
		return []ProductStock{}, errors.New(stripErrorMessage(xml))
	}
	members := listInventorySupplydecode(xml)
	return members, nil
}

func listInventorySupplydecode(xmlString string) []ProductStock {
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	var members = []ProductStock{}
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch tokenType := token.(type) {
		case xml.StartElement:
			if tokenType.Name.Local == "member" {
				var member ProductStock
				decoder.DecodeElement(&member, &tokenType)
				members = append(members, member)
			}
		}
	}
	return members
}
