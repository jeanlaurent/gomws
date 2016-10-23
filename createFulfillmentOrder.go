package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Create FulfillmentOrder
//https://docs.developer.amazonservices.com/en_US/fba_outbound/FBAOutbound_CreateFulfillmentOrder.html

//Address structure Input to fullfillment
type Address struct {
	Name        string
	Line1       string
	Line2       string
	City        string
	CountryCode string
	PostalCode  string
}

//FulfillmentItem Input param
type FulfillmentItem struct {
	SellerSKU string
	Quantity  int
}

//CreateFulfillmentOrder create a fullfillment order in MWS
func (api AmazonMWSAPI) CreateFulfillmentOrder(orderID string, items []FulfillmentItem, destinationAddress Address, orderComment string) (string, error) {
	params := make(map[string]string)

	params["MarketplaceId"] = marketPlaceFR
	params["SellerFulfillmentOrderId"] = orderID
	params["FulfillmentAction"] = "Hold"
	params["DisplayableOrderId"] = orderID
	params["DisplayableOrderDateTime"] = time.Now().UTC().Format(ISO8601)
	params["DisplayableOrderComment"] = orderComment
	params["ShippingSpeedCategory"] = "Standard"
	params["DestinationAddress.Name"] = destinationAddress.Name
	params["DestinationAddress.Line1"] = destinationAddress.Line1
	params["DestinationAddress.Line2"] = destinationAddress.Line2
	params["DestinationAddress.City"] = destinationAddress.City
	params["DestinationAddress.CountryCode"] = destinationAddress.CountryCode
	params["DestinationAddress.PostalCode"] = destinationAddress.PostalCode

	for index, item := range items {
		params[fmt.Sprintf("Items.member.%d.SellerSKU", (index+1))] = item.SellerSKU
		params[fmt.Sprintf("Items.member.%d.Quantity", (index+1))] = strconv.Itoa(item.Quantity)
		params[fmt.Sprintf("Items.member.%d.SellerFulfillmentOrderItemId", (index+1))] = fmt.Sprintf("%s.%d", orderID, (index + 1))
	}

	var xmlstring, err = api.mwsCall("/FulfillmentOutboundShipment", "CreateFulfillmentOrder", params)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(xmlstring, "<ErrorResponse") {
		return "", errors.New(stripErrorMessage(xmlstring))
	}
	requestID, err := stripSingleString(xmlstring, "RequestId")
	if err != nil {
		return "", errors.New("Can't find requestID")
	}
	return requestID, nil
}
