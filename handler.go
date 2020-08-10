package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/alessioVnt/frontend/pb"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type platformDetails struct {
	css      string
	provider string
}

var (
	templates = template.Must(template.New("").ParseGlob("templates/*.html"))
	plat      platformDetails
)

//User service handlers

func (fe *frontendServer) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("getUserByID request received\n")
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 32)
	id32 := int32(id)

	fe.getUserByID(r.Context(), id32)
	logrus.Infof("Request completed\n")
}

func (fe *frontendServer) updateUserPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("updateUserPreferences request received\n")
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 32)
	id32 := int32(id)

	newPref := r.FormValue("newPreference")

	fe.updatePreferences(r.Context(), id32, newPref)
	logrus.Infof("Request completed\n")
}

//Restaurant service handlers

func (fe *frontendServer) restaurantListHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("getAllRestaurants request received\n")
	fe.getAllRestaurants(r.Context())
	logrus.Infof("Request completed\n")
}

func (fe *frontendServer) addRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("addRestaurant request received\n")
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
	logrus.Infof("Request completed\n")

	//for key, value := range redTags {
	// Each value is an interface{} type, that is type asserted as a string
	//	fmt.Println(key, value.(string))
	//}

}

//Mail service handlers

func (fe *frontendServer) sendMailHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("sendMail request received\n")
	id := r.FormValue("id")
	tag := r.FormValue("tag")

	fe.sendMail(r.Context(), tag, id)
	logrus.Infof("Request completed\n")
}

//Recommendation service handlers

func (fe *frontendServer) getRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("getRecommendations request received\n")
	id := r.FormValue("id")
	recommendations := fe.getRecommendations(r.Context(), id)

	for _, recommended := range recommendations {
		print(recommended + "\n")
	}
	logrus.Infof("Request completed\n")
}

//Cart service handlers

func (fe *frontendServer) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("getOrder request received\n")
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
	logrus.Infof("Request completed\n")

}

func (fe *frontendServer) emptyCartHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("emptyCart request received\n")
	id := r.FormValue("id")
	fe.emptyCart(r.Context(), id)
	logrus.Infof("Request completed\n")
}

func (fe *frontendServer) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("addToCart request received\n")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		print("error in reading body of http message")
	}
	logrus.Infof("Request completed\n")

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
	logrus.Infof("Request completed\n")
}

//Checkout service handlers

func (fe *frontendServer) checkoutRequestHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("checkout request received\n")
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

	logrus.Infof("Request completed\n")
}

//Utils handler

func (fe *frontendServer) helpPageHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "help", map[string]interface{}{}); err != nil {
		print("Error executing template!")
	}
}
