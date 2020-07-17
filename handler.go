package main

import (
	"encoding/json"
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

func (fe *frontendServer) sendMailHandler(w http.ResponseWriter, r *http.Request) {

}

//Recommendation service handlers

func (fe *frontendServer) getRecommendationsHandler(w http.ResponseWriter, r *http.Request) {

}

//Cart service handlers

func (fe *frontendServer) getCartHandler(w http.ResponseWriter, r *http.Request) {

}

func (fe *frontendServer) emptyCartHandler(w http.ResponseWriter, r *http.Request) {

}

func (fe *frontendServer) addToCartHandler(w http.ResponseWriter, r *http.Request) {

}

//Checkout service handlers

func (fe *frontendServer) checkoutRequestHandler(w http.ResponseWriter, r *http.Request) {

}
