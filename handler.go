package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/alessioVnt/frontend/pb"
	"io/ioutil"
	"net/http"
	"strconv"
)

type platformDetails struct {
	css      string
	provider string
}

//User service handlers

func (fe *frontendServer) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 32)
	id32 := int32(id)

	fe.getUserByID(r.Context(), id32)
}

//Restaurant service handlers

func (fe *frontendServer) restaurantListHandler(w http.ResponseWriter, r *http.Request) {
	fe.getAllRestaurants(r.Context())
}

func (fe *frontendServer) addRestaurantHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		print("error")
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	name := result["name"]
	city := result["city"]
	address := result["address"]
	redTags := result["TAG"].(map[string]interface{})
	tags := &pb.TAG{Tag1: redTags["tag1"].(string), Tag2: redTags["tag2"].(string), Tag3: redTags["tag3"].(string)}

	fe.addRestaurant(r.Context(), name.(string), city.(string), address.(string), tags)

	//for key, value := range redTags {
	// Each value is an interface{} type, that is type asserted as a string
	//	fmt.Println(key, value.(string))
	//}

}

//Mail service handlers
//ToDo: sistemare il mail service, problema con le grpc

func (fe *frontendServer) sendMailHandler(w http.ResponseWriter, r *http.Request) {

}

//Recommendation service handlers

func (fe *frontendServer) getRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	recommendations := fe.getRecommendations(r.Context(), id)

	for _, recommended := range recommendations {
		print(recommended + "\n")
	}
}

//Cart service handlers

func (fe *frontendServer) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	order := fe.getOrder(r.Context(), id)

	//Print the result
	if order.UserId != "" {
		userID := order.UserId
		print("Order user ID: " + userID + " \n")
		restaurantID := order.RestaurantId
		print("From restaurant with ID: " + restaurantID + "\n")
		itemsInOrder := order.Items
		print("Items in the order: \n")
		for i, item := range itemsInOrder {
			print("ITEM NUMBER " + strconv.Itoa(i) + " IN THE ORDER\n")
			print("ProductID: " + item.ProductId + " x" + fmt.Sprint(item.Quantity) + "\n")
			print("With price: " + fmt.Sprint(item.Price) + "\n")
			print("\n")
		}
	}

}

func (fe *frontendServer) emptyCartHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	fe.emptyCart(r.Context(), id)
}

func (fe *frontendServer) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		print("error in reading body of http message")
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	userID := result["userid"].(string)
	restaurantID := result["restaurantid"].(string)

	itemToAdd := result["item"].(map[string]interface{})

	//Converting quantity to in32 from string
	quantity64, _ := strconv.ParseInt(itemToAdd["quantity"].(string), 10, 32)
	quantity := int32(quantity64)

	//Converting price to float32 from string
	price64, _ := strconv.ParseFloat(itemToAdd["price"].(string), 32)
	price := float32(price64)

	cItem := &pb.CartItem{
		ProductId: itemToAdd["productid"].(string),
		Quantity:  quantity,
		Price:     price,
	}

	fe.addToCart(r.Context(), userID, restaurantID, cItem)
}

//Checkout service handlers

func (fe *frontendServer) checkoutRequestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		print("error in reading body of http message")
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	userID := result["userid"].(string)
	restaurantID := result["restaurantid"].(string)

	menuItemsJSN := result["menuitems"].(map[string]interface{})
	var menuItems []*pb.MenuItem

	for _, value := range menuItemsJSN {

		item := value.(map[string]interface{})

		productid := item["productid"].(string)

		//Converting quantity to in32 from string
		quantity64, _ := strconv.ParseInt(item["quantity"].(string), 10, 32)
		quantity := int32(quantity64)

		//Converting price to float32 from string
		price64, _ := strconv.ParseFloat(item["price"].(string), 32)
		price := float32(price64)

		toAdd := &pb.MenuItem{
			ProductId: productid,
			Quantity:  quantity,
			Price:     price,
		}
		menuItems = append(menuItems, toAdd)
	}

	cardNumber := result["cardnumber"].(string)
	cvc := result["cvc"].(string)
	expiration := result["expiration"].(string)

	transactionResult := fe.executeCheckout(r.Context(), userID, restaurantID, menuItems, cardNumber, cvc, expiration)

	if transactionResult {
		print("Transaction successfull!")
	} else {
		print("Transaction failed!")
	}
}
