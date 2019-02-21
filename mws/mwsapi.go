package mws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Seller hold access credentials of the Amazon Seller
type Seller struct {
	ID        string
	AccessKey string
	SecretKey string
}

// to get tracking numbers :
//FulfillmentOutboundShipment
// _GET_AMAZON_FULFILLED_SHIPMENTS_DATA_

// AmazonMWSAPI holds all amazon MWS Api function
type AmazonMWSAPI struct {
	Seller   Seller
	EndPoint string
	Version  string

	client *http.Client
}

//NewAmazonMWSAPI creates an AmazonMWSAPI
func NewAmazonMWSAPI(seller Seller) *AmazonMWSAPI {
	return &AmazonMWSAPI{Seller: seller, EndPoint: "https://mws-eu.amazonservices.com", Version: "2010-10-01", client: http.DefaultClient}
}

func (api AmazonMWSAPI) mwsCall(path string, action string, Parameters map[string]string) (string, error) {
	amazonURL, err := generateAmazonURL(api, path, action, Parameters)
	if err != nil {
		return "", err
	}

	err = signAmazonURLMethod2(amazonURL, api)
	if err != nil {
		return "", err
	}

	fmt.Println(amazonURL)
	fmt.Println("")

	resp, err := api.client.Get(amazonURL.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyAsString := string(body)
	// fmt.Println()
	// fmt.Println(bodyAsString)
	// fmt.Println()
	return bodyAsString, nil
}

func generateAmazonURL(api AmazonMWSAPI, path string, action string, parameters map[string]string) (finalURL *url.URL, err error) {
	version := fmt.Sprintf("/%s", api.Version)
	if path == "" {
		version = ""
	}
	endPointURL, err := url.Parse(fmt.Sprintf("%s%s%s", api.EndPoint, path, version))
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("Action", action)
	values.Add("Version", api.Version)
	values.Add("AWSAccessKeyId", api.Seller.AccessKey)
	values.Add("SignatureVersion", "2")
	values.Add("SignatureMethod", "HmacSHA256")
	values.Add("SellerId", api.Seller.ID)
	values.Add("Timestamp", time.Now().UTC().Format(time.RFC3339))

	for k, v := range parameters {
		if v != "" {
			values.Set(k, v)
		}
	}

	endPointURL.RawQuery = values.Encode()

	return endPointURL, nil
}

func signAmazonURLMethod2(amazonURL *url.URL, api AmazonMWSAPI) error {
	params := amazonURL.Query()

	stringToSign := "GET\n" + amazonURL.Host + "\n" + amazonURL.Path + "\n" + encodeParameterAmazonStyle(params)
	hmac := hmac.New(sha256.New, []byte(api.Seller.SecretKey))
	_, err := hmac.Write([]byte(stringToSign))
	if err != nil {
		return err
	}
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))

	params.Set("Signature", signature)
	amazonURL.RawQuery = encodeParameterAmazonStyle(params)
	return nil
}

// Amazon wants space in query param to be encoded with %20 instead of +
// http://docs.aws.amazon.com/general/latest/gr/signature-version-2.html
// go urlencode does not do that and won't fix it...
// https://github.com/golang/go/issues/4013 Life is good
func encodeParameterAmazonStyle(params url.Values) string {
	encodedValues := params.Encode()
	return strings.Replace(encodedValues, "+", "%20", -1)
}
