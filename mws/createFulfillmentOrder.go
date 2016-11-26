package mws

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
	Name        string `json:"name"`
	Line1       string `json:"line2"`
	Line2       string `json:"line1"`
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
	PostalCode  string `json:"postalCode"`
}

//FulfillmentItem Input param
type FulfillmentItem struct {
	SellerSKU string `json:"sku"`
	Quantity  int    `json:"quantity"`
}

// Order input Param
type Order struct {
	ID              string            `json:"id"`
	Items           []FulfillmentItem `json:"items"`
	ShippingAddress Address           `json:"shippingAddress"`
	Comment         string            `json:"comment"`
}

type FullfillmentReturn struct {
	RequestID string `json:"requestID"`
}

//CreateFulfillmentOrder create a fullfillment order in MWS
func (api AmazonMWSAPI) CreateFulfillmentOrder(order Order) (FullfillmentReturn, error) {
	params := make(map[string]string)

	params["MarketplaceId"] = marketPlaceFR
	params["SellerFulfillmentOrderId"] = order.ID
	params["FulfillmentAction"] = "Hold"
	params["DisplayableOrderId"] = order.ID
	params["DisplayableOrderDateTime"] = time.Now().UTC().Format(ISO8601)
	params["DisplayableOrderComment"] = order.Comment
	params["ShippingSpeedCategory"] = "Standard"
	params["DestinationAddress.Name"] = order.ShippingAddress.Name
	params["DestinationAddress.Line1"] = order.ShippingAddress.Line1
	params["DestinationAddress.Line2"] = order.ShippingAddress.Line2
	params["DestinationAddress.City"] = order.ShippingAddress.City
	params["DestinationAddress.CountryCode"] = order.ShippingAddress.CountryCode
	params["DestinationAddress.PostalCode"] = order.ShippingAddress.PostalCode

	for index, item := range order.Items {
		params[fmt.Sprintf("Items.member.%d.SellerSKU", (index+1))] = item.SellerSKU
		params[fmt.Sprintf("Items.member.%d.Quantity", (index+1))] = strconv.Itoa(item.Quantity)
		params[fmt.Sprintf("Items.member.%d.SellerFulfillmentOrderItemId", (index+1))] = fmt.Sprintf("%s.%d", order.ID, (index + 1))
	}

	var xmlstring, err = api.mwsCall("/FulfillmentOutboundShipment", "CreateFulfillmentOrder", params)
	if err != nil {
		return FullfillmentReturn{}, err
	}
	if strings.HasPrefix(xmlstring, "<ErrorResponse") {
		return FullfillmentReturn{}, errors.New(stripErrorMessage(xmlstring))
	}
	requestID, err := stripSingleString(xmlstring, "RequestId")
	if err != nil {
		return FullfillmentReturn{}, errors.New("Can't find requestID")
	}
	return FullfillmentReturn{requestID}, nil
}
