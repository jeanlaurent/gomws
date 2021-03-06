package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeanlaurent/gomws/mws"
)

// Server is the http server instance
type Server struct {
	Port int
	mws  *mws.AmazonMWSAPI
}

// NewServer creates a new server
func NewServer(port int, seller mws.Seller) Server {
	return Server{Port: port, mws: mws.NewAmazonMWSAPI(seller)}
}

//Start the http server
func (server *Server) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/", errorHandler(server.root)).Methods("GET")
	router.HandleFunc("/stock", errorHandler(server.listInventorySupplyHandlerFunc)).Methods("GET")
	router.HandleFunc("/sentViaAmazon", errorHandler(server.createFulfillmentHandlerFunc)).Methods("POST")

	http.Handle("/", router)

	fmt.Println("listening on", server.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)
}

func (server *Server) root(response http.ResponseWriter, request *http.Request) error {
	response.Write([]byte("MWS bridge"))
	return nil
}

/*
Request:
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
Response:
200 OK
{"requestID":"skdjasfkjsaf"}
*/
func (server *Server) createFulfillmentHandlerFunc(response http.ResponseWriter, request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	var order mws.Order
	err := decoder.Decode(&order)
	if err != nil {
		return err
	}
	defer request.Body.Close()
	requestID, err := server.mws.CreateFulfillmentOrder(order)
	if err != nil {
		return err
	}
	json, err := json.Marshal(requestID)
	if err != nil {
		return err
	}
	response.Header().Set("Content-Type", "application/json")
	response.Write(json)
	return nil
}

/*
Request:
	GET stock?skus=a1,a2
Response:
[
	{
		"sellerSKU": "a1",
		"asin":"aaaaa",
		"totalSupplyQuantity":4,
		"inStockSupplyQuantity":4
	},
	{
		"sellerSKU": "a2",
		"asin":"aaaaa2",
		"totalSupplyQuantity":5,
		"inStockSupplyQuantity":3
	},
]
*/
func (server *Server) listInventorySupplyHandlerFunc(response http.ResponseWriter, request *http.Request) error {
	var skuString = request.URL.Query().Get("skus")
	fmt.Println("--->", skuString)
	var skus = strings.Split(skuString, ",")
	productStocks, err := server.mws.ListInventorySupply(skus)
	if err != nil {
		return err
	}
	fmt.Println(productStocks)
	json, err := json.Marshal(productStocks)
	if err != nil {
		return err
	}
	response.Header().Set("Content-Type", "application/json")
	response.Write(json)
	return nil
}

func errorHandler(handler func(response http.ResponseWriter, request *http.Request) error) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		err := handler(response, request)
		if err != nil {
			fmt.Println(err)
			response.WriteHeader(500)
			response.Header().Add("Content-Type", "application/json")
			response.Write([]byte(fmt.Sprintf("{\"error\":\"%v\"}", err)))
		}
	}
}
